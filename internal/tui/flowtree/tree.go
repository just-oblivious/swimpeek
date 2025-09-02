package flowtree

import (
	"strings"
	"swimpeek/internal/analyzer"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type flowTree struct {
	analyzer     *analyzer.Analyzer
	frame        *app.Frame
	flowNode     *flowNode
	breadcrumbs  []*graph.Node
	viewport     *viewport.Model
	selectionIdx int
	cursorMode   bool
	visibleCount int
}

// newFlowTree creates a new flow tree for the given root flow node.
func newFlowTree(analyzer *analyzer.Analyzer, frame *app.Frame, flowNode *flowNode, breadcrumbs []*graph.Node, vp *viewport.Model) *flowTree {
	flowNode.setFocus(true)
	visible := getVisibleNodeCount(flowNode)
	return &flowTree{
		analyzer:     analyzer,
		frame:        frame,
		flowNode:     flowNode,
		breadcrumbs:  breadcrumbs,
		viewport:     vp,
		selectionIdx: 0,
		cursorMode:   true,
		visibleCount: visible,
	}
}

func (m flowTree) Init() tea.Cmd {
	return nil
}

// toggleExpand expands or collapses the currently selected node. If applyToAll is true, it applies the action to all nodes.
func (m *flowTree) toggleExpand(state bool, applyToAll bool) {
	walkNodes(func(node *flowNode, idx int) bool {
		if idx == m.selectionIdx || applyToAll {
			node.setExpand(state)
			return applyToAll
		}
		return node.isExpanded
	}, m.flowNode)

	// recompute visible count after collapse/expand
	m.visibleCount = getVisibleNodeCount(m.flowNode)
}

// cursorStep moves the selection cursor up or down by the given step. If the step is 0 it just refreshes the focus state.
func (m *flowTree) cursorStep(step int, applyToAll bool) {
	m.selectionIdx = min(max(m.selectionIdx+step, 0), m.visibleCount-1)
	walkNodes(func(node *flowNode, idx int) bool {
		node.setFocus(idx == m.selectionIdx)
		return node.isExpanded || applyToAll
	}, m.flowNode)
}

func (m *flowTree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.NavCmd:
		if msg.NavEvent == app.NavNextTab {
			m.cursorMode = !m.cursorMode
			return m, nil
		}

		if !m.cursorMode {
			// In scroll mode, just pass the navigation commands to the viewport
			switch msg.NavEvent {
			case app.NavPageUp:
				m.viewport.PageUp()
			case app.NavPageDown:
				m.viewport.PageDown()
			case app.NavUp:
				m.viewport.ScrollUp(1)
			case app.NavDown:
				m.viewport.ScrollDown(1)
			case app.NavHome:
				m.viewport.GotoTop()
			case app.NavEnd:
				m.viewport.GotoBottom()
			case app.NavRight:
				m.viewport.ScrollRight(3)
			case app.NavLeft:
				m.viewport.ScrollLeft(3)
			}
			return m, nil
		}

		// In cursor mode, handle selection and expansion
		switch msg.NavEvent {
		case app.NavNextTab:
			m.cursorMode = !m.cursorMode
		case app.NavPageUp:
			m.cursorStep(-5, false)
		case app.NavPageDown:
			m.cursorStep(5, false)
		case app.NavUp:
			m.cursorStep(-1, false)
		case app.NavDown:
			m.cursorStep(1, false)
		case app.NavHome:
			m.selectionIdx = 0
			m.cursorStep(0, false)
		case app.NavEnd:
			m.selectionIdx = m.visibleCount - 1
			m.cursorStep(0, false)
		case app.NavRight:
			m.toggleExpand(true, false)
		case app.NavLeft:
			m.toggleExpand(false, false)
		case app.NavExpandAll:
			m.toggleExpand(true, true)
		case app.NavCollapseAll:
			m.toggleExpand(false, true)
			m.cursorStep(0, true) // reset focus to make sure we land on a visible node after collapsing all
		}
	}

	return m, nil
}

func (m flowTree) View() string {
	modeHint := "↑↓ SCROLL" // arrow up down symbol
	if m.cursorMode {
		modeHint = "↑↓ SELECT"
	}
	modeBlock := styles.ModeBlockStyle.Render(modeHint)
	modeHelp := styles.HelpDescStyle.Render("TAB to switch mode")
	header := lipgloss.JoinHorizontal(lipgloss.Right, modeBlock, " ", modeHelp)
	title := styles.BoldStyle.Render("Flow: ") + m.renderBreadcrumbs()
	headerTitle := lipgloss.JoinVertical(lipgloss.Left, header, title, "")

	m.viewport.SetContent(m.flowNode.render())
	m.viewport.Width = m.frame.Width - 2
	m.viewport.Height = m.frame.Height - lipgloss.Height(headerTitle)

	scrollBar := styles.RenderScrollBar(m.viewport)

	scrollPane := lipgloss.JoinHorizontal(lipgloss.Left, scrollBar, " ", m.viewport.View())

	return lipgloss.JoinVertical(lipgloss.Left, headerTitle, scrollPane)
}

func (m flowTree) renderBreadcrumbs() string {
	bc := make([]string, len(m.breadcrumbs))
	for i, b := range m.breadcrumbs {
		label, ok := app.NodeLabels[b.Meta.Type]
		if !ok {
			label = string(b.Meta.Type)
		}

		typeLabel := styles.ResTypeLabelStyle.Render("⟨" + label + "⟩")

		icon, ok := app.NodeIcons[b.Meta.Type]
		if ok {
			bc[i] = icon + " " + typeLabel + " " + b.Meta.Label
			continue
		}
		bc[i] = typeLabel + " " + b.Meta.Label

	}
	return strings.Join(bc, " ➜ ")
}
