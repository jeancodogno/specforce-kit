package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/project"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

// HandleImplementation dispatches to the correct implementation sub-command.
func (e *Executor) HandleImplementation(ctx context.Context, ui core.UI, args ...string) error {
	subCommand := ""
	if len(args) > 0 {
		subCommand = args[0]
	}

	switch subCommand {
	case "status":
		if len(args) < 2 {
			return fmt.Errorf("missing slug for implementation status. Usage: specforce implementation status <slug>")
		}
		jsonMode := false
		if len(args) > 2 && args[2] == "--json" {
			jsonMode = true
		}
		return e.HandleImplementationStatus(ctx, ui, args[1], jsonMode)
	case "update":
		// This is just for legacy switch, but we'll use Cobra for real work.
		return fmt.Errorf("update command requires specific flags, use Cobra-based CLI")
	default:
		fmt.Printf("Unknown implementation command: %s\n", subCommand)
		fmt.Println("Available commands: status, update")
		return nil
	}
}

// HandleImplementationStatus processes the 'implementation status' command.
func (e *Executor) HandleImplementationStatus(ctx context.Context, ui core.UI, slug string, jsonMode bool) error {
	projectRoot, err := spec.FindProjectRoot()
	if err != nil {
		return err
	}

	// Initialize SpecService lazily
	if e.SpecService == nil {
		artifactsFS, err := e.GetArtifactsFS(ui)
		if err != nil {
			return err
		}
		specFS, err := fs.Sub(artifactsFS, "spec")
		if err != nil {
			return fmt.Errorf("failed to load spec artifacts: %w", err)
		}
		registry, err := spec.NewRegistry(specFS)
		if err != nil {
			return fmt.Errorf("failed to initialize spec registry: %w", err)
		}

		// Initialize ProjectService if needed to pass as ConfigProvider
		if e.ProjectService == nil {
			kitFS, _ := e.GetKitFS(ui)
			e.ProjectService = project.NewService(kitFS, artifactsFS, ".")
		}
		e.SpecService = spec.NewService(registry, e.ProjectService)
	}

	report, err := e.SpecService.GetImplementationStatus(ctx, projectRoot, slug)
	if err != nil {
		if jsonMode {
			errData := map[string]string{"error": err.Error()}
			data, _ := json.MarshalIndent(errData, "", "  ")
			fmt.Println(string(data))
			os.Exit(1)
		}
		return fmt.Errorf("failed to get implementation status for %s: %w", slug, err)
	}

	if jsonMode {
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// TUI Mode
	return tui.RenderImplementationStatus(report)
}

// HandleImplementationUpdate processes the 'implementation update' command.
func (e *Executor) HandleImplementationUpdate(ctx context.Context, ui core.UI, slug, taskId, status string) error {
	projectRoot, err := spec.FindProjectRoot()
	if err != nil {
		return err
	}

	// Initialize SpecService lazily
	if e.SpecService == nil {
		artifactsFS, err := e.GetArtifactsFS(ui)
		if err != nil {
			return err
		}
		specFS, err := fs.Sub(artifactsFS, "spec")
		if err != nil {
			return fmt.Errorf("failed to load spec artifacts: %w", err)
		}
		registry, err := spec.NewRegistry(specFS)
		if err != nil {
			return fmt.Errorf("failed to initialize spec registry: %w", err)
		}

		// Initialize ProjectService if needed to pass as ConfigProvider
		if e.ProjectService == nil {
			kitFS, _ := e.GetKitFS(ui)
			e.ProjectService = project.NewService(kitFS, artifactsFS, ".")
		}
		e.SpecService = spec.NewService(registry, e.ProjectService)
	}

	if err := e.SpecService.UpdateTaskStatus(ctx, projectRoot, slug, taskId, status); err != nil {
		var hookErr *core.HookError
		if errors.As(err, &hookErr) {
			fmt.Println("\n[HOOK FAILURE] The following verification hooks failed:")
			for _, res := range hookErr.Results {
				if !res.Success {
					fmt.Printf("\n--- Command: %s ---\n", res.Command)
					if res.Stdout != "" {
						fmt.Printf("Stdout:\n%s\n", strings.TrimSpace(res.Stdout))
					}
					if res.Stderr != "" {
						fmt.Printf("Stderr:\n%s\n", strings.TrimSpace(res.Stderr))
					}
					fmt.Printf("Exit Code: %d\n", res.ExitCode)
				}
			}
			fmt.Println("\nUpdate aborted. Please fix the issues and try again.")
			return fmt.Errorf("task update blocked by hook failures")
		}
		return fmt.Errorf("failed to update task %s for %s: %w", taskId, slug, err)
	}

	fmt.Printf("[OK] Task %s status updated to %s in %s\n", taskId, status, slug)
	return nil
}
