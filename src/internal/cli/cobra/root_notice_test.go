package cobra

import (
	"bytes"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
	"github.com/spf13/cobra"
)

// setupNoticeTest isolates HOME, records a pending manual-update version and
// installs a service and TTY check for the post-run hook.
func setupNoticeTest(t *testing.T, tty bool) *upgrade.StateManager {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.Save(&upgrade.State{
		LatestVersion:       "v1.0.0-beta.3",
		ManualUpdateVersion: "v1.0.0-beta.3",
	}); err != nil {
		t.Fatal(err)
	}

	origSvc, origTTY, origVer := upgradeService, isTTY, appVersion
	t.Cleanup(func() { upgradeService, isTTY, appVersion = origSvc, origTTY, origVer })
	appVersion = "1.0.0-beta.2"
	upgradeService = upgrade.NewService(mgr, &upgrade.MockProvider{}, appVersion)
	isTTY = func() bool { return tty }
	return mgr
}

func newNoticeCmd(agent bool) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	parent := &cobra.Command{Use: "parent"}
	if agent {
		parent.Annotations = map[string]string{"IsAgentCommand": "true"}
	}
	child := &cobra.Command{Use: "child", Run: func(*cobra.Command, []string) {}}
	parent.AddCommand(child)
	var stdout, stderr bytes.Buffer
	child.SetOut(&stdout)
	child.SetErr(&stderr)
	return child, &stdout, &stderr
}

// setupStableNoticeTest is setupNoticeTest with a state holding only a pending
// stable version (no manual-update version, which would take precedence).
func setupStableNoticeTest(t *testing.T, tty bool) *upgrade.StateManager {
	t.Helper()
	mgr := setupNoticeTest(t, tty)
	if err := mgr.Save(&upgrade.State{
		LatestVersion: "v1.0.0-beta.2",
		StableVersion: "1.0.0",
	}); err != nil {
		t.Fatal(err)
	}
	return mgr
}

func TestRootUpdateNotice(t *testing.T) {
	t.Run("agent command writes nothing", func(t *testing.T) {
		mgr := setupNoticeTest(t, true)
		cmd, stdout, stderr := newNoticeCmd(true)
		rootCmd.PersistentPostRun(cmd, nil)
		if stdout.Len() != 0 || stderr.Len() != 0 {
			t.Errorf("expected no output, got stdout=%q stderr=%q", stdout, stderr)
		}
		if s, _ := mgr.Load(); s.ManualNoticeShownVersion != "" {
			t.Error("hint not written must not be marked as shown")
		}
	})

	t.Run("non-TTY writes nothing", func(t *testing.T) {
		mgr := setupNoticeTest(t, false)
		cmd, stdout, stderr := newNoticeCmd(false)
		rootCmd.PersistentPostRun(cmd, nil)
		if stdout.Len() != 0 || stderr.Len() != 0 {
			t.Errorf("expected no output, got stdout=%q stderr=%q", stdout, stderr)
		}
		if s, _ := mgr.Load(); s.ManualNoticeShownVersion != "" {
			t.Error("hint not written must not be marked as shown")
		}
	})

	t.Run("interactive non-agent writes once to stderr only", func(t *testing.T) {
		setupNoticeTest(t, true)
		cmd, stdout, stderr := newNoticeCmd(false)
		rootCmd.PersistentPostRun(cmd, nil)
		if stdout.Len() != 0 {
			t.Errorf("notice must never be written to stdout, got %q", stdout)
		}
		if !bytes.Contains(stderr.Bytes(), []byte("v1.0.0-beta.3")) {
			t.Errorf("expected manual hint naming v1.0.0-beta.3 on stderr, got %q", stderr)
		}
		if bytes.Contains(stderr.Bytes(), []byte("applied automatically")) {
			t.Errorf("manual hint must not claim automatic application: %q", stderr)
		}

		cmd2, stdout2, stderr2 := newNoticeCmd(false)
		rootCmd.PersistentPostRun(cmd2, nil)
		if stdout2.Len() != 0 || stderr2.Len() != 0 {
			t.Errorf("hint must be shown once per version, got stdout=%q stderr=%q", stdout2, stderr2)
		}
	})

	t.Run("stable interactive non-agent writes once to stderr and marks shown", testStableNoticeShownOnce)
	t.Run("stable agent command writes nothing", func(t *testing.T) { testStableNoticeSuppressed(t, true, true) })
	t.Run("stable non-TTY writes nothing", func(t *testing.T) { testStableNoticeSuppressed(t, false, false) })
}

// testStableNoticeShownOnce checks that a pending stable notice is written to
// stderr only, marked as shown, and not written again by the next command.
func testStableNoticeShownOnce(t *testing.T) {
	mgr := setupStableNoticeTest(t, true)
	cmd, stdout, stderr := newNoticeCmd(false)
	rootCmd.PersistentPostRun(cmd, nil)
	if stdout.Len() != 0 {
		t.Errorf("notice must never be written to stdout, got %q", stdout)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("STABLE RELEASE AVAILABLE")) ||
		!bytes.Contains(stderr.Bytes(), []byte("@jeancodogno/specforce-kit@latest")) {
		t.Errorf("expected stable notice on stderr, got %q", stderr)
	}
	if s, _ := mgr.Load(); s.StableNoticeShownVersion != "1.0.0" {
		t.Errorf("expected stable notice marked shown for 1.0.0, got %q", s.StableNoticeShownVersion)
	}

	cmd2, stdout2, stderr2 := newNoticeCmd(false)
	rootCmd.PersistentPostRun(cmd2, nil)
	if stdout2.Len() != 0 || stderr2.Len() != 0 {
		t.Errorf("stable notice must be shown once per version, got stdout=%q stderr=%q", stdout2, stderr2)
	}
}

// testStableNoticeSuppressed checks that a pending stable notice is neither
// written nor marked as shown for agent or non-TTY commands.
func testStableNoticeSuppressed(t *testing.T, tty, agent bool) {
	mgr := setupStableNoticeTest(t, tty)
	cmd, stdout, stderr := newNoticeCmd(agent)
	rootCmd.PersistentPostRun(cmd, nil)
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Errorf("expected no output, got stdout=%q stderr=%q", stdout, stderr)
	}
	if s, _ := mgr.Load(); s.StableNoticeShownVersion != "" {
		t.Error("stable notice not written must not be marked as shown")
	}
}
