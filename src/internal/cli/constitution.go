package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/constitution"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

// HandleConstitution dispatches to the correct constitution sub-command.
func (e *Executor) HandleConstitution(ctx context.Context, ui core.UI, args ...string) error {
	subCommand := ""
	if len(args) > 0 {
		subCommand = args[0]
	}

	switch subCommand {
	case "status":
		jsonMode := false
		if len(args) > 1 && args[1] == "--json" {
			jsonMode = true
		}
		return e.HandleConstitutionStatus(ctx, ui, jsonMode)
	case "artifact":
		slug := ""
		if len(args) > 1 {
			slug = args[1]
		}
		jsonMode := false
		if len(args) > 2 && args[2] == "--json" {
			jsonMode = true
		}
		return e.HandleConstitutionArtifact(ctx, ui, slug, jsonMode)
	default:
		fmt.Printf("Unknown constitution command: %s\n", subCommand)
		fmt.Println("Available commands: status, artifact")
		return nil
	}
}

// HandleConstitutionStatus processes the 'constitution status' command.
func (e *Executor) HandleConstitutionStatus(ctx context.Context, ui core.UI, jsonMode bool) error {
	artifactsFS, err := e.GetArtifactsFS(ui)
	if err != nil {
		return err
	}

	constitutionFS, err := fs.Sub(artifactsFS, "constitution")
	if err != nil {
		return fmt.Errorf("failed to sub-filesystem for constitution: %w", err)
	}

	registry, err := constitution.NewRegistry(constitutionFS)
	if err != nil {
		return fmt.Errorf("failed to initialize constitution registry: %w", err)
	}

	status, err := constitution.GetStatus(ctx, ".", registry)
	if err != nil {
		return fmt.Errorf("failed to scan constitution: %w", err)
	}

	if jsonMode {
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// TUI Mode
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	fmt.Println("\n" + tui.HeaderStyle.Render("CONSTITUTION COMPLETENESS REPORT"))
	fmt.Println(tui.SubtitleStyle.Render("Target: .specforce/docs/"))
	fmt.Println()

	fmt.Print(tui.RenderConstitutionStatus(status))
	fmt.Println()

	// Progress Bar (width 40 as per design or standard)
	fmt.Println(tui.RenderProgressBar(status.Progress, 40))
	fmt.Println()

	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}

	return nil
}

// HandleConstitutionArtifact processes the 'constitution artifact' command.
func (e *Executor) HandleConstitutionArtifact(ctx context.Context, ui core.UI, slug string, jsonMode bool) error {
	artifactsFS, err := e.GetArtifactsFS(ui)
	if err != nil {
		return err
	}

	constitutionFS, err := fs.Sub(artifactsFS, "constitution")
	if err != nil {
		return fmt.Errorf("failed to sub-filesystem for constitution: %w", err)
	}

	registry, err := constitution.NewRegistry(constitutionFS)
	if err != nil {
		return fmt.Errorf("failed to initialize constitution registry: %w", err)
	}

	if slug == "" {
		return e.ListArtifacts(ui, registry, jsonMode)
	}

	art, ok := registry.Get(slug)
	if !ok {
		fmt.Printf("Unknown artifact slug: %s\n", slug)
		fmt.Println("Use 'specforce constitution artifact' to list available slugs.")
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

	// TUI Mode
	if tui.IsTTY() {
		tui.PrintBranding()
	}

	fmt.Println("\n" + tui.HeaderStyle.Render("CONSTITUTION ARTIFACT: "+strings.ToUpper(art.Slug)))
	fmt.Println(tui.SubtitleStyle.Render("Description: " + art.Description))
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

	return nil
}

// ListArtifacts lists all available constitution artifacts.
func (e *Executor) ListArtifacts(ui core.UI, registry *constitution.Registry, jsonMode bool) error {
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
			ID:          art.Slug,
			Description: art.Description,
		}
	}

	tui.PrintArtifactList(
		"AVAILABLE CONSTITUTION ARTIFACTS",
		"Use 'specforce constitution artifact [slug]' to see full details.",
		displayItems,
		e.Version,
	)

	return nil
}
