package upgrade

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
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

	// Setup state: a fingerprinted staged update whose test run passes.
	sha, size := fingerprintOf([]byte("new content"))
	svc.runStaged = passingRunner("v1.1.0")
	_ = mgr.Save(&State{UpdateReady: true, StagedVersion: "v1.1.0", TargetPath: activePath, StagedSHA256: sha, StagedSize: size})

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

// fingerprintOf returns the staged-file fingerprint (sha256 hex, size) of content.
func fingerprintOf(content []byte) (string, int64) {
	sum := sha256.Sum256(content)
	return fmt.Sprintf("%x", sum), int64(len(content))
}

// passingRunner is a test-run stub whose staged binary reports version.
func passingRunner(version string) func(path, home string, timeout time.Duration) (string, error) {
	return func(path, home string, timeout time.Duration) (string, error) {
		return "specforce version " + strings.TrimPrefix(version, "v") + "\n", nil
	}
}

// --- Swap decision regression tests (FIX-2, FIX-3, FIX-4) ---

type swapFixture struct {
	mgr        *StateManager
	svc        *Service
	installDir string
	activePath string
	stagedPath string
}

// newSwapFixture builds an isolated HOME, a temp install dir holding the
// "running" binary, and a staged binary for stagedVer bound to target.
// An empty target means "this fixture's active binary".
func newSwapFixture(t *testing.T, runningVer, stagedVer, target string) *swapFixture {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	mgr, err := NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(mgr, &MockProvider{}, runningVer)

	installDir := filepath.Join(t.TempDir(), "install")
	if err := os.MkdirAll(installDir, 0750); err != nil {
		t.Fatal(err)
	}
	activePath := filepath.Join(installDir, "specforce")
	if err := os.WriteFile(activePath, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := mgr.EnsureStagedDir(); err != nil {
		t.Fatal(err)
	}
	stagedPath := filepath.Join(mgr.GetStagedDir(), "specforce_next")
	if err := os.WriteFile(stagedPath, []byte("new content"), 0600); err != nil {
		t.Fatal(err)
	}

	if target == "" {
		target = activePath
	}
	sha, size := fingerprintOf([]byte("new content"))
	svc.runStaged = passingRunner(stagedVer)
	if err := mgr.Save(&State{
		UpdateReady:   true,
		StagedVersion: stagedVer,
		LatestVersion: stagedVer,
		TargetPath:    target,
		StagedSHA256:  sha,
		StagedSize:    size,
	}); err != nil {
		t.Fatal(err)
	}
	return &swapFixture{mgr: mgr, svc: svc, installDir: installDir, activePath: activePath, stagedPath: stagedPath}
}

func (f *swapFixture) activeContent(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(f.activePath)
	if err != nil {
		t.Fatalf("read active binary: %v", err)
	}
	return string(b)
}

func (f *swapFixture) state(t *testing.T) *State {
	t.Helper()
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	return s
}

func assertDiscarded(t *testing.T, f *swapFixture, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected a non-nil error signalling the update was not applied")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("executable must not be replaced, got %q", got)
	}
	if _, statErr := os.Stat(f.stagedPath); !os.IsNotExist(statErr) {
		t.Errorf("expected staged binary to be deleted, stat err = %v", statErr)
	}
	if f.state(t).UpdateReady {
		t.Error("expected update to no longer be ready")
	}
}

func TestSwap_OlderStagedVersionIsDiscarded(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.3", "v1.0.0-beta.2", "")
	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
}

func TestSwap_EqualStagedVersionIsDiscarded(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.3", "v1.0.0-beta.3", "")
	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
}

func TestSwap_NewerStagedVersionForThisPathIsApplied(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected swap to succeed, got %v", err)
	}
	if got := f.activeContent(t); got != "new content" {
		t.Errorf("expected executable to be replaced, got %q", got)
	}
	if f.state(t).UpdateReady {
		t.Error("expected UpdateReady to be false after a successful swap")
	}
	if f.state(t).ManualUpdateVersion != "" {
		t.Error("a successful swap must not record a manual update")
	}
}

