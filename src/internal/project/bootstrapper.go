package project

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

	for _, dir := range dirs {
		if err := ctx.Err(); err != nil {
			return err
		}

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

		if ui != nil {
			// Format: ↳ .specforce/docs ................................. OK
			dotCount := 40 - len(dir)
			if dotCount < 1 {
				dotCount = 1
			}
			dots := strings.Repeat(".", dotCount)
			ui.LogSubTask(fmt.Sprintf("%s %s OK", dir, dots))
		}
	}

	return nil
}
