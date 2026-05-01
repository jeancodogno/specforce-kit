package upgrade

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewStateManager(t *testing.T) {
	// We need to set up HOME or XDG_CONFIG_HOME if we want to test the default path
	oldHome := os.Getenv("HOME")
	tempHome := t.TempDir()
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	mgr, err := NewStateManager()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mgr == nil {
		t.Fatal("expected non-nil manager")
	}
	if !strings.Contains(mgr.path, ".specforce") {
		t.Errorf("expected path to contain .specforce, got %s", mgr.path)
	}
}

func TestStateManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upgrade-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	statePath := filepath.Join(tempDir, "state.json")
	mgr := &StateManager{
		path: statePath,
	}

	// Test Load (Empty/Non-existent)
	state, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load initial state: %v", err)
	}
	if !state.LastCheckAt.IsZero() {
		t.Errorf("expected zero time for initial load, got %v", state.LastCheckAt)
	}

	// Test Save
	now := time.Now().Truncate(time.Second)
	state.LastCheckAt = now
	state.LatestVersion = "v1.0.0"
	state.IgnoredVersion = "v0.9.0"

	err = mgr.Save(state)
	if err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	// Test Load (Existing)
	loaded, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load saved state: %v", err)
	}

	if !loaded.LastCheckAt.Equal(now) {
		t.Errorf("expected LastCheckAt %v, got %v", now, loaded.LastCheckAt)
	}
	if loaded.LatestVersion != "v1.0.0" {
		t.Errorf("expected LatestVersion v1.0.0, got %s", loaded.LatestVersion)
	}
	if loaded.IgnoredVersion != "v0.9.0" {
		t.Errorf("expected IgnoredVersion v0.9.0, got %s", loaded.IgnoredVersion)
	}
}
