package analyzer

import (
	"slices"
	"swimpeek/internal/graph"
)

type findResult struct {
	nodes []*graph.Node
}

// FindAll finds all nodes reachable from the root node in the specified direction that match one of the nodeTypes.
// This function preserves the order of discovery.
func (a *Analyzer) FindAll(root *graph.Node, opts WalkOpts, nodeTypes ...graph.NodeType) []*graph.Node {
	result := findResult{}
	findFn(root, false, 0, &result, opts, nodeTypes...)
	return result.nodes
}

// FindUnique finds all unique nodes reachable from the root node in the specified direction that match one of the nodeTypes.
func (a *Analyzer) FindUnique(root *graph.Node, opts WalkOpts, nodeTypes ...graph.NodeType) map[*graph.Node]bool {
	result := findResult{}
	findFn(root, false, 0, &result, opts, nodeTypes...)
	return toSet(result.nodes)
}

// FindFirst finds the first node reachable from the root node in the specified direction that match one of the nodeTypes.
func (a *Analyzer) FindFirst(root *graph.Node, opts WalkOpts, nodeTypes ...graph.NodeType) *graph.Node {
	result := findResult{}
	findFn(root, true, 0, &result, opts, nodeTypes...)
	if len(result.nodes) == 0 {
		return nil
	}
	return result.nodes[0]
}

// findFn is a recursive function that performs the actual graph traversal based on the provided options.
func findFn(node *graph.Node, returnFirstResult bool, curDepth int, result *findResult, opts WalkOpts, nodeTypes ...graph.NodeType) {
	if opts.maxDepthReached(curDepth) || node == nil {
		return
	}
	edges := opts.filterEdges(opts.nextEdges(node))

	for _, edge := range edges {
		nextNode := opts.nextNode(edge)

		if len(nodeTypes) == 0 || slices.Contains(nodeTypes, nextNode.Meta.Type) {
			result.nodes = append(result.nodes, nextNode)
			if returnFirstResult {
				return
			}
		}
		if opts.shouldFollow(nextNode) {
			findFn(nextNode, returnFirstResult, curDepth+1, result, opts, nodeTypes...)
		}
	}
}

// toSet converts a slice of nodes to a set (map with bool values).
func toSet(nodes []*graph.Node) map[*graph.Node]bool {
	set := make(map[*graph.Node]bool, len(nodes))
	for _, n := range nodes {
		set[n] = true
	}
	return set
}
