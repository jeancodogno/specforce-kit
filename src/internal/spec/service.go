package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	auditor        Auditor
	projectRoot    string
}

// NewService creates a new instance of the spec service.
func NewService(registry *Registry, configProvider ConfigProvider) *Service {
	return &Service{
		registry:       registry,
		configProvider: configProvider,
		auditor:        NewAuditor("."), // Default project root
		projectRoot:    ".",
	}
}

// SetProjectRoot updates the project root for the service and its components.
func (s *Service) SetProjectRoot(root string) {
	s.projectRoot = root
	s.auditor = NewAuditor(root)
}

// ResizeSpec updates the size classification of an existing specification.
func (s *Service) ResizeSpec(ctx context.Context, slug string, size SpecSize) (*Metadata, error) {
	if !ValidateSize(size) {
		return nil, fmt.Errorf("invalid spec size: %s. Supported: small, medium, large, complex", size)
	}

	slug = ResolveSlug(s.projectRoot, slug)
	specDir := filepath.Join(s.projectRoot, ".specforce", "specs", slug)
	if fi, err := os.Stat(specDir); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("specification %q not found", slug)
	}

	meta, err := LoadMetadata(s.projectRoot, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec metadata: %w", err)
	}

	meta.Size = size
	if err := SaveMetadata(s.projectRoot, slug, meta); err != nil {
		return nil, fmt.Errorf("failed to save spec metadata: %w", err)
	}

	return meta, nil
}

// RefineSpec manages the automated refinement loop pass.
func (s *Service) RefineSpec(ctx context.Context, slug string, ui core.UI) error {
	slug = ResolveSlug(s.projectRoot, slug)
	meta, err := LoadMetadata(s.projectRoot, slug)
	if err != nil {
		return err
	}

	const maxIterations = 3
	if meta.Refinement.IterationCount >= maxIterations {
		if ui != nil {
			ui.Warn("Maximum refinement iterations reached.")
		}
		return nil
	}

	iter := meta.Refinement.IterationCount + 1
	if ui != nil {
		ui.SubTask(fmt.Sprintf("Refinement Pass %d/%d...", iter, maxIterations))
	}

	// 1. Run Audit
	errors, err := s.auditor.Audit(ctx, slug)
	if err != nil {
		return fmt.Errorf("audit failed: %w", err)
	}

	if len(errors) == 0 {
		meta.Refinement.IsValid = true
		meta.Refinement.Errors = nil
		return SaveMetadata(s.projectRoot, slug, meta)
	}

	// 2. Persist Errors
	meta.Refinement.Errors = nil
	for _, e := range errors {
		msg := fmt.Sprintf("[%s] %s: %s", strings.ToUpper(e.Artifact), e.Code, e.Message)
		meta.Refinement.Errors = append(meta.Refinement.Errors, msg)
	}
	meta.Refinement.IsValid = false

	// 3. Increment Counter
	meta.Refinement.IterationCount = iter
	meta.Refinement.LastAuditAt = time.Now().UTC()

	if err := SaveMetadata(s.projectRoot, slug, meta); err != nil {
		return err
	}

	if ui != nil {
		ui.Warn(fmt.Sprintf("Found %d coherence errors. Re-invoking agents...", len(errors)))
	}

	return nil
}

// BuildCorrectionPayload creates a surgical prompt for an agent to fix a specific coherence error.
func (s *Service) BuildCorrectionPayload(ctx context.Context, slug string, ce CoherenceError) (string, error) {
	slug = ResolveSlug(s.projectRoot, slug)
	
	// 1. Read the target artifact's current content
	targetPath := filepath.Join(s.projectRoot, ".specforce", "specs", slug, ce.Artifact+".md")
	// #nosec G304 - internal spec file
	targetContent, err := os.ReadFile(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to read target artifact %s: %w", ce.Artifact, err)
	}

	// 2. Read requirements for source of truth context if needed
	reqPath := filepath.Join(s.projectRoot, ".specforce", "specs", slug, "requirements.md")
	// #nosec G304 - internal spec file
	reqData, err := os.ReadFile(reqPath)
	if err != nil {
		return "", fmt.Errorf("failed to read requirements.md: %w", err)
	}

	// 3. Extract relevant requirement snippet if context is a US tag
	sourceContext := ce.Context
	if strings.HasPrefix(ce.Context, "## [US-") {
		lines := strings.Split(string(reqData), "\n")
		var snippet []string
		found := false
		for _, line := range lines {
			if strings.HasPrefix(line, ce.Context) {
				found = true
			} else if found && strings.HasPrefix(line, "## ") {
				break
			}
			if found {
				snippet = append(snippet, line)
			}
		}
		if len(snippet) > 0 {
			sourceContext = strings.Join(snippet, "\n")
		}
	}

	// 4. Build the payload
	payload := fmt.Sprintf(`### SURGICAL CORRECTION REQUIRED
Artifact: %s.md
Error Code: %s
Message: %s

### SOURCE OF TRUTH (CONTEXT)
%s

### CURRENT TARGET CONTENT
%s
`, ce.Artifact, ce.Code, ce.Message, sourceContext, string(targetContent))

	return payload, nil
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
			rules := s.resolveInstructions(conf, name)
			if len(rules) > 0 {
				custom := "## Project Specific Instructions\n- " + strings.Join(rules, "\n- ") + "\n\n"
				art.Instruction = custom + art.Instruction
			}
		}
	}

	return &art, nil
}