func TestSwap_UpdateStagedForAnotherExecutableIsLeftIntact(t *testing.T) {
	for _, otherExists := range []bool{true, false} {
		name := "other path no longer exists"
		if otherExists {
			name = "other path exists"
		}
		t.Run(name, func(t *testing.T) {
			other := filepath.Join(t.TempDir(), "specforce")
			if otherExists {
				if err := os.WriteFile(other, []byte("other"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", other)

			err := f.svc.PerformAtomicSwapAt(f.activePath)
			if err == nil {
				t.Error("expected a non-nil error signalling the update was not applied")
			}
			if got := f.activeContent(t); got != "old content" {
				t.Errorf("this executable must be untouched, got %q", got)
			}
			if _, statErr := os.Stat(f.stagedPath); statErr != nil {
				t.Errorf("staged binary for the other executable must remain: %v", statErr)
			}
			s := f.state(t)
			if !s.UpdateReady || s.StagedVersion != "v1.0.0-beta.3" || s.TargetPath != other {
				t.Errorf("staged update for the other executable must stay intact, got %+v", s)
			}
			if s.ManualUpdateVersion != "" {
				t.Error("no manual update must be recorded for another executable")
			}
		})
	}
}

func TestSwap_FailureIsTerminalForThatVersion(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")

	// #nosec G302 - test makes the install dir read-only on purpose
	if err := os.Chmod(f.installDir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.installDir, 0750) })

	if err := f.svc.PerformAtomicSwapAt(f.activePath); err == nil {
		t.Fatal("expected swap to fail on a read-only install dir")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("running binary must be unchanged, got %q", got)
	}
	if _, statErr := os.Stat(f.stagedPath); !os.IsNotExist(statErr) {
		t.Errorf("expected staged binary to be deleted, stat err = %v", statErr)
	}
	s := f.state(t)
	if s.UpdateReady {
		t.Error("expected update to no longer be ready after a failed swap")
	}
	if s.ManualUpdateVersion != "v1.0.0-beta.3" {
		t.Errorf("expected manual update version v1.0.0-beta.3, got %q", s.ManualUpdateVersion)
	}

	// Second call: make the dir writable again and put a staged file back;
	// the swap must not be attempted because the update is no longer ready.
	// #nosec G302 - restoring permissions in test
	if err := os.Chmod(f.installDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.stagedPath, []byte("new content"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err == nil {
		t.Error("expected second call to report that nothing was applied")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("second call must not attempt the swap, got %q", got)
	}
	if _, statErr := os.Stat(f.activePath + ".old"); !os.IsNotExist(statErr) {
		t.Error("second call must not create a backup (.old)")
	}
}

// --- Background check robustness tests (FIX-6, FIX-1) ---

type bgFixture struct {
	mgr        *StateManager
	svc        *Service
	installDir string
	exePath    string
	requests   *int
	launches   *int
	now        time.Time
}

func newBGFixture(t *testing.T, provider Provider) *bgFixture {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	mgr, err := NewStateManager()
	if err != nil {
		t.Fatal(err)
	}

	binaryContent := []byte("new specforce binary content")
	sum := sha256.Sum256(binaryContent)
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if strings.HasSuffix(r.URL.Path, "_checksums.txt") {
			assetName := fmt.Sprintf("specforce-kit_%s_%s", runtime.GOOS, runtime.GOARCH)
			if runtime.GOOS == "windows" {
				assetName += ".exe"
			}
			_, _ = fmt.Fprintf(w, "%x  %s\n", sum, assetName)
			return
		}
		_, _ = w.Write(binaryContent)
	}))
	t.Cleanup(ts.Close)

	installDir := filepath.Join(t.TempDir(), "install")
	if err := os.MkdirAll(installDir, 0750); err != nil {
		t.Fatal(err)
	}
	exePath := filepath.Join(installDir, "specforce")
	if err := os.WriteFile(exePath, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 11, 17, 10, 0, 0, 0, time.UTC)
	launches := 0
	svc := NewService(mgr, provider, "v1.0.0")
	svc.BaseURL = ts.URL
	svc.installer.Client = ts.Client()
	svc.now = func() time.Time { return now }
	svc.executable = func() (string, error) { return exePath, nil }
	svc.launch = func() error { launches++; return nil }

	return &bgFixture{mgr: mgr, svc: svc, installDir: installDir, exePath: exePath, requests: &requests, launches: &launches, now: now}
}

func TestBackgroundUpdate_ProviderErrorPersistsCheckTimeAndThrottles(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Err: errors.New("offline")})

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err == nil {
		t.Fatal("expected provider error to be returned")
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !s.LastCheckedAt.Equal(f.now) {
		t.Errorf("expected check time %v to be persisted, got %v", f.now, s.LastCheckedAt)
	}

	// One hour later: still inside the 6h window, no new check must be spawned.
	later := f.now.Add(1 * time.Hour)
	f.svc.now = func() time.Time { return later }
	f.svc.CheckForUpdate()
	if *f.launches != 0 {
		t.Errorf("expected throttle to skip the background check, got %d launches", *f.launches)
	}

	// After the window, the check is launched again.
	afterWindow := f.now.Add(7 * time.Hour)
	f.svc.now = func() time.Time { return afterWindow }
	f.svc.CheckForUpdate()
	if *f.launches != 1 {
		t.Errorf("expected one launch after the throttle window, got %d", *f.launches)
	}
}

func TestBackgroundUpdate_NonWritableInstallDirRecordsManualUpdate(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	// #nosec G302 - test makes the install dir read-only on purpose
	if err := os.Chmod(f.installDir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.installDir, 0750) })

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("expected no error for a non-writable install dir, got %v", err)
	}
	if *f.requests != 0 {
		t.Errorf("expected nothing to be downloaded, got %d requests", *f.requests)
	}
	if _, err := os.Stat(filepath.Join(f.mgr.GetStagedDir(), "specforce_next")); !os.IsNotExist(err) {
		t.Errorf("expected nothing staged, stat err = %v", err)
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.UpdateReady {
		t.Error("expected update not to be marked ready")
	}
	if s.ManualUpdateVersion != "v2.0.0" {
		t.Errorf("expected manual update version v2.0.0, got %q", s.ManualUpdateVersion)
	}
	if s.LatestVersion != "v2.0.0" || !s.LastCheckedAt.Equal(f.now) {
		t.Errorf("expected latest version and check time persisted, got %+v", s)
	}
}

