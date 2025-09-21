package layout

import (
	"github.com/just-oblivious/swimpeek/internal/tui/app"
	"github.com/just-oblivious/swimpeek/internal/tui/detailviews"
	"github.com/just-oblivious/swimpeek/internal/tui/flowtree"
	"github.com/just-oblivious/swimpeek/internal/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainView struct {
	keys         app.KeyMap
	title        string
	width        int
	height       int
	windowStack  []tea.Model
	help         help.Model
	contentFrame *app.Frame
	flowViews    *flowtree.FlowViews
	detailViews  *detailviews.DetailViews
}

func NewMainView(title string, windowStack []tea.Model, frame *app.Frame, flowViews *flowtree.FlowViews, detailViews *detailviews.DetailViews) tea.Model {
	h := help.New()
	h.Styles = styles.HelpStyles()

	return mainView{
		keys:         app.Keys,
		title:        title,
		windowStack:  windowStack,
		help:         h,
		contentFrame: frame,
		flowViews:    flowViews,
		detailViews:  detailViews,
	}
}

func (m mainView) Init() tea.Cmd {
	return nil
}

func (m mainView) getActiveContent() tea.Model {
	return m.windowStack[len(m.windowStack)-1]
}

func (m mainView) updateActiveContent(content tea.Model) mainView {
	m.windowStack[len(m.windowStack)-1] = content
	return m
}

func (m mainView) updateContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	cm, cmd := m.getActiveContent().Update(msg)
	return m.updateActiveContent(cm), cmd
}

func (m mainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.NextTab):
			return m.updateContent(app.NavCmdNextTab())
		case key.Matches(msg, m.keys.PrevTab):
			return m.updateContent(app.NavCmdPrevTab())
		case key.Matches(msg, m.keys.Up):
			return m.updateContent(app.NavCmdUp())
		case key.Matches(msg, m.keys.Down):
			return m.updateContent(app.NavCmdDown())
		case key.Matches(msg, m.keys.PageUp):
			return m.updateContent(app.NavCmdPageUp())
		case key.Matches(msg, m.keys.PageDown):
			return m.updateContent(app.NavCmdPageDown())
		case key.Matches(msg, m.keys.Home):
			return m.updateContent(app.NavCmdHome())
		case key.Matches(msg, m.keys.End):
			return m.updateContent(app.NavCmdEnd())
		case key.Matches(msg, m.keys.Expand):
			return m.updateContent(app.NavCmdRight())
		case key.Matches(msg, m.keys.Collapse):
			return m.updateContent(app.NavCmdLeft())
		case key.Matches(msg, m.keys.Select):
			return m.updateContent(app.NavCmdSelect())
		case key.Matches(msg, m.keys.ExpandAll):
			return m.updateContent(app.NavCmdExpandAll())
		case key.Matches(msg, m.keys.CollapseAll):
			return m.updateContent(app.NavCmdCollapseAll())

		// Quit active content first before quitting the application, this avoids the user from accidentally quitting
		// the app when they meant to go back to the previous view :)
		case key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Quit):
			if len(m.windowStack) > 1 {
				m.windowStack = m.windowStack[:len(m.windowStack)-1]
				return m, nil
			}
			if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}

	case app.ShowFlowCmd:
		m.windowStack = append(m.windowStack, m.flowViews.ShowFlow(msg.Node, msg.Breadcrumbs, msg.Highlight))
		return m, nil

	case app.ShowDetailsCmd:
		m.windowStack = append(m.windowStack, m.detailViews.ShowDetails(msg.Node))
		return m, nil

	case app.PushViewCmd:
		m.windowStack = append(m.windowStack, msg.View)
		return m, nil
	}

	return m.updateContent(msg)
}

func (m mainView) View() string {
	title := styles.TitleStyle.Render(m.title)
	usage := m.help.View(m.keys)

	m.contentFrame.Height = m.height - styles.WindowStyle.GetVerticalFrameSize() - lipgloss.Height(title) - lipgloss.Height(usage)
	m.contentFrame.Width = m.width - styles.WindowStyle.GetHorizontalFrameSize()
	content := m.getActiveContent().View()

	return styles.WindowStyle.
		Height(m.height).
		Width(m.width).
		Render(lipgloss.JoinVertical(lipgloss.Center, title, content, usage))
}
