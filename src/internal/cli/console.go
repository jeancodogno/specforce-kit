package cli

import (
	"context"
	"fmt"
	"io/fs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

// RunConsole launches the interactive TUI console.
func (e *Executor) RunConsole(ctx context.Context, ui core.UI) error {
	if !tui.IsTTY() {
		return fmt.Errorf("console requires an interactive TTY terminal")
	}

	artifactsFS, err := e.GetArtifactsFS(ui)
	if err != nil {
		return err
	}

	specFS, err := fs.Sub(artifactsFS, "spec")
	if err != nil {
		return fmt.Errorf("failed to sub-filesystem for spec: %w", err)
	}

	registry, err := spec.NewRegistry(specFS)
	if err != nil {
		return fmt.Errorf("failed to initialize spec registry: %w", err)
	}

	tree, err := spec.ScanProject(ctx, ".", registry)
	if err != nil {
		return fmt.Errorf("failed to scan project state: %w", err)
	}

	model := tui.NewConsoleModel(ctx, tree, registry, ".")
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
