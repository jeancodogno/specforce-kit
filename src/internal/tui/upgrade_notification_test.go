package tui

import (
	"strings"
	"testing"
)

func TestManualUpdateNotification_NamesVersionAndCommand(t *testing.T) {
	cases := []struct {
		name    string
		version string
		command string
	}{
		{"binary", "v1.0.0-beta.3", "Download v1.0.0-beta.3 from https://github.com/jeancodogno/specforce-kit/releases/tag/v1.0.0-beta.3 or run: sudo make install"},
		{"npm", "v1.0.0-beta.3", "npm install -g @jeancodogno/specforce-kit@1.0.0-beta.3"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := RenderManualUpdateNotification("1.0.0-beta.2", tc.version, tc.command)
			flat := strings.Join(strings.Fields(out), " ")
			if !strings.Contains(out, tc.version) {
				t.Errorf("manual notice must name the version %s, got:\n%s", tc.version, out)
			}
			if !strings.Contains(flat, strings.Join(strings.Fields(tc.command), " ")) {
				t.Errorf("manual notice must contain the exact command %q, got:\n%s", tc.command, out)
			}
			if strings.Contains(flat, "applied automatically") {
				t.Errorf("manual notice must never claim automatic application, got:\n%s", out)
			}
		})
	}
}

func TestUpdateAppliedNotification_NamesVersion(t *testing.T) {
	out := RenderUpdateAppliedNotification("v1.0.0-beta.3")
	if !strings.Contains(out, "v1.0.0-beta.3") {
		t.Errorf("applied notice must name the version, got %q", out)
	}
}

func TestUpdateNotification_Available(t *testing.T) {
	out := RenderUpdateNotification("1.0.0-beta.2", "v1.0.0-beta.3")
	if !strings.Contains(out, "v1.0.0-beta.3") {
		t.Errorf("available notice must name the latest version, got:\n%s", out)
	}
}
