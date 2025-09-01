package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	NextTab     key.Binding
	PrevTab     key.Binding
	Quit        key.Binding
	Back        key.Binding
	Select      key.Binding
	Filter      key.Binding
	Home        key.Binding
	End         key.Binding
	Help        key.Binding
	Expand      key.Binding
	Collapse    key.Binding
	ExpandAll   key.Binding
	CollapseAll key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Back}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.NextTab, k.PrevTab},
		{k.Expand, k.Collapse, k.ExpandAll, k.CollapseAll},
		{k.Back, k.Filter, k.Quit, k.Help},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "move down"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab", "shift+right"),
		key.WithHelp("tab", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab", "shift+left"),
		key.WithHelp("shift+tab", "previous tab"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "shift+up", "ctrl+b"),
		key.WithHelp("pgup", "rapid up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "shift+down", "ctrl+f"),
		key.WithHelp("pgdown", "rapid down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "jump to top"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "jump to bottom"),
	),
	Expand: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→", "expand node"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←", "collapse node"),
	),
	ExpandAll: key.NewBinding(
		key.WithKeys("X"),
		key.WithHelp("X", "expand all nodes"),
	),
	CollapseAll: key.NewBinding(
		key.WithKeys("Z"),
		key.WithHelp("Z", "collapse all nodes"),
	),
}