func TestBackgroundUpdate_WritableInstallDirStagesForResolvedTarget(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	if *f.requests == 0 {
		t.Error("expected the binary to be downloaded")
	}
	if _, err := os.Stat(filepath.Join(f.mgr.GetStagedDir(), "specforce_next")); err != nil {
		t.Errorf("expected staged binary: %v", err)
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !s.UpdateReady || s.StagedVersion != "v2.0.0" {
		t.Errorf("expected v2.0.0 staged and ready, got %+v", s)
	}
	if s.TargetPath != f.exePath {
		t.Errorf("expected target path %q, got %q", f.exePath, s.TargetPath)
	}
	if s.ManualUpdateVersion != "" {
		t.Errorf("expected no manual update recorded, got %q", s.ManualUpdateVersion)
	}
}

// --- Notice selection tests (FIX-5, FIX-1) ---

func newNoticeService(t *testing.T, runningVer string, state *State) (*StateManager, *Service) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	mgr, err := NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	if state != nil {
		if err := mgr.Save(state); err != nil {
			t.Fatal(err)
		}
	}
	return mgr, NewService(mgr, &MockProvider{}, runningVer)
}

func TestNotice_AppliedAfterSwapIsPendingExactlyOnce(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected swap to succeed, got %v", err)
	}

	// The re-executed process runs the new version.
	newSvc := NewService(f.mgr, &MockProvider{}, "1.0.0-beta.3")
	n := newSvc.PendingNotice()
	if n.Kind != NoticeApplied || n.Version != "v1.0.0-beta.3" {
		t.Fatalf("expected applied notice for v1.0.0-beta.3, got %+v", n)
	}
	if err := newSvc.MarkNoticeShown(n); err != nil {
		t.Fatalf("MarkNoticeShown: %v", err)
	}
	if n := NewService(f.mgr, &MockProvider{}, "1.0.0-beta.3").PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no pending notice after the applied notice was shown, got %+v", n)
	}
}

func TestNotice_AppliedIsNotClaimedByAnotherVersion(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected swap to succeed, got %v", err)
	}
	// A different build (another version) must not report that it was updated.
	if n := NewService(f.mgr, &MockProvider{}, "1.0.0-beta.1").PendingNotice(); n.Kind == NoticeApplied {
		t.Errorf("a different version must not get the applied notice, got %+v", n)
	}
}

func TestNotice_ManualShownOncePerVersion(t *testing.T) {
	mgr, svc := newNoticeService(t, "1.0.0-beta.2", &State{
		LatestVersion:       "v1.0.0-beta.3",
		ManualUpdateVersion: "v1.0.0-beta.3",
	})

	n := svc.PendingNotice()
	if n.Kind != NoticeManual || n.Version != "v1.0.0-beta.3" {
		t.Fatalf("expected manual notice for v1.0.0-beta.3, got %+v", n)
	}
	if err := svc.MarkNoticeShown(n); err != nil {
		t.Fatalf("MarkNoticeShown: %v", err)
	}

	// Persisted: a fresh service (next command) sees nothing for that version,
	// not even the "available" notice promising automatic application.
	if n := NewService(mgr, &MockProvider{}, "1.0.0-beta.2").PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no notice after the manual hint was shown, got %+v", n)
	}

	// A newer manual-update version is shown again, once.
	s, _ := mgr.Load()
	s.LatestVersion = "v1.0.0-beta.4"
	s.ManualUpdateVersion = "v1.0.0-beta.4"
	if err := mgr.Save(s); err != nil {
		t.Fatal(err)
	}
	n = NewService(mgr, &MockProvider{}, "1.0.0-beta.2").PendingNotice()
	if n.Kind != NoticeManual || n.Version != "v1.0.0-beta.4" {
		t.Errorf("expected manual notice for v1.0.0-beta.4, got %+v", n)
	}
}

func TestNotice_ManualVersionAlreadyRunningIsNotShown(t *testing.T) {
	_, svc := newNoticeService(t, "1.0.0-beta.3", &State{
		LatestVersion:       "v1.0.0-beta.3",
		ManualUpdateVersion: "v1.0.0-beta.3",
	})
	if n := svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no notice once the manual version is running, got %+v", n)
	}
}

func TestNotice_AvailableWhenStageable(t *testing.T) {
	_, svc := newNoticeService(t, "1.0.0-beta.2", &State{LatestVersion: "v1.0.0-beta.3"})
	n := svc.PendingNotice()
	if n.Kind != NoticeAvailable || n.Version != "v1.0.0-beta.3" {
		t.Errorf("expected available notice for v1.0.0-beta.3, got %+v", n)
	}
}

func TestNotice_NothingPendingOtherwise(t *testing.T) {
	_, svc := newNoticeService(t, "1.0.0-beta.3", nil)
	if n := svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no notice without state, got %+v", n)
	}
	_, svc = newNoticeService(t, "1.0.0-beta.3", &State{LatestVersion: "v1.0.0-beta.3"})
	if n := svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no notice when up to date, got %+v", n)
	}
}

