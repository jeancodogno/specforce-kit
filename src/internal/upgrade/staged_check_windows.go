//go:build windows

package upgrade

import "os/exec"

// isolateStagedRun keeps the default behavior on Windows: the timeout kills
// the test-run process.
func isolateStagedRun(*exec.Cmd) {}

// killStagedRunGroup is a no-op on Windows.
func killStagedRunGroup(*exec.Cmd) {}
