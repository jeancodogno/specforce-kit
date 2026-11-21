package upgrade

import (
	"runtime"
	"testing"
)

// setTestHome points both home variables at a fresh temporary directory, so
// per-user state never touches the real profile: HOME on Linux and macOS,
// USERPROFILE on Windows.
func setTestHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return home
}

// skipWithoutPermissionBits skips tests that rely on Unix permission bits or
// on read-only directories made with os.Chmod, which Windows does not enforce.
func skipWithoutPermissionBits(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("permission bits do not apply on Windows")
	}
}
