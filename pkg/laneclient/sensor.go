package laneclient

import (
	"context"
	"encoding/json"
	"net/http"
)

// Sensor represents a webhook or flow event listener.
type Sensor struct {
	Id   string `json:"id"`
	Meta struct {
		Enabled            bool     `json:"enabled"`
		Name               string   `json:"name"`
		Title              string   `json:"title"`
		EmittedByPlaybooks []string `json:"emittedByPlaybooks"`
		TriggeredPlaybooks []string `json:"triggeredPlaybooks"`
	} `json:"meta"`
	Sensor struct {
		Description string `json:"description"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Title       string `json:"title"`
	}
}

// GetSensors gets all 'sensors' in the tenant.
func (tc TenantClient) GetSensors(ctx context.Context) ([]Sensor, error) {
	url, err := tc.urlForTenantEndpoint("orchestration", "sensor/rql", 1)
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

	return decodeItems[Sensor](results...)
}
