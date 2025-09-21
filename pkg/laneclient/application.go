package laneclient

import (
	"context"
	"net/http"
	"time"
)

type Applications []Application

// Application represents a single application.
type Application struct {
	Type            string             `json:"$type"`
	Id              string             `json:"id"`
	Uid             string             `json:"uid"`
	Acronym         string             `json:"acronym"`
	Name            string             `json:"name"`
	CreatedDate     time.Time          `json:"createdDate"`
	ModifiedDate    time.Time          `json:"modifiedDate"`
	Workspaces      []string           `json:"workspaces"`
	TrackingFieldId string             `json:"trackingFieldId"`
	Version         int                `json:"version"`
	Fields          []ApplicationField `json:"fields"`
}

// ApplicationField describes a field in the application.
type ApplicationField struct {
	Type      string `json:"$type"`
	FieldType string `json:"fieldType"`
	Id        string `json:"id"`
	Key       string `json:"key"`
	InputType string `json:"inputType"`
	Name      string `json:"name"`
	Required  bool   `json:"required"`
	ReadOnly  bool   `json:"readonly"`
}

// GetApplications gets all applications in the tenant.
func (tc TenantClient) GetApplications(ctx context.Context) ([]Application, error) {
	url, err := tc.urlForTenantEndpoint("", "app", 0)
	if err != nil {
		return nil, err
	}

	req, err := tc.lc.prepareRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := tc.lc.sendRequest(req)
	if err != nil {
		return nil, err
	}

	apps, err := decodeItems[Applications](res)
	if err != nil {
		return nil, err
	}
	return apps[0], nil
}
