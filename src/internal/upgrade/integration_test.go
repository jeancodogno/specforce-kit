package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIntegration_UpdateFlow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upgrade-integration-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	statePath := filepath.Join(tempDir, "state.json")
	mgr := &StateManager{path: statePath}
	
	// v1.0.0 -> v2.0.0
	mockProvider := &MockProvider{Version: "v2.0.0"}
	svc := NewService(mgr, mockProvider, "v1.0.0")
	svc.now = func() time.Time { return time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC) }

	// 1. Check for update
	svc.CheckForUpdate(context.Background())
	time.Sleep(100 * time.Millisecond)

	// 2. Verify state updated
	state, err := mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if state.LatestVersion != "v2.0.0" {
		t.Errorf("expected latest version v2.0.0, got %s", state.LatestVersion)
	}

	// 3. Check if available
	available, latest := svc.IsUpdateAvailable()
	if !available || latest != "v2.0.0" {
		t.Errorf("expected update available v2.0.0, got %v, %s", available, latest)
	}
}

func TestIntegration_SilenceCheck(t *testing.T) {
	// This verifies that IsUpdateAvailable respects IgnoredVersion.
	tempDir, _ := os.MkdirTemp("", "silence-test-*")
	defer func() { _ = os.RemoveAll(tempDir) }()
	
	mgr := &StateManager{path: filepath.Join(tempDir, "state.json")}
	svc := NewService(mgr, &MockProvider{}, "v1.0.0")
	
	_ = mgr.Save(&State{LatestVersion: "v2.0.0", IgnoredVersion: "v2.0.0"})
	available, _ := svc.IsUpdateAvailable()
	if available {
		t.Error("expected update to be silent when version is ignored")
	}
}
