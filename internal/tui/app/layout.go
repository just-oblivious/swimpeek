package app

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

type Frame struct {
	Width  int
	Height int
}

func NewFrame() *Frame {
	return &Frame{}
}

func JoinVerticalNonEmpty(pos lipgloss.Position, itms ...string) string {
	return lipgloss.JoinVertical(pos, slices.DeleteFunc(itms, func(s string) bool { return s == "" })...)
}

func (f Frame) Render(content string) string {
	return lipgloss.NewStyle().Width(f.Width).Height(f.Height).Render(content)
}
