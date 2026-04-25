package spec

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// ConfigProvider defines the contract for project configuration access.
type ConfigProvider interface {
	GetConfig(ctx context.Context) (*core.ProjectConfig, error)
}

// Service orchestrates specification management logic.
type Service struct {
	registry       *Registry
	configProvider ConfigProvider
}

// NewService creates a new instance of the spec service.
func NewService(registry *Registry, configProvider ConfigProvider) *Service {
	return &Service{
		registry:       registry,
		configProvider: configProvider,
	}
}

// GetArtifact retrieves a specific artifact by name, injecting custom project instructions if present.
func (s *Service) GetArtifact(ctx context.Context, name string) (*Artifact, error) {
	art, ok := s.registry.Get(name)
	if !ok {
		return nil, fmt.Errorf("artifact %q not found", name)
	}

	if s.configProvider != nil {
		conf, err := s.configProvider.GetConfig(ctx)
		if err == nil && conf != nil {
			if rules, ok := conf.Instructions[name]; ok && len(rules) > 0 {
				custom := "\n\n## Project Specific Instructions\n- " + strings.Join(rules, "\n- ")
				art.Instruction += custom
			}
		}
	}

	return &art, nil
}

// GetImplementationStatus retrieves the status and details for implementing a specific feature.
func (s *Service) GetImplementationStatus(ctx context.Context, projectRoot, slug string) (*ImplementationReport, error) {
	// 1. Check artifacts
	ok, missing := CheckTriadArtifacts(projectRoot, slug)

	// 2. Parse tasks
	report, err := ParseTasks(ctx, projectRoot, slug)
	if err != nil {
		return nil, err
	}

	report.MissingArtifacts = missing
	if !ok {
		report.Status = "blocked"
	}

	// 3. Get context files
	contextFiles, err := GetContextFiles(projectRoot, slug)
	if err == nil {
		report.ContextFiles = contextFiles
	}

	// 4. Inject instructions
	if s.configProvider != nil {
		conf, err := s.configProvider.GetConfig(ctx)
		if err == nil && conf != nil {
			if rules, ok := conf.Instructions["implementation"]; ok {
				report.Instructions = rules
			}
		}
	}

	return report, nil
}

// GetStatus retrieves the status of a specific feature spec.
func (s *Service) GetStatus(ctx context.Context, projectRoot string, slug string) (SpecStatus, error) {
	return GetStatus(ctx, projectRoot, slug, s.registry)
}

// UpdateTaskStatus handles task status updates with event hooks.
func (s *Service) UpdateTaskStatus(ctx context.Context, projectRoot, slug, taskID, status string) error {
	// 1. Only run hooks if status is "finished"
	if strings.ToLower(status) != "finished" {
		return updateTaskStatusFile(projectRoot, slug, taskID, status)
	}

	// 2. Load Config
	var config *core.ProjectConfig
	if s.configProvider != nil {
		if c, err := s.configProvider.GetConfig(ctx); err == nil {
			config = c
		}
	}

	// 3. Collect Hooks
	hooks := s.collectHooksForTask(ctx, projectRoot, slug, taskID, config)

	// 4. Execute Hooks
	if len(hooks) > 0 {
		if _, err := core.ExecuteHooks(ctx, hooks); err != nil {
			return fmt.Errorf("hooks failed: %w", err)
		}
	}

	// 5. Update File
	return updateTaskStatusFile(projectRoot, slug, taskID, status)
}

func (s *Service) collectHooksForTask(ctx context.Context, projectRoot, slug, taskID string, config *core.ProjectConfig) []string {
	var hooks []string
	if config == nil {
		return hooks
	}

	// Always add task finished hook
	hooks = append(hooks, config.Hooks.OnTaskFinished...)

	// Use ParseTasks to determine if this is the last task in a phase or spec
	report, err := ParseTasks(ctx, projectRoot, slug)
	if err != nil {
		return hooks
	}

	isLastInPhase, isLastInSpec := s.checkTaskPosition(report, taskID)

	if isLastInPhase {
		hooks = append(hooks, config.Hooks.OnPhaseFinished...)
	}
	if isLastInSpec {
		hooks = append(hooks, config.Hooks.OnAllTasksFinished...)
	}

	return hooks
}

func (s *Service) checkTaskPosition(report *ImplementationReport, taskID string) (bool, bool) {
	for pi, p := range report.Phases {
		for ti, t := range p.Tasks {
			if t.ID == taskID {
				isLastInPhase := (ti == len(p.Tasks)-1)
				isLastInSpec := (pi == len(report.Phases)-1 && ti == len(p.Tasks)-1)
				return isLastInPhase, isLastInSpec
			}
		}
	}
	return false, false
}
