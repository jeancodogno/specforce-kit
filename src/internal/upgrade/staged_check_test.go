package upgrade

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

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

// --- Observability of state repair (correction cycle 1, item 1) ---

func TestStateRepair_DebugLogsParseErrorAndOrphanRemoval(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	writeRawState(t, f.mgr, "garbage{")
	t.Setenv("DEBUG", "1")

	out := captureStderr(t, func() {
		if _, err := f.mgr.Load(); err != nil {
			t.Errorf("expected the unparseable state to load as empty, got %v", err)
		}
	})
	if !strings.Contains(out, "state repair") {
		t.Errorf("expected a state repair line on stderr, got %q", out)
	}
	if !strings.Contains(out, "invalid character") {
		t.Errorf("expected the parse error as the reason, got %q", out)
	}
	if !strings.Contains(out, "removed the orphan staged binary") {
		t.Errorf("expected the orphan removal to be reported, got %q", out)
	}
}

func TestStateRepair_DebugLogsParseErrorWithoutOrphan(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := os.Remove(f.stagedPath); err != nil {
		t.Fatal(err)
	}
	writeRawState(t, f.mgr, `{"last_checked_at": `)
	t.Setenv("DEBUG", "1")

	out := captureStderr(t, func() { _, _ = f.mgr.Load() })
	if !strings.Contains(out, "state repair") || !strings.Contains(out, "unexpected end of JSON input") {
		t.Errorf("expected a state repair line with the parse error, got %q", out)
	}
	if !strings.Contains(out, "no orphan staged binary") {
		t.Errorf("expected the absence of an orphan to be reported, got %q", out)
	}
}

func TestStateRepair_NoDebugOutputWithoutDebug(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	writeRawState(t, f.mgr, "garbage{")
	t.Setenv("DEBUG", "")

	if out := captureStderr(t, func() { _, _ = f.mgr.Load() }); out != "" {
		t.Errorf("expected no stderr output without DEBUG=1, got %q", out)
	}
}

// --- Observability of the rollback reason (correction cycle 1, item 2) ---

func TestSwapRollback_DebugLineIncludesRestartError(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := f.svc.PerformAtomicSwapAt(f.activePath); err != nil {
		t.Fatalf("expected swap to succeed, got %v", err)
	}
	t.Setenv("DEBUG", "1")

	out := captureStderr(t, func() {
		if err := f.svc.RollbackFailedRestartAt(f.activePath, os.ErrNotExist); err != nil {
			t.Errorf("expected rollback to succeed, got %v", err)
		}
	})
	if !strings.Contains(out, "restored the previous binary") || !strings.Contains(out, os.ErrNotExist.Error()) {
		t.Errorf("expected the rollback line with the restart error, got %q", out)
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("expected the previous executable to be restored, got %q", got)
	}
	if s := f.state(t); s.ManualUpdateVersion != "v1.0.0-beta.3" || s.AppliedVersion != "" {
		t.Errorf("expected v1.0.0-beta.3 recorded as manual update and not applied, got %+v", s)
	}
}

// --- Installed bytes are the verified bytes (correction cycle 1, item 3) ---

// assertVerifyFailureCounted checks that one verification failure of version
// was recorded.
func assertVerifyFailureCounted(t *testing.T, f *swapFixture, version string) {
	t.Helper()
	s := f.state(t)
	if s.VerifyFailVersion != version || s.VerifyFailCount != 1 {
		t.Errorf("expected one verification failure of %s, got %q x%d", version, s.VerifyFailVersion, s.VerifyFailCount)
	}
	if s.AppliedVersion != "" {
		t.Errorf("a discarded update must not be recorded as applied, got %q", s.AppliedVersion)
	}
}

func TestIntegrity_StagedFileRewrittenDuringTestRunIsNotApplied(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	f.svc.runStaged = func(path, home string, timeout time.Duration) (string, error) {
		if err := os.WriteFile(path, []byte("rewritten by itself"), 0700); err != nil {
			t.Fatal(err)
		}
		return "specforce version 1.0.0-beta.3\n", nil
	}

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
	assertVerifyFailureCounted(t, f, "v1.0.0-beta.3")
}

func TestIntegrity_RealScriptRewritingItselfIsNotApplied(t *testing.T) {
	script := "#!/bin/sh\nprintf '#!/bin/sh\\necho evil\\n' > \"$0\"\necho \"specforce version 1.0.0-beta.3\"\n"
	f := newScriptSwapFixture(t, script)

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
	assertVerifyFailureCounted(t, f, "v1.0.0-beta.3")
}

// --- Symlinked staged file is rejected (correction cycle 1, item 4) ---

func TestIntegrity_SymlinkedStagedFileIsNotApplied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require privileges on Windows")
	}
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	linkTarget := filepath.Join(t.TempDir(), "elsewhere")
	if err := os.WriteFile(linkTarget, []byte("new content"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(f.stagedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(linkTarget, f.stagedPath); err != nil {
		t.Fatal(err)
	}
	calls := 0
	f.svc.runStaged = countingRunner("v1.0.0-beta.3", &calls)

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	assertDiscarded(t, f, err)
	assertVerifyFailureCounted(t, f, "v1.0.0-beta.3")
	if calls != 0 {
		t.Errorf("a symlinked staged file must never be run, got %d runs", calls)
	}
	info, statErr := os.Lstat(f.activePath)
	if statErr != nil || !info.Mode().IsRegular() {
		t.Errorf("the executable must stay a regular file, got %v (%v)", info, statErr)
	}
	info, statErr = os.Stat(linkTarget)
	if statErr != nil {
		t.Fatalf("the link target must be left in place: %v", statErr)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("the link target mode must be unchanged, got %o", perm)
	}
	if b, _ := os.ReadFile(linkTarget); string(b) != "new content" {
		t.Errorf("the link target content must be unchanged, got %q", string(b))
	}
}

func TestIntegrity_NonRegularStagedFileIsNotApplied(t *testing.T) {
	f := newSwapFixture(t, "1.0.0-beta.2", "v1.0.0-beta.3", "")
	if err := os.Remove(f.stagedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(f.stagedPath, 0700); err != nil {
		t.Fatal(err)
	}
	calls := 0
	f.svc.runStaged = countingRunner("v1.0.0-beta.3", &calls)

	err := f.svc.PerformAtomicSwapAt(f.activePath)
	if err == nil {
		t.Error("expected a non-regular staged path not to be applied")
	}
	if got := f.activeContent(t); got != "old content" {
		t.Errorf("executable must not be replaced, got %q", got)
	}
	if calls != 0 {
		t.Errorf("a non-regular staged path must never be run, got %d runs", calls)
	}
	assertVerifyFailureCounted(t, f, "v1.0.0-beta.3")
}