func TestNotice_ManualCommandForBinaryInstall(t *testing.T) {
	_, svc := newNoticeService(t, "1.0.0-beta.2", nil)
	svc.executable = func() (string, error) { return "/usr/local/bin/specforce", nil }
	cmd := svc.ManualUpdateCommand("v1.0.0-beta.3")
	if !strings.Contains(cmd, "v1.0.0-beta.3") {
		t.Errorf("manual command must name the version, got %q", cmd)
	}
	if !strings.Contains(cmd, "make install") {
		t.Errorf("binary manual command must explain how to reinstall, got %q", cmd)
	}
	if strings.Contains(cmd, "npm") {
		t.Errorf("binary manual command must not mention npm, got %q", cmd)
	}
}

// --- npm channel-aware updates (FIX-7) ---

type npmFixture struct {
	mgr  *StateManager
	svc  *Service
	exec *MockCommandExecutor
	now  time.Time
}

func newNPMFixture(t *testing.T, provider Provider) *npmFixture {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	mgr, err := NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	pkgBin := filepath.Join(t.TempDir(), "lib", "node_modules", "@jeancodogno", "specforce-kit-linux-x64", "bin")
	if err := os.MkdirAll(pkgBin, 0750); err != nil {
		t.Fatal(err)
	}
	exePath := filepath.Join(pkgBin, "specforce")
	if err := os.WriteFile(exePath, []byte("old content"), 0600); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 17, 10, 0, 0, 0, time.UTC)
	mock := &MockCommandExecutor{}
	svc := NewService(mgr, provider, "1.0.0-beta.2")
	svc.now = func() time.Time { return now }
	svc.executable = func() (string, error) { return exePath, nil }
	svc.launch = func() error { return nil }
	svc.npmInstaller = &NPMInstaller{Executor: mock}
	return &npmFixture{mgr: mgr, svc: svc, exec: mock, now: now}
}

func TestBackgroundUpdate_NPMInstallsExactVersion(t *testing.T) {
	f := newNPMFixture(t, &MockProvider{Version: "1.0.0-beta.3"})
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	want := []string{"install", "-g", "@jeancodogno/specforce-kit@1.0.0-beta.3"}
	if f.exec.LastCmd != "npm" || strings.Join(f.exec.LastArgs, " ") != strings.Join(want, " ") {
		t.Errorf("expected npm %v, got %s %v", want, f.exec.LastCmd, f.exec.LastArgs)
	}
	s, _ := f.mgr.Load()
	if s.UpdateReady {
		t.Error("npm installs must never mark a swap as ready")
	}
	if s.AppliedVersion != "1.0.0-beta.3" {
		t.Errorf("expected applied version 1.0.0-beta.3, got %q", s.AppliedVersion)
	}
	if _, err := os.Stat(filepath.Join(f.mgr.GetStagedDir(), "specforce_next")); !os.IsNotExist(err) {
		t.Errorf("npm installs must not stage a binary, stat err = %v", err)
	}
}

func TestBackgroundUpdate_NPMMissingChannelTagIsNoUpdate(t *testing.T) {
	f := newNPMFixture(t, &MockProvider{Err: fmt.Errorf("lookup: %w", ErrNoChannelTag)})
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("a missing channel tag is not a failure, got %v", err)
	}
	if f.exec.LastCmd != "" {
		t.Errorf("nothing must be installed, got %s %v", f.exec.LastCmd, f.exec.LastArgs)
	}
	s, _ := f.mgr.Load()
	if !s.LastCheckedAt.Equal(f.now) {
		t.Errorf("expected check time persisted, got %v", s.LastCheckedAt)
	}
	if s.ManualUpdateVersion != "" || s.AppliedVersion != "" || s.UpdateReady {
		t.Errorf("nothing must be recorded, got %+v", s)
	}
	if n := f.svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("nothing must be shown, got %+v", n)
	}
}

func TestBackgroundUpdate_NPMInstallFailureRecordsManualAndIsNotRetried(t *testing.T) {
	f := newNPMFixture(t, &MockProvider{Version: "1.0.0-beta.3"})
	f.exec.Err = errors.New("EACCES")
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err == nil {
		t.Error("expected the npm install error to be reported")
	}
	s, _ := f.mgr.Load()
	if s.ManualUpdateVersion != "1.0.0-beta.3" {
		t.Errorf("expected manual update version 1.0.0-beta.3, got %q", s.ManualUpdateVersion)
	}
	if s.AppliedVersion != "" {
		t.Errorf("a failed install must not be reported as applied, got %q", s.AppliedVersion)
	}

	// Next check (after the throttle window): the same version is not retried.
	f.exec.LastCmd, f.exec.LastArgs, f.exec.Err = "", nil, nil
	later := f.now.Add(7 * time.Hour)
	f.svc.now = func() time.Time { return later }
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("second check: %v", err)
	}
	if f.exec.LastCmd != "" {
		t.Errorf("install must not be retried for 1.0.0-beta.3, got %s %v", f.exec.LastCmd, f.exec.LastArgs)
	}
}

func TestSwap_NPMInstallNeverSwapsOrExecs(t *testing.T) {
	f := newNPMFixture(t, &MockProvider{})
	// Legacy state written by the old npm flow.
	if err := f.mgr.Save(&State{UpdateReady: true, StagedVersion: "1.0.0-beta.3"}); err != nil {
		t.Fatal(err)
	}
	exe, _ := f.svc.executable()
	err := f.svc.PerformAtomicSwapAt(exe)
	if !errors.Is(err, ErrSwapSkipped) {
		t.Errorf("expected ErrSwapSkipped for npm installs, got %v", err)
	}
	if s, _ := f.mgr.Load(); s.UpdateReady {
		t.Error("expected the legacy ready flag to be cleared")
	}
}

