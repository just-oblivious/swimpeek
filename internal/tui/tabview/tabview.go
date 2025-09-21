package tabview

import (
	"strings"
	"swimpeek/internal/tui/app"
	"swimpeek/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabView struct {
	tabLabels       []string
	tabViews        []tea.Model
	tabFrame        *app.Frame
	tabContentFrame *app.Frame
	activeTabIdx    int
}

func NewTabView(labels []string, views []tea.Model, tabFrame *app.Frame, tabContentFrame *app.Frame) tea.Model {
	return &tabView{
		tabLabels:       labels,
		tabViews:        views,
		tabFrame:        tabFrame,
		tabContentFrame: tabContentFrame,
	}
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func tabPaddingBorder() lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomRight = border.TopRight
	border.Right = ""
	return border
}

var (
	tabPaddingStyle   = lipgloss.NewStyle().Border(tabPaddingBorder()).BorderForeground(styles.FrameColor).UnsetBorderTop().UnsetBorderLeft()
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(styles.FrameColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true).Bold(true)
	tabContentStyle   = lipgloss.NewStyle().BorderForeground(styles.FrameColor).Padding(1).Border(lipgloss.RoundedBorder()).UnsetBorderTop()
)

func (m *tabView) Init() tea.Cmd {
	return nil
}

func (m *tabView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavNextTab:
			m.activeTabIdx = m.cycleTab(1)
			return m, nil
		case app.NavPrevTab:
			m.activeTabIdx = m.cycleTab(-1)
			return m, nil
		}
	}

	tm, cmd := m.tabViews[m.activeTabIdx].Update(msg)
	m.tabViews[m.activeTabIdx] = tm

	return m, cmd
}

func (m *tabView) View() string {
	tabRow := m.renderTabRow()
	m.tabContentFrame.Height = m.tabFrame.Height - lipgloss.Height(tabRow) - tabContentStyle.GetVerticalFrameSize()
	m.tabContentFrame.Width = m.tabFrame.Width - tabContentStyle.GetHorizontalFrameSize()

	content := lipgloss.NewStyle().
		Height(m.tabContentFrame.Height).
		MaxHeight(m.tabContentFrame.Height).
		Width(m.tabContentFrame.Width).
		MaxWidth(m.tabContentFrame.Width).
		Render(m.tabViews[m.activeTabIdx].View())

	return lipgloss.JoinVertical(lipgloss.Top,
		tabRow,
		tabContentStyle.Render(content),
	)
}

func (m *tabView) renderTabRow() string {
	var renderedTabs []string

	for i, t := range m.tabLabels {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabLabels)-1, i == m.activeTabIdx
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	tabPadding := ""
	tabPaddingWidth := m.tabFrame.Width - lipgloss.Width(tabRow) - 1
	if tabPaddingWidth > 0 {
		tabPadding = tabPaddingStyle.Render(strings.Repeat(" ", tabPaddingWidth))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabRow, tabPadding)

}

func (m *tabView) cycleTab(direction int) int {
	idx := m.activeTabIdx + direction
	if idx < 0 {
		return len(m.tabLabels) - 1
	}
	if idx >= len(m.tabLabels) {
		return 0
	}
	return idx
}
