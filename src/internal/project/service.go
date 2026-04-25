package project

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

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
	if err := BootstrapProject(ctx, config.ProjectRoot, s.kitFS, s.artifactsFS, ui); err != nil {
		if !errors.Is(err, core.ErrProjectAlreadyInitialized) {
			return err
		}
	}

	for _, a := range config.SelectedAgents {
		if err := agent.AdaptArtifacts(ctx, config.ProjectRoot, s.kitFS, a, ui, installer.Options{}); err != nil {
			return fmt.Errorf("failed to adapt artifacts for %s: %w", a, err)
		}
	}

	return nil
}

// UpdateTools refreshes agent tools and instructions while preserving the .specforce/ directory.
func (s *Service) UpdateTools(ctx context.Context, ui core.UI, selectedAgents []string) error {
	if ui != nil {
		ui.SubTask("Updating agent tools and instructions...")
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

	if ui != nil {
		ui.Success("Agent tools and instructions updated successfully.")
	}

	if err := EnsureAgentsMD(s.projectRoot, ui); err != nil {
		return fmt.Errorf("failed to update AGENTS.md: %w", err)
	}

	return nil
}
