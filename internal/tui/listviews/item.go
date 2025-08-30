package listviews

import (
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type simpleListItem struct {
	label    string
	resource *graph.Node
	analyzer *analyzer.Analyzer
	hasFocus bool
}

// NewSimpleListItem creates a basic selectable resource component for the given resource node.
func NewSimpleListItem(label string, res *graph.Node, analyzer *analyzer.Analyzer, focused bool) tea.Model {
	return simpleListItem{
		label:    label,
		resource: res,
		hasFocus: focused,
		analyzer: analyzer,
	}
}

func (m simpleListItem) Init() tea.Cmd {
	return nil
}

func (m simpleListItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.FocusCmd:
		m.hasFocus = msg.Focus
	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavSelect:
			return m, func() tea.Msg { return app.CmdSwitchView(m.resource) }
		}
	}

	return m, nil
}

func (m simpleListItem) View() string {
	if m.hasFocus {
		return styles.CursorStyle.Render(m.label)
	}
	return m.label
}
