package flowtree

import (
	"strings"

	"github.com/just-oblivious/swimpeek/internal/analyzer"
	"github.com/just-oblivious/swimpeek/internal/graph"
	"github.com/just-oblivious/swimpeek/internal/tui/app"
	"github.com/just-oblivious/swimpeek/internal/tui/styles"

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
	yOffsets     []int
	visibleNodes []*flowNode
	selectedNode *flowNode
}

// newFlowTree creates a new flow tree for the given root flow node.
func newFlowTree(analyzer *analyzer.Analyzer, frame *app.Frame, flowNode *flowNode, vp *viewport.Model) *flowTree {
	flowNode.setFocus(true)
	visible, yOffsets := getVisibleNodes(flowNode)
	return &flowTree{
		analyzer:     analyzer,
		frame:        frame,
		flowNode:     flowNode,
		viewport:     vp,
		cursorMode:   true,
		visibleNodes: visible,
		yOffsets:     yOffsets,
		selectedNode: flowNode,
		selectionIdx: 0,
	}
}

func (m flowTree) Init() tea.Cmd {
	return nil
}

// toggleExpand expands or collapses the currently selected node. If applyToAll is true, it applies the value to all visible nodes.
func (m *flowTree) toggleExpand(state bool, applyToAll bool) {
	m.selectedNode.setExpand(state)

	if applyToAll {
		for _, node := range m.visibleNodes {
			node.setExpand(state)
		}
	}

	// recompute visible nodes after a collapse/expand event
	m.visibleNodes, m.yOffsets = getVisibleNodes(m.flowNode)
}

// cursorStep moves the selection cursor up or down by the given step and scrolls the viewport to keep the selection in view.
// If the step is 0 it just resets the scroll position.
func (m *flowTree) cursorStep(step int) {
	m.selectionIdx = min(max(m.selectionIdx+step, 0), len(m.visibleNodes)-1)

	toNode := m.visibleNodes[m.selectionIdx]
	if toNode != m.selectedNode {
		m.selectedNode.setFocus(false)
		toNode.setFocus(true)
		m.selectedNode = toNode
	}

	m.scrollToSelection(m.yOffsets[m.selectionIdx])
}

// highlightNode highlights and scrolls to the given node if it is currently visible in the flow tree.
func (m *flowTree) highlightNode(node *graph.Node) {
	if node == nil || m.selectedNode.node == node {
		return
	}

	// Find the node in the visible nodes
	// TODO: auto-expand nodes until the target node becomes visible
	for idx, flowNode := range m.visibleNodes {
		if flowNode.node == node {
			m.selectedNode.setFocus(false)
			flowNode.setFocus(true)
			m.selectedNode = flowNode
			m.selectionIdx = idx
			m.scrollToSelection(m.yOffsets[m.selectionIdx])
			return
		}
	}
}

// setBreadcrumbs updates the breadcrumb trail for the current flow.
func (m *flowTree) setBreadcrumbs(breadcrumbs []*graph.Node) {
	m.breadcrumbs = breadcrumbs
}

// scrollToSelection scrolls the viewport to ensure the selected node remains visible given its Y-offset in the rendered tree.
func (m *flowTree) scrollToSelection(yOffset int) {
	visibleCount := m.viewport.VisibleLineCount()
	if m.viewport.TotalLineCount() <= visibleCount {
		m.viewport.GotoTop()
		return
	}

	// Bias the scroll position to show more context above the selection
	m.viewport.SetYOffset(yOffset - (visibleCount - visibleCount/4))
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
			m.cursorStep(-5)
		case app.NavPageDown:
			m.cursorStep(5)
		case app.NavUp:
			m.cursorStep(-1)
		case app.NavDown:
			m.cursorStep(1)
		case app.NavHome:
			m.selectionIdx = 0
			m.cursorStep(0)
		case app.NavEnd:
			m.selectionIdx = len(m.visibleNodes) - 1
			m.cursorStep(0)
		case app.NavRight:
			m.toggleExpand(true, false)
		case app.NavExpandAll:
			m.toggleExpand(true, true)
			m.cursorStep(0)
		case app.NavLeft:
			m.toggleExpand(false, false)
		case app.NavCollapseAll:
			m.toggleExpand(false, true)
			m.cursorStep(0)
		}
	}

	return m, nil
}

func (m flowTree) View() string {
	modeHint := "↑↓ SCROLL"
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

// renderBreadcrumbs renders the breadcrumb trail for the current flow.
func (m flowTree) renderBreadcrumbs() string {
	crumbs := make([]string, len(m.breadcrumbs))

	for idx, node := range m.breadcrumbs {
		label, ok := app.NodeLabels[node.Meta.Type]
		if !ok {
			label = string(node.Meta.Type)
		}

		typeLabel := styles.ResTypeLabelStyle.Render("⟨" + label + "⟩")

		icon, ok := app.NodeIcons[node.Meta.Type]
		if ok {
			crumbs[idx] = icon + " " + typeLabel + " " + node.Meta.Label
			continue
		}
		crumbs[idx] = typeLabel + " " + node.Meta.Label
	}

	return strings.Join(crumbs, " ➜ ")
}

// getVisibleNodes walks the nodes and returns a flat list of visible nodes and their Y-offsets.
// The Y-offsets are the cumulative line heights of the nodes as they would be rendered.
func getVisibleNodes(node *flowNode) ([]*flowNode, []int) {
	visibleNodes := make([]*flowNode, 0, 100)
	yOffsets := make([]int, 0, 100)
	yOffset := 0

	walkNodes(func(node *flowNode, idx int) bool {
		visibleNodes = append(visibleNodes, node)
		yOffset += node.lineHeight()
		yOffsets = append(yOffsets, yOffset)

		return node.isExpanded
	}, node)

	return visibleNodes, yOffsets
}
