package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
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
	WindowStyle    = lipgloss.NewStyle().Padding(0, 2)
	ErrorBoxStyle  = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.RoundedBorder()).BorderForeground(ErrorColor)
	ScrollBarStyle = lipgloss.NewStyle().Foreground(ScrollBarColor)
	ModeBlockStyle = lipgloss.NewStyle().Padding(0, 1).Background(FrameColor).Bold(true)
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

// Table Styles
var (
	TableHeaderStyle    = lipgloss.NewStyle().Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(TableHeaderColor).Foreground(TableHeaderColor)
	TableCellStyle      = lipgloss.NewStyle().Foreground(TableCellColor)
	TableSelectionStyle = lipgloss.NewStyle().Background(HighlightColor).Foreground(LightOnDarkBGColor).Bold(true)
	TableBorderStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(TableBorderColor)
	TableStyle          = table.Styles{
		Header:   TableHeaderStyle,
		Selected: TableSelectionStyle,
	}
)
