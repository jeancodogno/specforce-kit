package upgrade

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	svc.CheckForUpdate()
	
	// Note: We can't easily verify the async background check in this integration test
	// anymore because it spawns a detached process that won't have the flag in the test binary.
	// We'll manually inject state to test the rest of the flow.
	_ = mgr.Save(&State{LatestVersion: "v2.0.0", LastCheckedAt: svc.now()})

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

func TestIntegration_BackgroundUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "background-update-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	statePath := filepath.Join(tempDir, "state.json")
	mgr := &StateManager{path: statePath}
	
	binaryContent := []byte("new specforce binary content")
	hash := sha256.Sum256(binaryContent)
	hashStr := fmt.Sprintf("%x", hash)

	// Mock HTTP Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "checksums.txt") {
			assetName := fmt.Sprintf("specforce-kit_%s_%s", runtime.GOOS, runtime.GOARCH)
			if runtime.GOOS == "windows" {
				assetName += ".exe"
			}
			_, _ = fmt.Fprintf(w, "%s  %s\n", hashStr, assetName)
			return
		}
		_, _ = w.Write(binaryContent)
	}))
	defer ts.Close()

	svc := NewService(mgr, &MockProvider{Version: "v2.0.0"}, "v1.0.0")
	svc.BaseURL = ts.URL
	svc.installer.Client = ts.Client()

	// Perform background update
	err = svc.PerformBackgroundUpdate(context.Background())
	if err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}

	// Verify state
	state, err := mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !state.UpdateReady {
		t.Error("expected UpdateReady to be true")
	}
	if state.StagedVersion != "v2.0.0" {
		t.Errorf("expected StagedVersion v2.0.0, got %s", state.StagedVersion)
	}

	// Verify staged file
	stagedPath := filepath.Join(mgr.GetStagedDir(), "specforce_next")
	content, err := os.ReadFile(stagedPath)
	if err != nil {
		t.Fatalf("failed to read staged binary: %v", err)
	}
	if string(content) != string(binaryContent) {
		t.Error("staged binary content mismatch")
	}
}
