package laneclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OrchestrationSolution (playbook or component).
type OrchestrationSolution struct {
	Type                 string    `json:"$type"`
	Id                   string    `json:"id"`
	Uid                  string    `json:"uid"`
	PlaybookIds          []string  `json:"playbookIds"` // workflows referenced by a playbook
	PlaybookId           string    `json:"playbookId"`  // workflow referenced by a component
	ReferencedComponents []string  `json:"referencedComponents"`
	Disabled             bool      `json:"disabled"`
	Name                 string    `json:"name"`
	CreatedDate          time.Time `json:"createdDate"`
	ModifiedDate         time.Time `json:"modifiedDate"`
	Description          string    `json:"description"`
	Version              int       `json:"version"`
}

// OrchestrationTasks are tasks attached to a specific application (e.g. recordAction, playbookButton).
type OrchestrationTasks []OrchestrationTask

// OrchestrationTask describes an a application task.
type OrchestrationTask struct {
	Id            string        `json:"id"`
	Uid           string        `json:"uid"`
	ApplicationId string        `json:"applicationId"`
	Disabled      bool          `json:"disabled"`
	Name          string        `json:"name"`
	PlaybookId    string        `json:"playbookId"`
	Type          string        `json:"type"`
	Triggers      []TaskTrigger `json:"triggers"`
}

// TaskTrigger describes the trigger conditions of an OrchestrationTask.
type TaskTrigger struct {
	Type              string   `json:"type"`
	AvailableMappings []string `json:"availableMappings"`
	Conditions        []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"conditions"`
	OnCorrelationActionComplete bool `json:"onCorrelationActionComplete"`
	OnRecordCreate              bool `json:"onRecordCreate"`
	OnRecordUpdate              bool `json:"onRecordUpdate"`
	HasRecord                   bool `json:"hasRecord"`
}

// getOrchestrationSolutions gets orchestration solutions from the solution builder.
func (tc TenantClient) getOrchestrationSolutions(ctx context.Context, resourceType string) ([]OrchestrationSolution, error) {
	url, err := tc.urlForTenantEndpoint("", fmt.Sprintf("solution-builder/%s/filter", resourceType), 0)
	if err != nil {
		return nil, err
	}

	req, err := tc.lc.prepareRequest(ctx, http.MethodPost, url, nil, nil)
	if err != nil {
		return nil, err
	}

	var items []json.RawMessage
	if err := tc.lc.pagedRequest(ctx, req, &items); err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", resourceType, err)
	}

	return decodeItems[OrchestrationSolution](items...)
}

// GetSolutions gets all playbooks (solutions) in the tenant.
func (tc TenantClient) GetPlaybooks(ctx context.Context) ([]OrchestrationSolution, error) {
	return tc.getOrchestrationSolutions(ctx, "solutions")
}

// GetComponents gets all components in the tenant.
func (tc TenantClient) GetComponents(ctx context.Context) ([]OrchestrationSolution, error) {
	return tc.getOrchestrationSolutions(ctx, "components")
}

// GetOrchestrationTasks returns all orchestration tasks in the tenant.
func (tc TenantClient) GetOrchestrationTasks(ctx context.Context) ([]OrchestrationTask, error) {
	url, err := tc.urlForTenantEndpoint("", "orchestrationtask", 0)
	if err != nil {
		return nil, err
	}

	req, err := tc.lc.prepareRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := tc.lc.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orchestration tasks: %w", err)
	}

	tasks, err := decodeItems[OrchestrationTasks](resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode orchestration tasks: %w", err)
	}
	return tasks[0], nil
}
