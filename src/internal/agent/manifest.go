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
