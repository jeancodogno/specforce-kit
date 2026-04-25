package spec

import (
	"context"
	"os"
	"path/filepath"
)

// SpecInfo contains basic information about a specification.
type SpecInfo struct {
	Slug string `json:"slug"`
}

// ListActiveSpecs returns a list of all non-archived specs.
func ListActiveSpecs(ctx context.Context, projectRoot string) ([]SpecInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	specsDir := filepath.Join(projectRoot, ".specforce", "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SpecInfo{}, nil
		}
		return nil, err
	}

	var specs []SpecInfo
	for _, entry := range entries {
		// Ignore hidden directories and non-directories
		if entry.IsDir() && entry.Name()[0] != '.' {
			specs = append(specs, SpecInfo{Slug: entry.Name()})
		}
	}

	return specs, nil
}

// SpecExists checks if a spec slug already exists in either specs or archive.
func SpecExists(projectRoot string, slug string) (bool, string) {
	specsDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if _, err := os.Stat(specsDir); err == nil {
		return true, "active"
	}

	archiveDir := filepath.Join(projectRoot, ".specforce", "archive", slug)
	if _, err := os.Stat(archiveDir); err == nil {
		return true, "archived"
	}

	return false, ""
}
