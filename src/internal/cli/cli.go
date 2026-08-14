package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
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

func (e *Executor) HandleInit(ctx context.Context, ui core.UI, agents ...string) error {
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return err
	}

	selected, removed, err := e.ResolveSelectedAgents(ctx, ui, agents...)
	if err != nil {
		return err
	}

	// Auto-initialize config.yaml if it doesn't exist
	if err := core.EnsureConfigExists("."); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config.yaml: %v\n", err)
	}

	// Check if already initialized for update flow
	if project.IsInitialized(".") {
		return e.handleUpdateFlow(ctx, ui, kitFS, selected, removed)
	}

	return e.handleNewInitFlow(ctx, ui, kitFS, selected)
}

func (e *Executor) handleUpdateFlow(ctx context.Context, ui core.UI, kitFS fs.FS, selected []string, removed []string) error {
	// If the user already confirmed in the TUI (removed is non-empty and they passed),
	// or if it's a direct CLI command, we might not need another confirm.
	// However, if it's an interactive TUI session, we should trust the TUI result.
	
	// Initialize ProjectService lazily
	if e.ProjectService == nil {
		e.ProjectService = project.NewService(kitFS, nil, ".")
	}

	// Decommission removed agents
	if err := e.ProjectService.DecommissionAgents(ctx, ui, removed); err != nil {
		return err
	}

	if err := e.ProjectService.UpdateTools(ctx, ui, selected); err != nil {
		return err
	}
	if tui.IsTTY() {
		tui.PrintCompletionBox("UPDATE COMPLETE", "Agent tools and instructions updated successfully.")
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

	if tui.IsTTY() {
		tui.PrintCompletionBox("MISSION ACCOMPLISHED", "Specforce structure is live.\n\nNEXT: Run '/spf:discovery' to start the SDD cycle.")
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

func (e *Executor) ResolveSelectedAgents(ctx context.Context, ui core.UI, agents ...string) ([]string, []string, error) {
	// For HandleInit, we need a kitFS
	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return nil, nil, err
	}

	// Ensure registry is initialized
	if err := e.Registry.Initialize(kitFS, "."); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize agent registry: %w", err)
	}

	if len(agents) > 0 {
		return e.resolveDirectAgents(ctx, agents...)
	}

	if !tui.IsTTY() {
		return nil, nil, fmt.Errorf("no agents specified and no TTY detected")
	}

	existing := project.DetectExistingAgents(ctx, ".", e.Registry)
	selected, removed, err := tui.SelectAgents(e.Registry.GetAgents(), existing)
	if err != nil {
		if err.Error() == "aborted" {
			fmt.Println("Project initialization aborted.")
			os.Exit(0)
		}
		return nil, nil, fmt.Errorf("TUI failure: %w", err)
	}

	if len(selected) == 0 {
		return nil, nil, fmt.Errorf("no agents selected; project initialization aborted")
	}

	return selected, removed, nil
}

func (e *Executor) resolveDirectAgents(ctx context.Context, agents ...string) ([]string, []string, error) {
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
			return nil, nil, fmt.Errorf("agent %q is not supported. Use 'specforce init' without arguments to see available agents", a)
		}
		normalized[i] = n
	}

	// If project is initialized, we should also detect what's NOT in 'normalized' but is currently installed
	var removed []string
	if project.IsInitialized(".") {
		existing := project.DetectExistingAgents(ctx, ".", e.Registry)
		for _, ex := range existing {
			if !contains(normalized, ex) {
				removed = append(removed, ex)
			}
		}
	}

	return normalized, removed, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
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
	fmt.Println("  init          Project initialization with TUI")
	fmt.Println("  constitution  Manage project constitution docs")
	fmt.Println("  spec          Manage feature specification artifacts")
	fmt.Println("  implementation Task tracking and implementation status")
	fmt.Println("  console       Launch the Specforce Console TUI")
	fmt.Println("\nFlags:")
	// flag.PrintDefaults() - no longer using standard flags here.
}