// GetImplementationStatus retrieves the status and details for implementing a specific feature.
func (s *Service) GetImplementationStatus(ctx context.Context, projectRoot, slug string) (*ImplementationReport, error) {
	slug = ResolveSlug(projectRoot, slug)

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
func (s *Service) UpdateTaskStatus(ctx context.Context, projectRoot, slug string, taskIDs []string, status string) error {
	slug = ResolveSlug(projectRoot, slug)

	// 1. Only run hooks if status is "finished"
	if strings.ToLower(status) != "finished" {
		return updateTaskStatusesFile(projectRoot, slug, taskIDs, status)
	}

	// 2. Load Config
	var config *core.ProjectConfig
	if s.configProvider != nil {
		if c, err := s.configProvider.GetConfig(ctx); err == nil {
			config = c
		}
	}

	// 3. Collect Hooks
	hooks := s.collectHooksForBatch(ctx, projectRoot, slug, taskIDs, config)

	// 4. Execute Hooks
	if len(hooks) > 0 {
		if _, err := core.ExecuteHooks(ctx, hooks); err != nil {
			return fmt.Errorf("hooks failed: %w", err)
		}
	}

	// 5. Update File
	return updateTaskStatusesFile(projectRoot, slug, taskIDs, status)
}

func (s *Service) collectHooksForBatch(ctx context.Context, projectRoot, slug string, taskIDs []string, config *core.ProjectConfig) []string {
	if config == nil || len(taskIDs) == 0 {
		return nil
	}

	var rawHooks []string
	rawHooks = append(rawHooks, config.Hooks.OnTaskFinished...)

	report, err := ParseTasks(ctx, projectRoot, slug)
	if err == nil && report != nil {
		for _, taskID := range taskIDs {
			isLastInPhase, isLastInSpec := s.checkTaskPosition(report, taskID)
			if isLastInPhase {
				rawHooks = append(rawHooks, config.Hooks.OnPhaseFinished...)
			}
			if isLastInSpec {
				rawHooks = append(rawHooks, config.Hooks.OnAllTasksFinished...)
			}
		}
	}

	var deduplicated []string
	seen := make(map[string]bool, len(rawHooks))
	for _, hook := range rawHooks {
		if hook != "" && !seen[hook] {
			seen[hook] = true
			deduplicated = append(deduplicated, hook)
		}
	}

	return deduplicated
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

func (s *Service) inferBaseType(name string) string {
	baseTypes := []string{"requirements", "design", "tasks", "implementation", "archive"}
	bestMatch := ""
	bestIndex := -1

	for _, bt := range baseTypes {
		idx := strings.LastIndex(name, bt)
		if idx > bestIndex {
			bestIndex = idx
			bestMatch = bt
		}
	}

	return bestMatch
}

func (s *Service) resolveInstructions(conf *core.ProjectConfig, name string) []string {
	var combined []string
	seen := make(map[string]bool)

	// 1. Get Generic Instructions
	baseType := s.inferBaseType(name)
	if baseType != "" {
		if generic, ok := conf.Instructions[baseType]; ok {
			for _, rule := range generic {
				if !seen[rule] {
					combined = append(combined, rule)
					seen[rule] = true
				}
			}
		}
	}

	// 2. Get Specific Instructions (if name is different from baseType)
	if name != baseType {
		if specific, ok := conf.Instructions[name]; ok {
			for _, rule := range specific {
				if !seen[rule] {
					combined = append(combined, rule)
					seen[rule] = true
				}
			}
		}
	}

	return combined
}
