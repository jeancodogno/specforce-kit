package installer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func TestInstall_CancelledContext_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := Install(ctx, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

func TestCheckPermissions_WrapsWithDomainError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-installer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
		t.Fatalf("failed to create read-only dir: %v", err)
	}

	permErr := checkPermissions(nil, readOnlyDir)
	if permErr == nil {
		t.Skip("no permission error in this environment (possibly running as root)")
	}

	if !errors.Is(permErr, core.ErrInstallerPermissionDenied) {
		t.Errorf("expected error to wrap ErrInstallerPermissionDenied, got: %v", permErr)
	}

	if !strings.Contains(permErr.Error(), "failed to check") {
		t.Errorf("expected error message to contain wrapping context, got: %v", permErr)
	}
}

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
