package upgrade

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ErrSwapSkipped is returned (wrapped) by PerformAtomicSwapAt when no update
// was applied. PerformAtomicSwapAt returns nil only when the binary was replaced.
var ErrSwapSkipped = errors.New("update not applied")

// stagedFileName is the name of the staged binary inside the staging directory.
const stagedFileName = "specforce_next"

// Service coordinates the auto-update process.
type Service struct {
	stateManager *StateManager
	provider     Provider
	installer    *BinaryInstaller
	npmInstaller *NPMInstaller
	currentVer   string
	BaseURL      string

	// For testing
	now        func() time.Time
	executable func() (string, error)
	launch     func() error
	// runStaged executes the staged binary with --version, HOME set to home,
	// bounded by timeout, and returns its standard output.
	runStaged        func(path, home string, timeout time.Duration) (string, error)
	stagedRunTimeout time.Duration
}

// NewService creates a new upgrade service.
func NewService(mgr *StateManager, p Provider, currentVer string) *Service {
	svc := &Service{
		stateManager: mgr,
		provider:     p,
		installer:    NewBinaryInstaller(),
		npmInstaller: NewNPMInstaller(),
		currentVer:   currentVer,
		BaseURL:      "https://github.com",
		now:          time.Now,
		executable:   resolveExecutable,

		runStaged:        runStagedBinary,
		stagedRunTimeout: defaultStagedRunTimeout,
	}
	svc.launch = svc.LaunchBackgroundCheck
	return svc
}

// resolveExecutable returns the running executable path with symlinks resolved.
func resolveExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return resolvePath(exe), nil
}

// resolvePath returns an absolute, symlink-resolved form of p. When p cannot be
// resolved (e.g. it no longer exists), the cleaned absolute path is returned.
func resolvePath(p string) string {
	if resolved, err := filepath.EvalSymlinks(p); err == nil {
		p = resolved
	}
	if abs, err := filepath.Abs(p); err == nil {
		p = abs
	}
	return filepath.Clean(p)
}

// debugf writes an updater decision to stderr when DEBUG=1.
func debugf(format string, args ...interface{}) {
	if os.Getenv("DEBUG") == "1" {
		_, _ = fmt.Fprintf(os.Stderr, "[DEBUG] upgrade: "+format+"\n", args...)
	}
}

// isDirWritable reports whether a file can be created in dir. It creates and
// removes a temporary file instead of inspecting permission bits, which are
// unreliable with ownership and ACLs.
func isDirWritable(dir string) bool {
	f, err := os.CreateTemp(dir, ".specforce-write-check-*")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}

// CheckForUpdate initiates a detached background check if 6h have passed.
func (s *Service) CheckForUpdate() {
	state, err := s.stateManager.Load()
	if err != nil {
		return // Silent failure
	}

	// Throttle: only check every 6 hours
	if s.now().Sub(state.LastCheckedAt) < 6*time.Hour {
		debugf("check skipped: throttled (last check %s)", state.LastCheckedAt.Format(time.RFC3339))
		return
	}

	_ = s.launch()
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
	state, err := s.stateManager.Load()
	if err != nil {
		return err
	}

	latest, err := s.provider.GetLatestVersion(ctx)
	if errors.Is(err, ErrNoChannelTag) {
		// No release for this channel: nothing to do, but the throttle window starts.
		debugf("check skipped: %v", err)
		state.LastCheckedAt = s.now()
		return s.stateManager.Save(state)
	}
	if err != nil {
		// Persist the check time anyway so the throttle window starts even offline.
		debugf("check failed: %v", err)
		state.LastCheckedAt = s.now()
		_ = s.stateManager.Save(state)
		return err
	}

	// Update latest version and last check time
	state.LatestVersion = latest
	state.LastCheckedAt = s.now()

	if IsNewer(s.currentVer, latest) {
		if s.isNPM() {
			return s.installNPM(ctx, state, latest)
		}

		exe, err := s.executable()
		if err != nil {
			_ = s.stateManager.Save(state)
			return fmt.Errorf("failed to resolve executable path: %w", err)
		}
		targetPath := resolvePath(exe)
		if latest == state.ManualUpdateVersion {
			// A swap failure (or non-writable install dir) is terminal for
			// that version: never download and stage it again.
			debugf("staging skipped: %s already requires a manual update", latest)
			return s.stateManager.Save(state)
		}
		if !isDirWritable(filepath.Dir(targetPath)) {
			debugf("staging skipped: install directory %s is not writable; %s requires a manual update", filepath.Dir(targetPath), latest)
			state.ManualUpdateVersion = latest
			return s.stateManager.Save(state)
		}

		// If newer, download and stage
		if err := s.stageAndRecord(ctx, state, latest, targetPath); err != nil {
			// Log error but we must still save the LastCheckedAt to avoid infinite loops
			_ = s.stateManager.Save(state)
			return err
		}
	}

	return s.stateManager.Save(state)
}

