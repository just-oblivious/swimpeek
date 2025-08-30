package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	NeutralStyle  = lipgloss.NewStyle()
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(TitleColor)
	ErrorMsgStyle = lipgloss.NewStyle().Bold(true).Foreground(ErrorColor)
	CursorStyle   = lipgloss.NewStyle().Foreground(HighlightColor).Bold(true)
	HelpDescStyle = lipgloss.NewStyle().Foreground(HelpDescColor)
	HelpKeyStyle  = lipgloss.NewStyle().Foreground(HelpKeyColor)
	BoldStyle     = lipgloss.NewStyle().Bold(true)
)

// Help styles
func HelpStyles() help.Styles {
	return help.Styles{
		ShortKey:       HelpKeyStyle,
		FullKey:        HelpKeyStyle,
		FullDesc:       HelpDescStyle,
		ShortDesc:      HelpDescStyle,
		ShortSeparator: HelpDescStyle,
		FullSeparator:  HelpDescStyle,
	}
}

// Window styles
var (
	WindowStyle   = lipgloss.NewStyle().Padding(0, 2)
	ErrorBoxStyle = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.RoundedBorder()).BorderForeground(ErrorColor)
)

// Resource styles
var (
	ResDetailsStyle     = lipgloss.NewStyle().Padding(0, 0, 1, 0)
	ResDisabledStyle    = lipgloss.NewStyle().Foreground(DisabledColor)
	ResEnabledStyle     = lipgloss.NewStyle().Foreground(EnabledColor)
	ResDescriptionStyle = lipgloss.NewStyle().Italic(true).Foreground(MutedColor)
	ResTypeLabelStyle   = lipgloss.NewStyle().Foreground(TypeLabelColor)
	ResTriggerStyle     = lipgloss.NewStyle().Foreground(TriggerColor).Bold(true)
	ResReferenceStyle   = lipgloss.NewStyle().Foreground(ReferenceColor)
)

// Layout styles
func IndentLeft(n int) lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(n)
}
