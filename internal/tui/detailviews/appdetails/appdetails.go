package appdetails

import (
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/tabview"
	"swimpeek/pkg/laneclient"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// NewApplicationDetailsView creates a new tab view with application details.
func NewApplicationDetailsView(node *graph.Node, analyzer *analyzer.Analyzer, outerFrame *app.Frame, viewport *viewport.Model, appResource *laneclient.Application) tea.Model {
	innerFrame := app.NewFrame()

	appTriggers := analyzer.ApplicationTriggers(node)
	appAccessLocations := analyzer.ApplicationAccessedBy(node)

	labels := []string{"App Fields", "Record Actions", "Playbook Buttons", "Access Locations"}
	sections := []tea.Model{
		newApplicationFieldList(analyzer, innerFrame, outerFrame, appResource, node),
		newTriggerList(innerFrame, appTriggers.RecordEventTriggers, node),
		newTriggerList(innerFrame, appTriggers.ButtonTriggers, node),
		newAccessListView(analyzer, innerFrame, appAccessLocations, nil, node),
	}

	return tabview.NewTabView(labels, sections, outerFrame, innerFrame)
}
