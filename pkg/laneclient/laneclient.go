// LaneClient is a client for interacting with the SwimLane API.
// All responses were learned from observing the API and are not guaranteed to be complete or accurate.

package laneclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type LaneClient struct {
	domain      string
	accountId   string
	accessToken string
	client      *http.Client
	logger      *log.Logger
}

type TenantClient struct {
	lc     LaneClient
	Tenant Tenant
}

type Params *map[string]string

// NewLaneClient returns a new API client.
func NewLaneClient(domain string, accountId string, accessToken string, logger *log.Logger) LaneClient {
	return LaneClient{
		domain:      domain,
		accountId:   accountId,
		accessToken: accessToken,
		client:      &http.Client{Timeout: 2 * time.Minute},
		logger:      logger,
	}
}

// NewTenantClient returns a new TenantClient for the specified tenant.
func NewTenantClient(lc LaneClient, tenant Tenant) TenantClient {
	return TenantClient{
		lc:     lc,
		Tenant: tenant,
	}
}

// urlForAccountEndpoint returns the URL for an endpoint in the account context.
func (lc LaneClient) urlForAccountEndpoint(endpoint string) (string, error) {
	return url.JoinPath("https://", lc.domain, "tenant", "api", "accounts", lc.accountId, endpoint)
}

// urlForTenantEndpoint return the url for an endpoint in the tenant context. Versioned endpoints can be used by specifying a nonzero value for apiVer.
func (tc TenantClient) urlForTenantEndpoint(api string, endpoint string, apiVer uint8) (string, error) {
	if apiVer > 0 {
		return url.JoinPath("https://", tc.lc.domain, api, "api", "account", tc.lc.accountId, "tenant", tc.Tenant.Id, fmt.Sprintf("v%d", apiVer), endpoint)
	}
	return url.JoinPath("https://", tc.lc.domain, "api", "account", tc.lc.accountId, "tenant", tc.Tenant.Id, endpoint)
}

// prepareRequest prepares a new http.request.
func (lc LaneClient) prepareRequest(ctx context.Context, method string, url string, params Params, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add request headers
	req.Header.Add("Private-Token", lc.accessToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Add query parameters
	q := req.URL.Query()
	if params != nil {
		for k, v := range *params {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// sendRequest submits a request and checks the response.
func (lc LaneClient) sendRequest(req *http.Request) ([]byte, error) {
	lc.logger.Debug("Request", "method", req.Method, "url", req.URL.String())

	// Fire request
	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check response
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d for %s", resp.StatusCode, resp.Request.URL)
	}

	return data, nil
}

// pagedRequest incrementally requests data from a paged endpoint.
func (lc LaneClient) pagedRequest(ctx context.Context, req *http.Request, items *[]json.RawMessage) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	q := req.URL.Query()
	q.Set("size", "50")

	// Increment page number in query string
	if page := q.Get("page"); page != "" {
		pageNum, err := strconv.ParseUint(page, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse page number: %w", err)
		}
		if pageNum > 50 {
			return fmt.Errorf("recursion limit reached")
		}
		q.Set("page", strconv.FormatUint(pageNum+1, 10))
	} else {
		q.Add("page", "1")
	}
	req.URL.RawQuery = q.Encode()

	// Fire request
	resp, err := lc.sendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send paged request: %w", err)
	}

	// Decode page response
	page, err := decodeItems[ItemPage](resp)
	if err != nil {
		return fmt.Errorf("failed to decode paged response: %w", err)
	}
	*items = slices.Concat(*items, page[0].Items)

	// if the number of items returned is equal to the page size, assume there's another page and request it.
	if len(page[0].Items) == 50 {
		return lc.pagedRequest(ctx, req, items)
	}

	return nil
}

// rqlRequest queries an RQL endpoint.
func (lc LaneClient) rqlRequest(ctx context.Context, req *http.Request, results *[]json.RawMessage, cursor string, queryParts ...string) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	// Format RQL query string.
	rqlString := fmt.Sprintf("and(%s)", strings.Join(queryParts, ","))
	if cursor != "" {
		rqlString = fmt.Sprintf("and(%s)", strings.Join([]string{rqlString, fmt.Sprintf("after(%s)", cursor)}, ","))
	}

	// Format JSON
	rqlJSON, err := json.Marshal(struct {
		RQL string `json:"rql"`
	}{rqlString})
	if err != nil {
		return fmt.Errorf("failed to marshal RQL query: %w", err)
	}

	// Set request body
	req.Body = io.NopCloser(bytes.NewReader(rqlJSON))
	req.ContentLength = int64(len(rqlJSON))

	// Fire request
	resp, err := lc.sendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send RQL request: %w", err)
	}

	// Decode RQL page response
	page, err := decodeItems[RQLResult](resp)
	if err != nil {
		return fmt.Errorf("failed to decode RQL response: %w", err)
	}
	for _, item := range page[0].Items {
		*results = append(*results, item.Item)
	}

	// If there's a page cursor, request the next page
	if page[0].Meta.HasNextPage {
		return lc.rqlRequest(ctx, req, results, page[0].Meta.PageCursor.Next, queryParts...)
	}

	return nil
}

// decodeItems decodes a list of raw JSON items into a ResponseModel.
func decodeItems[T ResponseModel, I json.RawMessage | []byte](items ...I) ([]T, error) {
	var results []T

	for _, item := range items {
		res, err := decodeItem[T, I](item)
		if err != nil {
			return nil, fmt.Errorf("failed to decode item: %w", err)
		}
		results = append(results, res)
	}

	return results, nil
}

// decodeItem decodes a single raw JSON item into a ResponseModel.
func decodeItem[T ResponseModel, I json.RawMessage | []byte](item I) (T, error) {
	var res T
	err := json.Unmarshal(item, &res)
	if err != nil {
		return res, fmt.Errorf("failed to decode item: %w", err)
	}
	return res, nil
}
