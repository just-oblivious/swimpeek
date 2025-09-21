package analyzer

import (
	"github.com/just-oblivious/swimpeek/internal/graph"
)

type ComponentCalledByResult struct {
	Actions           map[*graph.Node]bool
	Components        map[*graph.Node]bool
	PlaybookWorkflows map[*graph.Node]map[*graph.Node]bool
}

// ComponentCalledBy analyzes which actions, components, and playbook-workflows call the given component.
func (a *Analyzer) ComponentCalledBy(compNode *graph.Node) *ComponentCalledByResult {
	// Find actions that call this component
	calledByActions := make(map[*graph.Node]bool)
	for _, edge := range compNode.Out {
		if edge.Type == graph.CalledByEdge {
			calledByActions[edge.Dst] = true
		}
	}

	// Find the workflows containing the call actions
	calledByWorkflows := make(map[*graph.Node]bool)
	for actionNode := range calledByActions {
		wfNode := a.GetWorkflowForAction(actionNode)
		if wfNode != nil {
			calledByWorkflows[wfNode] = true
		}
	}

	// Find components and playbooks containing the workflows
	calledByComponents := make(map[*graph.Node]bool)
	calledByPbWorkflows := make(map[*graph.Node]map[*graph.Node]bool)
	for wfNode := range calledByWorkflows {
		comp := a.GetComponentForWorkflow(wfNode)
		if comp != nil {
			calledByComponents[comp] = true
			delete(calledByWorkflows, wfNode)
		}

		pbNode := a.GetPlaybookForWorkflow(wfNode)
		if pbNode != nil {
			if _, exists := calledByPbWorkflows[pbNode]; !exists {
				calledByPbWorkflows[pbNode] = make(map[*graph.Node]bool)
			}
			calledByPbWorkflows[pbNode][wfNode] = true
			delete(calledByWorkflows, wfNode)
		}
	}

	return &ComponentCalledByResult{
		Actions:           calledByActions,
		Components:        calledByComponents,
		PlaybookWorkflows: calledByPbWorkflows,
	}
}

// ComponentCalls returns nodes that the given component node calls.
func (a *Analyzer) ComponentCalls(compNode *graph.Node) map[*graph.Node]bool {
	componentCalls := make(map[*graph.Node]bool)
	wfNode := a.GetWorkflowForComponent(compNode)
	if wfNode == nil {
		return componentCalls
	}

	componentCallActions := a.FindUnique(wfNode, NewWalkOpts(Descend), graph.ComponentActionNode)
	for _, actionNode := range SortSetByLabel(componentCallActions) {
		component := a.FindFirst(actionNode, NewWalkOpts(Ascend, WithFollowEdgeTypes(graph.CalledByEdge), WithMaxDepth(1)), graph.ComponentNode)
		if component != nil {
			componentCalls[component] = true
		}
	}

	return componentCalls
}
