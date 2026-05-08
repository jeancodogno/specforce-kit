package spec

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Metadata represents the core configuration of a specification.
type Metadata struct {
	Slug string `json:"slug" yaml:"slug"`
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"` // "feature" | "bug"
}

// LoadMetadata reads the spec.yaml from the specification directory.
// It returns a default "feature" metadata if the file does not exist.
func LoadMetadata(projectRoot, slug string) (*Metadata, error) {
	metaPath := filepath.Join(projectRoot, ".specforce", "specs", slug, "spec.yaml")
	
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

	return os.WriteFile(metaPath, data, 0644)
}
