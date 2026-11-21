package tui

import (
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// withColorState forces a color profile and background for one test and
// restores the global lipgloss state afterwards. The background is restored
// explicitly so no terminal query is ever sent.
func withColorState(t *testing.T, profile termenv.Profile, dark bool) {
	t.Helper()
	previous := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(profile)
	lipgloss.SetHasDarkBackground(dark)
	t.Cleanup(func() {
		lipgloss.SetColorProfile(previous)
		lipgloss.SetHasDarkBackground(true)
	})
}

func TestPaletteAdaptiveValues(t *testing.T) {
	cases := []struct {
		name        string
		color       lipgloss.AdaptiveColor
		light, dark string
	}{
		{"brandMint", brandMint, "#007A4D", "#00FA9A"},
		{"brandCyan", brandCyan, "#00727F", "#00FFFF"},
		{"brandIce", brandIce, "#3F5F6B", "#E0FFFF"},
		{"successGreen", successGreen, "#1E7A34", "#5FFF87"},
		{"warningYellow", warningYellow, "#8A6100", "#FFFFAF"},
		{"errorRed", errorRed, "#C62828", "#FF5F5F"},
		{"mutedGrey", mutedGrey, "#B0B0B0", "#444444"},
		{"textGrey", textGrey, "#5C5C5C", "#808080"},
		{"textPrimary", textPrimary, "", "#FFFFFF"},
		{"footerBg", footerBg, "#EEEEEE", "#111111"},
		{"badgeOnFill", badgeOnFill, "#FFFFFF", "#000000"},
		{"badgeOnError", badgeOnError, "#FFFFFF", "#FFFFFF"},
	}
	for _, tc := range cases {
		if tc.color.Light != tc.light || tc.color.Dark != tc.dark {
			t.Errorf("%s = {Light:%q Dark:%q}, want {Light:%q Dark:%q}",
				tc.name, tc.color.Light, tc.color.Dark, tc.light, tc.dark)
		}
	}
}

func TestPaletteAliases(t *testing.T) {
	if inProgressYellow != warningYellow {
		t.Errorf("inProgressYellow = %v, want warningYellow %v", inProgressYellow, warningYellow)
	}
	if finishedGreen != successGreen {
		t.Errorf("finishedGreen = %v, want successGreen %v", finishedGreen, successGreen)
	}
}

// seq returns the TrueColor SGR parameters termenv emits for hex.
func seq(hex string, background bool) string {
	return termenv.TrueColor.Color(hex).Sequence(background)
}

func TestPaletteLightBadgeAndBorderStyles(t *testing.T) {
	withColorState(t, termenv.TrueColor, false)
	out := SuccessBadgeStyle.Render("OK")
	if !strings.Contains(out, seq("#FFFFFF", false)) || !strings.Contains(out, seq("#1E7A34", true)) {
		t.Errorf("light success badge must be white on #1E7A34, got %q", out)
	}
	out = ErrorBadgeStyle.Render("ERR")
	if !strings.Contains(out, seq("#FFFFFF", false)) || !strings.Contains(out, seq("#C62828", true)) {
		t.Errorf("light error badge must be white on #C62828, got %q", out)
	}
	if out := FooterStyle.Render("f"); !strings.Contains(out, seq("#EEEEEE", true)) {
		t.Errorf("light footer background must be #EEEEEE, got %q", out)
	}
}

func TestPaletteDarkColorsUnchanged(t *testing.T) {
	withColorState(t, termenv.TrueColor, true)
	if out := HeaderStyle.Render("x"); !strings.Contains(out, seq("#00FA9A", false)) {
		t.Errorf("dark header must stay #00FA9A, got %q", out)
	}
	if out := BodyStyle.Render("x"); !strings.Contains(out, seq("#FFFFFF", false)) {
		t.Errorf("dark primary text must stay #FFFFFF, got %q", out)
	}
}

// relativeLuminance implements the WCAG 2.x relative luminance formula.
func relativeLuminance(t *testing.T, hex string) float64 {
	t.Helper()
	v, err := strconv.ParseUint(strings.TrimPrefix(hex, "#"), 16, 32)
	if err != nil || len(hex) != 7 {
		t.Fatalf("invalid hex color %q", hex)
	}
	channel := func(c uint64) float64 {
		s := float64(c) / 255
		if s <= 0.03928 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	r, g, b := channel(v>>16&0xFF), channel(v>>8&0xFF), channel(v&0xFF)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func contrastRatio(t *testing.T, a, b string) float64 {
	la, lb := relativeLuminance(t, a), relativeLuminance(t, b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func TestLightPaletteContrast(t *testing.T) {
	textColors := map[string]lipgloss.AdaptiveColor{
		"brandMint":     brandMint,
		"brandCyan":     brandCyan,
		"brandIce":      brandIce,
		"successGreen":  successGreen,
		"warningYellow": warningYellow,
		"errorRed":      errorRed,
		"textGrey":      textGrey,
	}
	for name, c := range textColors {
		if ratio := contrastRatio(t, c.Light, "#FFFFFF"); ratio < 4.5 {
			t.Errorf("%s light %s contrast on white = %.2f, want >= 4.5", name, c.Light, ratio)
		}
	}
	// Badge labels sit on the light fills.
	for name, fill := range map[string]lipgloss.AdaptiveColor{"successGreen": successGreen, "warningYellow": warningYellow, "errorRed": errorRed} {
		if ratio := contrastRatio(t, badgeOnFill.Light, fill.Light); ratio < 4.5 {
			t.Errorf("badge label on light %s = %.2f, want >= 4.5", name, ratio)
		}
	}
}

func TestNoStrayColorLiterals(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read tui dir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || name == "theme.go" {
			continue
		}
		data, err := os.ReadFile(name) // #nosec G304 -- test reads files of its own package directory
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, "lipgloss.Color(") {
				t.Errorf("%s:%d: stray lipgloss.Color literal outside theme.go: %s", name, i+1, strings.TrimSpace(line))
			}
		}
	}
}
