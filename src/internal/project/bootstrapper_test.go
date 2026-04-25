package project

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

type mockUI struct {
	warnCalled bool
}

func (m *mockUI) Log(msg string)          {}
func (m *mockUI) Success(msg string)      {}
func (m *mockUI) Warn(msg string)         { m.warnCalled = true }
func (m *mockUI) Error(msg string)        {}
func (m *mockUI) SubTask(msg string)      {}
func (m *mockUI) StartSpinner(msg string) {}
func (m *mockUI) StopSpinner()             {}
func (m *mockUI) Confirm(msg string) bool  { return true }

func TestBootstrapProject_Success(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Mock artifacts filesystem
	artifactsFS := fstest.MapFS{}

	err := BootstrapProject(ctx, tmpDir, nil, artifactsFS, &mockUI{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dirs := []string{
		".specforce/docs",
		".specforce/specs",
		".specforce/archive",
	}

	for _, dir := range dirs {
		path := filepath.Join(tmpDir, dir)
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist", path)
		}
		if _, err := os.Stat(filepath.Join(path, ".gitkeep")); os.IsNotExist(err) {
			t.Errorf("expected .gitkeep in %s", path)
		}
	}

	// Verify .specforce/templates does NOT exist
	templatesPath := filepath.Join(tmpDir, ".specforce/templates")
	if _, err := os.Stat(templatesPath); !os.IsNotExist(err) {
		t.Errorf("expected directory %s NOT to exist", templatesPath)
	}
}

func TestBootstrapProject_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()
	artifactsFS := fstest.MapFS{}

	t.Run("already initialized", func(t *testing.T) {
		if err := os.Mkdir(filepath.Join(tmpDir, ".specforce"), 0755); err != nil {
			t.Fatal(err)
		}
		err := BootstrapProject(ctx, tmpDir, nil, artifactsFS, nil)
		if err == nil {
			t.Fatal("expected error when already initialized")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		err := BootstrapProject(cancelCtx, t.TempDir(), nil, artifactsFS, nil)
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})

	t.Run("nil ui success", func(t *testing.T) {
		err := BootstrapProject(ctx, t.TempDir(), nil, artifactsFS, nil)
		if err != nil {
			t.Fatalf("BootstrapProject failed with nil UI: %v", err)
		}
	})
}
