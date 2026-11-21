package tui

import (
	"fmt"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// Tagline is the one-line product description shown with the brand mark.
const Tagline = "Spec-Driven Development for AI agents"

// gateMark is the 3-line Gate symbol; each line is split into the left
// bracket (Mint), the check mark (Ice) and the right bracket (Cyan).
var gateMark = [3][3]string{
	{"┌─", "    ╱", "─┐"},
	{"│ ", "╲  ╱ ", " │"},
	{"└─", " ╲╱  ", "─┘"},
}

// Layout of the full mark: text beside the symbol, or below it.
const (
	layoutSideBySide = iota
	layoutStacked
)

// sideBySideMinWidth is the narrowest terminal that fits mark + gap + tagline.
const sideBySideMinWidth = 48

// brandCache memoizes one rendered brand form for a given version.
type brandCache struct {
	key   string
	value string
}

var (
	cachedHeader   brandCache
	cachedFullMark [2]brandCache
)

// formatVersion normalizes a version to exactly one leading "v".
func formatVersion(version string) string {
	return "v" + strings.TrimPrefix(version, "v")
}

// renderWordmark renders the lowercase wordmark: "spec" primary, "force" Mint.
func renderWordmark() string {
	return WordmarkSpecStyle.Render("spec") + WordmarkForceStyle.Render("force")
}

// renderBrandLine renders "specforce v<version>" with the version dimmed.
func renderBrandLine(version string) string {
	return renderWordmark() + " " + DimmedStyle.Render(formatVersion(version))
}

// renderMarkLines renders the three colored lines of the Gate symbol.
func renderMarkLines() []string {
	lines := make([]string, 0, len(gateMark))
	for _, l := range gateMark {
		lines = append(lines, MarkLeftStyle.Render(l[0])+MarkCheckStyle.Render(l[1])+MarkRightStyle.Render(l[2]))
	}
	return lines
}

func layoutFor(width int) int {
	if width >= sideBySideMinWidth {
		return layoutSideBySide
	}
	return layoutStacked
}

// renderFullMark renders the full Gate mark with wordmark, version and
// tagline, preceded by one blank line of top padding.
func renderFullMark(version string, width int) string {
	mark := renderMarkLines()
	text := []string{renderBrandLine(version), DimmedStyle.Render(Tagline)}

	var lines []string
	if layoutFor(width) == layoutStacked {
		lines = append(mark, text...)
	} else {
		lines = []string{mark[0] + "  " + text[0], mark[1] + "  " + text[1], mark[2]}
	}
	return "\n" + strings.Join(lines, "\n")
}

// fullMark returns the full mark, rendered at most once per layout.
func fullMark(version string, width int) string {
	c := &cachedFullMark[layoutFor(width)]
	if c.value == "" || c.key != version {
		*c = brandCache{key: version, value: renderFullMark(version, width)}
	}
	return c.value
}

// RenderCompactHeader returns the one-line brand header "[✓] specforce v<version>".
func RenderCompactHeader(version string) string {
	if cachedHeader.value != "" && cachedHeader.key == version {
		return cachedHeader.value
	}
	symbol := MarkLeftStyle.Render("[") + MarkCheckStyle.Render("✓") + MarkRightStyle.Render("]")
	cachedHeader = brandCache{key: version, value: symbol + " " + renderBrandLine(version)}
	return cachedHeader.value
}

// PrintHeader prints the compact brand header used on working screens. The
// blank line that follows it comes from the leading newline every caller
// already prints before its title.
func PrintHeader() {
	fmt.Println(RenderCompactHeader(core.Version))
}

// PrintBranding prints the full Gate mark used on welcome surfaces.
func PrintBranding() {
	fmt.Println(fullMark(core.Version, GetTerminalWidth()))
}
