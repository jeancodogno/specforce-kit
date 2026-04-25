package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
)

// RenderSpecStatus renders a checklist of the spec artifacts for a given slug.
func RenderSpecStatus(status spec.SpecStatus) string {
	var builder strings.Builder

	// Styles
	successStyle := lipgloss.NewStyle().Foreground(successGreen)
	errorStyle := lipgloss.NewStyle().Foreground(errorRed)
	descStyle := lipgloss.NewStyle().Foreground(textGrey)
	nameStyle := lipgloss.NewStyle().Foreground(textWhite).Bold(true)

	for _, artifact := range status.Artifacts {
		glyph := errorStyle.Render(EmptyBulletGlyph)
		if artifact.Exists {
			glyph = successStyle.Render(BulletGlyph)
		}

		// Align columns: Glyph [Name] Description
		name := nameStyle.Render(fmt.Sprintf("%-16s", artifact.Name))
		description := descStyle.Render(artifact.Description)

		fmt.Fprintf(&builder, " %s %s %s\n", glyph, name, description)
	}

	return builder.String()
}
