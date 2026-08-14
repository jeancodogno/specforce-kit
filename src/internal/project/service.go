package project

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

// InitConfig contains project root, selected agents, and options for project initialization.
type InitConfig struct {
	ProjectRoot    string
	SelectedAgents []string
}

// Service orchestrates project-related domain logic.
type Service struct {
	kitFS       fs.FS
	artifactsFS fs.FS
	projectRoot string
	config      *core.ProjectConfig
	registry    *agent.Registry
}

// NewService creates a new instance of the project service.
func NewService(kitFS, artifactsFS fs.FS, projectRoot string) *Service {
	return &Service{
		kitFS:       kitFS,
		artifactsFS: artifactsFS,
		projectRoot: projectRoot,
	}
}

// GetConfig returns the project-specific configuration.
func (s *Service) GetConfig(ctx context.Context) (*core.ProjectConfig, error) {
	if s.config == nil {
		s.config = core.LoadConfig(s.projectRoot)
	}
	return s.config, nil
}

// InitializeProject orchestrates the project initialization flow.
func (s *Service) InitializeProject(ctx context.Context, ui core.UI, config InitConfig) error {
	if ui != nil {
		ui.LogSubTask("DEPLOYING INFRASTRUCTURE")
	}

	if err := BootstrapProject(ctx, config.ProjectRoot, s.kitFS, s.artifactsFS, ui); err != nil {
		if !errors.Is(err, core.ErrProjectAlreadyInitialized) {
			return err
		}
	}

	if err := EnsureAgentsMD(config.ProjectRoot, ui, config.SelectedAgents); err != nil {
		return fmt.Errorf("failed to finalize project with AGENTS.md: %w", err)
	}

	if ui != nil {
		ui.LogSubTask("SYNCING AGENT ARTIFACTS")
	}

	for _, a := range config.SelectedAgents {
		if err := agent.AdaptArtifacts(ctx, config.ProjectRoot, s.kitFS, a, ui, installer.Options{}); err != nil {
			return fmt.Errorf("failed to adapt artifacts for %s: %w", a, err)
		}
	}

	return nil
}

// DecommissionAgents removes agent directories that are no longer selected.
func (s *Service) DecommissionAgents(ctx context.Context, ui core.UI, agentsToRemove []string) error {
	if len(agentsToRemove) == 0 {
		return nil
	}

	if ui != nil {
		ui.LogSubTask("DECOMMISSIONING AGENTS")
	}

	for _, agentID := range agentsToRemove {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Use the agent registry to find the directory name
		agentMeta, ok := s.GetAgentMetadata(agentID)
		if !ok {
			if ui != nil {
				ui.Warn(fmt.Sprintf("Skip decommission: agent metadata not found for %s", agentID))
			}
			continue
		}

		path := filepath.Join(s.projectRoot, agentMeta.DirName)
		if _, err := os.Stat(path); err == nil {
			if err := os.RemoveAll(path); err != nil {
				if ui != nil {
					ui.Warn(fmt.Sprintf("Failed to remove agent directory %s: %v", path, err))
				}
			} else if ui != nil {
				// Format: ↳ .qwen ................................... DELETED
				dotCount := 35 - len(agentMeta.DirName)
				if dotCount < 1 {
					dotCount = 1
				}
				dots := strings.Repeat(".", dotCount)
				ui.LogSubTask(fmt.Sprintf("%s %s DELETED", agentMeta.DirName, dots))
			}
		}
	}

	return nil
}

func (s *Service) GetAgentMetadata(agentID string) (agent.AgentMetadata, bool) {
	if s.registry == nil {
		s.registry = &agent.Registry{}
		if err := s.registry.Initialize(s.kitFS, s.projectRoot); err != nil {
			return agent.AgentMetadata{}, false
		}
	}
	return s.registry.GetAgent(agentID)
}

// UpdateTools refreshes agent tools and instructions while preserving the .specforce/ directory.
func (s *Service) UpdateTools(ctx context.Context, ui core.UI, selectedAgents []string) error {
	if ui != nil {
		ui.LogSubTask("UPDATING AGENT ARTIFACTS")
	}

	if err := MigrateLegacyAgents(s.projectRoot, ui); err != nil {
		if ui != nil {
			ui.Warn(fmt.Sprintf("Migration failed: %v", err))
		}
	}

	opts := installer.Options{ToolsOnly: true}

	// Iterate through selected tools and update them
	for _, agentID := range selectedAgents {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := agent.AdaptArtifacts(ctx, s.projectRoot, s.kitFS, agentID, ui, opts); err != nil {
			// We log but continue if one agent fails to update
			if ui != nil {
				ui.Warn(fmt.Sprintf("Failed to update tools for %s: %v", agentID, err))
			}
		}
	}

	if err := EnsureAgentsMD(s.projectRoot, ui, selectedAgents); err != nil {
		return fmt.Errorf("failed to update AGENTS.md: %w", err)
	}

	if err := CleanupLegacySymlinks(s.projectRoot); err != nil {
		if ui != nil {
			ui.Warn(fmt.Sprintf("Failed to cleanup legacy symlinks: %v", err))
		}
	}

	return nil
}
