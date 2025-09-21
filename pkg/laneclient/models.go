package laneclient

import (
	"encoding/json"
)

// ResponseModel is a common interface for all response models.
type ResponseModel interface {
	ItemPage |
		RQLResult |
		OrchestrationSolution |
		Workflow |
		Applications |
		Connector |
		OrchestrationTasks |
		Sensor |
		TenantResponse
}

// ItemPage is a common structure for paginated responses.
type ItemPage struct {
	Items []json.RawMessage `json:"items"`
}

// RQLResult is the result of an RQL query.
type RQLResult struct {
	Items []struct {
		Item json.RawMessage `json:"item"`
	} `json:"items"`
	Meta struct {
		HasNextPage bool `json:"hasNextPage"`
		PageCursor  struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"pageCursor"`
		RQL string `json:"rql"`
	} `json:"meta"`
}
