package tui

import "github.com/charmbracelet/lipgloss"

// Ghost in the Machine Palette (adaptive: Light value for light terminal
// backgrounds, Dark value for dark ones; dark values are the original palette).
var (
	brandMint     = lipgloss.AdaptiveColor{Light: "#007A4D", Dark: "#00FA9A"} // Ghost Primary
	brandCyan     = lipgloss.AdaptiveColor{Light: "#00727F", Dark: "#00FFFF"} // Ghost Secondary
	brandIce      = lipgloss.AdaptiveColor{Light: "#3F5F6B", Dark: "#E0FFFF"} // Ghost Accent
	successGreen  = lipgloss.AdaptiveColor{Light: "#1E7A34", Dark: "#5FFF87"} // Success
	warningYellow = lipgloss.AdaptiveColor{Light: "#8A6100", Dark: "#FFFFAF"} // Warning
	errorRed      = lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#FF5F5F"} // Error
	mutedGrey     = lipgloss.AdaptiveColor{Light: "#B0B0B0", Dark: "#444444"} // Borders / Separators
	textGrey      = lipgloss.AdaptiveColor{Light: "#5C5C5C", Dark: "#808080"} // Secondary Text / Metadata
	textPrimary   = lipgloss.AdaptiveColor{Light: "", Dark: "#FFFFFF"}        // Primary Text (terminal default on light)
	footerBg      = lipgloss.AdaptiveColor{Light: "#EEEEEE", Dark: "#111111"} // Footer Bar Background
	badgeOnFill   = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"} // Label on success/warning badges
	badgeOnError  = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"} // Label on error badges (white in both modes)
)

// UI Refinement Colors (Alias for roles)
var (
	finishedGreen    = successGreen
	inProgressYellow = warningYellow
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
			Foreground(textPrimary)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(mutedDarkGrey)

	FinishedStyle = lipgloss.NewStyle().
			Foreground(finishedGreen)

	InProgressStyle = lipgloss.NewStyle().
			Foreground(inProgressYellow)

	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedDarkGrey)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successGreen)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningYellow)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorRed)

	ClockStyle = lipgloss.NewStyle().
			Foreground(brandMint)

	OvertimeStyle = lipgloss.NewStyle().
			Foreground(warningYellow)
)

// Badge Styles
var (
	SuccessBadgeStyle = lipgloss.NewStyle().
				Foreground(badgeOnFill).
				Background(successGreen).
				Padding(0, 1).
				Bold(true)

	WarningBadgeStyle = lipgloss.NewStyle().
				Foreground(badgeOnFill).
				Background(warningYellow).
				Padding(0, 1).
				Bold(true)

	ErrorBadgeStyle = lipgloss.NewStyle().
			Foreground(badgeOnError).
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
				Foreground(textPrimary)

	SelectionMarkerStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true)

	ActiveStatusStyle = lipgloss.NewStyle().
				Foreground(brandCyan).
				Bold(true)

	ReadyStatusStyle = lipgloss.NewStyle().
				Foreground(brandIce)
)

// Brand Styles (Gate symbol and wordmark)
var (
	MarkLeftStyle = lipgloss.NewStyle().
			Foreground(brandMint)

	MarkCheckStyle = lipgloss.NewStyle().
			Foreground(brandIce)

	MarkRightStyle = lipgloss.NewStyle().
			Foreground(brandCyan)

	WordmarkSpecStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				Bold(true)

	WordmarkForceStyle = lipgloss.NewStyle().
				Foreground(brandMint).
				Bold(true)
)
