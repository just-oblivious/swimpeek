package listviews

import (
	"fmt"
	"strings"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pbListItem struct {
	label       string
	playbook    *graph.Node
	pbWorkflows []map[*graph.Node][]*graph.Node
	analyzer    *analyzer.Analyzer
	hasFocus    bool
	isExpanded  bool
	selectedIdx int
}

// NewPbListItem creates a new expandable list item for a playbook node
func NewPbListItem(label string, pbNode *graph.Node, analyzer *analyzer.Analyzer, focused bool) tea.Model {
	return pbListItem{
		label:    label,
		playbook: pbNode,
		analyzer: analyzer,
		hasFocus: focused,
	}
}

func (m *pbListItem) expand() {
	if m.pbWorkflows == nil {
		pbWorkflows := m.analyzer.GetWorkflowsForPlaybook(m.playbook)
		m.pbWorkflows = make([]map[*graph.Node][]*graph.Node, 0, len(pbWorkflows))

		// Find triggers for each workflow
		for _, wf := range pbWorkflows {
			triggers := app.SortSetByLabel(m.analyzer.GetTriggersForWorkflow(wf))
			m.pbWorkflows = append(m.pbWorkflows, map[*graph.Node][]*graph.Node{wf: triggers})
		}
	}

	m.isExpanded = true
}

func (m *pbListItem) collapse() {
	m.isExpanded = false
}

func (m pbListItem) openWorkflow() tea.Msg {
	if len(m.pbWorkflows) == 0 || m.selectedIdx < 0 || m.selectedIdx >= len(m.pbWorkflows) {
		return nil
	}
	for wf, triggers := range m.pbWorkflows[m.selectedIdx] {
		for _, trig := range triggers {
			return app.CmdShowFlow(trig, m.playbook, wf)
		}
		return app.CmdShowFlow(wf, m.playbook, wf)
	}
	return nil
}

func (m pbListItem) Init() tea.Cmd {
	return nil
}

func (m pbListItem) cursorStepInBounds(step int) bool {
	newIdx := m.selectedIdx + step
	return newIdx >= 0 && newIdx < len(m.pbWorkflows)
}

func (m pbListItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.FocusCmd:
		m.hasFocus = msg.Focus

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavRight:
			m.expand()
		case app.NavLeft:
			m.collapse()
		case app.NavSelect:
			if m.isExpanded {
				return m, m.openWorkflow
			}
			m.expand()
		case app.NavUp:
			inBounds := m.cursorStepInBounds(-1)
			if m.isExpanded && inBounds {
				m.selectedIdx = max(m.selectedIdx-1, 0)
				return m, app.CmdCancelNav
			}
		case app.NavDown:
			inBounds := m.cursorStepInBounds(1)
			if m.isExpanded && inBounds {
				m.selectedIdx = max(0, min(m.selectedIdx+1, len(m.pbWorkflows)-1))
				return m, app.CmdCancelNav
			}
		}
	}
	return m, nil
}

func (m pbListItem) View() string {
	if m.isExpanded {
		title := m.compactView()
		detailBlock := styles.ResDetailsStyle.Render(m.detailedView())
		return lipgloss.JoinVertical(lipgloss.Left, title, detailBlock)

	}
	return m.compactView()
}

func (m pbListItem) compactView() string {
	if m.hasFocus {
		return styles.CursorStyle.Render(m.label)
	}
	return m.label
}

func (m pbListItem) detailedView() string {
	workflowDetails := m.renderWorkflowDetails()
	description := styles.FormatDescription(m.playbook.Meta.Description, false)
	return app.JoinVerticalNonEmpty(lipgloss.Top, description, workflowDetails)
}

func (m pbListItem) renderWorkflowDetails() string {
	workflows := make([]string, 0)

	for idx, pbWf := range m.pbWorkflows {
		for wfNode, wfTriggers := range pbWf {
			selected := m.hasFocus && idx == m.selectedIdx
			wf := m.analyzer.GetWorkflowResource(wfNode)

			style := styles.NeutralStyle
			if wf != nil {
				if wf.Meta.Enabled {
					style = styles.ResEnabledStyle
				} else {
					style = styles.ResDisabledStyle
				}
			}
			wfLabel := style.Bold(selected).Render(fmt.Sprintf("● %s", wfNode.Meta.Label))

			trigDetails := m.renderTriggerDetails(wfTriggers)
			wfDescription := styles.FormatDescription(wfNode.Meta.Description, false)
			wfDetails := lipgloss.JoinHorizontal(lipgloss.Left, "  ", app.JoinVerticalNonEmpty(lipgloss.Top, wfDescription, trigDetails))

			cursorPfx := "  "
			if m.hasFocus && idx == m.selectedIdx {
				cursorPfx = styles.CursorStyle.Render("❯ ")
			}
			wfBlock := lipgloss.JoinHorizontal(lipgloss.Top, cursorPfx, lipgloss.JoinVertical(lipgloss.Top, wfLabel, wfDetails))
			workflows = append(workflows, wfBlock)
		}
	}
	if len(workflows) == 0 {
		return styles.ErrorMsgStyle.Italic(true).Render("No workflows found")
	}

	for idx, wf := range workflows {
		if idx < len(workflows)-1 {
			workflows[idx] = lipgloss.JoinVertical(lipgloss.Left, wf, "")
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, workflows...)
}

func (m pbListItem) renderTriggerDetails(wfTriggers []*graph.Node) string {
	triggers := make([]string, 0)

	for _, trigNode := range wfTriggers {
		// Find the application related to the trigger
		refs := make([]string, 0)
		for refNode := range m.analyzer.GetReferences(trigNode) {
			refs = append(refs, refNode.Meta.Label)
		}

		triggerType := styles.ResTriggerStyle.Render("▶ " + string(trigNode.Meta.Type))
		triggerLabel := fmt.Sprintf("%s %s", triggerType, trigNode.Meta.Label)
		if len(refs) > 0 {
			triggerLabel += styles.ResReferenceStyle.Render(fmt.Sprintf(" ← %s", strings.Join(refs, " | ")))
		}

		triggers = append(triggers, triggerLabel)
	}
	if len(triggers) == 0 {
		return styles.ErrorMsgStyle.Italic(true).Render("No trigger found")
	}
	return lipgloss.JoinVertical(lipgloss.Left, triggers...)
}
