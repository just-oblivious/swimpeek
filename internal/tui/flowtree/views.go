package flowtree

import (
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type FlowViews struct {
	frame    *app.Frame
	analyzer *analyzer.Analyzer
	flows    map[*graph.Node]*flowTree
}

// NewFlowViews creates a manager for flow tree views, it handles dependency injection and caching of views.
func NewFlowViews(frame *app.Frame, analyzer *analyzer.Analyzer) *FlowViews {
	return &FlowViews{
		frame:    frame,
		analyzer: analyzer,
		flows:    make(map[*graph.Node]*flowTree),
	}
}

// ShowFlow returns a flow tree view for the given root node.
func (fv *FlowViews) ShowFlow(rootNode *graph.Node, breadcrumbs []*graph.Node, highlight *graph.Node) tea.Model {
	if flowView, exists := fv.flows[rootNode]; exists {
		flowView.highlightNode(highlight)
		flowView.setBreadcrumbs(breadcrumbs)
		return flowView
	}
	vp := viewport.New(fv.frame.Width, fv.frame.Height)

	flowNode := fv.createFlow(nil, rootNode)
	flowView := newFlowTree(fv.analyzer, fv.frame, flowNode, &vp)

	fv.flows[rootNode] = flowView
	flowView.highlightNode(highlight)
	flowView.setBreadcrumbs(breadcrumbs)
	return flowView
}

// createFlow creates a flow view for the given node and its branches.
func (fv FlowViews) createFlow(edge *graph.Edge, node *graph.Node) *flowNode {
	expanded := true

	// Find refs to components, applications, actions, etc.
	refs := fv.analyzer.GetReferences(node)

	// Render inner flows
	innerNodes := analyzer.NewWalkOpts(analyzer.Descend, analyzer.WithFollowEdgeTypes(graph.EntrypointEdge)).Next(node)

	// Integrate component workflow as inner nodes, but leave it collapsed by default
	if node.Meta.Type == graph.ComponentActionNode {
		expanded = false
		compWf := fv.analyzer.GetWorkflowForComponent(fv.analyzer.GetComponentForAction(node))
		if compWf != nil {
			innerNodes = append(innerNodes, analyzer.NewWalkOpts(analyzer.Descend).Next(compWf)...)
		}

	}
	innerActions := fv.createBranches(innerNodes)

	// Render branches
	branchNodes := analyzer.NewWalkOpts(analyzer.Descend, analyzer.WithSkipEdgeTypes(graph.EntrypointEdge)).Next(node)
	branchActions := fv.createBranches(branchNodes)

	return newFlowNode(node, edge, branchActions, innerActions, refs, expanded)
}

// createBranches creates flow views for a list of nodes connected by edges.
func (fv FlowViews) createBranches(nodes []map[*graph.Edge]*graph.Node) []*flowNode {
	branches := make([]*flowNode, 0, len(nodes))
	for _, n := range nodes {
		for edge, node := range n {
			branch := fv.createFlow(edge, node)
			branches = append(branches, branch)
		}
	}

	return branches
}
