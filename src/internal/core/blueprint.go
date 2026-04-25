package core

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// MappingConfig defines where and how a blueprint should be adapted for an agent.
type MappingConfig struct {
	Target string `yaml:"target,omitempty"`
	Path   string `yaml:"path"`
	Name   string `yaml:"name"`
	Ext    string `yaml:"ext"`
}

// BlueprintMetadata contains the structured metadata of a framework asset.
type BlueprintMetadata struct {
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Version     string                   `yaml:"version,omitempty"`
	Priority    string                   `yaml:"priority,omitempty"`
	Triggers    []string                 `yaml:"triggers,omitempty"`
	Mapping     map[string]MappingConfig `yaml:"mapping"`
	Content     string                   `yaml:"content,omitempty"`
}

// Blueprint represents a framework asset with metadata and content.
type Blueprint struct {
	ID       string
	Metadata BlueprintMetadata
	Content  string
}

// ParseBlueprint parses a file content (pure YAML) into a Blueprint struct.
func ParseBlueprint(id string, data []byte) (*Blueprint, error) {
	var metadata BlueprintMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse blueprint metadata (ID: %s): %w", id, err)
	}

	// For YAML blueprints, we expect both metadata and content to be in the same file structure.
	return &Blueprint{
		ID:       id,
		Metadata: metadata,
		Content:  strings.TrimSpace(metadata.Content),
	}, nil
}