// stageAndRecord stages version and records it in state (the caller saves it)
// as ready for targetPath with its fingerprint. A verification failure of the
// staged file counts toward the retry policy.
func (s *Service) stageAndRecord(ctx context.Context, state *State, version, targetPath string) error {
	fp, err := s.stageUpdate(ctx, version)
	if err != nil {
		debugf("staging %s failed: %v", version, err)
		if errors.Is(err, errStagedVerification) {
			s.recordVerifyFailure(state, version)
		}
		return err
	}
	state.StagedVersion = version
	state.UpdateReady = true
	state.TargetPath = targetPath
	state.StagedSHA256 = fp.SHA256
	state.StagedSize = fp.Size
	if state.VerifyFailVersion != version {
		// The failure counter belongs to one version only.
		state.VerifyFailVersion = ""
		state.VerifyFailCount = 0
	}
	return nil
}

// installNPM installs exactly version through npm in the background. There is
// no staged binary nor swap: on success the version is recorded as applied so
// the next interactive run announces it; on failure it is recorded as requiring
// a manual update and never retried.
func (s *Service) installNPM(ctx context.Context, state *State, version string) error {
	if version == state.ManualUpdateVersion {
		debugf("npm install skipped: %s already failed and requires a manual update", version)
		return s.stateManager.Save(state)
	}
	if err := s.npmInstaller.Install(ctx, version); err != nil {
		debugf("npm install of %s failed: %v; it requires a manual update", version, err)
		state.ManualUpdateVersion = version
		_ = s.stateManager.Save(state)
		return fmt.Errorf("failed to install %s with npm: %w", version, err)
	}
	state.AppliedVersion = version
	return s.stateManager.Save(state)
}

