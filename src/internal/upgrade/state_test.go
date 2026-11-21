package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewStateManager(t *testing.T) {
	// The default path lives under the home directory: isolate it.
	setTestHome(t)

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
	state.UpdateReady = true

	err = mgr.Save(state)
	if err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	// Verify file permissions
	t.Run("owner-only file mode", func(t *testing.T) {
		skipWithoutPermissionBits(t)
		info, err := os.Stat(statePath)
		if err != nil {
			t.Fatalf("failed to stat state file: %v", err)
		}
		expectedMode := os.FileMode(0600)
		if info.Mode().Perm() != expectedMode {
			t.Errorf("expected file mode %v, got %v", expectedMode, info.Mode().Perm())
		}
	})

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

// --- Legacy state tests (FIX-9) ---

func TestLegacyState_IgnoredKeyHasNoEffect(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	t.Setenv("TMPDIR", t.TempDir())

	// A state file written by an older release that still had the ignore feature.
	legacy := []byte(`{"latest_version": "v2.0.0", "ignored_version": "v2.0.0", "update_ready": false}`)
	if err := os.MkdirAll(filepath.Dir(f.mgr.path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.mgr.path, legacy, 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := f.mgr.Load(); err != nil {
		t.Fatalf("expected legacy state to load, got %v", err)
	}
	if available, ver := f.svc.IsUpdateAvailable(); !available || ver != "v2.0.0" {
		t.Errorf("expected v2.0.0 reported as available, got %v, %q", available, ver)
	}

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	if *f.requests == 0 {
		t.Error("expected the ignored version to be downloaded anyway")
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !s.UpdateReady || s.StagedVersion != "v2.0.0" {
		t.Errorf("expected v2.0.0 staged and ready, got %+v", s)
	}
}

// --- Unreadable state tests (FIX-4) ---

func TestStateRepair_UnparseableStateLoadsAsEmpty(t *testing.T) {
	cases := map[string]string{
		"garbage":    "not json at all \x00\xff",
		"truncated":  `{"latest_version": "v2.0.0", "update_ready": tr`,
		"wrongTyped": `{"update_ready": "yes", "last_checked_at": 42}`,
		"empty":      "",
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			setTestHome(t)
			mgr, err := NewStateManager()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(mgr.path), 0750); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(mgr.path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			if err := mgr.EnsureStagedDir(); err != nil {
				t.Fatal(err)
			}
			staged := filepath.Join(mgr.GetStagedDir(), "specforce_next")
			if err := os.WriteFile(staged, []byte("orphan"), 0600); err != nil {
				t.Fatal(err)
			}

			s, err := mgr.Load()
			if err != nil {
				t.Fatalf("expected unparseable state to load as empty, got %v", err)
			}
			if *s != (State{}) {
				t.Errorf("expected an empty state, got %+v", s)
			}
			if _, err := os.Stat(staged); !os.IsNotExist(err) {
				t.Errorf("expected the orphan staged binary to be deleted, stat err = %v", err)
			}
		})
	}
}

func TestStateRepair_ReadErrorIsReturned(t *testing.T) {
	setTestHome(t)
	mgr, err := NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(mgr.path, 0750); err != nil {
		t.Fatal(err)
	}
	if err := mgr.EnsureStagedDir(); err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(mgr.GetStagedDir(), "specforce_next")
	if err := os.WriteFile(staged, []byte("staged"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := mgr.Load(); err == nil {
		t.Error("expected a read error for a directory at the state path")
	}
	if _, err := os.Stat(staged); err != nil {
		t.Errorf("a read error must not delete the staged binary: %v", err)
	}
}
