package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TUI Glyphs
const (
	ArrowGlyph       = "›"
	BulletGlyph      = "◉"
	EmptyBulletGlyph = "○"
	SeparatorGlyph   = "─"
)

// Shared Component Styles
var (
	ActiveArrowStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true)

	SelectedBulletStyle = lipgloss.NewStyle().
				Foreground(brandCyan)

	UnselectedBulletStyle = lipgloss.NewStyle().
				Foreground(mutedGrey)

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(mutedGrey)

	FooterStyle = lipgloss.NewStyle().
			Foreground(textGrey).
			Background(lipgloss.Color("#111111")).
			Padding(0, 1)
)

// RenderSeparator returns a full-width dimmed separator
func RenderSeparator(width int) string {
	if width <= 0 {
		width = 80 // Default fallback
	}
	return SeparatorStyle.Render(strings.Repeat(SeparatorGlyph, width))
}

// PrintFooter renders a contextual information bar at the bottom.
func PrintFooter(version string) {
	cwd, _ := os.Getwd()
	path := cwd
	
	footerText := fmt.Sprintf("Specforce %s | Workspace: %s", version, path)
	fmt.Println("\n" + FooterStyle.Render(footerText))
}

// RenderBadge returns a high-contrast boxed status message.
func RenderBadge(status string, message string) string {
	var style lipgloss.Style
	var prefix string

	switch strings.ToLower(status) {
	case "success", "ok":
		style = SuccessBadgeStyle
		prefix = " SUCCESS "
	case "warning", "warn":
		style = WarningBadgeStyle
		prefix = " WARNING "
	case "error", "fail":
		style = ErrorBadgeStyle
		prefix = "  ERROR  "
	default:
		style = lipgloss.NewStyle().Foreground(textWhite).Background(mutedGrey).Padding(0, 1)
		prefix = "  INFO   "
	}

	badge := style.Render(prefix)
	return fmt.Sprintf("\n %s %s\n", badge, BodyStyle.Render(message))
}

// RenderErrorBadge returns a high-contrast boxed error message.
func RenderErrorBadge(message string) string {
	return RenderBadge("error", message)
}

// LogSubTask renders a styled sub-task log with an indented cyan arrow.
func LogSubTask(message string) {
	prefix := ActiveArrowStyle.Render(" › ")
	fmt.Println(SubTaskStyle.Render(prefix + message))
}

// PrintCompletionBox renders a high-fidelity summary box for successful operations.
func PrintCompletionBox(title, message string) {
	header := SuccessStyle.Bold(true).Render(title)
	body := BodyStyle.Render(message)
	fmt.Println("\n" + CompletionBoxStyle.Render(header+"\n\n"+body))
}

// ArtifactDisplay represents a generic artifact for TUI listing.
type ArtifactDisplay struct {
	ID          string
	Description string
}

// PrintArtifactList renders a standardized list of artifacts.
func PrintArtifactList(title, subtitle string, artifacts []ArtifactDisplay, version string) {
	if IsTTY() {
		PrintBranding()
	}

	fmt.Println("\n" + HeaderStyle.Render(title))
	fmt.Println(SubtitleStyle.Render(subtitle))
	fmt.Println()

	for _, art := range artifacts {
		fmt.Printf("  %s %s\n", HeaderStyle.Width(15).Render(art.ID), SubtitleStyle.Render(art.Description))
	}
	fmt.Println()

	if IsTTY() {
		PrintFooter(version)
	}
}

// Confirm asks a simple yes/no question in the terminal and returns true if yes.
func Confirm(question string) bool {
	fmt.Printf("\n %s %s (y/N): ", ActiveArrowStyle.Render("?"), BodyStyle.Render(question))
	var response string
	_, _ = fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
