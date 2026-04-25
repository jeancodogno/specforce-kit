package constitution

import (
	"context"
	"os"
	"path/filepath"
)

// ArtifactStatus represents the presence and description of a single constitution document.
type ArtifactStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Exists      bool   `json:"exists"`
}

// ConstitutionStatus represents the overall completion state of the project's documentation.
type ConstitutionStatus struct {
	Artifacts []ArtifactStatus `json:"artifacts"`
	Progress  int              `json:"progress"`
	Total     int              `json:"total"`
	Found     int              `json:"found"`
}

// GetStatus checks the filesystem for the required artifacts from the registry and returns a progress summary.
func GetStatus(ctx context.Context, projectRoot string, registry *Registry) (ConstitutionStatus, error) {
	artifacts := registry.List()
	
	// Filter for standard core artifacts (exclude index and module as they are meta/dynamic)
	// Actually, index IS core. memorial and current-state are also core.
	// Only 'module' is a generic template.
	
	coreArtifacts := make([]Artifact, 0)
	for _, art := range artifacts {
		if art.Slug == "module" {
			continue
		}
		coreArtifacts = append(coreArtifacts, art)
	}

	status := ConstitutionStatus{
		Artifacts: make([]ArtifactStatus, len(coreArtifacts)),
		Total:     len(coreArtifacts),
	}

	for i, art := range coreArtifacts {
		if err := ctx.Err(); err != nil {
			return status, err
		}
		fullPath := filepath.Join(projectRoot, art.Path)
		exists := false
		if _, err := os.Stat(fullPath); err == nil {
			exists = true
			status.Found++
		}

		status.Artifacts[i] = ArtifactStatus{
			Name:        art.Name,
			Description: art.Description,
			Path:        art.Path,
			Exists:      exists,
		}
	}

	if status.Total > 0 {
		status.Progress = (status.Found * 100) / status.Total
	}

	return status, nil
}
