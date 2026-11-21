package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// fakeEnv returns a getenv func backed by a map.
func fakeEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

// countingQuery returns a query func that reports result and counts calls.
func countingQuery(result bool, calls *int) func() bool {
	return func() bool {
		*calls++
		return result
	}
}

func TestDetectDarkBackground_COLORFGBG(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"0;15", false},
		{"15;0", true},
		{"0;7", false},
		{"0;8", true},
		{"0;9", false},
		{"0;16", true},
		{"0;default;15", false},
		{"15;default;0", true},
	}
	for _, tc := range cases {
		calls := 0
		// The query would answer the opposite; it must never be consulted.
		got := detectDarkBackground(fakeEnv(map[string]string{"COLORFGBG": tc.value}), true, countingQuery(!tc.want, &calls))
		if got != tc.want {
			t.Errorf("COLORFGBG=%q: dark = %v, want %v", tc.value, got, tc.want)
		}
		if calls != 0 {
			t.Errorf("COLORFGBG=%q: query called %d times, want 0", tc.value, calls)
		}
	}
}

func TestDetectDarkBackground_UnparseableFallsThrough(t *testing.T) {
	env := fakeEnv(map[string]string{"COLORFGBG": "default;default"})

	calls := 0
	if got := detectDarkBackground(env, true, countingQuery(false, &calls)); got != false || calls != 1 {
		t.Errorf("interactive fall-through: dark = %v (calls %d), want query result false with 1 call", got, calls)
	}

	calls = 0
	if got := detectDarkBackground(env, false, countingQuery(false, &calls)); got != true || calls != 0 {
		t.Errorf("non-interactive fall-through: dark = %v (calls %d), want true with 0 calls", got, calls)
	}
}

func TestDetectDarkBackground_NonInteractiveIsDarkWithoutQuery(t *testing.T) {
	calls := 0
	if got := detectDarkBackground(fakeEnv(nil), false, countingQuery(false, &calls)); !got {
		t.Error("non-interactive output must fall back to dark")
	}
	if calls != 0 {
		t.Errorf("non-interactive output must not query the terminal, got %d calls", calls)
	}
}

func TestDetectDarkBackground_InteractiveQueriesOnce(t *testing.T) {
	for _, answer := range []bool{true, false} {
		calls := 0
		if got := detectDarkBackground(fakeEnv(nil), true, countingQuery(answer, &calls)); got != answer {
			t.Errorf("query answered %v, got %v", answer, got)
		}
		if calls != 1 {
			t.Errorf("query called %d times, want 1", calls)
		}
	}
}

func TestConfigureColor_NoColorForcesAscii(t *testing.T) {
	withColorState(t, termenv.TrueColor, true)
	calls := 0
	ConfigureColor(fakeEnv(map[string]string{"NO_COLOR": "1", "COLORFGBG": "0;15"}), true, countingQuery(false, &calls))

	if p := lipgloss.ColorProfile(); p != termenv.Ascii {
		t.Errorf("NO_COLOR=1: profile = %v, want Ascii", p)
	}
	if calls != 0 {
		t.Errorf("NO_COLOR=1: query called %d times, want 0", calls)
	}
	if out := HeaderStyle.Render("x") + ErrorBadgeStyle.Render("ERROR"); strings.Contains(out, "\x1b[") {
		t.Errorf("NO_COLOR=1: output must contain no color sequence, got %q", out)
	}
}

func TestConfigureColor_EmptyNoColorKeepsColors(t *testing.T) {
	withColorState(t, termenv.TrueColor, true)
	ConfigureColor(fakeEnv(map[string]string{"NO_COLOR": "", "COLORFGBG": "15;0"}), true, func() bool {
		t.Error("query must not be called when COLORFGBG is parseable")
		return true
	})

	if p := lipgloss.ColorProfile(); p != termenv.TrueColor {
		t.Errorf("empty NO_COLOR: profile = %v, want TrueColor", p)
	}
	if out := HeaderStyle.Render("x"); !strings.Contains(out, seq("#00FA9A", false)) {
		t.Errorf("empty NO_COLOR on dark: header must use #00FA9A, got %q", out)
	}
}

func TestConfigureColor_LightBackground(t *testing.T) {
	withColorState(t, termenv.TrueColor, true)
	ConfigureColor(fakeEnv(map[string]string{"COLORFGBG": "0;15"}), true, func() bool { return true })

	if lipgloss.HasDarkBackground() {
		t.Fatal("COLORFGBG=0;15 must select the light background")
	}
	if out := BodyStyle.Render("x"); strings.Contains(out, "38;") {
		t.Errorf("light primary text must use the terminal default foreground, got %q", out)
	}
	if out := HeaderStyle.Render("x"); !strings.Contains(out, seq("#007A4D", false)) {
		t.Errorf("light header must use #007A4D, got %q", out)
	}
}
