package graph

import (
	"swimpeek/internal/lanedump"
)

type Graph struct {
	Resources *ResourceNodes
}

type Node struct {
	Meta Meta
	Out  []*Edge
	In   []*Edge
}

type Edge struct {
	Src  *Node
	Dst  *Node
	Type EdgeType
	Meta *Meta
}

type Meta struct {
	Id          string
	Type        NodeType
	Label       string
	Description string
}

// FromState creates a graph from a given LaneState.
func FromState(laneState *lanedump.LaneState) (*Graph, []error, error) {

	// Create nodes for top-level resources.
	resources := createNodes(laneState)
	g := newGraph(resources)

	// Run the linker to create edges between nodes.
	warns, err := linkGraph(g, laneState)

	return g, warns, err
}

// newGraph creates a new graph.
func newGraph(res *ResourceNodes) *Graph {
	g := &Graph{
		Resources: res,
	}
	return g
}

// newNode creates a new node with the given metadata.
func newNode(meta Meta) *Node {
	return &Node{
		Meta: meta,
		Out:  make([]*Edge, 0),
		In:   make([]*Edge, 0),
	}
}

// newEdge creates a new edge between two nodes and adds it to their edges list.
func newEdge(src *Node, dst *Node, edgeType EdgeType, meta *Meta) *Edge {
	edge := &Edge{
		Src:  src,
		Dst:  dst,
		Type: edgeType,
		Meta: meta,
	}

	src.Out = append(src.Out, edge)
	dst.In = append(dst.In, edge)

	return edge
}

// newMeta creates a new metadata object with the given properties.
func newMeta(id string, typ NodeType, label string, description string) Meta {
	return Meta{
		Id:          id,
		Type:        typ,
		Label:       label,
		Description: description,
	}
}
