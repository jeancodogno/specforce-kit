package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/project"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

// HandleSpec dispatches to the correct spec sub-command.
func (e *Executor) HandleSpec(ctx context.Context, ui core.UI, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Available commands: init, list, status, artifact, archive")
		return nil
	}

	subCommand := args[0]
	jsonMode := false
	for _, arg := range args {
		if arg == "--json" {
			jsonMode = true
			break
		}
	}

	switch subCommand {
	case "init":
		return e.handleSpecInitCmd(ctx, ui, args, jsonMode)
	case "list":
		return e.HandleSpecList(ctx, ui, jsonMode)
	case "status":
		return e.handleSpecStatusCmd(ctx, ui, args, jsonMode)
	case "artifact":
		return e.handleSpecArtifactCmd(ctx, ui, args, jsonMode)
	case "archive":
		return e.handleSpecArchiveCmd(ctx, ui, args)
	default:
		fmt.Printf("Unknown spec command: %s\n", subCommand)
		fmt.Println("Available commands: init, list, status, artifact, archive")
		return nil
	}
}

func (e *Executor) handleSpecInitCmd(ctx context.Context, ui core.UI, args []string, jsonMode bool) error {
	if len(args) < 2 {
		return fmt.Errorf("missing slug for spec init. Usage: specforce spec init <slug>")
	}
	return e.HandleSpecInit(ctx, ui, args[1], jsonMode)
}

func (e *Executor) handleSpecStatusCmd(ctx context.Context, ui core.UI, args []string, jsonMode bool) error {
	if len(args) < 2 {
		return fmt.Errorf("missing slug for spec status. Usage: specforce spec status <slug>")
	}
	return e.HandleSpecStatus(ctx, ui, args[1], jsonMode)
}

func (e *Executor) handleSpecArtifactCmd(ctx context.Context, ui core.UI, args []string, jsonMode bool) error {
	slug := ""
	if len(args) > 1 && args[1] != "--json" {
		slug = args[1]
	}
	return e.HandleSpecArtifact(ctx, ui, slug, jsonMode)
}

func (e *Executor) handleSpecArchiveCmd(ctx context.Context, ui core.UI, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing slug for spec archive. Usage: specforce spec archive <slug>")
	}
	force := false
	for _, arg := range args {
		if arg == "--force" {
			force = true
			break
		}
	}
	return e.HandleSpecArchive(ctx, ui, args[1], force)
}

