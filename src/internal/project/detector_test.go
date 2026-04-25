package project

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
)

func TestDetectExistingAgents(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	// Mock kitFS
	kitFS := fstest.MapFS{
		"kit.yaml": &fstest.MapFile{Data: []byte(`
tools:
  gemini:
    name: "Gemini"
    target: ".gemini/"
  claude:
    name: "Claude"
    target: ".claude/"
`)},
		"agents/gemini.yaml": &fstest.MapFile{Data: []byte("content: gemini")},
		"agents/claude.yaml": &fstest.MapFile{Data: []byte("content: claude")},
	}

	reg := &agent.Registry{}
	_ = reg.Initialize(kitFS)

	t.Run("no existing agents", func(t *testing.T) {
		got := DetectExistingAgents(ctx, tmpDir, reg)
		if len(got) != 0 {
			t.Errorf("expected no agents, got %v", got)
		}
	})

	t.Run("existing agents", func(t *testing.T) {
		if err := os.Mkdir(filepath.Join(tmpDir, ".gemini"), 0755); err != nil {
			t.Fatal(err)
		}
		got := DetectExistingAgents(ctx, tmpDir, reg)
		if len(got) != 1 || got[0] != "gemini" {
			t.Errorf("expected [gemini], got %v", got)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		got := DetectExistingAgents(ctx, tmpDir, nil)
		if len(got) != 0 {
			t.Errorf("expected no agents for nil registry, got %v", got)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		got := DetectExistingAgents(cancelCtx, tmpDir, reg)
		if len(got) != 0 {
			t.Errorf("expected no agents for cancelled context, got %v", got)
		}
	})
}

func TestIsInitialized(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("not initialized", func(t *testing.T) {
		if IsInitialized(tmpDir) {
			t.Fatal("expected IsInitialized to be false")
		}
	})

	t.Run("is initialized", func(t *testing.T) {
		if err := os.Mkdir(filepath.Join(tmpDir, ".specforce"), 0755); err != nil {
			t.Fatal(err)
		}
		if !IsInitialized(tmpDir) {
			t.Fatal("expected IsInitialized to be true")
		}
	})

	t.Run("is file not dir", func(t *testing.T) {
		anotherTmp := t.TempDir()
		if err := os.WriteFile(filepath.Join(anotherTmp, ".specforce"), []byte("not a dir"), 0644); err != nil {
			t.Fatal(err)
		}
		if IsInitialized(anotherTmp) {
			t.Fatal("expected IsInitialized to be false when .specforce is a file")
		}
	})
}