func TestNotice_ManualCommandForNPMInstall(t *testing.T) {
	f := newNPMFixture(t, &MockProvider{})
	if got := f.svc.ManualUpdateCommand("v1.0.0-beta.3"); got != "npm install -g @jeancodogno/specforce-kit@1.0.0-beta.3" {
		t.Errorf("unexpected npm manual command %q", got)
	}
}

// --- Staged file hygiene tests (FIX-8, FIX-2) ---

func assertStagedPrivate(t *testing.T, f *bgFixture) {
	t.Helper()
	info, err := os.Stat(filepath.Join(f.mgr.GetStagedDir(), stagedFileName))
	if err != nil {
		t.Fatalf("expected staged binary: %v", err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Errorf("expected staged file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestBackgroundUpdate_StagedPermIsOwnerOnly(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	t.Setenv("TMPDIR", t.TempDir())

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	assertStagedPrivate(t, f)
}

// The download temp dir lives on a different filesystem than HOME, so staging
// goes through the copy fallback instead of a rename.
func TestBackgroundUpdate_StagedPermIsOwnerOnlyAcrossFilesystems(t *testing.T) {
	const otherFS = "/dev/shm"
	if info, err := os.Stat(otherFS); err != nil || !info.IsDir() {
		t.Skip("no second filesystem available for a cross-device staging test")
	}
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	tmp, err := os.MkdirTemp(otherFS, "specforce-staged-perm-*")
	if err != nil {
		t.Skipf("cannot create a temp dir on %s: %v", otherFS, err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	t.Setenv("TMPDIR", tmp)

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	assertStagedPrivate(t, f)
	assertDirEmpty(t, tmp)
}

func TestBackgroundUpdate_BinaryManualVersionIsNotRestaged(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	if err := f.mgr.Save(&State{ManualUpdateVersion: "v2.0.0"}); err != nil {
		t.Fatal(err)
	}

	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *f.requests != 0 {
		t.Errorf("expected a manual-update version not to be downloaded again, got %d requests", *f.requests)
	}
	if _, err := os.Stat(filepath.Join(f.mgr.GetStagedDir(), stagedFileName)); !os.IsNotExist(err) {
		t.Errorf("expected nothing staged, stat err = %v", err)
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.UpdateReady {
		t.Error("expected update not to be marked ready")
	}
	if s.ManualUpdateVersion != "v2.0.0" {
		t.Errorf("expected manual update version kept, got %q", s.ManualUpdateVersion)
	}
	if s.LatestVersion != "v2.0.0" || !s.LastCheckedAt.Equal(f.now) {
		t.Errorf("expected latest version and check time persisted, got %+v", s)
	}
}

// --- Staged binary integrity tests (staged-binary-integrity FIX-1) ---

// bgBinaryContent is the release binary served by newBGFixture.
const bgBinaryContent = "new specforce binary content"

// countingRunner is a passing test-run stub that counts its calls.
func countingRunner(version string, calls *int) func(path, home string, timeout time.Duration) (string, error) {
	pass := passingRunner(version)
	return func(path, home string, timeout time.Duration) (string, error) {
		*calls++
		return pass(path, home, timeout)
	}
}

// stageForBG runs a background check that stages v2.0.0 for f.exePath.
func stageForBG(t *testing.T, f *bgFixture) string {
	t.Helper()
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !s.UpdateReady || s.StagedVersion != "v2.0.0" {
		t.Fatalf("expected v2.0.0 staged and ready, got %+v", s)
	}
	return filepath.Join(f.mgr.GetStagedDir(), stagedFileName)
}

func assertBGNotApplied(t *testing.T, f *bgFixture, stagedPath string, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected a non-nil error signalling the update was not applied")
	}
	b, readErr := os.ReadFile(f.exePath)
	if readErr != nil || string(b) != "old content" {
		t.Errorf("executable must not be modified, got %q (%v)", string(b), readErr)
	}
	if _, statErr := os.Stat(stagedPath); !os.IsNotExist(statErr) {
		t.Errorf("expected staged binary to be deleted, stat err = %v", statErr)
	}
	s, loadErr := f.mgr.Load()
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if s.UpdateReady {
		t.Error("expected update to no longer be ready")
	}
	if s.AppliedVersion != "" {
		t.Errorf("a discarded update must not be recorded as applied, got %q", s.AppliedVersion)
	}
}

func TestIntegrity_StagingRecordsReleaseChecksumAndSize(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	stageForBG(t, f)

	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	wantSHA, wantSize := fingerprintOf([]byte(bgBinaryContent))
	if s.StagedSHA256 != wantSHA || s.StagedSize != wantSize {
		t.Errorf("expected fingerprint %s/%d, got %q/%d", wantSHA, wantSize, s.StagedSHA256, s.StagedSize)
	}
}

func TestIntegrity_CorruptedStagedFileIsNeverApplied(t *testing.T) {
	cases := map[string][]byte{
		"empty":                 {},
		"truncated":             []byte(bgBinaryContent[:len(bgBinaryContent)/2]),
		"overwritten same size": []byte(strings.Repeat("x", len(bgBinaryContent))),
		"overwritten larger":    []byte(bgBinaryContent + " plus tampering"),
	}
	for name, corrupt := range cases {
		t.Run(name, func(t *testing.T) {
			f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
			stagedPath := stageForBG(t, f)
			if err := os.WriteFile(stagedPath, corrupt, 0600); err != nil {
				t.Fatal(err)
			}
			calls := 0
			f.svc.runStaged = countingRunner("v2.0.0", &calls)

			err := f.svc.PerformAtomicSwapAt(f.exePath)
			assertBGNotApplied(t, f, stagedPath, err)
			if calls != 0 {
				t.Errorf("a file failing its fingerprint must never be run, got %d runs", calls)
			}
		})
	}
}

func TestIntegrity_ReadyUpdateWithoutFingerprintIsDiscarded(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	s := f.state(t)
	s.StagedSHA256, s.StagedSize = "", 0
	if err := f.mgr.Save(s); err != nil {
		t.Fatal(err)
	}

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
	if !errors.Is(err, ErrSwapSkipped) {
		t.Errorf("expected ErrSwapSkipped, got %v", err)
	}
}

func TestIntegrity_IntactStagedFileIsApplied(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	stageForBG(t, f)
	f.svc.runStaged = passingRunner("v2.0.0")

	if err := f.svc.PerformAtomicSwapAt(f.exePath); err != nil {
		t.Fatalf("expected the intact staged update to be applied, got %v", err)
	}
	b, err := os.ReadFile(f.exePath)
	if err != nil || string(b) != bgBinaryContent {
		t.Errorf("expected executable replaced by the staged binary, got %q (%v)", string(b), err)
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if s.UpdateReady || s.AppliedVersion != "v2.0.0" {
		t.Errorf("expected v2.0.0 applied and not ready, got %+v", s)
	}
}

// --- Pre-swap test run of the staged binary (staged-binary-integrity FIX-1) ---

func TestStagedRun_ReportedVersionIsApplied(t *testing.T) {
	for name, out := range map[string]string{
		"with v":    "specforce version v2.0.0\n",
		"without v": "specforce version 2.0.0\n",
	} {
		t.Run(name, func(t *testing.T) {
			f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
			stageForBG(t, f)
			f.svc.runStaged = func(path, home string, timeout time.Duration) (string, error) { return out, nil }

			if err := f.svc.PerformAtomicSwapAt(f.exePath); err != nil {
				t.Fatalf("expected the update to be applied, got %v", err)
			}
			if b, _ := os.ReadFile(f.exePath); string(b) != bgBinaryContent {
				t.Errorf("expected executable replaced, got %q", string(b))
			}
		})
	}
}

func TestStagedRun_FailedRunIsDiscarded(t *testing.T) {
	cases := map[string]func(path, home string, timeout time.Duration) (string, error){
		"runner error": func(string, string, time.Duration) (string, error) {
			return "", errors.New("exec format error")
		},
		"timeout": func(string, string, time.Duration) (string, error) {
			return "", context.DeadlineExceeded
		},
		"different version": func(string, string, time.Duration) (string, error) {
			return "specforce version 1.9.0\n", nil
		},
		"unexpected output": func(string, string, time.Duration) (string, error) {
			return "hello\n", nil
		},
	}
	for name, runner := range cases {
		t.Run(name, func(t *testing.T) {
			f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
			stagedPath := stageForBG(t, f)
			f.svc.runStaged = runner

			err := f.svc.PerformAtomicSwapAt(f.exePath)
			assertBGNotApplied(t, f, stagedPath, err)
			if !errors.Is(err, ErrSwapSkipped) {
				t.Errorf("expected ErrSwapSkipped, got %v", err)
			}
		})
	}
}

func TestStagedRun_RunnerGetsStagedPathEmptyHomeAndTimeout(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	stagedPath := stageForBG(t, f)

	var gotPath, gotHome string
	var gotTimeout time.Duration
	f.svc.runStaged = func(path, home string, timeout time.Duration) (string, error) {
		gotPath, gotHome, gotTimeout = path, home, timeout
		entries, err := os.ReadDir(home)
		if err != nil || len(entries) != 0 {
			t.Errorf("expected an existing empty HOME, got %v entries (%v)", len(entries), err)
		}
		if runtime.GOOS != "windows" {
			info, err := os.Stat(path)
			if err != nil || info.Mode().Perm() != 0700 {
				t.Errorf("expected staged file mode 0700 during the test run, got %v (%v)", info, err)
			}
		}
		return "specforce version 2.0.0\n", nil
	}

	if err := f.svc.PerformAtomicSwapAt(f.exePath); err != nil {
		t.Fatalf("expected the update to be applied, got %v", err)
	}
	if gotPath != stagedPath {
		t.Errorf("expected runner to get staged path %q, got %q", stagedPath, gotPath)
	}
	if gotHome == "" || gotHome == os.Getenv("HOME") {
		t.Errorf("expected a fresh temporary HOME, got %q", gotHome)
	}
	if _, err := os.Stat(gotHome); !os.IsNotExist(err) {
		t.Errorf("expected the temporary HOME to be removed, stat err = %v", err)
	}
	if gotTimeout < time.Second || gotTimeout > 10*time.Second {
		t.Errorf("expected a timeout of a few seconds, got %v", gotTimeout)
	}
}

// newScriptSwapFixture stages a real shell script as the staged binary and
// uses the default (real) test runner.
func newScriptSwapFixture(t *testing.T, script string) *swapFixture {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the staged binary is a shell script")
	}
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	content := []byte(script)
	if err := os.WriteFile(f.stagedPath, content, 0600); err != nil {
		t.Fatal(err)
	}
	s := f.state(t)
	s.StagedSHA256, s.StagedSize = fingerprintOf(content)
	if err := f.mgr.Save(s); err != nil {
		t.Fatal(err)
	}
	f.svc.runStaged = NewService(f.mgr, &MockProvider{}, "1.0.0-beta.2").runStaged
	return f
}

func TestStagedRun_RealScriptReportingStagedVersionIsApplied(t *testing.T) {
	script := "#!/bin/sh\necho \"specforce version 1.0.0-beta.3\"\n"
	f := newScriptSwapFixture(t, script)
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected the update to be applied, got %v", err)
	}
	if got := f.activeContent(t); got != script {
		t.Errorf("expected executable replaced by the script, got %q", got)
	}
}

func TestStagedRun_RealScriptReportingAnotherVersionIsDiscarded(t *testing.T) {
	f := newScriptSwapFixture(t, "#!/bin/sh\necho \"specforce version 1.0.0-beta.2\"\n")
	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
}

func TestStagedRun_RealScriptThatHangsTimesOut(t *testing.T) {
	f := newScriptSwapFixture(t, "#!/bin/sh\nsleep 30\necho \"specforce version 1.0.0-beta.3\"\n")
	f.svc.stagedRunTimeout = 200 * time.Millisecond

	start := time.Now()
	err := f.svc.PerformAtomicSwapAt(f.activePath)
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("expected the test run to be bounded by its timeout, took %v", elapsed)
	}
	assertDiscarded(t, f, err)
}

// --- Verification-failure retry policy (staged-binary-integrity FIX-2) ---

func failingRunner(string, string, time.Duration) (string, error) {
	return "", errors.New("exec format error")
}

func TestVerifyRetry_FirstFailureRestagesSameVersion(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	stagedPath := stageForBG(t, f)
	f.svc.runStaged = failingRunner

	err := f.svc.PerformAtomicSwapAt(f.exePath)
	assertBGNotApplied(t, f, stagedPath, err)
	s, _ := f.mgr.Load()
	if s.ManualUpdateVersion != "" {
		t.Errorf("a first verification failure must not record a manual update, got %q", s.ManualUpdateVersion)
	}

	*f.requests = 0
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("second background check failed: %v", err)
	}
	if *f.requests == 0 {
		t.Error("expected v2.0.0 to be downloaded again after a first verification failure")
	}
	s, _ = f.mgr.Load()
	if !s.UpdateReady || s.StagedVersion != "v2.0.0" || s.StagedSHA256 == "" {
		t.Errorf("expected v2.0.0 staged again with a fingerprint, got %+v", s)
	}
	if _, err := os.Stat(stagedPath); err != nil {
		t.Errorf("expected staged binary again: %v", err)
	}
}

func TestVerifyRetry_SecondFailureRecordsManualAndStopsDownloads(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	f.svc.runStaged = failingRunner

	for attempt := 1; attempt <= 2; attempt++ {
		stagedPath := stageForBG(t, f)
		err := f.svc.PerformAtomicSwapAt(f.exePath)
		assertBGNotApplied(t, f, stagedPath, err)
	}

	s, _ := f.mgr.Load()
	if s.ManualUpdateVersion != "v2.0.0" {
		t.Fatalf("expected manual update version v2.0.0 after the second failure, got %q", s.ManualUpdateVersion)
	}
	n := f.svc.PendingNotice()
	if n.Kind != NoticeManual || n.Version != "v2.0.0" {
		t.Fatalf("expected a pending manual notice for v2.0.0, got %+v", n)
	}
	if err := f.svc.MarkNoticeShown(n); err != nil {
		t.Fatal(err)
	}
	if n := f.svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected the manual notice only once, got %+v", n)
	}

	*f.requests = 0
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("later background check failed: %v", err)
	}
	if *f.requests != 0 {
		t.Errorf("expected no download of a manual-update version, got %d requests", *f.requests)
	}
}

