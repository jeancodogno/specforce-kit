package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
	"github.com/jeancodogno/specforce-kit/src/internal/project"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

// Executor handles the CLI command orchestration and routing.
type Executor struct {
	Version        string
	DevMode        bool
	KitRoot        string
	ArtifactsRoot  string
	Registry       *agent.Registry
	ProjectService *project.Service
	SpecService    *spec.Service
}

// NewExecutor creates a new CLI executor.
func NewExecutor(Version string) *Executor {
	return &Executor{
		Version:       Version,
		KitRoot:       "src/internal/agent/kit",
		ArtifactsRoot: "src/internal/agent/artifacts",
		Registry:      &agent.Registry{},
	}
}

func (e *Executor) HandleInstall(ctx context.Context, ui core.UI) error {
	if tui.IsTTY() {
		tui.PrintBranding()
	}
	err := installer.Install(ctx, ui)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrInstallerPermissionDenied):
			// installer already printed the sudo suggestion; exit gracefully
			return nil
		}
	}
	return err
}

func (e *Executor) HandleInit(ctx context.Context, ui core.UI, agents ...string) error {
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return err
	}

	selected, err := e.ResolveSelectedAgents(ctx, ui, agents...)
	if err != nil {
		return err
	}

	// Check if already initialized for update flow
	if project.IsInitialized(".") {
		return e.handleUpdateFlow(ctx, ui, kitFS, selected)
	}

	return e.handleNewInitFlow(ctx, ui, kitFS, selected)
}

func (e *Executor) handleUpdateFlow(ctx context.Context, ui core.UI, kitFS fs.FS, selected []string) error {
	if ui != nil && ui.Confirm("Do you want to update agent tools and instructions (.gemini, .claude, etc)?") {
		// Initialize ProjectService lazily
		if e.ProjectService == nil {
			e.ProjectService = project.NewService(kitFS, nil, ".")
		}

		if err := e.ProjectService.UpdateTools(ctx, ui, selected); err != nil {
			return err
		}
		if tui.IsTTY() {
			tui.PrintCompletionBox("UPDATE COMPLETE", "Agent tools and instructions updated successfully.")
		}
		return nil
	}
	if ui != nil {
		ui.Warn("Initialization cancelled. Project remains unchanged.")
	}
	return nil
}

func (e *Executor) handleNewInitFlow(ctx context.Context, ui core.UI, kitFS fs.FS, selected []string) error {
	if tui.IsTTY() {
		ui.SubTask(fmt.Sprintf("Initializing project with agents: %s", strings.Join(selected, ", ")))
	}

	artifactsFS, err := e.GetArtifactsFS(ui)
	if err != nil {
		return err
	}

	// Initialize ProjectService lazily
	if e.ProjectService == nil {
		e.ProjectService = project.NewService(kitFS, artifactsFS, ".")
	}

	config := project.InitConfig{
		ProjectRoot:    ".",
		SelectedAgents: selected,
	}

	if err := e.ProjectService.InitializeProject(ctx, ui, config); err != nil {
		return err
	}

	// Auto-initialize config.yaml if it doesn't exist
	if err := core.EnsureConfigExists("."); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config.yaml: %v\n", err)
	}

	if tui.IsTTY() {
		tui.PrintCompletionBox("MISSION ACCOMPLISHED", "Specforce initialized successfully!\nYou can now start using SDD with your selected AI agents.")
	}
	return nil
}

func (e *Executor) GetKitFS(ui core.UI) (fs.FS, error) {
	if e.DevMode {
		if tui.IsTTY() {
			ui.Warn(fmt.Sprintf("DEVELOPMENT MODE ACTIVE: Using local %s", e.KitRoot))
		}
		return os.DirFS(e.KitRoot), nil
	}

	kitFS, err := agent.GetKitFS()
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded kit: %w", err)
	}
	return kitFS, nil
}

func (e *Executor) GetArtifactsFS(ui core.UI) (fs.FS, error) {
	if e.DevMode {
		if tui.IsTTY() {
			ui.Warn(fmt.Sprintf("DEVELOPMENT MODE ACTIVE: Using local %s", e.ArtifactsRoot))
		}
		return os.DirFS(e.ArtifactsRoot), nil
	}

	artifactsFS, err := agent.GetArtifactsFS()
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded artifacts: %w", err)
	}
	return artifactsFS, nil
}

func (e *Executor) ResolveSelectedAgents(ctx context.Context, ui core.UI, agents ...string) ([]string, error) {
	// For HandleInit, we need a kitFS
	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return nil, err
	}

	// Ensure registry is initialized
	if err := e.Registry.Initialize(kitFS, "."); err != nil {
		return nil, fmt.Errorf("failed to initialize agent registry: %w", err)
	}

	if len(agents) > 0 {
		// Normalize and validate provided agents
		normalized := make([]string, len(agents))
		for i, a := range agents {
			n := a
			// Aliases normalization
			switch a {
			case "opencode":
				n = "open-code"
			case "kilocode":
				n = "kilo-code"
			case "qwen-code":
				n = "qwen"
			case "kimicode":
				n = "kimi-code"
			}

			if _, ok := e.Registry.GetAgent(n); !ok {
				return nil, fmt.Errorf("agent %q is not supported. Use 'specforce init' without arguments to see available agents", a)
			}
			normalized[i] = n
		}
		return normalized, nil
	}

	if !tui.IsTTY() {
		return nil, fmt.Errorf("no agents specified and no TTY detected")
	}

	existing := project.DetectExistingAgents(ctx, ".", e.Registry)
	selected, err := tui.SelectAgents(e.Registry.GetAgents(), existing)
	if err != nil {
		if err.Error() == "aborted" {
			fmt.Println("Project initialization aborted.")
			os.Exit(0)
		}
		return nil, fmt.Errorf("TUI failure: %w", err)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no agents selected; project initialization aborted")
	}

	return selected, nil
}

func (e *Executor) HandleConsole(ctx context.Context, ui core.UI) error {
	return e.RunConsole(ctx, ui)
}

func (e *Executor) HandleHelp() {
	if tui.IsTTY() {
		tui.PrintBranding()
	}
	e.PrintUsage()
	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
}

func (e *Executor) HandleUnknown(command string) {
	if tui.IsTTY() {
		tui.PrintBranding()
	}
	fmt.Printf("Unknown command: %s\n", command)
	e.PrintUsage()
	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
}

func (e *Executor) PrintUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  specforce [flags] [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  install       Global framework installation")
	fmt.Println("  init          Project initialization with TUI")
	fmt.Println("  constitution  Manage project constitution docs")
	fmt.Println("  spec          Manage feature specification artifacts")
	fmt.Println("  implementation Task tracking and implementation status")
	fmt.Println("  console       Launch the Specforce Console TUI")
	fmt.Println("\nFlags:")
	// flag.PrintDefaults() - no longer using standard flags here.
}
