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

type walkOpts struct {
	direction       WalkDirection
	followNodeTypes []graph.NodeType
	skipNodeTypes   []graph.NodeType
	followEdgeTypes []graph.EdgeType
	skipEdgeTypes   []graph.EdgeType
	skipNodes       []*graph.Node
	maxDepth        int
}

func NewWalkOpts(direction WalkDirection, options ...func(*walkOpts)) walkOpts {
	opts := walkOpts{
		direction: direction,
	}
	for _, option := range options {
		option(&opts)
	}
	return opts
}

func WithFollowNodeTypes(types ...graph.NodeType) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.followNodeTypes = append(opts.followNodeTypes, types...)
	}
}

func WithSkipNodeTypes(types ...graph.NodeType) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.skipNodeTypes = append(opts.skipNodeTypes, types...)
	}
}

func WithFollowEdgeTypes(types ...graph.EdgeType) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.followEdgeTypes = append(opts.followEdgeTypes, types...)
	}
}

func WithSkipEdgeTypes(types ...graph.EdgeType) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.skipEdgeTypes = append(opts.skipEdgeTypes, types...)
	}
}

func WithSkipNodes(nodes ...*graph.Node) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.skipNodes = append(opts.skipNodes, nodes...)
	}
}

func WithMaxDepth(depth int) func(*walkOpts) {
	return func(opts *walkOpts) {
		opts.maxDepth = depth
	}
}

// nextEdges returns the next edges to traverse based on the trace direction.
func (t walkOpts) nextEdges(node *graph.Node) []*graph.Edge {
	if t.direction == Descend {
		return node.Out
	}
	return node.In
}

// nextNode returns the next node to traverse based on the trace direction.
func (t walkOpts) nextNode(edge *graph.Edge) *graph.Node {
	if t.direction == Descend {
		return edge.Dst
	}
	return edge.Src
}

// shouldFollow checks if a node should be followed based on the trace options.
func (t walkOpts) shouldFollow(node *graph.Node) bool {
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
func (t walkOpts) filterEdges(edges []*graph.Edge) []*graph.Edge {
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

func (t walkOpts) maxDepthReached(currentDepth int) bool {
	return t.maxDepth > 0 && currentDepth >= t.maxDepth
}
