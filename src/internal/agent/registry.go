package agent

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"gopkg.in/yaml.v3"
)

// Registry stores and manages discovered agent metadata.
type Registry struct {
	agents map[string]AgentMetadata
	skills map[string]SkillMetadata
}

// Initialize parses kit.yaml from the provided filesystem and optionally from the local
// .specforce directory to populate agent and skill metadata.
func (r *Registry) Initialize(kitFS fs.FS, rootDir string) error {
	r.agents = make(map[string]AgentMetadata)
	r.skills = make(map[string]SkillMetadata)

	// 1. Load embedded agents
	if err := r.loadAgents(kitFS); err != nil {
		return err
	}

	// 2. Load embedded skills
	if err := r.loadSkills(kitFS); err != nil {
		return err
	}

	// 3. Load local overrides (optional)
	if rootDir != "" {
		localPath := filepath.Join(rootDir, ".specforce", "kit.yaml")
		if _, err := os.Stat(localPath); err == nil {
			data, err := os.ReadFile(localPath)
			if err == nil {
				if err := r.parseAgents(data); err != nil {
					return fmt.Errorf("failed to parse local kit.yaml: %w", err)
				}
			}
		}

		// Local skills
		localSkillsPath := filepath.Join(rootDir, ".specforce", "skills")
		if info, err := os.Stat(localSkillsPath); err == nil && info.IsDir() {
			if err := r.scanSkills(os.DirFS(localSkillsPath)); err != nil {
				return fmt.Errorf("failed to scan local skills: %w", err)
			}
		}
	}

	return nil
}

func (r *Registry) loadAgents(kitFS fs.FS) error {
	data, err := fs.ReadFile(kitFS, "kit.yaml")
	if err != nil {
		return fmt.Errorf("failed to read kit.yaml: %w", err)
	}
	return r.parseAgents(data)
}

func (r *Registry) parseAgents(data []byte) error {
	var config core.KitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse kit.yaml: %w", err)
	}

	for id, route := range config.Tools {
		name := route.Name
		if name == "" {
			name = id
		}

		metadata := AgentMetadata{
			ID:          id,
			Name:        name,
			Description: route.Description,
			DirName:     route.Target,
			Version:     "1.0.0",
		}

		r.agents[id] = metadata
	}
	return nil
}

func (r *Registry) loadSkills(kitFS fs.FS) error {
	skillsFS, err := fs.Sub(kitFS, "skills")
	if err != nil {
		return nil // No skills subfolder
	}
	return r.scanSkills(skillsFS)
}

func (r *Registry) scanSkills(skillsFS fs.FS) error {
	return fs.WalkDir(skillsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != "." {
			skillID := d.Name()
			r.skills[skillID] = r.loadSkillMetadata(skillsFS, path, skillID)
		}
		return nil
	})
}

func (r *Registry) loadSkillMetadata(skillsFS fs.FS, path, skillID string) SkillMetadata {
	metadata := SkillMetadata{
		ID:      skillID,
		Path:    path,
		Version: "1.0.0", // Default
	}

	// Try to load detailed metadata from SKILL.yaml
	if data, err := fs.ReadFile(skillsFS, filepath.Join(path, "SKILL.yaml")); err == nil {
		var skill SkillMetadata
		if err := yaml.Unmarshal(data, &skill); err == nil {
			if skill.Name != "" {
				metadata.Name = skill.Name
			}
			if skill.Description != "" {
				metadata.Description = skill.Description
			}
			if skill.Version != "" {
				metadata.Version = skill.Version
			}
		}
		return metadata
	}

	// Try to load from SKILL.md
	if data, err := fs.ReadFile(skillsFS, filepath.Join(path, "SKILL.md")); err == nil {
		// Try to parse version from Markdown frontmatter
		bp, err := core.ParseBlueprint(path, data)
		if err == nil {
			if bp.Metadata.Version != "" {
				metadata.Version = bp.Metadata.Version
			}
			if bp.Metadata.Description != "" {
				metadata.Description = bp.Metadata.Description
			}
		}
	}

	return metadata
}

// GetAgents returns a list of all discovered agent metadata, sorted by ID.
func (r *Registry) GetAgents() []AgentMetadata {
	agents := make([]AgentMetadata, 0, len(r.agents))
	for _, a := range r.agents {
		agents = append(agents, a)
	}

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].ID < agents[j].ID
	})

	return agents
}

// GetSkills returns a list of all discovered skill metadata, sorted by ID.
func (r *Registry) GetSkills() []SkillMetadata {
	skills := make([]SkillMetadata, 0, len(r.skills))
	for _, s := range r.skills {
		skills = append(skills, s)
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].ID < skills[j].ID
	})

	return skills
}

// GetAgent returns a specific agent's metadata by ID.
func (r *Registry) GetAgent(id string) (AgentMetadata, bool) {
	a, ok := r.agents[id]
	return a, ok
}
