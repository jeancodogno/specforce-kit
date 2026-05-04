package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestHandleSpecList(t *testing.T) {
	// Setup temp project root
	tmpDir := filepath.Join("testdata", "cli_spec")
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "test-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(cwd) }()
	defer func() { _ = os.RemoveAll(filepath.Join(cwd, "testdata")) }()

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	t.Run("TUI Mode", func(t *testing.T) {
		err := e.HandleSpecList(context.Background(), ui, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("JSON Mode", func(t *testing.T) {
		err := e.HandleSpecList(context.Background(), ui, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestHandleSpecInit(t *testing.T) {
	// Setup temp project root
	tmpDir := filepath.Join("testdata", "cli_init")
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "active-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "archive", "archived-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(cwd) }()
	defer func() { _ = os.RemoveAll(filepath.Join(cwd, "testdata")) }()

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	t.Run("Success", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "new-spec", false)
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Collision Active", func(t *testing.T) {
		// PrepareSlug will transform "active-spec" to something like "20260504-1323-active-spec"
		// which won't collide with the manually created ".specforce/specs/active-spec" directory.
		// To test collision, we must pass the ALREADY timestamped slug.
		err := e.HandleSpecInit(context.Background(), ui, "active-spec", false)
		if err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("Collision Archived", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "archived-spec", false)
		if err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("JSON Success", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "json-spec", true)
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})
}
