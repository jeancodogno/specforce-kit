package agent

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func setupMockKitFS() fstest.MapFS {
	return fstest.MapFS{
		"kit.yaml": {
			Data: []byte(`
tools:
  claude:
    name: "Claude"
    description: "Claude agent"
    target: ".claude"
`),
		},
		"skills/tdd/SKILL.yaml": {
			Data: []byte("name: TDD"),
		},
	}
}

func TestRegistry_Initialize_Embedded(t *testing.T) {
	mockFS := setupMockKitFS()
	registry := &Registry{}
	err := registry.Initialize(mockFS, "")
	if err != nil {
		t.Fatalf("failed to initialize registry: %v", err)
	}

	agents := registry.GetAgents()
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}

	skills := registry.GetSkills()
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}

func TestRegistry_Initialize_LocalOverrides(t *testing.T) {
	mockFS := setupMockKitFS()
	registry := &Registry{}

	tempDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	specDir := filepath.Join(tempDir, ".specforce")
	if err := os.MkdirAll(filepath.Join(specDir, "skills", "local-skill"), 0755); err != nil {
		t.Fatal(err)
	}

	localKit := `
tools:
  claude:
    name: "Local Claude"
  new-agent:
    name: "New Agent"
`
	if err := os.WriteFile(filepath.Join(specDir, "kit.yaml"), []byte(localKit), 0644); err != nil {
		t.Fatal(err)
	}

	err = registry.Initialize(mockFS, tempDir)
	if err != nil {
		t.Fatalf("failed to initialize registry with overrides: %v", err)
	}

	claude, _ := registry.GetAgent("claude")
	if claude.Name != "Local Claude" {
		t.Errorf("expected Name 'Local Claude', got '%s'", claude.Name)
	}

	if _, ok := registry.GetAgent("new-agent"); !ok {
		t.Error("expected to find 'new-agent'")
	}

	skills := registry.GetSkills()
	foundLocal := false
	for _, s := range skills {
		if s.ID == "local-skill" {
			foundLocal = true
			break
		}
	}
	if !foundLocal {
		t.Error("expected to find 'local-skill'")
	}
}

func TestRegistry_ScanSkills_Metadata(t *testing.T) {
	mockFS := fstest.MapFS{
		"skill-yaml/SKILL.yaml": {Data: []byte("name: YAML Skill\nversion: 2.0.0\ndescription: A YAML skill")},
		"skill-md/SKILL.md":     {Data: []byte("---\nversion: 3.1.4\ndescription: A MD skill\n---\nContent")},
	}
	registry := &Registry{
		skills: make(map[string]SkillMetadata),
	}
	err := registry.scanSkills(mockFS)
	if err != nil {
		t.Fatal(err)
	}

	yamlSkill := registry.skills["skill-yaml"]
	if yamlSkill.Version != "2.0.0" || yamlSkill.Name != "YAML Skill" {
		t.Errorf("YAML skill metadata incorrect: %+v", yamlSkill)
	}

	mdSkill := registry.skills["skill-md"]
	if mdSkill.Version != "3.1.4" || mdSkill.Description != "A MD skill" {
		t.Errorf("MD skill metadata incorrect: %+v", mdSkill)
	}
}
