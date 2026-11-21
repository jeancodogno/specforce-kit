package spec

import (
	"runtime"
	"testing"
)

// skipWithoutPermissionBits skips assertions on Unix permission bits, which
// Windows does not enforce.
func skipWithoutPermissionBits(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("permission bits do not apply on Windows")
	}
}
