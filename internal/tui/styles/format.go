package styles

import (
	"strings"
	"swimpeek/internal/tui/app"

	"github.com/charmbracelet/lipgloss"
)

func FormatErrMsg(msg string) string {
	return ErrorMsgStyle.Render(msg)
}

func FrameContent(content string, frame *app.Frame) string {
	return lipgloss.NewStyle().Width(frame.Width).Height(frame.Height).Render(content)
}

func FormatDescription(description string, addNewline bool) string {
	trimmed := strings.Trim(description, "\n ")
	if trimmed == "" {
		return ""
	}

	if addNewline {
		trimmed += "\n"
	}
	return ResDescriptionStyle.Render(trimmed)
}
