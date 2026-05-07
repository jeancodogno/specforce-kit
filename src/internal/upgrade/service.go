package upgrade

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Service coordinates the auto-update process.
type Service struct {
	stateManager *StateManager
	provider     Provider
	installer    *BinaryInstaller
	currentVer   string
	BaseURL      string
	
	// For testing
	now func() time.Time
}

// NewService creates a new upgrade service.
func NewService(mgr *StateManager, p Provider, currentVer string) *Service {
	return &Service{
		stateManager: mgr,
		provider:     p,
		installer:    NewBinaryInstaller(),
		currentVer:   currentVer,
		BaseURL:      "https://github.com",
		now:          time.Now,
	}
}

// CheckForUpdate initiates a detached background check if 6h have passed.
func (s *Service) CheckForUpdate() {
	state, err := s.stateManager.Load()
	if err != nil {
		return // Silent failure
	}

	// Throttle: only check every 6 hours
	if s.now().Sub(state.LastCheckedAt) < 6*time.Hour {
		return
	}

	_ = s.LaunchBackgroundCheck()
}

// LaunchBackgroundCheck spawns a detached specforce process to perform the update check.
func (s *Service) LaunchBackgroundCheck() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare detached process attributes
	attr := &os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			nil, // stdin
			nil, // stdout
			nil, // stderr
		},
		Sys: getDetachedSysProcAttr(),
	}

	// Spawn specforce with the internal check flag
	// #nosec G204,G702 - exe is obtained via os.Executable() which is safe
	process, err := os.StartProcess(exe, []string{exe, "--internal-upgrade-check"}, attr)
	if err != nil {
		return err
	}

	// Release the process immediately so it becomes a daemon/detached
	return process.Release()
}

// PerformBackgroundUpdate performs the check, download and staging silently.
func (s *Service) PerformBackgroundUpdate(ctx context.Context) error {
	latest, err := s.provider.GetLatestVersion(ctx)
	if err != nil {
		return err
	}

	state, err := s.stateManager.Load()
	if err != nil {
		return err
	}

	// Update latest version and last check time
	state.LatestVersion = latest
	state.LastCheckedAt = s.now()

	if IsNewer(s.currentVer, latest) && latest != state.IgnoredVersion {
		// If newer, download and stage
		if err := s.stageUpdate(ctx, latest); err != nil {
			// Log error but we must still save the LastCheckedAt to avoid infinite loops
			_ = s.stateManager.Save(state)
			return err
		}
		state.StagedVersion = latest
		state.UpdateReady = true
	}

	return s.stateManager.Save(state)
}

func (s *Service) stageUpdate(ctx context.Context, version string) error {
	if s.isNPM() {
		// For NPM, staging IS the installation, but we'll do it in background.
		// Actually, REQ-4 says background process triggers npm install -g.
		installer := NewNPMInstaller()
		return installer.Install(ctx)
	}

	if err := s.stateManager.EnsureStagedDir(); err != nil {
		return err
	}

	tmpPath, err := s.installer.DownloadAndVerify(ctx, version, s.BaseURL)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmpPath) }()

	stagedPath := filepath.Join(s.stateManager.GetStagedDir(), "specforce_next")
	return s.moveFile(tmpPath, stagedPath)
}

// PerformAtomicSwap replaces the active binary with the staged one.
func (s *Service) PerformAtomicSwap() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	return s.PerformAtomicSwapAt(exePath)
}

// PerformAtomicSwapAt replaces the binary at targetPath with the staged one.
func (s *Service) PerformAtomicSwapAt(targetPath string) error {
	// Cleanup any old backup from previous runs
	_ = os.Remove(targetPath + ".old")

	if s.isNPM() {
		// For NPM, the upgrade happened in background via npm install -g
		// We just need to reset the state.
		state, _ := s.stateManager.Load()
		if state != nil {
			state.UpdateReady = false
			_ = s.stateManager.Save(state)
		}
		return nil
	}

	state, err := s.stateManager.Load()
	if err != nil {
		return err
	}

	if !state.UpdateReady || state.StagedVersion == "" {
		return nil
	}

	stagedPath := filepath.Join(s.stateManager.GetStagedDir(), "specforce_next")
	if _, err := os.Stat(stagedPath); err != nil {
		return fmt.Errorf("staged binary not found: %w", err)
	}

	// 1. Rename current to .old
	oldPath := targetPath + ".old"
	if err := os.Rename(targetPath, oldPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// 2. Move staged to target
	if err := s.moveFile(stagedPath, targetPath); err != nil {
		// Rollback
		_ = os.Rename(oldPath, targetPath)
		return fmt.Errorf("failed to move staged binary: %w", err)
	}

	// 3. Set permissions
	// #nosec G302 - Binary must be executable
	if err := os.Chmod(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// 4. Update state
	state.UpdateReady = false
	return s.stateManager.Save(state)
}

func (s *Service) moveFile(src, dst string) error {
	// Try renaming first (same filesystem)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fallback to copy + delete (different filesystems)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	// #nosec G304 - src and dst are internal paths managed by upgrade service
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	_ = os.Remove(src)
	return nil
}

// ExecuteNewBinary replaces the current process with the new binary.
func (s *Service) ExecuteNewBinary() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	return s.ExecuteBinary(exePath)
}

// ExecuteBinary replaces the current process with the binary at the given path.
func (s *Service) ExecuteBinary(path string) error {
	return executeNewBinary(path)
}

// IsUpdateAvailable returns true if a newer version is available and not ignored.
func (s *Service) IsUpdateAvailable() (bool, string) {
	state, err := s.stateManager.Load()
	if err != nil {
		return false, ""
	}

	if state.LatestVersion == "" {
		return false, ""
	}

	if state.LatestVersion == state.IgnoredVersion {
		return false, ""
	}

	if IsNewer(s.currentVer, state.LatestVersion) {
		return true, state.LatestVersion
	}

	return false, ""
}

// PerformUpgrade executes the update strategy (NPM vs Binary).
func (s *Service) PerformUpgrade(ctx context.Context, version string) error {
	if s.isNPM() {
		installer := NewNPMInstaller()
		return installer.Install(ctx)
	}

	installer := NewBinaryInstaller()
	// GitHub API base URL is hardcoded in provider but we need it here for download.
	// For production, it's https://github.com
	baseURL := "https://github.com"

	tmpPath, err := installer.DownloadAndVerify(ctx, version, baseURL)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmpPath) }()

	return installer.Replace(tmpPath)
}

func (s *Service) isNPM() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	// Simple detection: if the executable name is 'node' or path contains 'node_modules'
	// or if we're running via the npm-distributed index.js (which usually results in a symlink)
	return strings.Contains(exePath, "node") || strings.Contains(exePath, "npm")
}
