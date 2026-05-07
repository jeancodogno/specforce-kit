package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestService_PerformUpgrade(t *testing.T) {
	tempDir := t.TempDir()
	mgr := &StateManager{path: filepath.Join(tempDir, "state.json")}
	svc := NewService(mgr, &MockProvider{}, "v1.0.0")

	// We can't easily test the full success path without a real server or lots of mocking,
	// but we can test that it correctly identifies NPM or tries binary.
	// We'll just call it and expect it to fail (as it will try to download from github).
	err := svc.PerformUpgrade(context.Background(), "v2.0.0")
	if err == nil {
		t.Error("expected error from PerformUpgrade in test environment, got nil")
	}
}

func TestServiceCheckForUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upgrade-svc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	mgr := &StateManager{path: filepath.Join(tempDir, "state.json")}
	mockProvider := &MockProvider{Version: "v2.0.0"}
	
	fakeNow := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	svc := NewService(mgr, mockProvider, "v1.0.0")
	svc.now = func() time.Time { return fakeNow }

	// 1. First check - should trigger
	// 1. Initial check should trigger background process
	svc.CheckForUpdate()
	
	// Since we can't easily verify the detached process in unit test without complex mocking,
	// we assume it doesn't crash and returns.
	
	// 2. Second check immediately - should be throttled
	// No way to verify easily without checking state, but CheckForUpdate doesn't update state itself anymore,
	// it's the spawned process that does.
}

func TestServiceIsUpdateAvailable(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upgrade-svc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	mgr := &StateManager{path: filepath.Join(tempDir, "state.json")}
	svc := NewService(mgr, &MockProvider{}, "v1.0.0")

	// No state yet
	available, _ := svc.IsUpdateAvailable()
	if available {
		t.Error("expected no update available when no state exists")
	}

	// Newer version available
	_ = mgr.Save(&State{LatestVersion: "v1.1.0"})
	available, ver := svc.IsUpdateAvailable()
	if !available || ver != "v1.1.0" {
		t.Errorf("expected update v1.1.0 available, got %v, %s", available, ver)
	}

	// Newer version is ignored
	_ = mgr.Save(&State{LatestVersion: "v1.1.0", IgnoredVersion: "v1.1.0"})
	available, _ = svc.IsUpdateAvailable()
	if available {
		t.Error("expected update to be ignored")
	}

	// Current version is same or newer
	_ = mgr.Save(&State{LatestVersion: "v1.0.0"})
	available, _ = svc.IsUpdateAvailable()
	if available {
		t.Error("expected no update when versions are equal")
	}
}

func TestService_IsNPM_Detection(t *testing.T) {
	svc := &Service{}
	// This will call os.Executable() which is likely 'specforce.test' or similar.
	// It should return false unless the test binary name matches NPM criteria.
	_ = svc.isNPM()
}

func TestService_PerformAtomicSwap(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "swap-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	mgr := &StateManager{path: filepath.Join(tempDir, "state.json")}
	svc := NewService(mgr, &MockProvider{}, "v1.0.0")

	// Create a dummy staged binary
	_ = mgr.EnsureStagedDir()
	stagedPath := filepath.Join(mgr.GetStagedDir(), "specforce_next")
	_ = os.WriteFile(stagedPath, []byte("new content"), 0644)

	// Create a dummy active binary
	activePath := filepath.Join(tempDir, "specforce")
	_ = os.WriteFile(activePath, []byte("old content"), 0755)

	// Setup state
	_ = mgr.Save(&State{UpdateReady: true, StagedVersion: "v1.1.0"})

	// Perform swap
	err = svc.PerformAtomicSwapAt(activePath)
	if err != nil {
		t.Fatalf("PerformAtomicSwapAt failed: %v", err)
	}

	// Verify active binary is now new content
	content, _ := os.ReadFile(activePath)
	if string(content) != "new content" {
		t.Errorf("expected active binary to have new content, got %s", string(content))
	}

	// Verify .old exists
	if _, err := os.Stat(activePath + ".old"); err != nil {
		t.Error("expected .old backup to exist")
	}

	// Verify state reset
	state, _ := mgr.Load()
	if state.UpdateReady {
		t.Error("expected UpdateReady to be false after swap")
	}
}
