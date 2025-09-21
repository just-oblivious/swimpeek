package detailviews

import (
	"fmt"

	"github.com/just-oblivious/swimpeek/internal/graph"
	"github.com/just-oblivious/swimpeek/internal/tui/app"
	"github.com/just-oblivious/swimpeek/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type FallbackDetailsView struct {
	node   *graph.Node
	frame  *app.Frame
	errMsg string
}

// NewFallbackDetailsView creates a new fallback detail view for unsupported node types or error cases.
func NewFallbackDetailsView(node *graph.Node, frame *app.Frame, errMsg string) *FallbackDetailsView {
	return &FallbackDetailsView{
		node:   node,
		frame:  frame,
		errMsg: errMsg,
	}
}

func (m *FallbackDetailsView) Init() tea.Cmd {
	return nil
}

func (m *FallbackDetailsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *FallbackDetailsView) View() string {
	msg := fmt.Sprintf("%s (label: %s, type: %s, id: %s)", m.errMsg, m.node.Meta.Label, m.node.Meta.Type, m.node.Meta.Id)

	return m.frame.Render(styles.ErrorMsgStyle.Render(msg))
}
