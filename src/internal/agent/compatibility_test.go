package agent

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

// List returns all discovered agents (alias for GetAgents).
func (r *Registry) List() []AgentMetadata {
	return r.GetAgents()
}

func getTestKitConfig() *core.KitConfig {
	return &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"claude": {
				Target: ".claude",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"qwen": {
				Target: ".qwen",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"open-code": {
				Target: ".opencode",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"kilo-code": {
				Target: ".kilocode",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"codex": {
				Target: ".codex",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"antigravity": {
				Target: ".agents",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"cursor": {
				Target: ".cursor",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"kimi-code": {
				Target: ".kimi",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
		},
	}
}

func TestNewAgentMappings(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-compat-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectRoot := filepath.Join(tmpDir, "project")
	kitFS, err := GetKitFS()
	if err != nil {
		t.Fatalf("failed to get kit FS: %v", err)
	}

	kitConfig := getTestKitConfig()
	agents := []struct {
		id     string
		folder string
	}{
		{"claude", ".claude"},
		{"qwen", ".qwen"},
		{"open-code", ".opencode"},
		{"kilo-code", ".kilocode"},
		{"codex", ".codex"},
		{"antigravity", ".agents"},
		{"cursor", ".cursor"},
		{"kimi-code", ".kimi"},
	}

	for _, agent := range agents {
		t.Run(agent.id, func(t *testing.T) {
			testAgentMappings(t, projectRoot, kitFS, kitConfig, agent.id, agent.folder)
		})
	}

	// Ensure gemini-cli is no longer in the kit config
	if _, ok := kitConfig.Tools["gemini-cli"]; ok {
		t.Errorf("gemini-cli should not be in the kit config")
	}

	// Verify Gemini CLI is not in registry.List()
	registry := &Registry{}
	if err := registry.Initialize(kitFS, ""); err != nil {
		t.Fatalf("failed to initialize registry: %v", err)
	}
	for _, a := range registry.List() {
		if a.ID == "gemini-cli" || a.Name == "Gemini CLI" {
			t.Errorf("gemini-cli should not be in registry.List(), found: %+v", a)
		}
	}
	if _, ok := registry.GetAgent("gemini-cli"); ok {
		t.Errorf("gemini-cli should not be in registry")
	}
}

func testAgentMappings(t *testing.T, projectRoot string, kitFS fs.FS, kitConfig *core.KitConfig, agentID, folder string) {
	skills := []string{
		"spf-archive",
		"spf-constitution",
		"spf-discovery",
		"spf-implement",
		"spf-spec",
	}

	for _, skill := range skills {
		skillPath := fmt.Sprintf("skills/%s/SKILL.yaml", skill)
		err := processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, skillPath, agentID, installer.Options{})
		if err != nil {
			t.Fatalf("processBlueprint failed for skill %s on %s: %v", skill, agentID, err)
		}
		expectedSkillPath := filepath.Join(projectRoot, folder, "skills", skill, "SKILL.md")
		if _, err := os.Stat(expectedSkillPath); os.IsNotExist(err) {
			t.Errorf("expected skill file %s to exist", expectedSkillPath)
		}
	}

	// Add negative assertions ensuring that .agents/workflows/ and <tool>/commands/ directories are never created
	badPaths := []string{
		filepath.Join(projectRoot, folder, "workflows"),
		filepath.Join(projectRoot, folder, "commands"),
	}
	for _, badPath := range badPaths {
		if _, err := os.Stat(badPath); !os.IsNotExist(err) {
			t.Errorf("directory %s should not exist", badPath)
		}
	}
}
