package tui

import "github.com/charmbracelet/lipgloss"

// Ghost in the Machine Palette
var (
	brandMint     = lipgloss.Color("#00FA9A") // Ghost Primary
	brandCyan     = lipgloss.Color("#00FFFF") // Ghost Secondary
	brandIce      = lipgloss.Color("#E0FFFF") // Ghost Accent
	successGreen  = lipgloss.Color("#5FFF87") // Success
	warningYellow = lipgloss.Color("#FFFFAF") // Warning
	errorRed      = lipgloss.Color("#FF5F5F") // Error
	mutedGrey     = lipgloss.Color("#444444") // Standard UI Border
	textGrey      = lipgloss.Color("#808080") // Secondary Text / Metadata
	textWhite     = lipgloss.Color("#FFFFFF") // Primary Text
	black         = lipgloss.Color("#000000") // Canvas
)

// UI Refinement Colors (Alias for roles)
var (
	finishedGreen    = successGreen
	inProgressOrange = warningYellow
	mutedDarkGrey    = textGrey
)

// Base Styles
var (
	HeaderStyle = lipgloss.NewStyle().
			Foreground(brandMint).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(brandIce)

	BodyStyle = lipgloss.NewStyle().
			Foreground(textWhite)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(mutedDarkGrey)

	FinishedStyle = lipgloss.NewStyle().
			Foreground(finishedGreen)

	InProgressStyle = lipgloss.NewStyle().
			Foreground(inProgressOrange)

	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedDarkGrey)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successGreen)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningYellow)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorRed)
)

// Badge Styles
var (
	SuccessBadgeStyle = lipgloss.NewStyle().
				Foreground(black).
				Background(successGreen).
				Padding(0, 1).
				Bold(true)

	WarningBadgeStyle = lipgloss.NewStyle().
				Foreground(black).
				Background(warningYellow).
				Padding(0, 1).
				Bold(true)

	ErrorBadgeStyle = lipgloss.NewStyle().
				Foreground(textWhite).
				Background(errorRed).
				Padding(0, 1).
				Bold(true)

	// Clean Separator Border (replacing thick borders)
	CleanBorder = lipgloss.Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
	}

	// Completion and Sub-task Styles
	CompletionBoxStyle = lipgloss.NewStyle().
				Border(CleanBorder).
				BorderForeground(mutedGrey).
				Padding(1, 4).
				Margin(1, 0)

	SubTaskStyle = lipgloss.NewStyle().
			Foreground(textGrey).
			PaddingLeft(2)

	// Console specific styles
	PaneBorderStyle = lipgloss.NewStyle().
			Border(CleanBorder).
			BorderForeground(mutedGrey).
			Padding(1)

	CategoryHeaderStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true).
				Underline(true).
				MarginBottom(1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(textWhite)

	SelectionMarkerStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true)
)