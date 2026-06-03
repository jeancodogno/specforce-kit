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
	warnStyle := lipgloss.NewStyle().Foreground(brandCyan)

	// Render Refinement State if active
	if status.RefinementIteration > 0 {
		iterText := fmt.Sprintf("REFINEMENT: Iteration %d/3", status.RefinementIteration)
		fmt.Fprintf(&builder, " %s\n\n", warnStyle.Bold(true).Render(iterText))
	}

	for _, artifact := range status.Artifacts {
		glyph := errorStyle.Render(EmptyBulletGlyph)
		if artifact.Exists {
			glyph = successStyle.Render(BulletGlyph)
		}

		// Align columns: Glyph [Name] Description
		name := nameStyle.Render(fmt.Sprintf("%-16s", artifact.Name))
		description := descStyle.Render(artifact.Description)

		fmt.Fprintf(&builder, " %s %s %s\n", glyph, name, description)

		// Render validation errors if they exist
		for _, err := range artifact.ValidationErrors {
			indent := strings.Repeat(" ", 4)
			fmt.Fprintf(&builder, "%s%s\n", indent, errorStyle.Render("- "+err))
		}
	}

	// Render Coherence Errors
	if len(status.RefinementErrors) > 0 {
		fmt.Fprintf(&builder, "\n %s\n", errorStyle.Bold(true).Render("COHERENCE ERRORS:"))
		for _, err := range status.RefinementErrors {
			fmt.Fprintf(&builder, "  %s\n", errorStyle.Render("↳ "+err))
		}
		
		if status.RefinementIteration >= 3 {
			fmt.Fprintf(&builder, "\n %s\n", warnStyle.Render("Automatic refinement reached limit. Manual intervention required."))
		}
	}

	return builder.String()
}
