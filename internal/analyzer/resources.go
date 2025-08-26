package analyzer

import (
	"swimpeek/internal/graph"
	"swimpeek/pkg/laneclient"
)

// GetWorkflowResource returns the workflow resource associated with the given workflow node, if it exists.
func (a *Analyzer) GetWorkflowResource(wfNode *graph.Node) *laneclient.Workflow {
	if wf, exists := a.Lanestate.WorkflowsById[wfNode.Meta.Id]; exists {
		return &wf
	}
	return nil
}
