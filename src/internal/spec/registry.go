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
	order     []string // Cached topological order of artifact names
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

	// Resolve topological order and check for cycles
	if err := r.topologicalSort(); err != nil {
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

// Get returns an artifact by its name. It supports prefixed lookups (e.g., "bug-requirements")
// and right-to-left inference for base types (e.g., "tasks-for-design" -> "design").
func (r *Registry) Get(name string) (Artifact, bool) {
	// 1. Try exact match (could be base name or already a full prefixed name in the registry)
	if art, ok := r.artifacts[name]; ok {
		return art, true
	}

	// 2. Try to resolve if it has a known prefix
	for _, t := range []string{"bug", "feature"} {
		prefix := t + "-"
		if strings.HasPrefix(name, prefix) {
			baseName := strings.TrimPrefix(name, prefix)
			return r.GetForType(t, baseName)
		}
	}

	// 3. Right-to-left inference for base types
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
	if bestMatch != "" {
		if art, ok := r.artifacts[bestMatch]; ok {
			return art, true
		}
	}

	return Artifact{}, false
}

// List returns all loaded artifacts in topological order.
func (r *Registry) List() []Artifact {
	list := make([]Artifact, 0, len(r.order))
	for _, name := range r.order {
		if art, ok := r.artifacts[name]; ok {
			list = append(list, art)
		}
	}

	// Catch any artifacts not in the dependency graph (if any were somehow missed by topologicalSort)
	for name, art := range r.artifacts {
		found := false
		for _, orderedName := range r.order {
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

// GetForType returns an artifact by its name, checking for typed versions first.
func (r *Registry) GetForType(typeName, artifactName string) (Artifact, bool) {
	// Try typed version: "{type}-{name}"
	typedName := fmt.Sprintf("%s-%s", typeName, artifactName)
	if art, ok := r.artifacts[typedName]; ok {
		// Return with the original artifact name to keep filenames consistent
		art.Name = artifactName
		return art, true
	}

	// Fallback to default
	art, ok := r.artifacts[artifactName]
	return art, ok
}

// ListForType returns all loaded artifacts resolved for a specific type, respecting topological order.
func (r *Registry) ListForType(typeName string) []Artifact {
	list := make([]Artifact, 0, len(r.order))
	seenBase := make(map[string]bool)

	for _, name := range r.order {
		if art, ok := r.GetForType(typeName, name); ok {
			list = append(list, art)
			seenBase[name] = true
		}
	}

	// Catch any artifacts not in the explicit order
	for name := range r.artifacts {
		// Extract base name if it's a typed artifact (e.g., "bug-requirements" -> "requirements")
		baseName := name
		for _, t := range []string{"bug", "feature"} {
			prefix := t + "-"
			if strings.HasPrefix(name, prefix) {
				baseName = strings.TrimPrefix(name, prefix)
				break
			}
		}

		if !seenBase[baseName] {
			if art, ok := r.GetForType(typeName, baseName); ok {
				list = append(list, art)
				seenBase[baseName] = true
			}
		}
	}

	return list
}

func (r *Registry) topologicalSort() error {
	var sorted []string
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	// We only want to sort base artifacts and their dependencies.
	// Prefixed artifacts (bug-*) should not be in the base dependency graph
	// as they are resolved via GetForType.
	var baseArtifacts []string
	for name := range r.artifacts {
		isPrefixed := false
		for _, p := range []string{"bug-", "feature-"} {
			if strings.HasPrefix(name, p) {
				isPrefixed = true
				break
			}
		}
		if !isPrefixed {
			baseArtifacts = append(baseArtifacts, name)
		}
	}

	var visit func(string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("circular dependency detected involving artifact: %s", name)
		}
		if !visited[name] {
			temp[name] = true
			art, ok := r.artifacts[name]
			if ok && art.Dependency != "" {
				if err := visit(art.Dependency); err != nil {
					return err
				}
			}
			visited[name] = true
			delete(temp, name)
			sorted = append(sorted, name)
		}
		return nil
	}

	for _, name := range baseArtifacts {
		if !visited[name] {
			if err := visit(name); err != nil {
				return err
			}
		}
	}

	r.order = sorted
	return nil
}
