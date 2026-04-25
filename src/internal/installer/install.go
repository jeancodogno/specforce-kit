package installer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// Install handles the global installation of the specforce binary.
func Install(ctx context.Context, ui core.UI) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	targetDir := "/usr/local/bin"
	targetPath := filepath.Join(targetDir, "specforce")

	if ui != nil {
		ui.SubTask(fmt.Sprintf("Attempting to install specforce to %s...", targetPath))
	}

	if err := checkPermissions(ui, targetDir); err != nil {
		handlePermissionError(ui, targetDir, exePath, targetPath)
		return err
	}

	if err := createSymlink(exePath, targetPath); err != nil {
		return err
	}

	if ui != nil {
		ui.SubTask(fmt.Sprintf("specforce installed successfully to %s", targetPath))
	}
	return nil
}

func checkPermissions(ui core.UI, dir string) error {
	tmpFile := filepath.Join(filepath.Clean(dir), ".specforce_write_test")
	// #nosec G304 - This is a write test for the installer in a known system directory
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to check write permissions in %s: %w", dir, errors.Join(core.ErrInstallerPermissionDenied, err))
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}
	if err := os.Remove(tmpFile); err != nil && ui != nil {
		ui.Warn(fmt.Sprintf("Failed to remove temporary file: %s", tmpFile))
	}
	return nil
}

func handlePermissionError(ui core.UI, targetDir, exePath, targetPath string) {
	msg := fmt.Sprintf("%s is not writable without elevated permissions.", targetDir)
	sudoCmd := "Please run: sudo ln -sf " + exePath + " " + targetPath

	if ui != nil {
		ui.Warn(msg)
		fmt.Println(sudoCmd)
	} else {
		fmt.Printf("[WARN] %s\n", msg)
		fmt.Println(sudoCmd)
	}
}

func createSymlink(exePath, targetPath string) error {
	if err := os.Symlink(exePath, targetPath); err != nil {
		if os.IsExist(err) {
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("failed to remove existing symlink: %w", err)
			}
			if err := os.Symlink(exePath, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	}
	return nil
}
