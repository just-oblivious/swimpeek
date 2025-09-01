package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	NoColor        = lipgloss.AdaptiveColor{}
	ErrorColor     = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#ff6200"}
	TitleColor     = lipgloss.AdaptiveColor{Light: "#0400ff", Dark: "#00ffee"}
	FrameColor     = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	HighlightColor = lipgloss.AdaptiveColor{Light: "#0400ff", Dark: "#00ffee"}
	HelpDescColor  = lipgloss.AdaptiveColor{Light: "#555454", Dark: "#7D7D7D"}
	HelpKeyColor   = lipgloss.AdaptiveColor{Light: "#202020", Dark: "#ADB5BD"}

	DisabledColor  = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#ff6200"}
	EnabledColor   = lipgloss.AdaptiveColor{Light: "#197300", Dark: "#00ff04"}
	TriggerColor   = lipgloss.AdaptiveColor{Light: "#b64c00", Dark: "#ffc400"}
	ReferenceColor = lipgloss.AdaptiveColor{Light: "#0400ff", Dark: "#00aeff"}
	MutedColor     = lipgloss.AdaptiveColor{Light: "#6C757D", Dark: "#ADB5BD"}

	DarkOnLightBGColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	LightOnDarkBGColor = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}
	DarkColor          = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#000000"}
	LightColor         = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}

	SuccessColor = lipgloss.AdaptiveColor{Light: "#197300", Dark: "#00ff04"}
	FailureColor = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#ff6200"}

	FlowLineSegmentColor         = lipgloss.AdaptiveColor{Light: "#6C757D", Dark: "#ADB5BD"}
	FlowLineSegmentFocussedColor = lipgloss.AdaptiveColor{Light: "#0400ff", Dark: "#ffff00"}
	FlowNodeHighlightColor       = lipgloss.AdaptiveColor{Light: "#0400ff", Dark: "#ffff00"}
	TypeLabelColor               = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#ff6ae1"}

	ScrollBarColor = lipgloss.AdaptiveColor{Light: "#6C757D", Dark: "#626262"}
)