func TestVerifyRetry_FailuresOfAnotherVersionDoNotCount(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	f.svc.runStaged = failingRunner

	stagedPath := stageForBG(t, f)
	assertBGNotApplied(t, f, stagedPath, f.svc.PerformAtomicSwapAt(f.exePath))

	// A newer release appears: its first failure is a first failure, not a second.
	f.svc.provider = &MockProvider{Version: "v2.1.0"}
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err != nil {
		t.Fatalf("PerformBackgroundUpdate failed: %v", err)
	}
	assertBGNotApplied(t, f, stagedPath, f.svc.PerformAtomicSwapAt(f.exePath))
	s, _ := f.mgr.Load()
	if s.ManualUpdateVersion != "" {
		t.Errorf("a first failure of v2.1.0 must not record a manual update, got %q", s.ManualUpdateVersion)
	}
}

func TestVerifyRetry_SwapFailureStaysTerminal(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	f := newBGFixture(t, &MockProvider{Version: "v2.0.0"})
	stageForBG(t, f)
	f.svc.runStaged = passingRunner("v2.0.0")

	// #nosec G302 - test makes the install dir read-only on purpose
	if err := os.Chmod(f.installDir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.installDir, 0750) })

	if err := f.svc.PerformAtomicSwapAt(f.exePath); err == nil {
		t.Fatal("expected the swap to fail on a read-only install dir")
	}
	s, _ := f.mgr.Load()
	if s.ManualUpdateVersion != "v2.0.0" {
		t.Errorf("a swap failure must be terminal on its first occurrence, got manual %q", s.ManualUpdateVersion)
	}
}

