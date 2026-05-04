package agent

// AgentMetadata contains the descriptive metadata for an agent.
type AgentMetadata struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	DirName     string `yaml:"dir_name"`
	Version     string `yaml:"version"`
}

// AgentManifest is the structure of the manifest.yaml file found in each agent directory.
type AgentManifest struct {
	AgentMetadata `yaml:",inline"`
}

// SkillMetadata contains the descriptive metadata for a skill.
type SkillMetadata struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Version     string `yaml:"version"`
}
