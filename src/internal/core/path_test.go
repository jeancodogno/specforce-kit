package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath_Tilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"tilde expansion", "~/test", filepath.Join(home, "test")},
		{"tilde dot expansion", "~/.codex", filepath.Join(home, ".codex")},
		{"no tilde", "/absolute/path", "/absolute/path"},
		{"relative path", "relative/path", "relative/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandPath(tt.input)
			if got != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExpandPath_EnvVars(t *testing.T) {
	if err := os.Setenv("SET_VAR", "value"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() { _ = os.Unsetenv("SET_VAR") }()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"env var expansion (set)", "${SET_VAR:-default}", "value"},
		{"env var expansion (unset)", "${UNSET_VAR:-default}", "default"},
		{"env var expansion (no fallback, set)", "${SET_VAR}", "value"},
		{"env var expansion (no fallback, unset)", "${UNSET_VAR}", ""},
		{"nested env var expansion", "${UNSET_VAR:-${SET_VAR:-fallback}}", "value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandPath(tt.input)
			if got != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExpandPath_Complex(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir: %v", err)
	}

	if err := os.Setenv("PROJECT_NAME", "my-project"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() { _ = os.Unsetenv("PROJECT_NAME") }()

	input := "~/prompts/${PROJECT_NAME:-default}/spf"
	expected := filepath.Join(home, "prompts/my-project/spf")
	got := ExpandPath(input)
	if got != expected {
		t.Errorf("ExpandPath(%q) = %q, want %q", input, got, expected)
	}
}

func TestExpandPath_Malformed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"missing closing brace", "${VAR", "${VAR"},
		{"unbalanced nested", "${VAR:-${NESTED", "${VAR:-${NESTED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandPath(tt.input)
			if got != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
