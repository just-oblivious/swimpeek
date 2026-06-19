package filter

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/just-oblivious/swimpeek/internal/tui/app"
)

type FilterInput struct {
	textInput textinput.Model
	items     []FilterAble
}

func NewFilterInput(items []FilterAble) *FilterInput {
	ti := textinput.New()
	ti.Placeholder = "/ to filter"
	ti.Blur()
	ti.CharLimit = 256
	ti.Width = 20

	return &FilterInput{
		textInput: ti,
		items:     items,
	}
}

func (fi *FilterInput) Init() tea.Cmd {
	return textinput.Blink
}

func (fi *FilterInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case app.FocusCmd:
		if msg.Focus {
			fi.textInput.Focus()
		} else {
			fi.textInput.Blur()
		}
	case app.NavCmd:
		switch msg.NavEvent {
		case app.NavSelect, app.NavDown:
			if fi.textInput.Focused() {
				fi.textInput.Blur()
				return fi, nil
			}
		}
	}

	fi.textInput, cmd = fi.textInput.Update(msg)
	return fi, cmd
}

func (fi *FilterInput) View() string {
	return fi.textInput.View()
}

func (fi *FilterInput) GetFilteredItems() ([]tea.Model, []string) {
	return fi.items.Match(fi.textInput.Value())
}

func (fi *FilterInput) IsFocused() bool {
	return fi.textInput.Focused()
}
