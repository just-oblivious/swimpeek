package laneclient

import (
	"context"
	"fmt"
	"time"
)

type Tenant struct {
	Name            string    `json:"name"`
	Id              string    `json:"id"`
	UserCount       int       `json:"userCount"`
	CreatedDateTime time.Time `json:"createdDateTime"`
}

type TenantResponse struct {
	Tenants    []Tenant `json:"viewModels"`
	TotalCount int      `json:"totalCount"`
}

// GetTenants retrieves the list of tenants for the account.
func (lc LaneClient) GetTenants(ctx context.Context) (TenantResponse, error) {
	url, err := lc.urlForAccountEndpoint("tenants")
	if err != nil {
		return TenantResponse{}, fmt.Errorf("failed to construct tenants URL: %w", err)
	}
	req, err := lc.prepareRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return TenantResponse{}, fmt.Errorf("failed to prepare request for tenants: %w", err)
	}
	resp, err := lc.sendRequest(req)
	if err != nil {
		return TenantResponse{}, fmt.Errorf("failed to request tenants: %w", err)
	}
	return decodeItem[TenantResponse](resp)
}
