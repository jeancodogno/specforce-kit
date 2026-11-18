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
	content += "  It will be downloaded in the background and applied automatically."

	return boxStyle.Render(content)
}

// RenderManualUpdateNotification returns a notice for a newer version that cannot
// be applied automatically, with the exact command to update by hand.
func RenderManualUpdateNotification(current, latest, command string) string {
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

	content := fmt.Sprintf("%s Specforce %s cannot be installed automatically.\n\n", titleStyle.Render("MANUAL UPDATE REQUIRED"), versionStyle.Render(latest))
	content += fmt.Sprintf("  Current: %s\n", current)
	content += fmt.Sprintf("  Latest:  %s\n\n", versionStyle.Render(latest))
	content += "  To update, run:\n"
	content += fmt.Sprintf("  %s", command)

	return boxStyle.Render(content) + "\n"
}

// RenderUpdateAppliedNotification returns a subtle line announcing that the
// running binary was updated automatically.
func RenderUpdateAppliedNotification(version string) string {
	return fmt.Sprintf("\n\x1b[38;5;245m(Specforce updated to %s)\x1b[0m\n", version)
}
