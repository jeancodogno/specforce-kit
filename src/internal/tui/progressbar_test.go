package tui

import (
	"strings"
	"testing"
)

func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		percent  int
		width    int
		contains string
	}{
		{"0 percent", 0, 10, "0%"},
		{"50 percent", 50, 10, "50%"},
		{"100 percent", 100, 10, "100%"},
		{"Underflow", -10, 10, "0%"},
		{"Overflow", 110, 10, "100%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderProgressBar(tt.percent, tt.width)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("RenderProgressBar() = %q, want it to contain %q", got, tt.contains)
			}
		})
	}
}
