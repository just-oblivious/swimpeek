package analyzer

import (
	"slices"
	"swimpeek/internal/graph"
)

type WalkDirection int

const (
	Descend WalkDirection = iota
	Ascend
)

type WalkOpts struct {
	direction       WalkDirection
	followNodeTypes []graph.NodeType
	skipNodeTypes   []graph.NodeType
	followEdgeTypes []graph.EdgeType
	skipEdgeTypes   []graph.EdgeType
	skipNodes       []*graph.Node
	maxDepth        int
}

func NewWalkOpts(direction WalkDirection, options ...func(*WalkOpts)) WalkOpts {
	opts := WalkOpts{
		direction: direction,
	}
	for _, option := range options {
		option(&opts)
	}
	return opts
}

func WithFollowNodeTypes(types ...graph.NodeType) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.followNodeTypes = append(opts.followNodeTypes, types...)
	}
}

func WithSkipNodeTypes(types ...graph.NodeType) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.skipNodeTypes = append(opts.skipNodeTypes, types...)
	}
}

func WithFollowEdgeTypes(types ...graph.EdgeType) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.followEdgeTypes = append(opts.followEdgeTypes, types...)
	}
}

func WithSkipEdgeTypes(types ...graph.EdgeType) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.skipEdgeTypes = append(opts.skipEdgeTypes, types...)
	}
}

func WithSkipNodes(nodes ...*graph.Node) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.skipNodes = append(opts.skipNodes, nodes...)
	}
}

func WithMaxDepth(depth int) func(*WalkOpts) {
	return func(opts *WalkOpts) {
		opts.maxDepth = depth
	}
}

// Next returns a slice of maps representing the next edges and nodes matching the filtering criteria.
func (t WalkOpts) Next(node *graph.Node) []map[*graph.Edge]*graph.Node {
	edges := t.filterEdges(t.nextEdges(node))
	next := make([]map[*graph.Edge]*graph.Node, 0, len(edges))

	for _, edge := range edges {
		nextNode := t.nextNode(edge)
		if t.shouldFollow(nextNode) {
			next = append(next, map[*graph.Edge]*graph.Node{edge: nextNode})
		}
	}
	return next
}

// nextEdges returns the next edges to traverse based on the trace direction.
func (t WalkOpts) nextEdges(node *graph.Node) []*graph.Edge {
	if t.direction == Descend {
		return node.Out
	}
	return node.In
}

// nextNode returns the next node to traverse based on the trace direction.
func (t WalkOpts) nextNode(edge *graph.Edge) *graph.Node {
	if t.direction == Descend {
		return edge.Dst
	}
	return edge.Src
}

// shouldFollow checks if a node should be followed based on the trace options.
func (t WalkOpts) shouldFollow(node *graph.Node) bool {
	if slices.Contains(t.skipNodes, node) {
		return false
	}
	if len(t.skipNodeTypes) > 0 && slices.Contains(t.skipNodeTypes, node.Meta.Type) {
		return false
	}
	if len(t.followNodeTypes) > 0 && !slices.Contains(t.followNodeTypes, node.Meta.Type) {
		return false
	}
	return true
}

// filterEdges filters edges based on types to skip or follow.
func (t WalkOpts) filterEdges(edges []*graph.Edge) []*graph.Edge {
	if len(t.skipEdgeTypes) == 0 && len(t.followEdgeTypes) == 0 {
		return edges
	}
	filtered := make([]*graph.Edge, 0, len(edges))
	for _, edge := range edges {
		if slices.Contains(t.skipEdgeTypes, edge.Type) {
			continue
		}
		if len(t.followEdgeTypes) > 0 && !slices.Contains(t.followEdgeTypes, edge.Type) {
			continue
		}
		filtered = append(filtered, edge)
	}
	return filtered
}

func (t WalkOpts) maxDepthReached(currentDepth int) bool {
	return t.maxDepth > 0 && currentDepth >= t.maxDepth
}
