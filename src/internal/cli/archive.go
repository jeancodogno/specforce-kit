package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
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
	default:
		fmt.Printf("Unknown archive command: %s\n", subCommand)
		fmt.Println("Available commands: instructions")
		return nil
	}
}

// HandleArchiveInstructions processes the 'archive instructions' command.
func (e *Executor) HandleArchiveInstructions(ctx context.Context, ui core.UI) error {
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

	mgr := agent.NewInstructionManager(kitFS, config)
	finalInstructions, err := mgr.GetInstructions("archive")
	if err != nil {
		return fmt.Errorf("failed to get archive instructions: %w", err)
	}

	e.printArchiveInstructions(finalInstructions)
	return nil
}

func (e *Executor) printArchiveInstructions(instructions string) {
	fmt.Println(instructions)
	fmt.Println()

	if tui.IsTTY() {
		tui.PrintFooter(e.Version)
	}
}
