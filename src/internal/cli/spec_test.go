package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestHandleSpecList(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	t.Run("Populated List", func(t *testing.T) {
		tmpDir := filepath.Join(cwd, "testdata", "cli_spec_populated")
		if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "test-spec"), 0755); err != nil {
			t.Fatal(err)
		}
		_ = os.Chdir(tmpDir)

		e := NewExecutor("1.0.0")
		ui := tui.NewUI()

		err := e.HandleSpecList(context.Background(), ui, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = e.HandleSpecList(context.Background(), ui, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty List JSON", func(t *testing.T) {
		tmpDir := filepath.Join(cwd, "testdata", "cli_spec_empty")
		if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs"), 0755); err != nil {
			t.Fatal(err)
		}
		_ = os.Chdir(tmpDir)

		e := NewExecutor("1.0.0")
		ui := tui.NewUI()

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecList(context.Background(), ui, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := strings.TrimSpace(buf.String())
		if output != "[]" {
			t.Errorf("expected empty JSON array [], got %q", output)
		}
	})

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func TestHandleSpecInit(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	// Setup temp project root
	tmpDir := filepath.Join(cwd, "testdata", "cli_init")
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "active-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "archive", "archived-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.Chdir(tmpDir)

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	t.Run("Success", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "new-spec", false, "")
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Collision Active", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "active-spec", false, "")
		if err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("Collision Archived", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "archived-spec", false, "")
		if err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("JSON Success", func(t *testing.T) {
		err := e.HandleSpecInit(context.Background(), ui, "json-spec", true, "")
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}
