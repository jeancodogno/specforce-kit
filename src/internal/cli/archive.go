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
	default:
		fmt.Printf("Unknown archive command: %s\n", subCommand)
		fmt.Println("Available commands: instructions, memorial")
		return nil
	}
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
	config.Context["MEMORIAL_FRAGMENTS"] = e.getMemorialList()

	mgr := agent.NewInstructionManager(kitFS, config)
	finalInstructions, err := mgr.GetInstructions("archive")
	if err != nil {
		return fmt.Errorf("failed to get archive instructions: %w", err)
	}

	e.printArchiveInstructions(&status, finalInstructions)
	return nil
}

func (e *Executor) getMemorialList() string {
	memorialFiles := []string{}
	entries, _ := os.ReadDir(filepath.Join(".specforce", "memorial"))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			memorialFiles = append(memorialFiles, entry.Name())
		}
	}
	if len(memorialFiles) == 0 {
		return "None"
	}
	return strings.Join(memorialFiles, ", ")
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
