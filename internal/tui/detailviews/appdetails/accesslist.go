package appdetails

import (
	"fmt"
	"slices"
	"strings"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"
	"swimpeek/pkg/laneclient"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type accessListView struct {
	analyzer      *analyzer.Analyzer
	frame         *app.Frame
	accessActions []analyzer.AccessAction
	cursorIdx     int
	app           *graph.Node
	viewport      viewport.Model
	field         *laneclient.ApplicationField
}

// newAccessListView creates a list view for displaying components and playbook-workflows that access records in this application.
func newAccessListView(analyzer *analyzer.Analyzer, frame *app.Frame, accessActions []analyzer.AccessAction, field *laneclient.ApplicationField, app *graph.Node) tea.Model {
	sortAccessActions(accessActions)

	return &accessListView{
		analyzer:      analyzer,
		frame:         frame,
		accessActions: accessActions,
		app:           app,
		viewport:      viewport.New(frame.Width-2, frame.Height),
		field:         field,
	}
}

// sortAccessActions sorts access actions first by access type (create, update, delete, etc.), then by playbook or component label, and finally by workflow name.
func sortAccessActions(actions []analyzer.AccessAction) {
	actionRankFn := func(action analyzer.AccessAction) int {
		switch action.Action.Meta.Type {
		case graph.RecordCreateActionNode:
			return 1
		case graph.RecordUpsertActionNode:
			return 2
		case graph.RecordUpdateActionNode:
			return 3
		case graph.RecordDeleteActionNode:
			return 4
		case graph.RecordSearchActionNode:
			return 5
		case graph.RecordExportActionNode:
			return 6
		default:
			return 7
		}
	}

	pbLabelFn := func(action analyzer.AccessAction) string {
		if action.Playbook != nil {
			return action.Playbook.Meta.Label
		} else if action.Component != nil {
			return action.Component.Meta.Label
		}
		return ""
	}

	wfLabelFn := func(action analyzer.AccessAction) string {
		if action.Workflow != nil {
			return action.Workflow.Meta.Label
		}
		return ""
	}

	slices.SortStableFunc(actions, func(a, b analyzer.AccessAction) int {
		aRank, bRank := actionRankFn(a), actionRankFn(b)
		aPbLabel, bPbLabel := pbLabelFn(a), pbLabelFn(b)
		aWfLabel, bWfLabel := wfLabelFn(a), wfLabelFn(b)

		if aRank != bRank {
			return aRank - bRank
		}
		if aPbLabel != bPbLabel {
			return strings.Compare(aPbLabel, bPbLabel)
		}
		return strings.Compare(aWfLabel, bWfLabel)
	})
}

func (m *accessListView) openWorkflow() tea.Msg {
	if len(m.accessActions) == 0 {
		return nil
	}
	act := m.accessActions[m.cursorIdx]
	if act.IsComponentAction() {
		for ep := range m.analyzer.GetEntrypointsForWorkflow(act.Workflow) {
			return app.CmdShowFlowWithHighlight(ep, act.Action, m.app, act.Component)
		}
		return app.CmdShowFlowWithHighlight(act.Workflow, act.Action, m.app, act.Component)
	}

	for trig := range m.analyzer.GetTriggersForWorkflow(act.Workflow) {
		return app.CmdShowFlowWithHighlight(trig, act.Action, m.app, act.Playbook, act.Workflow)
	}
	return app.CmdShowFlowWithHighlight(act.Workflow, act.Action, m.app, act.Playbook, act.Workflow)
}

func (m *accessListView) Init() tea.Cmd {
	return nil
}

func (m *accessListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavUp:
			m.cursorIdx = max(0, m.cursorIdx-1)
		case app.NavDown:
			m.cursorIdx = min(len(m.accessActions)-1, m.cursorIdx+1)
		case app.NavPageUp:
			m.cursorIdx = max(0, m.cursorIdx-5)
		case app.NavPageDown:
			m.cursorIdx = min(len(m.accessActions)-1, m.cursorIdx+5)
		case app.NavHome:
			m.cursorIdx = 0
		case app.NavEnd:
			m.cursorIdx = len(m.accessActions) - 1
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

func (m *accessListView) View() string {
	content := m.renderAccessList()
	title := styles.TitleStyle.Render(m.app.Meta.Label+fmt.Sprintf(" - %d Access Locations", len(m.accessActions))) + "\n"
	if m.field != nil {
		title = styles.TitleStyle.Render(m.app.Meta.Label+fmt.Sprintf(" - Field: %s (%s) - %d Access Locations", m.field.Name, m.field.Key, len(m.accessActions))) + "\n"
	}

	m.viewport.SetContent(content)
	m.viewport.Width = m.frame.Width - 2
	m.viewport.Height = m.frame.Height - lipgloss.Height(title)
	m.viewport.SetYOffset(m.cursorIdx)

	scrollBar := styles.RenderScrollBar(&m.viewport)
	contentPane := lipgloss.JoinHorizontal(lipgloss.Left, scrollBar, " ", m.viewport.View())

	return lipgloss.JoinVertical(lipgloss.Left, title, contentPane)
}

// renderAccessList renders the access list as access type, access location, and action.
func (m *accessListView) renderAccessList() string {
	if len(m.accessActions) == 0 {
		return styles.ResDescriptionStyle.Render("No actions found")
	}

	formatAccessTypeFn := func(node *graph.Node) string {
		style := styles.ResReferenceStyle

		switch node.Meta.Type {
		case graph.RecordCreateActionNode:
			style = styles.ResEnabledStyle
		case graph.RecordUpdateActionNode:
			style = styles.ResTypeLabelStyle
		case graph.RecordUpsertActionNode:
			style = styles.ResTriggerStyle
		case graph.RecordDeleteActionNode:
			style = styles.ResDisabledStyle
		}

		label := app.NodeLabels[node.Meta.Type]
		if label == "" {
			label = string(node.Meta.Type)
		}

		icon := app.NodeIcons[node.Meta.Type]
		if icon != "" {
			label = icon + " " + label
		}
		return style.Bold(true).Render(label)
	}

	playbookIcon := app.NodeIcons[graph.WorkflowNode]
	componentIcon := app.NodeIcons[graph.ComponentNode]

	formatActionContainerFn := func(action analyzer.AccessAction, highlighted bool) string {
		style := styles.TableCellStyle
		icnStyle := styles.HelpDescStyle
		pfx := "   "
		if highlighted {
			style = styles.CursorStyle
			icnStyle = styles.CursorStyle
			pfx = styles.CursorStyle.Render(" ➜ ")
		}
		if action.IsComponentAction() {
			return fmt.Sprintf("%s%s %s", pfx, icnStyle.Render(componentIcon), style.Render(action.Component.Meta.Label))
		}
		wfStyle := styles.ResDisabledStyle
		if action.Enabled {
			wfStyle = styles.ResEnabledStyle
		}
		return fmt.Sprintf("%s%s %s (%s)", pfx, icnStyle.Render(playbookIcon), style.Render(action.Playbook.Meta.Label), wfStyle.Render(action.Workflow.Meta.Label))
	}

	formatActionFn := func(action analyzer.AccessAction) string {
		return "  " + styles.ResDescriptionStyle.Render(action.Action.Meta.Label)
	}

	formatInspectionErrFn := func(action analyzer.AccessAction) string {
		if action.InspectionErr == nil {
			return ""
		}
		return styles.ErrorMsgStyle.Render(" ! " + action.InspectionErr.Error())
	}

	accessTypeCol := make([]string, len(m.accessActions))
	accessLocationCol := make([]string, len(m.accessActions))
	actionCol := make([]string, len(m.accessActions))
	inspectionErrCol := make([]string, len(m.accessActions))

	for idx, action := range m.accessActions {
		highlighted := idx == m.cursorIdx
		accessTypeCol[idx] = formatAccessTypeFn(action.Action)
		accessLocationCol[idx] = formatActionContainerFn(action, highlighted)
		actionCol[idx] = formatActionFn(action)
		inspectionErrCol[idx] = formatInspectionErrFn(action)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Right, accessTypeCol...),
		lipgloss.JoinVertical(lipgloss.Left, accessLocationCol...),
		lipgloss.JoinVertical(lipgloss.Left, actionCol...),
		lipgloss.JoinVertical(lipgloss.Left, inspectionErrCol...),
	)
}
