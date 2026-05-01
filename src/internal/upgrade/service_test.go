package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
	svc.CheckForUpdate(context.Background())
	
	// Wait a bit for background goroutine
	time.Sleep(100 * time.Millisecond)

	state, _ := mgr.Load()
	if state.LatestVersion != "v2.0.0" {
		t.Errorf("expected LatestVersion v2.0.0, got %s", state.LatestVersion)
	}
	if !state.LastCheckAt.Equal(fakeNow) {
		t.Errorf("expected LastCheckAt %v, got %v", fakeNow, state.LastCheckAt)
	}

	// 2. Second check immediately - should be throttled
	mockProvider.Version = "v3.0.0"
	svc.CheckForUpdate(context.Background())
	time.Sleep(100 * time.Millisecond)

	state, _ = mgr.Load()
	if state.LatestVersion != "v2.0.0" {
		t.Errorf("expected LatestVersion v2.0.0 (throttled), got %s", state.LatestVersion)
	}

	// 3. Third check after 25h - should trigger
	svc.now = func() time.Time { return fakeNow.Add(25 * time.Hour) }
	svc.CheckForUpdate(context.Background())
	time.Sleep(100 * time.Millisecond)

	state, _ = mgr.Load()
	if state.LatestVersion != "v3.0.0" {
		t.Errorf("expected LatestVersion v3.0.0, got %s", state.LatestVersion)
	}
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

func TestServiceIsNPM(t *testing.T) {
	svc := &Service{}
	// This will depend on the environment, but we can verify it doesn't crash
	_ = svc.isNPM()
}
