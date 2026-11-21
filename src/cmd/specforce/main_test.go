package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
)

func TestMain(t *testing.T) {
	// Dummy test to avoid 'no such tool covdata' error in Go 1.25+
}

// captureStdout runs fn while os.Stdout is redirected to a pipe and returns
// everything written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	fn()
	_ = w.Close()
	return <-done
}

// stageScript stages a real shell script reporting version as a fingerprinted
// update ready for exePath, so the pre-swap test run passes. It returns the
// staged bytes.
func stageScript(t *testing.T, mgr *upgrade.StateManager, exePath, version string) []byte {
	t.Helper()
	if err := mgr.EnsureStagedDir(); err != nil {
		t.Fatal(err)
	}
	staged := []byte("#!/bin/sh\necho 'specforce version " + version + "'\n")
	if err := os.WriteFile(filepath.Join(mgr.GetStagedDir(), "specforce_next"), staged, 0600); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Save(&upgrade.State{
		UpdateReady:   true,
		StagedVersion: version,
		LatestVersion: version,
		TargetPath:    exePath,
		StagedSHA256:  fmt.Sprintf("%x", sha256.Sum256(staged)),
		StagedSize:    int64(len(staged)),
	}); err != nil {
		t.Fatal(err)
	}
	return staged
}

func TestSwapCheck_WritesNothingToStdout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the staged binary is a shell script")
	}
	t.Setenv("HOME", t.TempDir())
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		t.Fatal(err)
	}

	installDir := filepath.Join(t.TempDir(), "install")
	if err := os.MkdirAll(installDir, 0750); err != nil {
		t.Fatal(err)
	}
	exePath := filepath.Join(installDir, "specforce")
	if err := os.WriteFile(exePath, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}
	staged := stageScript(t, mgr, exePath, "v1.0.0-beta.3")

	var execed string
	out := captureStdout(t, func() {
		checkForSwap(swapEnv{
			version:    "1.0.0-beta.2",
			executable: func() (string, error) { return exePath, nil },
			exec:       func(p string) error { execed = p; return nil },
		})
	})

	if out != "" {
		t.Errorf("pre-command swap must write nothing to stdout, got %q", out)
	}
	b, err := os.ReadFile(exePath)
	if err != nil || string(b) != string(staged) {
		t.Fatalf("expected the executable to be replaced, got %q (%v)", string(b), err)
	}
	if execed != exePath {
		t.Errorf("expected re-exec of %q, got %q", exePath, execed)
	}
	s, err := mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.AppliedVersion != "v1.0.0-beta.3" {
		t.Errorf("expected applied version recorded in state, got %q", s.AppliedVersion)
	}
}

// newSwapCheckFixture isolates HOME and stages a healthy update for a fresh
// executable holding "old content". It returns the state manager and the path.
func newSwapCheckFixture(t *testing.T, version string) (*upgrade.StateManager, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the staged binary is a shell script")
	}
	t.Setenv("HOME", t.TempDir())
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	installDir := filepath.Join(t.TempDir(), "install")
	if err := os.MkdirAll(installDir, 0750); err != nil {
		t.Fatal(err)
	}
	exePath := filepath.Join(installDir, "specforce")
	if err := os.WriteFile(exePath, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}
	stageScript(t, mgr, exePath, version)
	return mgr, exePath
}

func TestSwapCheck_RollbackOnFailedRestart(t *testing.T) {
	mgr, exePath := newSwapCheckFixture(t, "v1.0.0-beta.3")

	execCalls := 0
	checkForSwap(swapEnv{
		version:    "1.0.0-beta.2",
		executable: func() (string, error) { return exePath, nil },
		exec:       func(string) error { execCalls++; return errors.New("exec format error") },
	})

	if execCalls != 1 {
		t.Fatalf("expected one restart attempt, got %d", execCalls)
	}
	b, err := os.ReadFile(exePath)
	if err != nil || string(b) != "old content" {
		t.Fatalf("expected the previous executable restored, got %q (%v)", string(b), err)
	}
	if _, err := os.Stat(exePath + ".old"); !os.IsNotExist(err) {
		t.Errorf("expected no backup (.old) to remain, stat err = %v", err)
	}
	s, err := mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.AppliedVersion != "" {
		t.Errorf("a failed restart must not record the version as applied, got %q", s.AppliedVersion)
	}
	if s.ManualUpdateVersion != "v1.0.0-beta.3" {
		t.Errorf("expected v1.0.0-beta.3 recorded as manual update, got %q", s.ManualUpdateVersion)
	}
	if s.UpdateReady {
		t.Error("expected the update to no longer be ready")
	}
}

func TestSwapCheck_NoRollbackOnSuccessfulRestart(t *testing.T) {
	mgr, exePath := newSwapCheckFixture(t, "v1.0.0-beta.3")

	checkForSwap(swapEnv{
		version:    "1.0.0-beta.2",
		executable: func() (string, error) { return exePath, nil },
		exec:       func(string) error { return nil },
	})

	b, err := os.ReadFile(exePath)
	if err != nil || string(b) == "old content" {
		t.Fatalf("expected the new executable to stay in place, got %q (%v)", string(b), err)
	}
	s, err := mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.AppliedVersion != "v1.0.0-beta.3" {
		t.Errorf("expected applied version recorded, got %q", s.AppliedVersion)
	}
	if s.ManualUpdateVersion != "" {
		t.Errorf("a successful restart must not record a manual update, got %q", s.ManualUpdateVersion)
	}
}

// captureStderr runs fn while os.Stderr is redirected to a pipe and returns
// everything written to it.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	fn()
	_ = w.Close()
	return <-done
}

func TestSwapCheck_RollbackDebugLineIncludesRestartError(t *testing.T) {
	mgr, exePath := newSwapCheckFixture(t, "v1.0.0-beta.3")
	t.Setenv("DEBUG", "1")

	out := captureStderr(t, func() {
		checkForSwap(swapEnv{
			version:    "1.0.0-beta.2",
			executable: func() (string, error) { return exePath, nil },
			exec:       func(string) error { return os.ErrNotExist },
		})
	})

	if !strings.Contains(out, "restored the previous binary") || !strings.Contains(out, os.ErrNotExist.Error()) {
		t.Errorf("expected the rollback line with the restart error, got %q", out)
	}
	b, err := os.ReadFile(exePath)
	if err != nil || string(b) != "old content" {
		t.Fatalf("expected the previous executable restored, got %q (%v)", string(b), err)
	}
	if s, err := mgr.Load(); err != nil || s.ManualUpdateVersion != "v1.0.0-beta.3" {
		t.Errorf("expected v1.0.0-beta.3 recorded as manual update, got %+v (%v)", s, err)
	}
}
