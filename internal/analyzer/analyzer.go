package analyzer

import (
	"github.com/just-oblivious/swimpeek/internal/graph"
	"github.com/just-oblivious/swimpeek/internal/lanedump"
)

type Analyzer struct {
	Lanestate *lanedump.LaneState
	Graph     *graph.Graph
}

// NewAnalyzer creates a new analyzer instance with the given state and graph.
func NewAnalyzer(lanestate *lanedump.LaneState, graph *graph.Graph) *Analyzer {
	return &Analyzer{
		Lanestate: lanestate,
		Graph:     graph,
	}
}
