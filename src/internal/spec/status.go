package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ArtifactStatus represents the presence and description of a single spec document.
type ArtifactStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Exists      bool   `json:"exists"`
	Blocked     bool   `json:"blocked"`
	Dependency  string `json:"dependency"`
}

// SpecStatus represents the overall completion state of a specific feature spec.
type SpecStatus struct {
	Slug      string           `json:"slug"`
	Artifacts []ArtifactStatus `json:"artifacts"`
	Progress  int              `json:"progress"`
	Total     int              `json:"total"`
	Found     int              `json:"found"`
}

// GetStatus checks the filesystem for the required artifacts from the registry and returns a progress summary.
func GetStatus(ctx context.Context, projectRoot string, slug string, registry *Registry) (SpecStatus, error) {
	artifacts := registry.List()

	status := SpecStatus{
		Slug:      slug,
		Artifacts: make([]ArtifactStatus, 0, len(artifacts)),
		Total:     len(artifacts),
	}

	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if _, err := os.Stat(specDir); os.IsNotExist(err) {
		return SpecStatus{}, fmt.Errorf("feature directory not found: %s", specDir)
	}

	// Pre-scan for existing files to resolve dependencies
	existsMap := make(map[string]bool)
	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return SpecStatus{}, err
		}
		fullPath := filepath.Join(specDir, art.Name+".md")
		if _, err := os.Stat(fullPath); err == nil {
			existsMap[art.Name] = true
		}
	}

	for _, art := range artifacts {
		if err := ctx.Err(); err != nil {
			return status, err
		}
		// Artifact paths within the spec directory
		fileName := art.Name + ".md"
		relPath := filepath.Join(".specforce", "specs", slug, fileName)

		exists := existsMap[art.Name]
		if exists {
			status.Found++
		}

		blocked := false
		if art.Dependency != "" {
			// Check if dependency exists in registry first (as per REQ-1 scenario 4)
			if _, depInRegistry := registry.Get(art.Dependency); !depInRegistry {
				blocked = true
			} else if !existsMap[art.Dependency] {
				blocked = true
			}
		}

		status.Artifacts = append(status.Artifacts, ArtifactStatus{
			Name:        art.Name,
			Description: art.Description,
			Path:        relPath,
			Exists:      exists,
			Blocked:     blocked,
			Dependency:  art.Dependency,
		})
	}

	if status.Total > 0 {
		status.Progress = (status.Found * 100) / status.Total
	}

	return status, nil
}
