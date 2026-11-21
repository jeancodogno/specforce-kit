package tui

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

// captureStdout runs fn and returns what it wrote to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()
	_ = w.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(data)
}

func TestFooter_Text(t *testing.T) {
	withColorState(t, termenv.Ascii, true)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for _, v := range []string{"1.0.0", "v1.0.0"} {
		out := strings.TrimSpace(captureStdout(t, func() { PrintFooter(v) }))
		if !strings.HasPrefix(out, "specforce v") {
			t.Errorf("footer must start with %q, got %q", "specforce v", out)
		}
		if want := "specforce v1.0.0 | Workspace: " + cwd; out != want {
			t.Errorf("footer = %q, want %q", out, want)
		}
	}
}

func TestFooter_BadgeLabelsSurviveNoColor(t *testing.T) {
	withColorState(t, termenv.Ascii, true)
	for status, label := range map[string]string{"success": "SUCCESS", "warning": "WARNING", "error": "ERROR", "info": "INFO"} {
		out := RenderBadge(status, "msg")
		if !strings.Contains(out, label) || strings.Contains(out, "\x1b[") {
			t.Errorf("badge %s must show %s without escape sequences, got %q", status, label, out)
		}
	}
}
