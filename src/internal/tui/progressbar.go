package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderProgressBar returns a stylized progress bar string.
func RenderProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// Calculate filled width
	filledWidth := (percent * width) / 100
	emptyWidth := width - filledWidth

	// Styles
	barColor := brandCyan
	if percent == 100 {
		barColor = successGreen
	}

	filledStyle := lipgloss.NewStyle().Foreground(barColor)
	emptyStyle := lipgloss.NewStyle().Foreground(mutedGrey)

	// Use solid blocks for the neon look
	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty := emptyStyle.Render(strings.Repeat("░", emptyWidth))

	return fmt.Sprintf("%s%s %d%%", filled, empty, percent)
}
