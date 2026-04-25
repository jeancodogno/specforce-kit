package project

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// BootstrapProject creates the basic Specforce directory structure and initializes standard templates.
func BootstrapProject(ctx context.Context, root string, kitFS fs.FS, artifactsFS fs.FS, ui core.UI) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(root, ".specforce")); err == nil {
		return fmt.Errorf("cannot initialize project: %w", core.ErrProjectAlreadyInitialized)
	}

	dirs := []string{
		".specforce/docs",
		".specforce/specs",
		".specforce/archive",
	}

	if ui != nil {
		ui.StartSpinner("Creating directories...")
	}
	if err := createDirectories(root, dirs); err != nil {
		return err
	}

	if ui != nil {
		ui.StopSpinner()
		ui.Success("Specforce directory structure initialized successfully.")
	}

	if err := EnsureAgentsMD(root, ui); err != nil {
		return fmt.Errorf("failed to ensure AGENTS.md: %w", err)
	}

	return nil
}

func createDirectories(root string, dirs []string) error {
	for _, dir := range dirs {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}

		// Add .gitkeep to ensure directories are tracked if empty
		gitkeep := filepath.Join(path, ".gitkeep")
		if _, err := os.Stat(gitkeep); os.IsNotExist(err) {
			if err := os.WriteFile(gitkeep, []byte(""), 0600); err != nil {
				return fmt.Errorf("failed to create .gitkeep in %s: %w", path, err)
			}
		}
	}
	return nil
}