// stageUpdate downloads and verifies version, moves it to the staging
// directory and checks the staged file against the verified fingerprint, which
// it returns. A staged file that no longer matches is deleted immediately.
func (s *Service) stageUpdate(ctx context.Context, version string) (Fingerprint, error) {
	if err := s.stateManager.EnsureStagedDir(); err != nil {
		return Fingerprint{}, err
	}

	tmpPath, fp, err := s.installer.DownloadAndVerify(ctx, version, s.BaseURL)
	if err != nil {
		return Fingerprint{}, err
	}
	defer func() { _ = os.Remove(tmpPath) }()

	// The staged file stays owner-only (0600) until the swap makes it executable.
	stagedPath := filepath.Join(s.stateManager.GetStagedDir(), stagedFileName)
	if err := s.moveFile(tmpPath, stagedPath); err != nil {
		_ = os.Remove(stagedPath)
		return Fingerprint{}, err
	}
	if err := os.Chmod(stagedPath, 0600); err != nil {
		_ = os.Remove(stagedPath)
		return Fingerprint{}, err
	}
	if err := verifyFingerprint(stagedPath, fp); err != nil {
		_ = os.Remove(stagedPath)
		return Fingerprint{}, err
	}
	return fp, nil
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
// It returns nil only when the binary was replaced; every other outcome returns
// an error (wrapping ErrSwapSkipped when the update was deliberately not applied).
// The swap is applied only when the update is ready, was staged for this exact
// path and the staged version is strictly newer than the running version.
// Any failure while swapping is terminal for that version: the binary is rolled
// back, the staged file deleted and the version recorded as a manual update.
func (s *Service) PerformAtomicSwapAt(targetPath string) error {
	// Cleanup any old backup from previous runs
	_ = os.Remove(targetPath + ".old")

	if s.isNPM() {
		// For NPM, the upgrade happens in background via npm install -g: there
		// is never anything to swap. Clear a ready flag left by older versions.
		if state, err := s.stateManager.Load(); err == nil && state.UpdateReady {
			state.UpdateReady = false
			_ = s.stateManager.Save(state)
		}
		return fmt.Errorf("%w: npm installations are updated by npm", ErrSwapSkipped)
	}

	state, err := s.stateManager.Load()
	if err != nil {
		return err
	}

	if !state.UpdateReady || state.StagedVersion == "" {
		return fmt.Errorf("%w: no update ready", ErrSwapSkipped)
	}

	stagedPath := filepath.Join(s.stateManager.GetStagedDir(), stagedFileName)

	if state.TargetPath == "" {
		// Legacy state without a bound target: on doubt, discard.
		debugf("swap skipped: staged %s has no target path; discarding", state.StagedVersion)
		return s.discardStaged(state, stagedPath, "staged update has no target path")
	}

	if resolvePath(state.TargetPath) != resolvePath(targetPath) {
		debugf("swap skipped: staged %s targets %s, not %s", state.StagedVersion, state.TargetPath, targetPath)
		return fmt.Errorf("%w: staged for another executable", ErrSwapSkipped)
	}

	if !IsNewer(s.currentVer, state.StagedVersion) {
		debugf("swap skipped: staged %s is not newer than running %s; discarding", state.StagedVersion, s.currentVer)
		return s.discardStaged(state, stagedPath, "staged version is not newer than running version")
	}

	if err := s.verifyStaged(state, stagedPath); err != nil {
		return err
	}

	if err := s.swapBinary(stagedPath, targetPath); err != nil {
		debugf("swap of %s failed: %v; %s requires a manual update", targetPath, err, state.StagedVersion)
		_ = os.Remove(stagedPath)
		state.UpdateReady = false
		state.ManualUpdateVersion = state.StagedVersion
		_ = s.stateManager.Save(state)
		return err
	}

	state.UpdateReady = false
	state.AppliedVersion = state.StagedVersion
	return s.stateManager.Save(state)
}

// verifyStaged checks the staged file against its recorded fingerprint and
// test-runs it. On any failure the update is discarded (the returned error
// wraps ErrSwapSkipped) and verification failures count toward the retry policy.
func (s *Service) verifyStaged(state *State, stagedPath string) error {
	if state.StagedSHA256 == "" {
		// Staged by an earlier build: nothing proves the file is the verified one.
		debugf("swap skipped: staged %s has no fingerprint; discarding", state.StagedVersion)
		return s.discardStaged(state, stagedPath, "staged update has no fingerprint")
	}

	if err := verifyFingerprint(stagedPath, Fingerprint{SHA256: state.StagedSHA256, Size: state.StagedSize}); err != nil {
		debugf("swap skipped: staged %s failed verification: %v; discarding", state.StagedVersion, err)
		s.recordVerifyFailure(state, state.StagedVersion)
		return s.discardStaged(state, stagedPath, err.Error())
	}

	if err := s.verifyStagedRun(stagedPath, state.StagedVersion); err != nil {
		debugf("swap skipped: staged %s failed its test run: %v; discarding", state.StagedVersion, err)
		s.recordVerifyFailure(state, state.StagedVersion)
		return s.discardStaged(state, stagedPath, err.Error())
	}

	// The installed bytes must be the ones that matched and passed the test
	// run: re-check right before the swap (the run could have changed them).
	if err := verifyFingerprint(stagedPath, Fingerprint{SHA256: state.StagedSHA256, Size: state.StagedSize}); err != nil {
		debugf("swap skipped: staged %s changed during its test run: %v; discarding", state.StagedVersion, err)
		s.recordVerifyFailure(state, state.StagedVersion)
		return s.discardStaged(state, stagedPath, err.Error())
	}
	return nil
}

// maxVerifyFailures is the number of verification failures of one version
// after which it is recorded as requiring a manual update.
const maxVerifyFailures = 2

// recordVerifyFailure counts a verification failure of version in state (the
// caller saves it). The first failure allows re-staging; the second records the
// version as a manual update so it is never downloaded again.
func (s *Service) recordVerifyFailure(state *State, version string) {
	if state.VerifyFailVersion != version {
		state.VerifyFailVersion = version
		state.VerifyFailCount = 0
	}
	state.VerifyFailCount++
	if state.VerifyFailCount >= maxVerifyFailures {
		debugf("%s failed verification %d times; it requires a manual update", version, state.VerifyFailCount)
		state.ManualUpdateVersion = version
	}
}

// discardStaged deletes the staged binary and clears the ready flag.
func (s *Service) discardStaged(state *State, stagedPath, reason string) error {
	_ = os.Remove(stagedPath)
	state.UpdateReady = false
	if err := s.stateManager.Save(state); err != nil {
		return err
	}
	return fmt.Errorf("%w: %s", ErrSwapSkipped, reason)
}

// swapBinary replaces targetPath with stagedPath, rolling back on failure.
func (s *Service) swapBinary(stagedPath, targetPath string) error {
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
		_ = os.Remove(targetPath)
		_ = os.Rename(oldPath, targetPath)
		return fmt.Errorf("failed to move staged binary: %w", err)
	}

	// 3. Set permissions
	// #nosec G302 - Binary must be executable
	if err := os.Chmod(targetPath, 0755); err != nil {
		// Rollback
		_ = os.Remove(targetPath)
		_ = os.Rename(oldPath, targetPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}

// RollbackSwapAt restores the backup (targetPath+".old") left by a successful
// swap over targetPath. It is used when the swapped executable fails to start:
// the version is no longer recorded as applied and is recorded as requiring a
// manual update instead. It returns an error when no backup exists.
func (s *Service) RollbackSwapAt(targetPath string) error {
	return s.rollbackSwap(targetPath, nil)
}

// RollbackFailedRestartAt is RollbackSwapAt for a restart that failed with
// restartErr; the error is included in the DEBUG rollback line as its reason.
func (s *Service) RollbackFailedRestartAt(targetPath string, restartErr error) error {
	return s.rollbackSwap(targetPath, restartErr)
}

// rollbackSwap implements RollbackSwapAt; reason, when non-nil, is why the
// restart failed.
func (s *Service) rollbackSwap(targetPath string, reason error) error {
	oldPath := targetPath + ".old"
	if _, err := os.Stat(oldPath); err != nil {
		return fmt.Errorf("no backup to restore: %w", err)
	}
	if err := os.Rename(oldPath, targetPath); err != nil {
		return fmt.Errorf("failed to restore previous binary: %w", err)
	}

	state, err := s.stateManager.Load()
	if err != nil {
		return err
	}
	version := state.AppliedVersion
	if version == "" {
		version = state.StagedVersion
	}
	cause := ""
	if reason != nil {
		cause = fmt.Sprintf(" (%v)", reason)
	}
	debugf("restart of %s failed%s; restored the previous binary, %s requires a manual update", targetPath, cause, version)
	state.AppliedVersion = ""
	state.UpdateReady = false
	if version != "" {
		state.ManualUpdateVersion = version
	}
	return s.stateManager.Save(state)
}

func (s *Service) moveFile(src, dst string) error {
	// Try renaming first (same filesystem)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fallback to copy + delete (different filesystems)
	// #nosec G304 - src is a temporary file managed by the upgrade service
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	// Created owner-only; swapBinary sets the executable bits explicitly.
	// #nosec G304 - dst is an internal path managed by the upgrade service
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

// IsUpdateAvailable returns true if a newer version is available.
func (s *Service) IsUpdateAvailable() (bool, string) {
	state, err := s.stateManager.Load()
	if err != nil {
		return false, ""
	}

	if state.LatestVersion == "" {
		return false, ""
	}

	if IsNewer(s.currentVer, state.LatestVersion) {
		return true, state.LatestVersion
	}

	return false, ""
}

// isNPM reports whether the running executable is the Specforce npm package.
func (s *Service) isNPM() bool {
	if s.executable == nil {
		return false
	}
	exePath, err := s.executable()
	if err != nil {
		return false
	}
	return IsNPMInstallPath(exePath)
}
