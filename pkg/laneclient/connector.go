package laneclient

import (
	"context"
	"encoding/json"
	"net/http"
)

// Connector holds metadata for a connector in the platform.
type Connector struct {
	Id   string `json:"id"`
	Meta struct {
		IsSystem      bool                       `json:"isSystem"`
		Actions       map[string]ConnectorAction `json:"actions"`
		InstallSource struct {
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"installSource"`
		Manifest struct {
			Author  string `json:"author"`
			Name    string `json:"name"`
			Product string `json:"product"`
			Title   string `json:"title"`
			Version string `json:"version"`
		} `json:"manifest"`
	} `json:"meta"`
}

// ConnectorAction is an action that the connector can perform.
type ConnectorAction struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Title       string               `json:"title"`
	Inputs      ConnectorInputOutput `json:"inputs"`
	Outputs     ConnectorInputOutput `json:"outputs"`
}

// ConnectorInputOutput describes the inputs and outputs of a connector action.
type ConnectorInputOutput struct {
	Properties map[string]InputOutputProperty `json:"properties"`
	Type       string                         `json:"type"`
	Required   []string                       `json:"required"`
}

// InputOutputProperty describes a single input or output property of a connector action.
type InputOutputProperty struct {
	Description          string `json:"description"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	AdditionalProperties bool   `json:"additionalProperties"`
}

// GetConnectors gets all connectors available in the tenant.
func (tc TenantClient) GetConnectors(ctx context.Context) ([]Connector, error) {
	url, err := tc.urlForTenantEndpoint("orchestration", "connector/rql", 1)
	if err != nil {
		return nil, err
	}

	req, err := tc.lc.prepareRequest(ctx, http.MethodPost, url, nil, nil)
	if err != nil {
		return nil, err
	}

	var results []json.RawMessage
	if err := tc.lc.rqlRequest(ctx, req, &results, "", "limit(100)"); err != nil {
		return nil, err
	}

	return decodeItems[Connector](results...)
}
