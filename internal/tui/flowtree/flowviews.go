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
	flows    map[*graph.Node]*flowView
}

// NewFlowViews creates a manager for flow tree views, it handles dependency injection and caching of views.
func NewFlowViews(frame *app.Frame, analyzer *analyzer.Analyzer) *FlowViews {
	return &FlowViews{
		frame:    frame,
		analyzer: analyzer,
		flows:    make(map[*graph.Node]*flowView),
	}
}

// ShowFlow returns a flow tree view for the given root node.
func (fv *FlowViews) ShowFlow(rootNode *graph.Node) tea.Model {
	if flow, exists := fv.flows[rootNode]; exists {
		return flow
	}
	vp := viewport.New(fv.frame.Width, fv.frame.Height)
	newFlow := &flowView{
		frame:    fv.frame,
		analyzer: fv.analyzer,
		root:     rootNode,
		workflow: fv.createFlow(nil, rootNode),
		viewport: &vp,
	}
	fv.flows[rootNode] = newFlow
	return newFlow
}

func (fv FlowViews) createFlow(edge *graph.Edge, node *graph.Node) tea.Model {

	// Render inner flows
	innerNodes := analyzer.NewWalkOpts(analyzer.Descend, analyzer.WithFollowEdgeTypes(graph.EntrypointEdge)).Next(node)
	innerActions := fv.createBranches(innerNodes)

	// Render branches
	branchNodes := analyzer.NewWalkOpts(analyzer.Descend, analyzer.WithSkipEdgeTypes(graph.EntrypointEdge)).Next(node)
	branchActions := fv.createBranches(branchNodes)

	// Find refs to components, applications, actions, etc.
	refs := fv.analyzer.GetReferences(node)

	return newActionNode(node, edge, branchActions, innerActions, refs, true)
}

func (fv FlowViews) createBranches(nodes []map[*graph.Edge]*graph.Node) []tea.Model {
	branches := make([]tea.Model, 0, len(nodes))
	for _, n := range nodes {
		for edge, node := range n {
			branch := fv.createFlow(edge, node)
			branches = append(branches, branch)
		}
	}

	return branches
}
