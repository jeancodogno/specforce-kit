// +build windows

package upgrade

import (
	"os"
	"os/exec"
	"syscall"
)

func getDetachedSysProcAttr() *syscall.SysProcAttr {
	// On Windows, there is no direct equivalent to Setsid, 
	// but we can use flags to detach from the parent console.
	return &syscall.SysProcAttr{
		CreationFlags: 0x00000010, // CREATE_NEW_CONSOLE
	}
}

func executeNewBinary(binaryPath string) error {
	// Windows fallback: spawn a new process and exit the current one.
	// #nosec G204,G702 - binaryPath is obtained via os.Executable() or is verified.
	cmd := exec.Command(binaryPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
