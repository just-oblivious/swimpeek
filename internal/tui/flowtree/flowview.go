package flowtree

import (
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type flowView struct {
	analyzer *analyzer.Analyzer
	frame    *app.Frame
	root     *graph.Node
	workflow tea.Model
	viewport *viewport.Model
}

func (m flowView) Init() tea.Cmd {
	return nil
}

func (m flowView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavUp:
			m.viewport.ScrollUp(5)
		case app.NavDown:
			m.viewport.ScrollDown(5)
		case app.NavRight:
			m.viewport.ScrollRight(5)
		case app.NavLeft:
			m.viewport.ScrollLeft(5)
		case app.NavPageUp:
			m.viewport.ScrollUp(m.frame.Height / 2)
		case app.NavPageDown:
			m.viewport.ScrollDown(m.frame.Height / 2)
		case app.NavHome:
			m.viewport.GotoTop()
		case app.NavEnd:
			m.viewport.GotoBottom()
		}
	}

	return m, nil
}

func (m flowView) View() string {
	title := styles.TitleStyle.Render("Exploring: " + m.root.Meta.Label)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title, "",
		m.workflow.View(),
	)

	m.viewport.SetContent(content)
	m.viewport.Width = m.frame.Width
	m.viewport.Height = m.frame.Height

	return m.viewport.View()
}
