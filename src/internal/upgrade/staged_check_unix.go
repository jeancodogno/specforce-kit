//go:build !windows

package upgrade

import (
	"os/exec"
	"syscall"
)

// isolateStagedRun starts the test run in its own process group and makes the
// timeout kill the whole group, so children of the staged binary die with it.
func isolateStagedRun(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}

// killStagedRunGroup kills whatever is left of the test run's process group.
func killStagedRunGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
