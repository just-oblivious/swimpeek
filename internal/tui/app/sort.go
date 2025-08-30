package app

import (
	"slices"
	"strings"
	"swimpeek/internal/graph"
)

// SortNodesByLabel sorts a slice of nodes in place by their label.
func SortNodesByLabel(nodes []*graph.Node) {
	slices.SortFunc(nodes, func(a, b *graph.Node) int { return strings.Compare(a.Meta.Label, b.Meta.Label) })
}

// SortSetByLabel takes a map of nodes and returns a sorted slice of nodes by their label.
func SortSetByLabel[T any](nodes map[*graph.Node]T) []*graph.Node {
	nodeList := make([]*graph.Node, 0, len(nodes))
	for node := range nodes {
		nodeList = append(nodeList, node)
	}
	SortNodesByLabel(nodeList)
	return nodeList
}
