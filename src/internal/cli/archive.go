package cli

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jeancodogno/specforce-kit/src/internal/constitution"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// HandleArchive dispatches to the correct archive sub-command.
func (e *Executor) HandleArchive(ctx context.Context, ui core.UI, args ...string) error {
	subCommand := ""
	if len(args) > 0 {
		subCommand = args[0]
	}

	switch subCommand {
	case "instructions":
		return e.HandleArchiveInstructions(ctx, ui)
	default:
		fmt.Printf("Unknown archive command: %s\n", subCommand)
		fmt.Println("Available commands: instructions")
		return nil
	}
}

// HandleArchiveInstructions processes the 'archive instructions' command.
func (e *Executor) HandleArchiveInstructions(ctx context.Context, ui core.UI) error {
	// 1. Fetch Constitution Status
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

	// 2. Fetch Kit Instructions
	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return err
	}
	kitInstructions, err := fs.ReadFile(kitFS, "instructions/archive.md")
	if err != nil {
		// Fallback if file doesn't exist
		kitInstructions = []byte("Follow the standard archiving procedure.")
	}

	// 3. Fetch Config Instructions
	config := core.LoadConfig(".")
	customInstructions := config.Instructions["archive"]

	// 4. Print Combined Output
	fmt.Println("# ARCHIVE INSTRUCTIONS")
	fmt.Println()
	fmt.Println("## 1. Project Constitution Context")
	fmt.Println("These are the global standards of the project:")
	for _, art := range status.Artifacts {
		existsStr := "[MISSING]"
		if art.Exists {
			existsStr = "[EXISTS]"
		}
		fmt.Printf("- %s %s: %s (Path: %s)\n", existsStr, art.Name, art.Description, art.Path)
	}
	fmt.Println()

	fmt.Println("## 2. Core Archiving Rules")
	fmt.Println(string(kitInstructions))
	fmt.Println()

	if len(customInstructions) > 0 {
		fmt.Println("## 3. Project-Specific Rules (config.yaml)")
		for _, instr := range customInstructions {
			fmt.Printf("- %s\n", instr)
		}
		fmt.Println()
	}

	return nil
}
