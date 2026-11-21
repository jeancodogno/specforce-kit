// +build !windows

package upgrade

import (
	"os"
	"syscall"
)

func getDetachedSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func executeNewBinary(binaryPath string) error {
	// #nosec G204,G702 - binaryPath is obtained via os.Executable() or is verified.
	return syscall.Exec(binaryPath, os.Args, os.Environ())
}
