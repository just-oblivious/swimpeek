package detailviews

import (
	"github.com/just-oblivious/swimpeek/internal/analyzer"
	"github.com/just-oblivious/swimpeek/internal/graph"
	"github.com/just-oblivious/swimpeek/internal/tui/app"
	"github.com/just-oblivious/swimpeek/internal/tui/detailviews/appdetails"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type DetailViews struct {
	frame    *app.Frame
	analyzer *analyzer.Analyzer
	views    map[*graph.Node]tea.Model
}

// NewDetailViews creates a manager for detail views, it handles dependency injection and caching of views.
func NewDetailViews(frame *app.Frame, analyzer *analyzer.Analyzer) *DetailViews {
	return &DetailViews{
		frame:    frame,
		analyzer: analyzer,
		views:    make(map[*graph.Node]tea.Model),
	}
}

// ShowDetails returns a detail view for the given node.
func (dv *DetailViews) ShowDetails(node *graph.Node) tea.Model {
	if view, exists := dv.views[node]; exists {
		return view
	}

	vp := viewport.New(dv.frame.Width, dv.frame.Height)

	var detailView tea.Model

	switch node.Meta.Type {
	case graph.ApplicationNode:
		app := dv.analyzer.GetApplicationResource(node)
		if app == nil {
			detailView = NewFallbackDetailsView(node, dv.frame, "Application data not found")
			break
		}
		detailView = appdetails.NewApplicationDetailsView(node, dv.analyzer, dv.frame, &vp, app)
	default:
		detailView = NewFallbackDetailsView(node, dv.frame, "No detail view available for this resource type")
	}

	dv.views[node] = detailView

	return detailView
}
