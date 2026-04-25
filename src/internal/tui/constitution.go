package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeancodogno/specforce-kit/src/internal/constitution"
)

// RenderConstitutionStatus renders a checklist of the constitution artifacts.
func RenderConstitutionStatus(status constitution.ConstitutionStatus) string {
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
		// Assume name max length around 15-20 characters for padding
		name := nameStyle.Render(fmt.Sprintf("%-16s", artifact.Name))
		description := descStyle.Render(artifact.Description)

		fmt.Fprintf(&builder, " %s %s %s\n", glyph, name, description)
	}

	return builder.String()
}
