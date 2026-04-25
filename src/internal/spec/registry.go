package spec

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

// Artifact represents a specification document type.
type Artifact struct {
	Name        string `json:"name" yaml:"-"`
	Description string `json:"description" yaml:"description"`
	Instruction string `json:"instruction" yaml:"instruction"`
	Template    string `json:"template" yaml:"template"`
	Dependency  string `json:"dependency" yaml:"dependency"`
}

// Registry manages the collection of spec artifacts.
type Registry struct {
	artifacts map[string]Artifact
}

// NewRegistry initializes a new registry by loading YAML files from the provided filesystem.
func NewRegistry(artifactsFS fs.FS) (*Registry, error) {
	r := &Registry{
		artifacts: make(map[string]Artifact),
	}

	artifactRoot := "."

	err := fs.WalkDir(artifactsFS, artifactRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}

		art, err := loadArtifact(artifactsFS, path, d.Name())
		if err != nil {
			return err
		}

		r.artifacts[art.Name] = art
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load spec artifacts: %w", err)
	}

	// Circular dependency check
	if err := r.checkCircularDependencies(); err != nil {
		return nil, err
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

	// Name is derived from the filename without extension
	art.Name = strings.TrimSuffix(fileName, ".yaml")

	return art, nil
}

// Get returns an artifact by its name.
func (r *Registry) Get(name string) (Artifact, bool) {
	art, ok := r.artifacts[name]
	return art, ok
}

// List returns all loaded artifacts.
func (r *Registry) List() []Artifact {
	list := make([]Artifact, 0, len(r.artifacts))
	// Standard order: requirements, design, tasks
	order := []string{"requirements", "design", "tasks"}

	for _, name := range order {
		if art, ok := r.artifacts[name]; ok {
			list = append(list, art)
		}
	}

	// Catch any artifacts not in the explicit order
	for name, art := range r.artifacts {
		found := false
		for _, orderedName := range order {
			if name == orderedName {
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

func (r *Registry) checkCircularDependencies() error {
	for name := range r.artifacts {
		visited := make(map[string]bool)
		curr := name
		for curr != "" {
			if visited[curr] {
				return fmt.Errorf("circular dependency detected involving artifact: %s", curr)
			}
			visited[curr] = true
			art, ok := r.artifacts[curr]
			if !ok {
				break
			}
			curr = art.Dependency
		}
	}
	return nil
}
