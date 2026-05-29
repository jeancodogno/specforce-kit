package project

import (
	"fmt"
	"os"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// MigrateLegacyAgents renames the legacy .agent/ directory to the new .agents/ standard.
func MigrateLegacyAgents(root string, ui core.UI) error {
	legacyPath, err := core.SecurePath(root, ".agent")
	if err != nil {
		return err
	}

	newPath, err := core.SecurePath(root, ".agents")
	if err != nil {
		return err
	}

	// Check if legacy directory exists
	info, err := os.Stat(legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to migrate
		}
		return fmt.Errorf("failed to stat legacy .agent directory: %w", err)
	}

	if !info.IsDir() {
		return nil
	}

	// Check if target directory already exists
	if _, err := os.Stat(newPath); err == nil {
		// If .agents already exists, we do not overwrite to prevent data loss.
		return nil
	}

	if ui != nil {
		ui.SubTask("Renaming legacy .agent/ to .agents/...")
	}

	if err := os.Rename(legacyPath, newPath); err != nil {
		return fmt.Errorf("failed to migrate .agent to .agents: %w", err)
	}

	return nil
}
