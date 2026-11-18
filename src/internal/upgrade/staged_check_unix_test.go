//go:build !windows

package upgrade

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// --- Timed-out test run kills the whole process group (correction cycle 1, item 5) ---

func TestStagedRun_TimeoutKillsBackgroundChildren(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	t.Setenv("SPF_TEST_CHILD_PIDFILE", pidFile)
	script := "#!/bin/sh\nsleep 30 &\necho $! > \"$SPF_TEST_CHILD_PIDFILE\"\nsleep 30\necho \"specforce version 1.0.0-beta.3\"\n"
	f := newScriptSwapFixture(t, script)
	f.svc.stagedRunTimeout = 500 * time.Millisecond

	start := time.Now()
	err := f.svc.PerformAtomicSwapAt(f.activePath)
	if elapsed := time.Since(start); elapsed > 7*time.Second {
		t.Errorf("expected the test run to be bounded by its timeout, took %v", elapsed)
	}
	assertDiscarded(t, f, err)

	b, readErr := os.ReadFile(pidFile)
	if readErr != nil {
		t.Fatalf("the script did not record its background child: %v", readErr)
	}
	pid, convErr := strconv.Atoi(strings.TrimSpace(string(b)))
	if convErr != nil || pid <= 0 {
		t.Fatalf("invalid child pid %q", string(b))
	}
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })

	// The killed child is reparented and reaped by init; allow a short delay.
	deadline := time.Now().Add(3 * time.Second)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("background child %d survived the timed-out test run", pid)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// processAlive reports whether pid exists and is not a zombie.
func processAlive(pid int) bool {
	if err := syscall.Kill(pid, 0); err != nil {
		return false
	}
	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return true // no procfs: rely on kill(0)
	}
	// Format: pid (comm) state ...
	s := string(stat)
	if i := strings.LastIndex(s, ")"); i >= 0 && i+2 < len(s) {
		return s[i+2] != 'Z'
	}
	return true
}
