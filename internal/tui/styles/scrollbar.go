package styles

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// RenderScrollBar renders a scrollbar for the given viewport.
func RenderScrollBar(vp *viewport.Model) string {
	lineCount := vp.TotalLineCount()
	if lineCount <= vp.Height || vp.Height <= 0 {
		vp.SetYOffset(0) // reset offset if content fits in viewport
		return " "
	}
	sbHeight := max(1, vp.Height*vp.Height/lineCount)
	sbOffset := 1 + vp.YOffset*(vp.Height-sbHeight-2)/(lineCount-vp.Height)

	// Force scrollbar to scroll when not at the edges to hint that there's more content
	if vp.YOffset != 0 {
		sbOffset = max(2, sbOffset)
	}
	if vp.YOffset+vp.Height < lineCount {
		sbOffset = min(vp.Height-sbHeight-1, sbOffset)
	}

	sb := make([]string, vp.Height)
	for i := range sb {
		if i == 0 {
			sb[i] = "⌃"
			continue
		}
		if i == vp.Height-1 {
			sb[i] = "⌄"
			continue
		}
		if i >= sbOffset && i < sbOffset+sbHeight {
			sb[i] = "█"
			continue
		}
		sb[i] = "░"
	}
	return ScrollBarStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sb...))
}
