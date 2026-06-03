package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// TimeSession represents a single period of work on a task.
type TimeSession struct {
	StartedAt   time.Time  `json:"started_at" yaml:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

// TaskTimeLog stores all sessions for a specific task.
type TaskTimeLog struct {
	Sessions []TimeSession `json:"sessions" yaml:"sessions"`
}

// RefinementMetadata tracks the state of the automated refinement loop.
type RefinementMetadata struct {
	IterationCount int       `json:"iteration_count" yaml:"iteration_count"`
	LastAuditAt    time.Time `json:"last_audit_at,omitempty" yaml:"last_audit_at,omitempty"`
	IsValid        bool      `json:"is_valid" yaml:"is_valid"`
	Errors         []string  `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// Metadata represents the core configuration of a specification.
type Metadata struct {
	Slug       string                 `json:"slug" yaml:"slug"`
	Name       string                 `json:"name" yaml:"name"`
	Type       string                 `json:"type" yaml:"type"` // "feature" | "bug"
	TimeLogs   map[string]TaskTimeLog `json:"time_logs,omitempty" yaml:"time_logs,omitempty"`
	Refinement RefinementMetadata     `json:"refinement,omitempty" yaml:"refinement,omitempty"`
}

// LoadMetadata reads the spec.yaml from the specification directory.
// It returns a default "feature" metadata if the file does not exist.
func LoadMetadata(projectRoot, slug string) (*Metadata, error) {
	metaPath := filepath.Join(projectRoot, ".specforce", "specs", slug, "spec.yaml")
	
	// #nosec G304 - metaPath is constructed from projectRoot and points to an internal spec file
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Backward compatibility: default to feature
			return &Metadata{
				Slug: slug,
				Name: slug,
				Type: "feature",
			}, nil
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta Metadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Default to feature if not specified
	if meta.Type == "" {
		meta.Type = "feature"
	}

	return &meta, nil
}

// SaveMetadata writes the metadata to spec.yaml in the specification directory.
func SaveMetadata(projectRoot, slug string, meta *Metadata) error {
	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	metaPath := filepath.Join(specDir, "spec.yaml")

	data, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return os.WriteFile(metaPath, data, 0600)
}

// StartSession appends a new session with StartedAt = time.Now().UTC().
// It ensures any previously open sessions for the task are closed.
func (m *Metadata) StartSession(taskID string) {
	if m.TimeLogs == nil {
		m.TimeLogs = make(map[string]TaskTimeLog)
	}

	log := m.TimeLogs[taskID]

	// Close any open sessions first
	for i := range log.Sessions {
		if log.Sessions[i].CompletedAt == nil {
			now := time.Now().UTC()
			log.Sessions[i].CompletedAt = &now
		}
	}

	now := time.Now().UTC()
	log.Sessions = append(log.Sessions, TimeSession{
		StartedAt: now,
	})

	m.TimeLogs[taskID] = log
}

// EndSession finds the active session (CompletedAt == nil) and sets it to UTC now.
func (m *Metadata) EndSession(taskID string) {
	if m.TimeLogs == nil {
		return
	}

	log, ok := m.TimeLogs[taskID]
	if !ok {
		return
	}

	changed := false
	for i := range log.Sessions {
		if log.Sessions[i].CompletedAt == nil {
			now := time.Now().UTC()
			log.Sessions[i].CompletedAt = &now
			changed = true
		}
	}

	if changed {
		m.TimeLogs[taskID] = log
	}
}

// GetTaskDuration calculates the cumulative duration for a task.
func (m *Metadata) GetTaskDuration(taskID string) time.Duration {
	if m.TimeLogs == nil {
		return 0
	}

	log, ok := m.TimeLogs[taskID]
	if !ok {
		return 0
	}

	var total time.Duration
	for _, s := range log.Sessions {
		if s.CompletedAt != nil {
			total += s.CompletedAt.Sub(s.StartedAt)
		} else {
			total += time.Since(s.StartedAt)
		}
	}

	return total
}
