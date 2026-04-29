package core_test

import (
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"gopkg.in/yaml.v3"
)

func TestKitConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
tools:
  gemini-cli:
    target: ".gemini/"
    mappings:
      skills:
        path: "skills"
        ext: ".md"
      agents:
        path: "prompts"
        ext: ".md"
  claude-code:
    target: ".claude/"
    mappings:
      skills:
        path: "instructions"
        ext: ".md"
`)

	var config core.KitConfig
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}

	if len(config.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(config.Tools))
	}

	gemini, ok := config.Tools["gemini-cli"]
	if !ok {
		t.Fatal("expected 'gemini-cli' tool")
	}

	if gemini.Target != ".gemini/" {
		t.Errorf("expected target '.gemini/', got %q", gemini.Target)
	}

	if len(gemini.Mappings) != 2 {
		t.Fatalf("expected 2 mappings for gemini-cli, got %d", len(gemini.Mappings))
	}

	skillsMap := gemini.Mappings["skills"][0]
	if skillsMap.Path != "skills" || skillsMap.Ext != ".md" {
		t.Errorf("unexpected skills mapping: %+v", skillsMap)
	}

	claude, ok := config.Tools["claude-code"]
	if !ok {
		t.Fatal("expected 'claude-code' tool")
	}

	if claude.Target != ".claude/" {
		t.Errorf("expected target '.claude/', got %q", claude.Target)
	}
}
