package core

// ToolRoute represents the target routing path and specific directory mappings for an agent tool.
type ToolRoute struct {
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Target      string                   `yaml:"target"`
	Mappings    map[string]MappingConfig `yaml:"mappings"`
}

// KitConfig represents the root configuration structure for kit.yaml,
// defining mapping rules for all supported AI agents.
type KitConfig struct {
	Tools map[string]ToolRoute `yaml:"tools"`
}
