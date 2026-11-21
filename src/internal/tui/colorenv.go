package tui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// ConfigureColor decides, once at process start and before any output, how
// colors are rendered. A non-empty NO_COLOR disables color entirely; otherwise
// the terminal background (light or dark) is detected and fixed so lipgloss
// never queries the terminal again.
func ConfigureColor(getenv func(string) string, interactive bool, query func() bool) {
	if getenv("NO_COLOR") != "" {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}
	lipgloss.SetHasDarkBackground(detectDarkBackground(getenv, interactive, query))
}

// detectDarkBackground reports whether the terminal background is dark.
// Precedence: COLORFGBG, then dark for non-interactive output (no query),
// then the terminal background query.
func detectDarkBackground(getenv func(string) string, interactive bool, query func() bool) bool {
	if dark, ok := parseColorFGBG(getenv("COLORFGBG")); ok {
		return dark
	}
	if !interactive {
		return true
	}
	return query()
}

// parseColorFGBG reads the background index from the last ";"-separated
// field of COLORFGBG. Indexes 7 and 9-15 are light backgrounds.
func parseColorFGBG(value string) (dark bool, ok bool) {
	if value == "" {
		return false, false
	}
	fields := strings.Split(value, ";")
	bg, err := strconv.Atoi(strings.TrimSpace(fields[len(fields)-1]))
	if err != nil {
		return false, false
	}
	light := bg == 7 || (bg >= 9 && bg <= 15)
	return !light, true
}
