package installer

import (
	"testing"
)

func TestShouldInstall_ToolsOnly(t *testing.T) {
	opts := Options{ToolsOnly: true}

	tests := []struct {
		path     string
		expected bool
	}{
		{".gemini/agents/spf.toml", true},
		{".claude/commands/spf.md", true},
		{".opencode/skills/tdd.md", true},
		{".specforce/config.yaml", false},
		{".specforce/docs/architecture.md", false},
		{"README.md", false},
		{"go.mod", false},
		{".agent/workflows/spf.md", true},
	}

	for _, tt := range tests {
		got := ShouldInstall(tt.path, opts)
		if got != tt.expected {
			t.Errorf("ShouldInstall(%q, ToolsOnly: true) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}

func TestShouldInstall_All(t *testing.T) {
	opts := Options{ToolsOnly: false}

	tests := []struct {
		path     string
		expected bool
	}{
		{".gemini/agents/spf.toml", true},
		{".specforce/config.yaml", true},
		{"README.md", true},
	}

	for _, tt := range tests {
		got := ShouldInstall(tt.path, opts)
		if got != tt.expected {
			t.Errorf("ShouldInstall(%q, ToolsOnly: false) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}
