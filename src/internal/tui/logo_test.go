package tui

import (
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

// asciiBrand forces the Ascii profile (plain glyphs) and clears brand caches.
func asciiBrand(t *testing.T) {
	t.Helper()
	withColorState(t, termenv.Ascii, true)
	resetBrandCaches()
	t.Cleanup(resetBrandCaches)
}

func resetBrandCaches() {
	cachedHeader = brandCache{}
	cachedFullMark = [2]brandCache{}
}

func trimmedLines(s string) []string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " ")
	}
	return lines
}

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFullMark_SideBySide(t *testing.T) {
	asciiBrand(t)
	assertLines(t, trimmedLines(renderFullMark("1.0.0", 80)), []string{
		"",
		"┌─    ╱─┐  specforce v1.0.0",
		"│ ╲  ╱  │  Spec-Driven Development for AI agents",
		"└─ ╲╱  ─┘",
	})
}

func TestFullMark_SideBySideAtThreshold(t *testing.T) {
	asciiBrand(t)
	lines := trimmedLines(renderFullMark("1.0.0", 48))
	if len(lines) != 4 || !strings.HasSuffix(lines[2], Tagline) {
		t.Errorf("width 48 must render side by side, got %q", lines)
	}
}

func TestFullMark_StackedWhenNarrow(t *testing.T) {
	asciiBrand(t)
	assertLines(t, trimmedLines(renderFullMark("1.0.0", 40)), []string{
		"",
		"┌─    ╱─┐",
		"│ ╲  ╱  │",
		"└─ ╲╱  ─┘",
		"specforce v1.0.0",
		"Spec-Driven Development for AI agents",
	})
}

func TestFullMark_VersionNormalization(t *testing.T) {
	for _, v := range []string{"1.0.0", "v1.0.0"} {
		if got := formatVersion(v); got != "v1.0.0" {
			t.Errorf("formatVersion(%q) = %q, want v1.0.0", v, got)
		}
	}
	asciiBrand(t)
	out := renderFullMark("v1.0.0", 80)
	if !strings.Contains(out, "specforce v1.0.0") || strings.Contains(out, "vv") {
		t.Errorf("full mark must show exactly one leading v, got:\n%s", out)
	}
}

func TestFullMark_CachedPerLayout(t *testing.T) {
	asciiBrand(t)
	wide := fullMark("1.0.0", 80)
	narrow := fullMark("1.0.0", 40)
	if wide == narrow {
		t.Fatal("side-by-side and stacked layouts must differ")
	}
	if fullMark("1.0.0", 100) != wide || fullMark("1.0.0", 30) != narrow {
		t.Error("each layout must be reused from its cache")
	}
}

func TestFullMark_ColoredSegments(t *testing.T) {
	withColorState(t, termenv.TrueColor, true)
	resetBrandCaches()
	t.Cleanup(resetBrandCaches)
	out := renderFullMark("1.0.0", 80)
	for _, hex := range []string{"#00FA9A", "#E0FFFF", "#00FFFF"} {
		if !strings.Contains(out, seq(hex, false)) {
			t.Errorf("full mark must use %s, got %q", hex, out)
		}
	}
}

func TestCompactHeader(t *testing.T) {
	asciiBrand(t)
	for _, v := range []string{"1.0.0", "v1.0.0"} {
		resetBrandCaches()
		if got := RenderCompactHeader(v); got != "[✓] specforce v1.0.0" {
			t.Errorf("RenderCompactHeader(%q) = %q, want %q", v, got, "[✓] specforce v1.0.0")
		}
	}
}

func TestBrandOutputsUseLowercaseWordmark(t *testing.T) {
	asciiBrand(t)
	outputs := map[string]string{
		"full mark wide":   renderFullMark("1.0.0", 80),
		"full mark narrow": renderFullMark("1.0.0", 40),
		"compact header":   RenderCompactHeader("1.0.0"),
		"footer":           captureStdout(t, func() { PrintFooter("1.0.0") }),
	}
	for name, out := range outputs {
		if strings.Contains(out, "SPECFORCE") || strings.Contains(out, "Specforce") {
			t.Errorf("%s must not contain an uppercase wordmark, got:\n%s", name, out)
		}
	}
}
