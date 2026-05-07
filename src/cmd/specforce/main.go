package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/cli/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
)

var version = "1.0.0"

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
	checkForSwap()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cobra.Execute(ctx, version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func checkForSwap() {
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		return
	}

	state, err := mgr.Load()
	if err != nil || !state.UpdateReady {
		return
	}

	// Create service to perform swap
	// We use a mock provider here because we won't need it for swap
	svc := upgrade.NewService(mgr, &upgrade.MockProvider{}, version)
	
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	stagedVer := state.StagedVersion
	if err := svc.PerformAtomicSwapAt(exePath); err != nil {
		// Failed to swap, just continue
		return
	}

	// Subtle notification
	fmt.Printf("\n\x1b[38;5;245m(Specforce updated to %s)\x1b[0m\n", stagedVer)

	// Restart with new binary
	_ = svc.ExecuteBinary(exePath)
	// If it fails to restart, we just continue with the old binary in memory
}

func runInternalUpgradeCheck() {
	mgr, err := upgrade.NewStateManager()
	if err != nil {
		os.Exit(1)
	}

	// TODO: Proper provider detection. For now, default to GitHub.
	provider := upgrade.NewGitHubProvider()
	svc := upgrade.NewService(mgr, provider, version)

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
