package agent

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

func getTestKitConfig() *core.KitConfig {
	return &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"qwen": {
				Target: ".qwen",
				Mappings: map[string]core.MappingConfig{
					"agents":   {Path: "agents", Ext: ".md"},
					"skills":   {Path: "skills", Ext: ".md"},
					"commands": {Path: "commands/spf", Ext: ".md"},
				},
			},
			"open-code": {
				Target: ".opencode",
				Mappings: map[string]core.MappingConfig{
					"agents":   {Path: "agents", Ext: ".md"},
					"skills":   {Path: "skills", Ext: ".md"},
					"commands": {Path: "commands", Ext: ".md"},
				},
			},
			"kilo-code": {
				Target: ".kilocode",
				Mappings: map[string]core.MappingConfig{
					"agents":   {Path: "agents", Ext: ".md"},
					"skills":   {Path: "skills", Ext: ".md"},
					"commands": {Path: "commands", Ext: ".md"},
				},
			},
			"codex": {
				Target: ".codex",
				Mappings: map[string]core.MappingConfig{
					"agents":   {Path: "agents", Ext: ".md"},
					"skills":   {Path: "skills", Ext: ".md"},
					"commands": {Path: "commands/spf", Ext: ".md"},
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
		{"qwen", ".qwen"},
		{"open-code", ".opencode"},
		{"kilo-code", ".kilocode"},
		{"codex", ".codex"},
	}

	for _, agent := range agents {
		t.Run(agent.id, func(t *testing.T) {
			testAgentMappings(t, projectRoot, kitFS, kitConfig, agent.id, agent.folder)
		})
	}
}

func testAgentMappings(t *testing.T, projectRoot string, kitFS fs.FS, kitConfig *core.KitConfig, agentID, folder string) {
	// Test agent mapping
	err := processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/technical-developer.yaml", agentID, installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for %s: %v", agentID, err)
	}

	expectedPath := filepath.Join(projectRoot, folder, "agents/technical-developer.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expectedPath)
	}

	// Test skill mapping
	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "skills/task-atomic-decomposition/SKILL.yaml", agentID, installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for skill on %s: %v", agentID, err)
	}
	expectedSkillPath := filepath.Join(projectRoot, folder, "skills/task-atomic-decomposition/SKILL.md")
	if _, err := os.Stat(expectedSkillPath); os.IsNotExist(err) {
		t.Errorf("expected skill file %s to exist", expectedSkillPath)
	}

	// Test command mapping
	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "commands/archive.yaml", agentID, installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for command on %s: %v", agentID, err)
	}

	expectedCmdPath := filepath.Join(projectRoot, folder, "commands/spf/archive.md")
	if agentID == "open-code" || agentID == "kilo-code" {
		expectedCmdPath = filepath.Join(projectRoot, folder, "commands/spf.archive.md")
	}

	if _, err := os.Stat(expectedCmdPath); os.IsNotExist(err) {
		t.Errorf("expected command file %s to exist", expectedCmdPath)
	}
}