// --- Restart failure rollback tests (FIX-3) ---

func TestSwapRollback_RestoresBackupAndRecordsManualUpdate(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected swap to succeed, got %v", err)
	}

	if err := f.svc.RollbackSwapAt(f.activePath); err != nil {
		t.Fatalf("expected rollback to succeed, got %v", err)
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("expected the previous executable to be restored, got %q", got)
	}
	if _, err := os.Stat(f.activePath + ".old"); !os.IsNotExist(err) {
		t.Errorf("expected no backup (.old) to remain, stat err = %v", err)
	}
	s := f.state(t)
	if s.AppliedVersion != "" {
		t.Errorf("a rolled back version must not be recorded as applied, got %q", s.AppliedVersion)
	}
	if s.ManualUpdateVersion != "v1.0.0-beta.3" {
		t.Errorf("expected v1.0.0-beta.3 recorded as manual update, got %q", s.ManualUpdateVersion)
	}
	if n := f.svc.PendingNotice(); n.Kind != NoticeManual || n.Version != "v1.0.0-beta.3" {
		t.Errorf("expected a manual-update notice for v1.0.0-beta.3, got %+v", n)
	}
}

func TestSwapRollback_ErrorsWithoutBackup(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")

	if err := f.svc.RollbackSwapAt(f.activePath); err == nil {
		t.Fatal("expected an error when no backup exists")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("executable must be left untouched, got %q", got)
	}
	if s := f.state(t); s.ManualUpdateVersion != "" || !s.UpdateReady {
		t.Errorf("state must be left untouched without a backup, got %+v", s)
	}
}

