package analyzer

import (
	"github.com/just-oblivious/swimpeek/internal/graph"
	"github.com/just-oblivious/swimpeek/pkg/laneclient"
)

// GetWorkflowResource returns the workflow resource associated with the given workflow node, if it exists.
func (a *Analyzer) GetWorkflowResource(wfNode *graph.Node) *laneclient.Workflow {
	if wf, exists := a.Lanestate.WorkflowsById[wfNode.Meta.Id]; exists {
		return &wf
	}
	return nil
}

// GetApplicationResource returns the application resource associated with the given application node, if it exists.
func (a *Analyzer) GetApplicationResource(appNode *graph.Node) *laneclient.Application {
	if app, exists := a.Lanestate.ApplicationsById[appNode.Meta.Id]; exists {
		return &app
	}
	return nil
}

// GetActionResource returns the playbook action associated with the given action node within the specified workflow, if it exists.
func (a *Analyzer) GetActionResource(wfNode *graph.Node, actNode *graph.Node) *laneclient.PlaybookAction {
	wfResource := a.GetWorkflowResource(wfNode)
	if wfResource == nil {
		return nil
	}

	// Recursively search for the action within the workflow actions.
	var findActFn func(actions map[string]laneclient.PlaybookAction) *laneclient.PlaybookAction
	findActFn = func(actions map[string]laneclient.PlaybookAction) *laneclient.PlaybookAction {
		for actId, act := range actions {
			if actId == actNode.Meta.Id {
				return &act
			}
			if innerAct := findActFn(act.Actions); innerAct != nil {
				return innerAct
			}
		}
		return nil
	}

	return findActFn(wfResource.Playbook.Actions)
}
