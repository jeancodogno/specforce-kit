package project

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
)

// DetectExistingAgents checks for existing agent directories in the current project using the registry.
func DetectExistingAgents(ctx context.Context, root string, reg *agent.Registry) []string {
	var existing []string
	if reg == nil {
		return existing
	}

	for _, agent := range reg.GetAgents() {
		if err := ctx.Err(); err != nil {
			return existing
		}
		path := filepath.Join(root, agent.DirName)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			existing = append(existing, agent.ID)
		}
	}

	return existing
}

// IsInitialized checks if the project has been initialized by looking for the .specforce directory.
func IsInitialized(root string) bool {
	info, err := os.Stat(filepath.Join(root, ".specforce"))
	return err == nil && info.IsDir()
}