// --- Unreadable state tests (FIX-4) ---

// writeRawState writes data verbatim at the fixture's state path.
func writeRawState(t *testing.T, mgr *StateManager, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(mgr.path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mgr.path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestStateRepair_OrphanStagedFileIsDeleted(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	writeRawState(t, f.mgr, "garbage{")

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	if err == nil {
		t.Error("expected the swap not to be applied with an unparseable state")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("executable must not be replaced, got %q", got)
	}
	if _, statErr := os.Stat(f.stagedPath); !os.IsNotExist(statErr) {
		t.Errorf("expected the orphan staged binary to be deleted, stat err = %v", statErr)
	}
}

func TestStateRepair_BackgroundCheckWritesValidState(t *testing.T) {
	f := newBGFixture(t, &MockProvider{Err: errors.New("offline")})
	writeRawState(t, f.mgr, `{"last_checked_at": `)

	f.svc.CheckForUpdate()
	if *f.launches != 1 {
		t.Errorf("expected the background check to be launched, got %d launches", *f.launches)
	}

	_ = f.svc.PerformBackgroundUpdate(context.Background())
	data, err := os.ReadFile(f.mgr.path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("expected a valid JSON state, got %q: %v", string(data), err)
	}
	if _, ok := raw["last_checked_at"]; !ok {
		t.Errorf("expected last_checked_at in the repaired state, got %s", string(data))
	}
	s, err := f.mgr.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !s.LastCheckedAt.Equal(f.now) {
		t.Errorf("expected check time %v, got %v", f.now, s.LastCheckedAt)
	}
}

func TestStateRepair_DirectoryAtStatePathDoesNothing(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := os.Remove(f.mgr.path); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(f.mgr.path, 0750); err != nil {
		t.Fatal(err)
	}
	launches := 0
	f.svc.launch = func() error { launches++; return nil }

	if _, err := f.mgr.Load(); err == nil {
		t.Error("expected a read error for a directory at the state path")
	}
	f.svc.CheckForUpdate()
	if launches != 0 {
		t.Errorf("expected no background check, got %d launches", launches)
	}
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err == nil {
		t.Error("expected the swap not to be applied")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("executable must not be replaced, got %q", got)
	}
	if _, err := os.Stat(f.stagedPath); err != nil {
		t.Errorf("staged file must be left alone when the state cannot be read: %v", err)
	}
	if err := f.svc.PerformBackgroundUpdate(context.Background()); err == nil {
		t.Error("expected the background update to report the read error")
	}
	if info, err := os.Stat(f.mgr.path); err != nil || !info.IsDir() {
		t.Errorf("the state path must not be overwritten, stat = %v, %v", info, err)
	}
	if n := f.svc.PendingNotice(); n.Kind != NoticeNone {
		t.Errorf("expected no notice, got %+v", n)
	}
}
