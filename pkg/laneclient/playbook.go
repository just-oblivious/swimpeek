package laneclient

import (
	"context"
	"encoding/json"
	"net/http"
)

// Workflow is the top-level container of a playbook.
type Workflow struct {
	Id   string `json:"id"`
	Meta struct {
		Enabled    bool   `json:"enabled"`
		SolutionId string `json:"solutionId"`
		Validation struct {
			Valid bool `json:"valid"`
		} `json:"validation"`
	} `json:"meta"`
	Playbook Playbook `json:"playbook"`
}

// Playbook describes the chain of actions in a workflow.
type Playbook struct {
	Title       string                    `json:"title"`
	Schema      string                    `json:"schema"`
	Description string                    `json:"description"`
	Name        string                    `json:"name"`
	Entrypoints []string                  `json:"entrypoints"`
	Triggers    map[string]any            `json:"triggers"`
	Actions     map[string]PlaybookAction `json:"actions"`
	Meta        struct {
		Enabled bool `json:"enabled"`
	} `json:"meta"`
}

// PlaybookAction is a single action in a playbook.
type PlaybookAction struct {
	Title            string                    `json:"title"`
	Type             string                    `json:"type"`
	Inputs           any                       `json:"inputs"`
	Description      string                    `json:"description"`
	OnComplete       []map[string]any          `json:"on-complete"`
	OnFailure        []map[string]any          `json:"on-failure"`
	OnSuccess        []map[string]any          `json:"on-success"`
	Conditions       []ActionCondition         `json:"conditions"`
	Actions          map[string]PlaybookAction `json:"actions"`
	Action           string                    `json:"action"`
	Asset            string                    `json:"asset"`
	Else             string                    `json:"else"`
	Entrypoint       string                    `json:"entrypoint"`
	Entrypoints      []string                  `json:"entrypoints"`
	Loop             ActionLoop                `json:"loop"`
	Transformations  map[string]any            `json:"transformations"`
	RecordActionType string                    `json:"recordActionType"`
}

// ActionCondition describes the conditions for a conditional action node.
type ActionCondition struct {
	Action    string           `json:"action"`
	Condition map[string][]any `json:"condition"`
}

// ActionLoop describes the configuration of a loop action node.
type ActionLoop struct {
	Type     string `json:"type"`
	Each     any    `json:"each"` // can be a slice or map depending on the loop type
	Parallel bool   `json:"parallel"`
}

// GetPlaybookWorkflows gets all playbook workflows in the tenant.
func (tc TenantClient) GetPlaybookWorkflows(ctx context.Context) ([]Workflow, error) {
	url, err := tc.urlForTenantEndpoint("orchestration", "playbook/rql", 1)
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

	return decodeItems[Workflow](results...)
}
