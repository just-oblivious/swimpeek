package app

import (
	"swimpeek/internal/graph"

	tea "github.com/charmbracelet/bubbletea"
)

type NavDirection int

const (
	NavUp NavDirection = iota
	NavDown
	NavNextTab
	NavPrevTab
	NavPageUp
	NavPageDown
	NavHome
	NavEnd
	NavRight
	NavLeft
	NavSelect
	NavExpandAll
	NavCollapseAll
)

type NavCmd struct {
	NavEvent NavDirection
}

func NavCmdUp() tea.Msg {
	return NavCmd{NavEvent: NavUp}
}
func NavCmdDown() tea.Msg {
	return NavCmd{NavEvent: NavDown}
}
func NavCmdNextTab() tea.Msg {
	return NavCmd{NavEvent: NavNextTab}
}
func NavCmdPrevTab() tea.Msg {
	return NavCmd{NavEvent: NavPrevTab}
}
func NavCmdPageUp() tea.Msg {
	return NavCmd{NavEvent: NavPageUp}
}
func NavCmdPageDown() tea.Msg {
	return NavCmd{NavEvent: NavPageDown}
}
func NavCmdHome() tea.Msg {
	return NavCmd{NavEvent: NavHome}
}
func NavCmdEnd() tea.Msg {
	return NavCmd{NavEvent: NavEnd}
}
func NavCmdRight() tea.Msg {
	return NavCmd{NavEvent: NavRight}
}
func NavCmdLeft() tea.Msg {
	return NavCmd{NavEvent: NavLeft}
}
func NavCmdSelect() tea.Msg {
	return NavCmd{NavEvent: NavSelect}
}
func NavCmdExpandAll() tea.Msg {
	return NavCmd{NavEvent: NavExpandAll}
}
func NavCmdCollapseAll() tea.Msg {
	return NavCmd{NavEvent: NavCollapseAll}
}

type FocusCmd struct {
	Focus bool
}

func CmdFocus() tea.Msg {
	return FocusCmd{Focus: true}
}
func CmdUnfocus() tea.Msg {
	return FocusCmd{Focus: false}
}

type CancelNavCmd struct{}

func CmdCancelNav() tea.Msg {
	return CancelNavCmd{}
}

type ShowFlowCmd struct {
	Node        *graph.Node
	Highlight   *graph.Node
	Breadcrumbs []*graph.Node
}

func CmdShowFlow(node *graph.Node, breadcrumbs ...*graph.Node) tea.Msg {
	return ShowFlowCmd{Node: node, Breadcrumbs: breadcrumbs}
}

func CmdShowFlowWithHighlight(node, highlight *graph.Node, breadcrumbs ...*graph.Node) tea.Msg {
	return ShowFlowCmd{Node: node, Highlight: highlight, Breadcrumbs: breadcrumbs}
}

type ShowDetailsCmd struct {
	Node *graph.Node
}

func CmdShowDetails(node *graph.Node) tea.Msg {
	return ShowDetailsCmd{Node: node}
}

type PushViewCmd struct {
	View tea.Model
}

func CmdPushView(view tea.Model) tea.Msg {
	return PushViewCmd{View: view}
}
