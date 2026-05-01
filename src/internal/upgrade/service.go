package upgrade

import (
	"context"
	"os"
	"strings"
	"time"
)

// Service coordinates the auto-update process.
type Service struct {
	stateManager *StateManager
	provider     Provider
	currentVer   string
	
	// For testing
	now func() time.Time
}

// NewService creates a new upgrade service.
func NewService(mgr *StateManager, p Provider, currentVer string) *Service {
	return &Service{
		stateManager: mgr,
		provider:     p,
		currentVer:   currentVer,
		now:          time.Now,
	}
}

// CheckForUpdate initiates an async background check if 24h have passed.
func (s *Service) CheckForUpdate(ctx context.Context) {
	state, err := s.stateManager.Load()
	if err != nil {
		return // Silent failure for background check
	}

	// Throttle: only check every 24 hours
	if s.now().Sub(state.LastCheckAt) < 24*time.Hour {
		return
	}

	// Run in background
	go func() {
		// Create a detached context with a timeout for the background task
		bgCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()

		latest, err := s.provider.GetLatestVersion(bgCtx)
		if err != nil {
			return // Silent failure
		}

		// Update state
		state.LatestVersion = latest
		state.LastCheckAt = s.now()
		_ = s.stateManager.Save(state)
	}()
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
