package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ArtifactStatus represents the presence and description of a single spec document.
type ArtifactStatus struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Path             string   `json:"path"`
	Exists           bool     `json:"exists"`
	Blocked          bool     `json:"blocked"`
	Dependency       string   `json:"dependency"`
	ValidationErrors []string `json:"validation_errors,omitempty"`
}

// SpecStatus represents the overall completion state of a specific feature spec.
type SpecStatus struct {
	Slug      string           `json:"slug"`
	Type      string           `json:"type"`
	Artifacts []ArtifactStatus `json:"artifacts"`
	Progress  int              `json:"progress"`
	Total     int              `json:"total"`
	Found     int              `json:"found"`
	IsValid   bool             `json:"is_valid"`
}

// GetStatus checks the filesystem for the required artifacts from the registry and returns a progress summary.
func GetStatus(ctx context.Context, projectRoot string, slug string, registry *Registry) (SpecStatus, error) {
	// Load Metadata to determine spec type
	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return SpecStatus{}, fmt.Errorf("failed to load metadata for %s: %w", slug, err)
	}

	artifacts := registry.ListForType(meta.Type)

	status := SpecStatus{
		Slug:      slug,
		Type:      meta.Type,
		Artifacts: make([]ArtifactStatus, 0, len(artifacts)),
		Total:     len(artifacts),
		IsValid:   true,
	}

	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if _, err := os.Stat(specDir); os.IsNotExist(err) {
		return SpecStatus{}, fmt.Errorf("feature directory not found: %s", specDir)
	}

	existsMap, err := scanArtifactExistence(ctx, specDir, artifacts)
	if err != nil {
		return SpecStatus{}, err
	}

	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return status, err
		}

		artStatus, err := processArtifactStatus(ctx, projectRoot, slug, meta.Type, art, existsMap, registry)
		if err != nil {
			return status, err
		}

		if len(artStatus.ValidationErrors) > 0 {
			status.IsValid = false
		}
		if artStatus.Exists {
			status.Found++
		}

		status.Artifacts = append(status.Artifacts, artStatus)
	}

	if status.Total > 0 {
		status.Progress = (status.Found * 100) / status.Total
	}

	return status, nil
}

func scanArtifactExistence(ctx context.Context, specDir string, artifacts []Artifact) (map[string]bool, error) {
	existsMap := make(map[string]bool)
	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		fullPath := filepath.Join(specDir, art.Name+".md")
		if _, err := os.Stat(fullPath); err == nil {
			existsMap[art.Name] = true
		}
	}
	return existsMap, nil
}

func processArtifactStatus(ctx context.Context, projectRoot, slug, specType string, art Artifact, existsMap map[string]bool, registry *Registry) (ArtifactStatus, error) {
	fileName := art.Name + ".md"
	relPath := filepath.Join(".specforce", "specs", slug, fileName)
	exists := existsMap[art.Name]

	prefixedName := fmt.Sprintf("%s-%s", specType, art.Name)

	blocked := false
	if art.Dependency != "" {
		if _, depInRegistry := registry.Get(art.Dependency); !depInRegistry {
			blocked = true
		} else if !existsMap[art.Dependency] {
			blocked = true
		}
	}

	var validationErrors []string
	if art.Name == "tasks" && exists {
		var err error
		validationErrors, err = ValidateTasks(ctx, projectRoot, slug)
		if err != nil {
			return ArtifactStatus{}, fmt.Errorf("failed to validate tasks.md: %w", err)
		}
	}

	return ArtifactStatus{
		Name:             prefixedName,
		Description:      art.Description,
		Path:             relPath,
		Exists:           exists,
		Blocked:          blocked,
		Dependency:       art.Dependency,
		ValidationErrors: validationErrors,
	}, nil
}
