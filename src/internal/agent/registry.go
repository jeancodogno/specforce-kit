package agent

import (
	"fmt"
	"io/fs"
	"sort"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"gopkg.in/yaml.v3"
)

// Registry stores and manages discovered agent metadata.
type Registry struct {
	agents map[string]AgentMetadata
}

// Initialize parses kit.yaml from the provided filesystem to populate agent metadata.
func (r *Registry) Initialize(kitFS fs.FS) error {
	r.agents = make(map[string]AgentMetadata)

	data, err := fs.ReadFile(kitFS, "kit.yaml")
	if err != nil {
		return fmt.Errorf("failed to read kit.yaml: %w", err)
	}

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
			Version:     "1.0.0", // default since we removed version
		}

		r.agents[id] = metadata
	}

	return nil
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

// GetAgent returns a specific agent's metadata by ID.
func (r *Registry) GetAgent(id string) (AgentMetadata, bool) {
	a, ok := r.agents[id]
	return a, ok
}
