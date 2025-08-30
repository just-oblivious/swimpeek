package layout

import (
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listView struct {
	items      []tea.Model
	listFrame  *app.Frame
	cursorIdx  int
	focusedIdx int
	viewport   viewport.Model
}

func NewListView(items []tea.Model, listFrame *app.Frame) tea.Model {
	return listView{
		items:     items,
		listFrame: listFrame,
	}
}

func (m listView) Init() tea.Cmd {
	return nil
}

func (m listView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if len(m.items) == 0 {
		return m, nil
	}

	im, cmd := m.items[m.focusedIdx].Update(msg)
	m.items[m.focusedIdx] = im

	// Cancel navigation if the inner item handles the nav event
	if cmd != nil {
		switch cmd().(type) {
		case app.CancelNavCmd:
			return m, nil
		}
	}

	switch msg := msg.(type) {

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavUp:
			m.cursorIdx = max(m.cursorIdx-1, 0)
		case app.NavDown:
			m.cursorIdx = min(m.cursorIdx+1, len(m.items)-1)
		case app.NavPageUp:
			m.cursorIdx = max(m.cursorIdx-m.listFrame.Height/3, 0)
		case app.NavPageDown:
			m.cursorIdx = min(m.cursorIdx+m.listFrame.Height/3, len(m.items)-1)
		case app.NavHome:
			m.cursorIdx = 0
		case app.NavEnd:
			m.cursorIdx = len(m.items) - 1
		}
	}

	if m.cursorIdx != m.focusedIdx {
		m.items[m.focusedIdx], _ = m.items[m.focusedIdx].Update(app.CmdUnfocus())
		m.items[m.cursorIdx], _ = m.items[m.cursorIdx].Update(app.CmdFocus())

		m.focusedIdx = m.cursorIdx
	}

	return m, cmd
}

func (m listView) View() string {
	renderedItems := make([]string, len(m.items))
	yOffset := 0
	hWrapContainer := lipgloss.NewStyle().Width(m.listFrame.Width)
	for i, itm := range m.items {
		itmView := hWrapContainer.Render(itm.View())
		itmHeight := lipgloss.Height(itmView)
		if i < m.cursorIdx {
			yOffset += itmHeight
		}

		pfx := "  "
		if i == m.cursorIdx {
			pfx = styles.CursorStyle.Render("❯ ")
		}
		renderedItems[i] = lipgloss.JoinHorizontal(lipgloss.Top, pfx, itmView)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, renderedItems...)

	m.viewport.SetContent(content)
	m.viewport.Width = m.listFrame.Width
	m.viewport.Height = m.listFrame.Height
	m.viewport.SetYOffset(yOffset)

	return m.viewport.View()

}
