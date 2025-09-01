package analyzer

import (
	"swimpeek/internal/graph"
)

// GetPlaybookForWorkflow returns the playbook associated with the given workflow node
func (a *Analyzer) GetPlaybookForWorkflow(wfNode *graph.Node) *graph.Node {
	return a.FindFirst(wfNode, NewWalkOpts(Ascend, WithMaxDepth(1), WithFollowEdgeTypes(graph.WorkflowEdge)), graph.PlaybookNode)
}

// GetComponentForWorkflow returns the component associated with the given workflow node
func (a *Analyzer) GetComponentForWorkflow(wfNode *graph.Node) *graph.Node {
	return a.FindFirst(wfNode, NewWalkOpts(Ascend, WithMaxDepth(1), WithFollowEdgeTypes(graph.WorkflowEdge)), graph.ComponentNode)
}

// GetWorkflowForComponent returns the workflow associated with the given component node
func (a *Analyzer) GetWorkflowForComponent(compNode *graph.Node) *graph.Node {
	return a.FindFirst(compNode, NewWalkOpts(Descend, WithMaxDepth(1), WithFollowEdgeTypes(graph.WorkflowEdge)), graph.WorkflowNode)
}

// GetWorkflowsForPlaybook returns the workflows contained in the given playbook node
func (a *Analyzer) GetWorkflowsForPlaybook(pbNode *graph.Node) []*graph.Node {
	// Use FindAll to maintain the order of workflows as defined in the playbook
	return a.FindAll(pbNode, NewWalkOpts(Descend, WithMaxDepth(1), WithFollowEdgeTypes(graph.WorkflowEdge)), graph.WorkflowNode)
}

// GetWorkflowForAction returns the workflow in which the given action node is contained
func (a *Analyzer) GetWorkflowForAction(actionNode *graph.Node) *graph.Node {
	return a.FindFirst(actionNode, NewWalkOpts(Ascend), graph.WorkflowNode)
}

// GetTriggersForWorkflow returns the triggers associated with the given workflow node.
func (a *Analyzer) GetTriggersForWorkflow(wfNode *graph.Node) map[*graph.Node]bool {
	return a.FindUnique(wfNode, NewWalkOpts(Ascend, WithFollowEdgeTypes(graph.TriggersWorkflowEdge), WithMaxDepth(1)))
}

// GetEntrypointsForWorkflow returns the entrypoints of a workflow node.
func (a *Analyzer) GetEntrypointsForWorkflow(wfNode *graph.Node) map[*graph.Node]bool {
	return a.FindUnique(wfNode, NewWalkOpts(Descend, WithFollowEdgeTypes(graph.EntrypointEdge), WithMaxDepth(1)))
}

// GetComponentForAction returns the component associated with the given action node
func (a *Analyzer) GetComponentForAction(actionNode *graph.Node) *graph.Node {
	return a.FindFirst(actionNode, NewWalkOpts(Ascend, WithFollowEdgeTypes(graph.CalledByEdge), WithMaxDepth(1)), graph.ComponentNode)
}

// GetReferences returns all nodes that reference the given node.
func (a *Analyzer) GetReferences(node *graph.Node) map[*graph.Node]bool {
	refEdges := []graph.EdgeType{
		graph.AccessedByEdge,
		graph.HasActionEdge,
		graph.CalledByEdge,
		graph.EmittedByEdge,
		graph.HasEventEdge,
	}
	return a.FindUnique(node, NewWalkOpts(Ascend, WithMaxDepth(1), WithFollowEdgeTypes(refEdges...)))
}
