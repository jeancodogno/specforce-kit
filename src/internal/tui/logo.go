package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var AppVersion = "v0.1.8"

var brailleLines = []string{
	`    ⢠⣶⣶⡄     `,
	`  ⢠⣿⣿⣿⣿⡄    `,
	` ⢠⣿⣿⣿⣿⣿⣿⡄   `,
	` ⠻⣿⣿⣿⣿⣿⣿⠟   `,
	`   ⠻⣿⣿⣿⠟    `,
	`     ⠻⠟       `,
}

var gradientColors = []lipgloss.Color{
	lipgloss.Color("#00FA9A"),
	lipgloss.Color("#00FCBB"),
	lipgloss.Color("#00FDEE"),
	lipgloss.Color("#00FFFF"),
	lipgloss.Color("#70FFFF"),
	lipgloss.Color("#E0FFFF"),
}

var (
	cachedLogo     string
	cachedBranding string
)

// GenerateLogo returns the styled lipgloss string for the new Ghost in the Machine logo.
func GenerateLogo(withSubtitle bool) string {
	if withSubtitle && cachedBranding != "" {
		return cachedBranding
	}
	if !withSubtitle && cachedLogo != "" {
		return cachedLogo
	}

	var renderedBraille []string
	for i, line := range brailleLines {
		style := lipgloss.NewStyle().Foreground(gradientColors[i]).Bold(true)
		renderedBraille = append(renderedBraille, style.Render(line))
	}
	brailleBlock := lipgloss.JoinVertical(lipgloss.Left, renderedBraille...)

	var textLines []string
	textLines = append(textLines, "") // Alignment for Line 0

	specforceText := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("SPECFORCE ")
	versionStr := AppVersion
	if !strings.HasPrefix(versionStr, "v") {
		versionStr = "v" + versionStr
	}
	versionText := DimmedStyle.Render(versionStr)
	textLines = append(textLines, specforceText+versionText)

	if withSubtitle {
		textLines = append(textLines, DimmedStyle.Render("Spec-Driven Development (SDD)"))
		textLines = append(textLines, DimmedStyle.Render("Ecosystem for AI-assisted development"))
	}

	textBlock := lipgloss.JoinVertical(lipgloss.Left, textLines...)

	finalStr := lipgloss.JoinHorizontal(lipgloss.Top, brailleBlock, textBlock)
	finalStr = lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(finalStr)

	if withSubtitle {
		cachedBranding = finalStr
	} else {
		cachedLogo = finalStr
	}

	return finalStr
}

// PrintLogo renders the Specforce logo.
func PrintLogo() {
	fmt.Println(GenerateLogo(false))
}

// PrintBranding renders both the logo and the subtitle.
func PrintBranding() {
	fmt.Println(GenerateLogo(true))
}
