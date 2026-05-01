package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderUpdateNotification returns a stylized notification box for an available update.
func RenderUpdateNotification(current, latest string) string {
	// Neon Theme (Cyan/Magenta)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")). // Cyan
		Padding(0, 1)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")). // Magenta
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFFF")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	content := fmt.Sprintf("%s A new version of Specforce is available!\n\n", titleStyle.Render("UPDATE AVAILABLE"))
	content += fmt.Sprintf("  Current: %s\n", current)
	content += fmt.Sprintf("  Latest:  %s\n\n", versionStyle.Render(latest))
	content += "  Run " + lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("specforce install") + " to upgrade."

	return boxStyle.Render(content)
}
