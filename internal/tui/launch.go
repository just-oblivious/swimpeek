package tui

import (
	"fmt"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/lanedump"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/flowtree"
	"swimpeek/internal/tui/layout"
	"swimpeek/internal/tui/listviews"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type createResViewFn func(string, *graph.Node, *analyzer.Analyzer, bool) tea.Model

// LaunchExplorer launches the TUI resource explorer application
func LaunchExplorer(laneState *lanedump.LaneState, graph *graph.Graph) error {

	windowStack := make([]tea.Model, 1)
	analyzer := analyzer.NewAnalyzer(laneState, graph)

	tabLabels := []string{"Playbooks", "Components", "Applications", "Triggers"}
	windowFrame := app.NewFrame()
	tabContentFrame := app.NewFrame()

	tabViews := []tea.Model{
		layout.NewListView(createListItemViews(graph.Resources.PlaybooksById, analyzer, listviews.NewPbListItem), tabContentFrame),
		layout.NewListView(createListItemViews(graph.Resources.ComponentsById, analyzer, listviews.NewCompListItem), tabContentFrame),
		layout.NewListView(createListItemViews(graph.Resources.AppsById, analyzer, listviews.NewSimpleListItem), tabContentFrame),
		layout.NewListView(createListItemViews(graph.Resources.TriggersById, analyzer, listviews.NewSimpleListItem), tabContentFrame),
	}

	flowViews := flowtree.NewFlowViews(windowFrame, analyzer)

	windowStack[0] = layout.NewTabView(tabLabels, tabViews, windowFrame, tabContentFrame)
	windowTitle := fmt.Sprintf("SwimPeek - %s (%s)", laneState.Tenant.Name, laneState.TimeStamp.Format(time.DateTime))
	mainView := layout.NewMainView(windowTitle, windowStack, windowFrame, flowViews)

	if _, err := tea.NewProgram(mainView, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

// flattenAndSort flattens a map of graph nodes into a sorted slice
func flattenAndSort(nodes map[string]*graph.Node) []*graph.Node {
	nodeList := make([]*graph.Node, 0, len(nodes))
	for _, node := range nodes {
		nodeList = append(nodeList, node)
	}

	app.SortNodesByLabel(nodeList)
	return nodeList
}

// createListItemViews creates a list of resource components from the given graph nodes
func createListItemViews(resNodes map[string]*graph.Node, analyzer *analyzer.Analyzer, fn createResViewFn) []tea.Model {
	nodes := flattenAndSort(resNodes)

	views := make([]tea.Model, 0, len(nodes))
	for idx, node := range nodes {
		views = append(views, fn(node.Meta.Label, node, analyzer, idx == 0))
	}

	return views
}
