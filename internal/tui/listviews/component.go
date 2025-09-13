package listviews

import (
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type compListItem struct {
	label      string
	component  *graph.Node
	analyzer   *analyzer.Analyzer
	hasFocus   bool
	isExpanded bool
	calledBy   *analyzer.ComponentCalledByResult
	calls      map[*graph.Node]bool
}

// NewCompListItem creates a new expandable list item for a component node
func NewCompListItem(label string, compNode *graph.Node, analyzer *analyzer.Analyzer, focused bool) tea.Model {
	return compListItem{
		label:     label,
		component: compNode,
		analyzer:  analyzer,
		hasFocus:  focused,
	}
}

func (m *compListItem) expand() {
	if m.calledBy == nil {
		m.calledBy = m.analyzer.ComponentCalledBy(m.component)
	}
	if m.calls == nil {
		m.calls = m.analyzer.ComponentCalls(m.component)
	}
	m.isExpanded = true
}

func (m *compListItem) collapse() {
	m.isExpanded = false
}

func (m compListItem) openComponent() tea.Msg {
	wf := m.analyzer.GetWorkflowForComponent(m.component)
	for ep := range m.analyzer.GetEntrypointsForWorkflow(wf) {
		return app.CmdShowFlow(ep, m.component)
	}
	return app.CmdShowFlow(wf, m.component)
}

func (m compListItem) Init() tea.Cmd {
	return nil
}

func (m compListItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, m.openComponent
		}
	}

	return m, nil
}

func (m compListItem) View() string {
	if m.isExpanded {
		title := m.compactView()
		detailView := m.detailedView()
		detailBlock := styles.ResDetailsStyle.Padding(0, 0, 1, 2).Render(detailView)
		return lipgloss.JoinVertical(lipgloss.Left, title, detailBlock)
	}
	return m.compactView()
}

func (m compListItem) compactView() string {
	if m.hasFocus {
		return styles.CursorStyle.Render(m.label)
	}
	return m.label
}

func (m compListItem) detailedView() string {
	description := styles.FormatDescription(m.component.Meta.Description, true)
	callLocations := m.renderCallLocations()
	return app.JoinVerticalNonEmpty(lipgloss.Top, description, callLocations)
}

func (m compListItem) renderCallLocations() string {

	// Format the components section
	calledByComponents := make([]string, 0, len(m.calledBy.Components))
	for _, comp := range app.SortSetByLabel(m.calledBy.Components) {
		calledByComponents = append(calledByComponents, comp.Meta.Label)
	}

	calledByComponentsList := styles.ResDescriptionStyle.Render("None")
	if len(calledByComponents) > 0 {
		calledByComponentsList = styles.ResReferenceStyle.Render(lipgloss.JoinVertical(lipgloss.Left, calledByComponents...))
	}
	calledByComponentsSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoldStyle.Render("Called by components:"),
		styles.IndentLeft(1).Render(calledByComponentsList),
	)

	// Format the playbook-workflows section
	calledByPlaybooks := make([]string, 0, len(m.calledBy.PlaybookWorkflows))
	for pbIdx, pb := range app.SortSetByLabel(m.calledBy.PlaybookWorkflows) {
		wfs := app.SortSetByLabel(m.calledBy.PlaybookWorkflows[pb])
		wfLabels := make([]string, 0, len(wfs))
		for wfIdx, wf := range wfs {
			pfx := "├─"
			if wfIdx == len(wfs)-1 {
				pfx = "╰─"
			}
			wfLabels = append(wfLabels, styles.ResReferenceStyle.Render(pfx+wf.Meta.Label))
		}

		sfx := "\n"
		if pbIdx == len(m.calledBy.PlaybookWorkflows)-1 {
			sfx = ""
		}
		wfList := lipgloss.JoinVertical(lipgloss.Left, wfLabels...)
		pbSection := lipgloss.JoinVertical(lipgloss.Left,
			styles.ResReferenceStyle.Bold(true).Render(pb.Meta.Label),
			wfList+sfx,
		)

		calledByPlaybooks = append(calledByPlaybooks, pbSection)
	}

	calledByPlaybookList := styles.ResDescriptionStyle.Render("None")
	if len(calledByPlaybooks) > 0 {
		calledByPlaybookList = lipgloss.JoinVertical(lipgloss.Left, calledByPlaybooks...)
	}

	calledByPlaybookSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoldStyle.Render("Called by playbooks:"),
		styles.IndentLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, calledByPlaybookList)),
	)

	// Format the calls section
	callLabels := make([]string, 0, len(m.calls))
	for _, call := range app.SortSetByLabel(m.calls) {
		callLabels = append(callLabels, call.Meta.Label)
	}

	callList := styles.ResDescriptionStyle.Render("None")
	if len(callLabels) > 0 {
		callList = styles.ResReferenceStyle.Render(lipgloss.JoinVertical(lipgloss.Left, callLabels...))
	}
	callsComponentsSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoldStyle.Render("Calls components:"),
		styles.IndentLeft(1).Render(callList),
	)

	return lipgloss.JoinVertical(lipgloss.Left, calledByComponentsSection, "", calledByPlaybookSection, "", callsComponentsSection)
}
