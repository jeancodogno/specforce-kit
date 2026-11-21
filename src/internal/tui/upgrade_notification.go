package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderUpdateNotification returns a stylized notification box for an available update.
func RenderUpdateNotification(current, latest string) string {
	// Brand palette (Cyan/Mint)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(brandCyan).
		Padding(0, 1)

	versionStyle := lipgloss.NewStyle().
		Foreground(brandMint).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(brandCyan).
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
		Foreground(brandCyan).
		Padding(0, 1)

	versionStyle := lipgloss.NewStyle().
		Foreground(brandMint).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(brandCyan).
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

// RenderStableNotification returns a notice telling a pre-release user that a
// newer stable version exists, with the exact command to switch channels by hand.
func RenderStableNotification(stable, channel, command string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(brandCyan).
		Padding(0, 1)

	versionStyle := lipgloss.NewStyle().
		Foreground(brandMint).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(brandCyan).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	plain := strings.TrimPrefix(strings.TrimSpace(stable), "v")
	content := fmt.Sprintf("%s\n\n", titleStyle.Render("STABLE RELEASE AVAILABLE"))
	content += fmt.Sprintf("  Specforce %s (stable) is available.\n", versionStyle.Render(plain))
	content += fmt.Sprintf("  You are on the %s channel; to switch, run:\n", channel)
	content += fmt.Sprintf("  %s", command)

	return boxStyle.Render(content) + "\n"
}

// RenderUpdateAppliedNotification returns a subtle line announcing that the
// running binary was updated automatically.
func RenderUpdateAppliedNotification(version string) string {
	style := lipgloss.NewStyle().Foreground(textGrey)
	return "\n" + style.Render(fmt.Sprintf("(Specforce updated to %s)", version)) + "\n"
}
