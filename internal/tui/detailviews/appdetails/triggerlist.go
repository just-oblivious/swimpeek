package appdetails

import (
	"fmt"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appTriggerList struct {
	frame          *app.Frame
	triggerActions []analyzer.TriggerAction
	app            *graph.Node
	cursorIdx      int
	viewport       viewport.Model
}

// newTriggerList creates a list view for displaying trigger actions associated with an application.
func newTriggerList(frame *app.Frame, triggerActions []analyzer.TriggerAction, app *graph.Node) tea.Model {
	return &appTriggerList{
		frame:          frame,
		triggerActions: triggerActions,
		app:            app,
		viewport:       viewport.New(frame.Width-2, frame.Height),
	}
}

func (m *appTriggerList) openWorkflow() tea.Msg {
	if len(m.triggerActions) == 0 {
		return nil
	}
	trig := m.triggerActions[m.cursorIdx]
	return app.CmdShowFlow(trig.Trigger, m.app, trig.Playbook, trig.Workflow)
}

func (m *appTriggerList) Init() tea.Cmd {
	return nil
}

func (m *appTriggerList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavUp:
			m.cursorIdx = max(0, m.cursorIdx-1)
		case app.NavDown:
			m.cursorIdx = min(len(m.triggerActions)-1, m.cursorIdx+1)
		case app.NavPageUp:
			m.cursorIdx = max(0, m.cursorIdx-5)
		case app.NavPageDown:
			m.cursorIdx = min(len(m.triggerActions)-1, m.cursorIdx+5)
		case app.NavHome:
			m.cursorIdx = 0
		case app.NavEnd:
			m.cursorIdx = len(m.triggerActions) - 1
		case app.NavSelect:
			return m, m.openWorkflow
		case app.NavLeft:
			m.viewport.ScrollLeft(5)
		case app.NavRight:
			m.viewport.ScrollRight(5)
		}
	}

	return m, nil
}

func (m *appTriggerList) View() string {
	content := m.renderTriggerList()
	title := styles.TitleStyle.Render(m.app.Meta.Label+fmt.Sprintf(" - %d Trigger Actions", len(m.triggerActions))) + "\n"

	m.viewport.SetContent(content)
	m.viewport.Width = m.frame.Width - 2
	m.viewport.Height = m.frame.Height - lipgloss.Height(title)
	m.viewport.SetYOffset(m.cursorIdx)

	scrollBar := styles.RenderScrollBar(&m.viewport)
	contentPane := lipgloss.JoinHorizontal(lipgloss.Left, scrollBar, " ", m.viewport.View())

	return lipgloss.JoinVertical(lipgloss.Left, title, contentPane)

}

// renderTriggerList renders the trigger list as label -> playbook (workflow)
func (m *appTriggerList) renderTriggerList() string {
	if len(m.triggerActions) == 0 {
		return styles.ResDescriptionStyle.Render("No actions found")
	}

	wfLabels := make([]string, 0, len(m.triggerActions))
	wfTriggers := make([]string, 0, len(m.triggerActions))

	for idx, trig := range m.triggerActions {
		wfStyle := styles.ResDisabledStyle
		if trig.Enabled {
			wfStyle = styles.ResEnabledStyle
		}

		trigStyle := styles.ResTriggerStyle
		pbStyle := styles.TableCellStyle
		sepStyle := styles.HelpDescStyle

		if m.cursorIdx == idx {
			pbStyle = styles.CursorStyle
			sepStyle = styles.CursorStyle
		}
		sep := sepStyle.Render(" ➜ ")

		wfLabels = append(wfLabels, trigStyle.Render(trig.Trigger.Meta.Label))
		wfTriggers = append(wfTriggers, fmt.Sprintf("%s%s (%s)", sep, pbStyle.Render(trig.Playbook.Meta.Label), wfStyle.Render(trig.Workflow.Meta.Label)))
	}

	trigLabels := lipgloss.JoinVertical(lipgloss.Right, wfLabels...)
	trigWorkflow := lipgloss.JoinVertical(lipgloss.Left, wfTriggers...)
	return lipgloss.JoinHorizontal(lipgloss.Left, trigLabels, trigWorkflow)
}
