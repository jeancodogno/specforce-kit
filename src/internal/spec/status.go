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
	ValidationGuide  string   `json:"validation_guide,omitempty"`
}

// SpecStatus represents the overall completion state of a specific feature spec.
type SpecStatus struct {
	Slug         string           `json:"slug"`
	Type         string           `json:"type"`
	Artifacts    []ArtifactStatus `json:"artifacts"`
	Progress     int              `json:"progress"`
	Total        int              `json:"total"`
	Found        int              `json:"found"`
	IsValid      bool             `json:"is_valid"`
	ContextFiles []string         `json:"context_files,omitempty"`
}

// GetStatus checks the filesystem for the required artifacts from the registry and returns a progress summary.
func GetStatus(ctx context.Context, projectRoot string, slug string, registry *Registry) (SpecStatus, error) {
	slug = ResolveSlug(projectRoot, slug)

	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return SpecStatus{}, fmt.Errorf("failed to load metadata for %s: %w", slug, err)
	}

	specDir, err := resolveSpecDir(projectRoot, slug)
	if err != nil {
		return SpecStatus{}, err
	}

	artifacts := registry.ListForType(meta.Type)
	status := SpecStatus{
		Slug:      slug,
		Type:      meta.Type,
		Artifacts: make([]ArtifactStatus, 0, len(artifacts)),
		Total:     len(artifacts),
		IsValid:   true,
	}

	detectProposal(projectRoot, specDir, &status)

	existsMap, foundCount, err := scanArtifactExistence(ctx, specDir, artifacts)
	if err != nil {
		return SpecStatus{}, err
	}

	if err := processAllArtifacts(ctx, projectRoot, slug, meta.Type, artifacts, existsMap, registry, foundCount == len(artifacts), &status); err != nil {
		return status, err
	}

	if status.Total > 0 {
		status.Progress = (status.Found * 100) / status.Total
	}

	return status, nil
}

func resolveSpecDir(projectRoot, slug string) (string, error) {
	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if _, err := os.Stat(specDir); os.IsNotExist(err) {
		specDir = filepath.Join(projectRoot, ".specforce", "archive", slug)
		if _, err := os.Stat(specDir); os.IsNotExist(err) {
			return "", fmt.Errorf("feature directory not found: %s", slug)
		}
	}
	return specDir, nil
}

func detectProposal(projectRoot, specDir string, status *SpecStatus) {
	proposalPath := filepath.Join(specDir, "proposal.md")
	if _, err := os.Stat(proposalPath); err == nil {
		if rel, err := filepath.Rel(projectRoot, proposalPath); err == nil {
			status.ContextFiles = append(status.ContextFiles, rel)
		}
	}
}

func processAllArtifacts(ctx context.Context, projectRoot, slug, specType string, artifacts []Artifact, existsMap map[string]bool, registry *Registry, shouldValidate bool, status *SpecStatus) error {
	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return err
		}

		artStatus, err := processArtifactStatus(ctx, projectRoot, slug, specType, art, existsMap, registry, shouldValidate)
		if err != nil {
			return err
		}

		if len(artStatus.ValidationErrors) > 0 {
			status.IsValid = false
		}
		if artStatus.Exists {
			status.Found++
		}

		status.Artifacts = append(status.Artifacts, artStatus)
	}
	return nil
}

func scanArtifactExistence(ctx context.Context, specDir string, artifacts []Artifact) (map[string]bool, int, error) {
	existsMap := make(map[string]bool)
	foundCount := 0
	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		fullPath := filepath.Join(specDir, art.Name+".md")
		if _, err := os.Stat(fullPath); err == nil {
			existsMap[art.Name] = true
			foundCount++
		}
	}
	return existsMap, foundCount, nil
}

func processArtifactStatus(ctx context.Context, projectRoot, slug, specType string, art Artifact, existsMap map[string]bool, registry *Registry, validate bool) (ArtifactStatus, error) {
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
	var validationGuide string
	if art.Name == "tasks" && exists && validate {
		var err error
		validationErrors, err = ValidateTasks(ctx, projectRoot, slug)
		if err != nil {
			return ArtifactStatus{}, fmt.Errorf("failed to validate tasks.md: %w", err)
		}

		if len(validationErrors) > 0 {
			validationGuide = "### Phase 1: Example Phase\n- [ ] T1.1: Example Task\n**Target:** `path/to/file.go`\n**Context:** [US-1]\n**Action Steps:**\n- Step 1\n**Acceptance Check:**\n- run test"
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
		ValidationGuide:  validationGuide,
	}, nil
}
