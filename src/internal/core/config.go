package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// HooksConfig contains the commands to run for specific implementation events.
type HooksConfig struct {
	OnTaskFinished     []string `yaml:"on_task_finished"`
	OnPhaseFinished    []string `yaml:"on_phase_finished"`
	OnAllTasksFinished []string `yaml:"on_all_tasks_finished"`
}

// ProjectConfig contains the project-specific configuration.
type ProjectConfig struct {
	Instructions map[string][]string `yaml:"instructions"`
	Hooks        HooksConfig         `yaml:"hooks"`
}

// DefaultConfigContent is the default content for the .specforce/config.yaml file.
const DefaultConfigContent = `instructions:
  # Example: Project-wide instructions for all requirements artifacts
  # requirements:
  #   - "Always use BDD GIVEN/WHEN/THEN syntax"
  #   - "Ensure accessibility is mentioned for UI components"

  # Example: Project-wide instructions for all design artifacts
  # design:
  #   - "Use Mermaid.js for architecture diagrams"
  #   - "Include a detailed component inventory"

  # Example: Project-wide instructions for all tasks artifacts
  # tasks:
  #   - "Each task must have a clear verification step"
  #   - "Group tasks by implementation phases"

  # Example: Project-wide instructions for the implementation phase
  # implementation:
  #   - "Always run 'go fmt' before finishing a task"
  #   - "Use explicit type casts instead of interfaces when possible"

  # Example: Project-wide instructions for the archive phase
  # archive:
  #   - "Always update the project memorial with lessons learned"
  #   - "Ensure all temporary artifacts are cleaned up"

hooks:
  # Example: Run linting and tests automatically when finishing tasks
  # on_task_finished:
  #   - "golangci-lint run"
  # on_phase_finished:
  #   - "go test ./src/internal/..."
  # on_all_tasks_finished:
  #   - "go test ./..."
`

// EnsureConfigExists checks if the configuration file exists, and creates it with default content if it doesn't.
func EnsureConfigExists(root string) error {
	specforceDir, err := SecurePath(root, ".specforce")
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	configPath, err := SecurePath(root, filepath.Join(".specforce", "config.yaml"))
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // Already exists, do nothing
	}

	// Create .specforce directory if it doesn't exist
	if err := os.MkdirAll(specforceDir, 0750); err != nil {
		return fmt.Errorf("failed to create .specforce directory: %w", err)
	}

	// Write default content
	// #nosec G306 - Config path is secured by SecurePath
	if err := os.WriteFile(configPath, []byte(DefaultConfigContent), 0600); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}

// LoadConfig reads the .specforce/config.yaml file from the project root.
func LoadConfig(root string) *ProjectConfig {
	config := &ProjectConfig{
		Instructions: make(map[string][]string),
	}

	configPath, err := SecurePath(root, filepath.Join(".specforce", "config.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Security: %v\n", err)
		return config
	}
	
	// #nosec G304 - Path is secured by SecurePath
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Warning: Failed to read config file at %s: %v\n", configPath, err)
		}
		return config
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Malformed config file at %s: %v\n", configPath, err)
		// Ensure we return an empty but non-nil Instructions map
		config.Instructions = make(map[string][]string)
		return config
	}

	if config.Instructions == nil {
		config.Instructions = make(map[string][]string)
	}

	return config
}
