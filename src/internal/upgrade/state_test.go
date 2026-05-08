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

func TestStateManager_Persistence(t *testing.T) {
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
	if !state.LastCheckedAt.IsZero() {
		t.Errorf("expected zero time for initial load, got %v", state.LastCheckedAt)
	}

	// Test Save
	now := time.Now().Truncate(time.Second)
	state.LastCheckedAt = now
	state.LatestVersion = "v1.0.0"
	state.StagedVersion = "v1.0.0"
	state.IgnoredVersion = "v0.9.0"
	state.UpdateReady = true

	err = mgr.Save(state)
	if err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(statePath)
	if err != nil {
		t.Fatalf("failed to stat state file: %v", err)
	}
	expectedMode := os.FileMode(0600)
	if info.Mode().Perm() != expectedMode {
		t.Errorf("expected file mode %v, got %v", expectedMode, info.Mode().Perm())
	}

	// Test Load (Existing)
	loaded, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load saved state: %v", err)
	}

	if !loaded.LastCheckedAt.Equal(now) {
		t.Errorf("expected LastCheckedAt %v, got %v", now, loaded.LastCheckedAt)
	}
	if loaded.LatestVersion != "v1.0.0" {
		t.Errorf("expected LatestVersion v1.0.0, got %s", loaded.LatestVersion)
	}
	if loaded.StagedVersion != "v1.0.0" {
		t.Errorf("expected StagedVersion v1.0.0, got %s", loaded.StagedVersion)
	}
	if loaded.IgnoredVersion != "v0.9.0" {
		t.Errorf("expected IgnoredVersion v0.9.0, got %s", loaded.IgnoredVersion)
	}
	if !loaded.UpdateReady {
		t.Error("expected UpdateReady to be true")
	}
}

func TestStateManager_Staging(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upgrade-test-staging-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	statePath := filepath.Join(tempDir, "state.json")
	mgr := &StateManager{
		path: statePath,
	}

	// Test EnsureStagedDir
	err = mgr.EnsureStagedDir()
	if err != nil {
		t.Fatalf("failed to ensure staged dir: %v", err)
	}
	stagedDir := mgr.GetStagedDir()
	info, err := os.Stat(stagedDir)
	if err != nil {
		t.Fatalf("staged dir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected staged path to be a directory")
	}
}
