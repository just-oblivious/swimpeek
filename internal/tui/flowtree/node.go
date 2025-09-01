package flowtree

import (
	"strings"
	"swimpeek/internal/graph"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

type flowNode struct {
	node         *graph.Node
	edge         *graph.Edge
	branches     []*flowNode
	innerActions []*flowNode
	references   map[*graph.Node]bool
	hasFocus     bool
	isExpanded   bool
}

func newFlowNode(node *graph.Node, edge *graph.Edge, branches []*flowNode, innerActions []*flowNode, refs map[*graph.Node]bool, expanded bool) *flowNode {
	return &flowNode{
		node:         node,
		edge:         edge,
		branches:     branches,
		innerActions: innerActions,
		references:   refs,
		hasFocus:     false,
		isExpanded:   expanded,
	}
}

func (m *flowNode) setFocus(focussed bool) {
	m.hasFocus = focussed
}

func (m *flowNode) setExpand(expanded bool) {
	m.isExpanded = expanded
}

func (m *flowNode) render() string {
	icon, ok := app.NodeIcons[m.node.Meta.Type]
	if !ok {
		icon = "●"
	}
	label := m.renderEdge(icon) + m.renderNodeLabel() + m.renderReferences()

	blocks := make([]string, 0, 2)

	// Render inner actions
	innerActions := make([]string, 0, len(m.innerActions))
	if m.isExpanded {
		for _, ia := range m.innerActions {
			innerActions = append(innerActions, ia.render())
		}
	} else if len(m.innerActions) > 0 {
		innerActions = append(innerActions, styles.ResDescriptionStyle.Render("■➜ expand..."))
	}
	innerActionsBlock := lipgloss.JoinVertical(lipgloss.Left, innerActions...)

	innerActionLines := m.renderLineSegments(innerActions, 0)
	if len(innerActions) > 0 {
		blocks = append(blocks, lipgloss.JoinHorizontal(lipgloss.Left, innerActionLines, innerActionsBlock))
	}

	// Render branches
	branches := make([]string, 0, len(m.branches))
	for _, b := range m.branches {
		branches = append(branches, b.render())
	}

	branchBlock := lipgloss.JoinVertical(lipgloss.Left, branches...)
	branchLines := ""
	if len(branches) > 0 {
		blocks = append(blocks, branchBlock)
		offset := 0
		if len(innerActions) > 0 {
			offset = lipgloss.Height(innerActionsBlock)
		}
		branchLines = m.renderLineSegments(branches, offset)
	}
	actionBlocks := lipgloss.JoinVertical(lipgloss.Left, blocks...)

	renderedAction := lipgloss.JoinVertical(lipgloss.Left, label, lipgloss.JoinHorizontal(lipgloss.Left, branchLines, actionBlocks))
	return renderedAction
}

func (m flowNode) renderEdge(icon string) string {
	if m.edge == nil {
		return icon + " "
	}
	color := styles.NoColor

	switch m.edge.Type {
	case graph.OnSuccessEdge, graph.IfEdge:
		color = styles.SuccessColor
	case graph.OnFailureEdge, graph.ElseEdge, graph.UnreachableEdge:
		color = styles.FailureColor
	case graph.EntrypointEdge:
		return icon + "➜ "
	}

	label, ok := app.EdgeLabels[m.edge.Type]
	if !ok {
		label = string(m.edge.Type)
	}

	style := lipgloss.NewStyle().Bold(m.hasFocus).Foreground(color)
	return style.Render(icon + " " + label + " ")
}

func (m flowNode) renderNodeLabel() string {
	label, ok := app.NodeLabels[m.node.Meta.Type]
	if !ok {
		label = string(m.node.Meta.Type)
	}
	typeLabel := styles.ResTypeLabelStyle.Render("⟨" + label + "⟩")
	style := lipgloss.NewStyle()
	if m.hasFocus {
		style = style.Foreground(styles.FlowNodeHighlightColor).Bold(true)
	}

	return typeLabel + " " + style.Render(m.node.Meta.Label)
}

func (m flowNode) renderReferences() string {
	if len(m.references) == 0 {
		return ""
	}
	refs := make([]string, 0, len(m.references))
	for ref := range m.references {
		icon, ok := app.NodeIcons[ref.Meta.Type]
		if !ok {
			icon = "?"
		}
		refs = append(refs, icon+" "+ref.Meta.Label)
	}
	return styles.ResReferenceStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, " ➜ ", strings.Join(refs, " · ")))
}

func (m flowNode) renderLineSegments(blocks []string, offset int) string {
	border := lipgloss.RoundedBorder()
	color := styles.FlowLineSegmentColor
	if m.hasFocus {
		color = styles.FlowLineSegmentFocussedColor
	}
	lineSegments := make([]string, 0, len(blocks)+offset)
	for range offset {
		lineSegments = append(lineSegments, "┆ ")
	}

	for idx, itm := range blocks {
		if idx == len(blocks)-1 {
			lineSegments = append(lineSegments, border.BottomLeft+border.Bottom)
			break
		}
		lineSegments = append(lineSegments, "├"+border.Bottom)
		for range lipgloss.Height(itm) - 1 {
			lineSegments = append(lineSegments, "│ ")
		}
	}

	colorStyle := lipgloss.NewStyle().Foreground(color)

	return colorStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lineSegments...))
}
