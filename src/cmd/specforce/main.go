package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/cli/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
)

func main() {
	// 1. Check for internal upgrade check flag (must be first to be silent and fast)
	for _, arg := range os.Args {
		if arg == "--internal-upgrade-check" {
			runInternalUpgradeCheck()
			return
		}
	}

	// 2. Check if an update is ready to be swapped
	// This happens before any command execution to ensure we use the new binary
	checkForSwap(defaultSwapEnv())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cobra.Execute(ctx, core.Version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// swapEnv holds the process-level dependencies of the pre-command swap so it can
// be exercised without replacing the test process.
type swapEnv struct {
	version    string
	executable func() (string, error)
	exec       func(path string) error
}

func defaultSwapEnv() swapEnv {
	return swapEnv{
		version:    core.Version,
		executable: os.Executable,
		exec:       (&upgrade.Service{}).ExecuteBinary,
	}
}

// checkForSwap applies a staged update before any command runs (agent or not).
// It never writes to stdout or stderr: the applied version is recorded in state
// and announced later, only for interactive non-agent commands.
func checkForSwap(env swapEnv) {
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		return
	}

	state, err := mgr.Load()
	if err != nil || !state.UpdateReady {
		return
	}

	// We use a mock provider here because we won't need it for swap
	svc := upgrade.NewService(mgr, &upgrade.MockProvider{}, env.version)

	exePath, err := env.executable()
	if err != nil {
		return
	}

	if err := svc.PerformAtomicSwapAt(exePath); err != nil {
		// Not applied (skipped or failed): just continue on the current binary
		return
	}

	// Restart with new binary. The .old backup is kept until the outcome is
	// known: when the restart fails, the previous executable is restored and the
	// command continues on the old binary in memory.
	if err := env.exec(exePath); err != nil {
		_ = svc.RollbackFailedRestartAt(exePath, err)
	}
}

func runInternalUpgradeCheck() {
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		os.Exit(1)
	}

	// Version source follows the installation method (npm channel or GitHub)
	exe, _ := os.Executable()
	svc := upgrade.NewService(mgr, upgrade.SelectProvider(exe, core.Version), core.Version)

	// Perform the background update (check + stage)
	// Timeout after 1 minute to avoid zombie processes
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := svc.PerformBackgroundUpdate(ctx); err != nil {
		// Silent failure
		os.Exit(1)
	}
	os.Exit(0)
}
