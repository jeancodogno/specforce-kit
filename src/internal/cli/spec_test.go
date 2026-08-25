package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
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

	testSpecInitValidation(t, e, ui)
	testSpecInitCollisions(t, e, ui)
	testSpecInitJSON(t, e, ui)

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func testSpecInitValidation(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Success Default Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "new-spec", false, "", ""); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Success With Custom Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "small-spec", false, "feature", "small"); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Invalid Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "invalid-size-spec", false, "feature", "gigantic"); err == nil {
			t.Errorf("expected error for invalid size, got nil")
		}
	})
}

func testSpecInitCollisions(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Collision Active", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "active-spec", false, "", ""); err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("Collision Archived", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "archived-spec", false, "", ""); err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})
}

func testSpecInitJSON(t *testing.T, e *Executor, ui core.UI) {
	t.Run("JSON Success With Size", func(t *testing.T) {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecInit(context.Background(), ui, "json-spec", true, "feature", "large")
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := buf.String()
		if !strings.Contains(output, `"size": "large"`) {
			t.Errorf("expected size large in json output, got: %s", output)
		}
	})
}

func TestHandleSpecResize(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmpDir := filepath.Join(cwd, "testdata", "cli_resize")
	specDir := filepath.Join(tmpDir, ".specforce", "specs", "resize-target")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.Chdir(tmpDir)

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	if err := e.HandleSpecInit(context.Background(), ui, "resize-target", false, "feature", "small"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	testSpecResizeTransitions(t, e, ui)
	testSpecResizeErrors(t, e, ui)
	testSpecResizeJSON(t, e, ui)
	testSpecResizeDispatch(t, e, ui)

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func testSpecResizeTransitions(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Promotion Small to Large", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "large", false); err != nil {
			t.Fatalf("HandleSpecResize failed: %v", err)
		}
	})

	t.Run("Demotion Large to Small", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "small", false); err != nil {
			t.Fatalf("HandleSpecResize failed: %v", err)
		}
	})
}

func testSpecResizeErrors(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Invalid Size", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "huge", false); err == nil {
			t.Fatal("expected error for invalid size, got nil")
		}
	})

	t.Run("Non-existent Spec", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "non-existent-spec-xyz", "large", false); err == nil {
			t.Fatal("expected error for non-existent spec, got nil")
		}
	})
}

func testSpecResizeJSON(t *testing.T, e *Executor, ui core.UI) {
	t.Run("JSON Output", func(t *testing.T) {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecResize(context.Background(), ui, "resize-target", "complex", true)
		if err != nil {
			t.Errorf("HandleSpecResize failed: %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := buf.String()
		if !strings.Contains(output, `"status": "ok"`) ||
			!strings.Contains(output, `"previous_size": "small"`) ||
			!strings.Contains(output, `"size": "complex"`) {
			t.Errorf("unexpected json output: %s", output)
		}
	})
}

func testSpecResizeDispatch(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Dispatch Resize via HandleSpec", func(t *testing.T) {
		if err := e.HandleSpec(context.Background(), ui, "resize", "resize-target", "--size", "medium"); err != nil {
			t.Fatalf("HandleSpec resize dispatch failed: %v", err)
		}
	})
}