// HandleSpecInit processes the 'spec init' command.
func (e *Executor) HandleSpecInit(ctx context.Context, ui core.UI, slug string, jsonMode bool) error {
	projectRoot, err := spec.FindProjectRoot()
	if err != nil {
		return err
	}

	// [REQ-1] Auto-timestamped slugs
	slug = spec.PrepareSlug(slug)

	exists, status := spec.SpecExists(projectRoot, slug)
	if exists {
		var errObj error
		if status == "active" {
			errObj = core.ErrSpecAlreadyActive
		} else {
			errObj = core.ErrSpecAlreadyArchived
		}

		if jsonMode {
			res := map[string]string{
				"status": "error",
				"error":  errObj.Error(),
			}
			data, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		return errObj
	}

	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0750); err != nil {
		return fmt.Errorf("failed to create spec directory %s: %w", specDir, err)
	}

	if jsonMode {
		res := map[string]string{
			"status":  "ok",
			"message": fmt.Sprintf("Spec directory initialized: .specforce/specs/%s", slug),
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("[OK] Spec directory initialized: .specforce/specs/%s\n", slug)
	return nil
}

// HandleSpecList processes the 'spec list' command.
func (e *Executor) HandleSpecList(ctx context.Context, ui core.UI, jsonMode bool) error {
	projectRoot, err := spec.FindProjectRoot()
	if err != nil {
		return err
	}

	specs, err := spec.ListActiveSpecs(ctx, projectRoot)
	if err != nil {
		return fmt.Errorf("failed to list active specs: %w", err)
	}

	if jsonMode {
		data, err := json.MarshalIndent(specs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if len(specs) == 0 {
		fmt.Println("No active specs found.")
		return nil
	}

	for _, s := range specs {
		fmt.Println(s.Slug)
	}

	return nil
}

// HandleSpecStatus processes the 'spec status' command.
func (e *Executor) HandleSpecStatus(ctx context.Context, ui core.UI, slug string, jsonMode bool) error {
	artifactsFS, err := e.GetArtifactsFS(ui)
	if err != nil {
		return err
	}
	specFS, err := fs.Sub(artifactsFS, "spec")
	if err != nil {
		return fmt.Errorf("failed to load spec artifacts: %w", err)
	}

	// Initialize SpecService lazily
	if e.SpecService == nil {
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

	status, err := e.SpecService.GetStatus(ctx, ".", slug)
	if err != nil {
		if jsonMode {
			errData := map[string]string{"error": err.Error()}
			data, _ := json.MarshalIndent(errData, "", "  ")
			fmt.Println(string(data))
			os.Exit(1)
		}
		return err
	}

	if jsonMode {
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal status to JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	e.renderSpecStatusTUI(slug, status)
	return nil
}

func (e *Executor) renderSpecStatusTUI(slug string, status spec.SpecStatus) {
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	fmt.Println("\n" + tui.HeaderStyle.Render("SPEC COMPLETENESS REPORT: "+strings.ToUpper(slug)))
	fmt.Println(tui.SubtitleStyle.Render("Target: .specforce/specs/"+slug+"/"))
	fmt.Println()

	fmt.Print(tui.RenderSpecStatus(status))
	fmt.Println()

	// Progress Bar
	fmt.Println(tui.RenderProgressBar(status.Progress, 40))
	fmt.Println()

	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
}

// HandleSpecArtifact processes the 'spec artifact' command.
func (e *Executor) HandleSpecArtifact(ctx context.Context, ui core.UI, slug string, jsonMode bool) error {
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

	if slug == "" {
		return e.ListSpecArtifacts(ui, registry, jsonMode)
	}

	// Initialize SpecService lazily
	if e.SpecService == nil {
		// Initialize ProjectService if needed to pass as ConfigProvider
		if e.ProjectService == nil {
			kitFS, _ := e.GetKitFS(ui)
			e.ProjectService = project.NewService(kitFS, artifactsFS, ".")
		}
		e.SpecService = spec.NewService(registry, e.ProjectService)
	}

	art, err := e.SpecService.GetArtifact(ctx, slug)
	if err != nil {
		fmt.Printf("Unknown spec artifact slug: %s\n", slug)
		fmt.Println("Use 'specforce spec artifact' to list available slugs.")
		return nil
	}

	if jsonMode {
		data, err := json.MarshalIndent(art, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	e.renderSpecArtifactTUI(art)
	return nil
	}

	func (e *Executor) renderSpecArtifactTUI(art *spec.Artifact) {
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	fmt.Println("\n" + tui.HeaderStyle.Render("SPEC ARTIFACT: "+strings.ToUpper(art.Name)))
	fmt.Println(tui.SubtitleStyle.Render("Description: " + art.Description))
	if art.Dependency != "" {
		fmt.Println(tui.SubtitleStyle.Render("Dependency:  " + art.Dependency))
	}
	fmt.Println()

	fmt.Println(tui.HeaderStyle.Render("> INSTRUCTIONS"))
	fmt.Println(tui.BodyStyle.Render(art.Instruction))
	fmt.Println()

	fmt.Println(tui.HeaderStyle.Render("> TEMPLATE"))
	fmt.Println(tui.BodyStyle.Render(art.Template))
	fmt.Println()

	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
	}


// ListSpecArtifacts lists all available spec artifacts.
func (e *Executor) ListSpecArtifacts(ui core.UI, registry *spec.Registry, jsonMode bool) error {
	artifacts := registry.List()

	if jsonMode {
		data, err := json.MarshalIndent(artifacts, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	displayItems := make([]tui.ArtifactDisplay, len(artifacts))
	for i, art := range artifacts {
		displayItems[i] = tui.ArtifactDisplay{
			ID:          art.Name,
			Description: art.Description,
		}
	}

	tui.PrintArtifactList(
		"AVAILABLE SPEC ARTIFACTS",
		"Use 'specforce spec artifact [slug]' to see full details.",
		displayItems,
		e.Version,
	)

	return nil
}

// HandleSpecArchive processes the 'spec archive' command.
func (e *Executor) HandleSpecArchive(ctx context.Context, ui core.UI, slug string, force bool) error {
	projectRoot, err := spec.FindProjectRoot()
	if err != nil {
		return err
	}

	// [REQ-4] Check task completion
	tasksPath := filepath.Join(projectRoot, ".specforce", "specs", slug, "tasks.md")
	if _, err := os.Stat(tasksPath); err == nil {
		report, err := spec.ParseTasks(ctx, projectRoot, slug)
		if err == nil && report.Status != "finished" && !force {
			err := fmt.Errorf("cannot archive '%s': pending tasks remain. Use --force to override", slug)
			ui.Error(err.Error() + "\n")
			return err
		}
	}

	if err := spec.ArchiveSpec(ctx, projectRoot, slug); err != nil {
		ui.Error(fmt.Sprintf("Archive failed: %v\n", err))
		return err
	}

	ui.Success(fmt.Sprintf("Spec successfully archived: %s\n", slug))
	return nil
}
