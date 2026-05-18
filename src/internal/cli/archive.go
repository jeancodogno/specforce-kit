package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
	"github.com/jeancodogno/specforce-kit/src/internal/constitution"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/project"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
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
	case "memorial":
		if len(args) < 5 {
			return fmt.Errorf("missing arguments for archive memorial. Usage: specforce archive memorial <slug> <type> <title> <content> [author]")
		}
		author := "agent"
		if len(args) > 5 {
			author = args[5]
		}
		return e.HandleArchiveMemorial(ctx, ui, args[1], args[2], args[3], args[4], author)
	case "distill":
		if len(args) < 3 {
			return fmt.Errorf("missing arguments for archive distill. Usage: specforce archive distill <slugs> <summary> [author]")
		}
		author := "agent"
		if len(args) > 3 {
			author = args[3]
		}
		slugs := strings.Split(args[1], ",")
		return e.HandleArchiveDistill(ctx, ui, slugs, args[2], author)
	default:
		fmt.Printf("Unknown archive command: %s\n", subCommand)
		fmt.Println("Available commands: instructions, memorial, distill")
		return nil
	}
}

// HandleArchiveDistill processes the 'archive distill' command.
func (e *Executor) HandleArchiveDistill(ctx context.Context, ui core.UI, slugs []string, summary string, author string) error {
	memSvc := project.NewMemorialService(".")
	if err := memSvc.Distill(ctx, slugs, summary, author); err != nil {
		return fmt.Errorf("failed to distill memorial: %w", err)
	}

	ui.Success(fmt.Sprintf("Memorial fragments %s distilled into DISTILLED.md", strings.Join(slugs, ", ")))
	return nil
}

// HandleArchiveMemorial processes the 'archive memorial' command.
func (e *Executor) HandleArchiveMemorial(ctx context.Context, ui core.UI, slug string, ftype string, title string, content string, author string) error {
	var fragmentType project.FragmentType
	switch strings.ToLower(ftype) {
	case "action":
		fragmentType = project.FragmentAction
	case "lesson":
		fragmentType = project.FragmentLesson
	case "decision":
		fragmentType = project.FragmentDecision
	case "context":
		fragmentType = project.FragmentContext
	default:
		return fmt.Errorf("invalid fragment type: %s. Must be one of: action, lesson, decision, context", ftype)
	}

	memSvc := project.NewMemorialService(".")
	f := project.Fragment{
		Date:    time.Now(),
		Scope:   slug,
		Author:  author,
		Type:    fragmentType,
		Title:   title,
		Content: content,
	}

	if err := memSvc.Record(ctx, f); err != nil {
		return fmt.Errorf("failed to record memorial: %w", err)
	}

	ui.Success(fmt.Sprintf("Memorial fragment recorded for %s", slug))
	return nil
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

	// [REQ-3] Legacy Memorial Cleanup
	legacyPath := filepath.Join(".specforce", "docs", "memorial.md")
	if _, err := os.Stat(legacyPath); err == nil {
		_ = os.Remove(legacyPath)
		// Re-scan status to reflect the deletion
		status, _ = constitution.GetStatus(ctx, ".", registry)
	}

	// 2. Fetch Kit Instructions & 3. Fetch Config Instructions
	kitFS, err := e.GetKitFS(ui)
	if err != nil {
		return err
	}
	config := core.LoadConfig(".")
	if config.Context == nil {
		config.Context = make(map[string]string)
	}

	// Inject dynamic context for archiving
	now := time.Now()
	config.Context["CURRENT_DATE"] = now.Format("20060102")
	config.Context["CURRENT_TIME"] = now.Format("1504")
	config.Context["CURRENT_DATETIME"] = now.Format("20060102-1504")

	memSvc := project.NewMemorialService(".")
	fragments, err := memSvc.Consolidate(ctx, 10)
	if err != nil {
		fragments = "None"
	}
	config.Context["MEMORIAL_FRAGMENTS"] = fragments

	mgr := agent.NewInstructionManager(kitFS, config)
	finalInstructions, err := mgr.GetInstructions("archive")
	if err != nil {
		return fmt.Errorf("failed to get archive instructions: %w", err)
	}

	e.printArchiveInstructions(&status, finalInstructions)
	return nil
}

func (e *Executor) printArchiveInstructions(status *constitution.ConstitutionStatus, instructions string) {
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
	fmt.Println(instructions)
	fmt.Println()

	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
}
