package constitution

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Artifact represents a constitution document specification.
type Artifact struct {
	Slug        string `json:"slug" yaml:"-"`
	Name        string `json:"name" yaml:"-"`
	Description string `json:"description" yaml:"description"`
	Instruction string `json:"instruction" yaml:"instruction"`
	Template    string `json:"template" yaml:"template"`
	Path        string `json:"path" yaml:"-"`
}

// Registry manages the collection of constitution artifacts.
type Registry struct {
	artifacts map[string]Artifact
}

// NewRegistry initializes a new registry by loading YAML files from the provided filesystem.
func NewRegistry(artifactsFS fs.FS) (*Registry, error) {
	r := &Registry{
		artifacts: make(map[string]Artifact),
	}

	// The artifacts are located at the root of the provided filesystem
	// (usually sub-filesystem of internal/agent/artifacts/constitution)
	artifactRoot := "."
	
	err := fs.WalkDir(artifactsFS, artifactRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".yaml") || d.Name() == "mapping.yaml" {
			return nil
		}

		art, err := loadArtifact(artifactsFS, path, d.Name())
		if err != nil {
			return err
		}

		r.artifacts[art.Slug] = art
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load constitution artifacts: %w", err)
	}

	return r, nil
}

func loadArtifact(artifactsFS fs.FS, path, fileName string) (Artifact, error) {
	data, err := fs.ReadFile(artifactsFS, path)
	if err != nil {
		return Artifact{}, fmt.Errorf("failed to read artifact %s: %w", path, err)
	}

	var art Artifact
	if err := yaml.Unmarshal(data, &art); err != nil {
		return Artifact{}, fmt.Errorf("failed to parse artifact %s: %w", path, err)
	}

	// Validation: Check required fields
	if art.Description == "" {
		return Artifact{}, fmt.Errorf("artifact %s is missing 'description'", path)
	}
	if art.Instruction == "" {
		return Artifact{}, fmt.Errorf("artifact %s is missing 'instruction'", path)
	}
	if art.Template == "" {
		return Artifact{}, fmt.Errorf("artifact %s is missing 'template'", path)
	}

	// Slug is the filename without extension
	art.Slug = strings.TrimSuffix(fileName, ".yaml")
	
	// Handle "index" special case (remove leading underscore for Slug/Name)
	if art.Slug == "_index" {
		art.Slug = "index"
	}

	// Name is the clean artifact name (slug)
	art.Name = art.Slug

	// Path is the expected project-relative path (with extension)
	mdFileName := art.Slug + ".md"
	if art.Slug == "index" {
		mdFileName = "_index.md"
	}

	art.Path = filepath.Join(".specforce", "docs", mdFileName)

	return art, nil
}

// Get returns an artifact by its slug.
func (r *Registry) Get(slug string) (Artifact, bool) {
	art, ok := r.artifacts[slug]
	return art, ok
}

// List returns all loaded artifacts.
func (r *Registry) List() []Artifact {
	list := make([]Artifact, 0, len(r.artifacts))
	// We want a stable order, but for now simple slice is fine.
	// Common order: principles, architecture, ui-ux, security, engineering, governance, memorial
	order := []string{"principles", "architecture", "ui-ux", "security", "engineering", "governance", "memorial"}
	
	for _, slug := range order {
		if art, ok := r.artifacts[slug]; ok {
			list = append(list, art)
		}
	}

	// Catch any artifacts not in the explicit order
	for slug, art := range r.artifacts {
		found := false
		for _, orderedSlug := range order {
			if slug == orderedSlug {
				found = true
				break
			}
		}
		if !found {
			list = append(list, art)
		}
	}

	return list
}
