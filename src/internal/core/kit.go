package core

import "gopkg.in/yaml.v3"

// MappingConfigs represents one or more mapping configurations.
type MappingConfigs []MappingConfig

// UnmarshalYAML implements custom unmarshaling to handle both single object and array of objects.
func (mc *MappingConfigs) UnmarshalYAML(value *yaml.Node) error {
	// Try single object first
	var single MappingConfig
	if err := value.Decode(&single); err == nil {
		*mc = []MappingConfig{single}
		return nil
	}

	// Try slice of objects
	var slice []MappingConfig
	if err := value.Decode(&slice); err != nil {
		return err
	}
	*mc = slice
	return nil
}

// ToolRoute represents the target routing path and specific directory mappings for an agent tool.
type ToolRoute struct {
	Name        string                    `yaml:"name"`
	Description string                    `yaml:"description"`
	Target      string                    `yaml:"target"`
	Mappings    map[string]MappingConfigs `yaml:"mappings"`
}

// KitConfig represents the root configuration structure for kit.yaml,
// defining mapping rules for all supported AI agents.
type KitConfig struct {
	Tools map[string]ToolRoute `yaml:"tools"`
}
