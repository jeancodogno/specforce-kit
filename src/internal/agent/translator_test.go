package agent

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

func getMockKitConfig() *core.KitConfig {
	return &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"claude": {
				Target: ".claude/",
				Mappings: map[string]core.MappingConfigs{
					"agents":   {core.MappingConfig{Path: "agents", Ext: ".md"}},
					"skills":   {core.MappingConfig{Path: "skills", Ext: ".md"}},
					"commands": {core.MappingConfig{Path: "commands/spf", Ext: ".md"}},
				},
			},
			"gemini-cli": {
				Target: ".gemini/",
				Mappings: map[string]core.MappingConfigs{
					"agents": {core.MappingConfig{Path: "agents", Ext: ".toml"}},
				},
			},
			"kimi-code": {
				Target: ".kimi/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "skills/spf-*", Name: "SKILL", Ext: ".md"}},
					"skills":   {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
			"wildcard-tool": {
				Target: ".wildcard/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "tools", Name: "*-kit", Ext: ".md"}},
				},
			},
		},
	}
}

func TestResolveMapping_Wildcards(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-wildcards-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "commands"), 0755)

	blueprintYAML := "description: Test Command\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(kitDir, "commands/archive.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"generic-tool": {
				Target: ".generic/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "skills/spf-*", Name: "SKILL", Ext: ".md"}},
				},
			},
			"wildcard-tool": {
				Target: ".wildcard/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "tools", Name: "*-kit", Ext: ".md"}},
				},
			},
		},
	}
	projectRoot := filepath.Join(tmpDir, "project")

	// Test case 1: Path expansion + fixed name
	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "commands/archive.yaml", "generic-tool", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for generic-tool: %v", err)
	}
	expectedGeneric := filepath.Join(projectRoot, ".generic/skills/spf-archive/SKILL.md")
	if _, err := os.Stat(expectedGeneric); os.IsNotExist(err) {
		t.Errorf("expected expanded path %s to exist", expectedGeneric)
	}

	// Test case 2: Fixed path + name expansion
	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "commands/archive.yaml", "wildcard-tool", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for wildcard-tool: %v", err)
	}
	expectedWildcard := filepath.Join(projectRoot, ".wildcard/tools/archive-kit.md")
	if _, err := os.Stat(expectedWildcard); os.IsNotExist(err) {
		t.Errorf("expected expanded name %s to exist", expectedWildcard)
	}
}

func TestResolveMapping_Whitelisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-whitelisting-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "agents"), 0755)
	_ = os.WriteFile(filepath.Join(kitDir, "agents/test.yaml"), []byte("description: test"), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"minimal": {
				Target: ".minimal/",
				Mappings: map[string]core.MappingConfigs{
					"skills": {core.MappingConfig{Path: "skills", Ext: ".md"}},
				},
			},
		},
	}
	projectRoot := filepath.Join(tmpDir, "project")

	// Category 'agents' is NOT in 'minimal' mappings
	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/test.yaml", "minimal", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Should NOT exist
	expectedPath := filepath.Join(projectRoot, ".minimal/agents/test.md")
	if _, err := os.Stat(expectedPath); err == nil {
		t.Errorf("expected file %s to NOT exist (should be whitelisted out)", expectedPath)
	}
}

func TestResolveMapping_Defaulting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-defaulting-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "agents"), 0755)
	_ = os.WriteFile(filepath.Join(kitDir, "agents/my-slug.yaml"), []byte("description: test"), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"test-tool": {
				Target: ".test/",
				Mappings: map[string]core.MappingConfigs{
					"agents": {core.MappingConfig{Path: "dir", Name: "", Ext: ".md"}}, // Empty Name
				},
			},
		},
	}
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/my-slug.yaml", "test-tool", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Should default to 'my-slug'
	expectedPath := filepath.Join(projectRoot, ".test/dir/my-slug.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist (defaulted to slug)", expectedPath)
	}
}

func TestKimiExclusion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-kimi-exclusion-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	agentsDir := filepath.Join(kitDir, "agents")
	_ = os.MkdirAll(agentsDir, 0755)

	blueprintYAML := "description: Test Agent\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(agentsDir, "test-agent.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/test-agent.yaml", "kimi-code", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Should NOT exist
	expectedPath := filepath.Join(projectRoot, ".kimi/agents/test-agent.md")
	if _, err := os.Stat(expectedPath); err == nil {
		t.Errorf("expected file %s to NOT exist for kimi-code", expectedPath)
	}
}

func TestKimiCommandToSkill(t *testing.T) {
	runTranslatorTest(t, translatorTestParams{
		blueprintPath: "commands/archive.yaml",
		blueprintYAML: "description: Test Command\ncontent: |\n  # Command Content\n",
		targetTool:    "kimi-code",
		expectedPath:  ".kimi/skills/spf-archive/SKILL.md",
		checkContent: func(t *testing.T, content string) {
			if !strings.Contains(content, "description: Test Command") {
				t.Error("Command-to-Skill transformation missing header")
			}
		},
	})
}

func TestAdaptArtifacts_Success(t *testing.T) {
	runTranslatorTest(t, translatorTestParams{
		blueprintPath: "agents/test-agent.yaml",
		blueprintYAML: "description: Test Agent\ncontent: |\n  # Hello\n",
		targetTool:    "claude",
		expectedPath:  ".claude/agents/test-agent.md",
		checkContent: func(t *testing.T, content string) {
			if !strings.HasPrefix(content, "---") {
				t.Errorf("expected YAML frontmatter header, got: %q", content)
			}
			if !strings.Contains(content, "name: test-agent") {
				t.Error("header missing name")
			}
		},
	})
}

func TestUnmappedAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-unmapped-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	agentsDir := filepath.Join(kitDir, "agents")
	_ = os.MkdirAll(agentsDir, 0755)

	blueprintYAML := "description: Test Agent\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(agentsDir, "test-agent.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/test-agent.yaml", "unknown-agent", installer.Options{})
	if !errors.Is(err, core.ErrToolMappingNotFound) {
		t.Errorf("expected ErrToolMappingNotFound, got %v", err)
	}
}

func TestLegacyMappingIgnored(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-legacy-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	agentsDir := filepath.Join(kitDir, "agents")
	_ = os.MkdirAll(agentsDir, 0755)

	// Write a legacy mapping.yaml with completely different settings
	legacyMapping := "claude:\n  path: .legacy/path\n  ext: .legacy\n"
	_ = os.WriteFile(filepath.Join(agentsDir, "mapping.yaml"), []byte(legacyMapping), 0644)

	blueprintYAML := "description: Test Agent\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(agentsDir, "test-agent.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/test-agent.yaml", "claude", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Should ignore legacy mapping and use kitConfig
	expectedPath := filepath.Join(projectRoot, ".claude/agents/test-agent.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist (legacy mapping was not ignored)", expectedPath)
	}
}

func TestMappingRegistry_Subdirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-mapping-sub-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	agentsDir := filepath.Join(kitDir, "agents")
	subDir := filepath.Join(agentsDir, "special")
	_ = os.MkdirAll(subDir, 0755)

	blueprintYAML := "description: Test Agent\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(subDir, "sub-agent.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project-rel")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/special/sub-agent.yaml", "claude", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Expected: .claude/agents/special/sub-agent.md
	expectedPath := filepath.Join(projectRoot, ".claude/agents/special/sub-agent.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expectedPath)
	}
}

func TestMappingGemini(t *testing.T) {
	runTranslatorTest(t, translatorTestParams{
		blueprintPath: "agents/test-agent.yaml",
		blueprintYAML: "description: Test Agent\ncontent: |\n  # Hello\n",
		targetTool:    "gemini-cli",
		expectedPath:  ".gemini/agents/test-agent.toml",
		checkContent: func(t *testing.T, content string) {
			if !strings.Contains(content, `description = "Test Agent"`) {
				t.Error("missing description in TOML")
			}
		},
	})
}

type translatorTestParams struct {
	blueprintPath string
	blueprintYAML string
	targetTool    string
	expectedPath  string
	checkContent  func(t *testing.T, content string)
}

func runTranslatorTest(t *testing.T, params translatorTestParams) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-translator-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	blueprintFile := filepath.Join(kitDir, params.blueprintPath)
	_ = os.MkdirAll(filepath.Dir(blueprintFile), 0755)
	_ = os.WriteFile(blueprintFile, []byte(params.blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, params.blueprintPath, params.targetTool, installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	fullExpectedPath := filepath.Join(projectRoot, params.expectedPath)
	if _, err := os.Stat(fullExpectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", fullExpectedPath)
	}

	if params.checkContent != nil {
		data, _ := os.ReadFile(fullExpectedPath)
		params.checkContent(t, string(data))
	}
}

func TestMappingOverrides(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-overrides-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	agentsDir := filepath.Join(kitDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}

	// Blueprint with explicit override
	blueprintYAML := "description: Override Agent\nmapping:\n  claude:\n    path: custom/path\n    name: overridden-name\ncontent: |\n  # Hello\n"
	if err := os.WriteFile(filepath.Join(agentsDir, "override-agent.yaml"), []byte(blueprintYAML), 0644); err != nil {
		t.Fatalf("failed to write blueprint: %v", err)
	}

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "agents/override-agent.yaml", "claude", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Overrides should be relative to the target directory of the tool
	// Expected: .claude/custom/path/overridden-name.md
	expectedPath := filepath.Join(projectRoot, ".claude/custom/path/overridden-name.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expectedPath)
	}
}

func TestSkillHeader(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-skills-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	skillsDir := filepath.Join(kitDir, "skills/test-skill")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}

	skillYAML := "name: test-skill\ndescription: Test Skill Description\nversion: \"2.1\"\npriority: HIGH\ncontent: |\n  # Skill Content\n"
	skillYAMLPath := filepath.Join(skillsDir, "SKILL.yaml")
	if err := os.WriteFile(skillYAMLPath, []byte(skillYAML), 0644); err != nil {
		t.Fatalf("failed to write skill: %v", err)
	}

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "skills/test-skill/SKILL.yaml", "claude", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for SKILL.yaml: %v", err)
	}

	skillMD := filepath.Join(projectRoot, ".claude/skills/test-skill/SKILL.md")
	data, err := os.ReadFile(skillMD)
	if err != nil {
		t.Fatalf("failed to read expected skill file %s: %v", skillMD, err)
	}
	content := string(data)
	if !strings.Contains(content, "name: test-skill") {
		t.Errorf("SKILL.md header missing name from metadata. Content: %q", content)
	}
	if !strings.Contains(content, "version: 2.1") || !strings.Contains(content, "priority: HIGH") {
		t.Errorf("SKILL.md missing version or priority. Content: %q", content)
	}
}

func TestCommandHeader(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-commands-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	commandsDir := filepath.Join(kitDir, "commands")
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		t.Fatalf("failed to create commands dir: %v", err)
	}

	commandYAML := "description: Test Command Description\ncontent: |\n  # Command Content\n"
	if err := os.WriteFile(filepath.Join(commandsDir, "test-cmd.yaml"), []byte(commandYAML), 0644); err != nil {
		t.Fatalf("failed to write command: %v", err)
	}

	kitFS := os.DirFS(kitDir)
	kitConfig := getMockKitConfig()
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "commands/test-cmd.yaml", "claude", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	cmdMD := filepath.Join(projectRoot, ".claude/commands/spf/test-cmd.md")
	data, _ := os.ReadFile(cmdMD)
	content := string(data)
	if !strings.HasPrefix(content, "---") || !strings.Contains(content, "description: Test Command Description") {
		t.Errorf("Command missing description in header. Content: %q", content)
	}
}

func TestStandardizedNaming(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-naming-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	commandsDir := filepath.Join(kitDir, "commands")
	_ = os.MkdirAll(commandsDir, 0755)

	blueprintYAML := "description: Test Command\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(commandsDir, "archive.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"codex": {
				Target: ".codex/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "commands/spf", Name: "spf-*", Ext: ".md"}},
				},
			},
		},
	}
	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfig, "commands/archive.yaml", "codex", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	expectedPath := filepath.Join(projectRoot, ".codex/commands/spf/spf-archive.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expectedPath)
	}

	// Test case 3: Absolute path for global-enabled agent
	globalTmp, _ := os.MkdirTemp("", "specforce-global-*")
	defer func() { _ = os.RemoveAll(globalTmp) }()

	kitConfigGlobal := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"codex": {
				Target: globalTmp,
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "spf", Name: "spf-*", Ext: ".md"}},
				},
			},
		},
	}

	err = processBlueprint(context.Background(), projectRoot, kitFS, kitConfigGlobal, "commands/archive.yaml", "codex", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed for global absolute path: %v", err)
	}

	expectedGlobalPath := filepath.Join(globalTmp, "spf/spf-archive.md")
	if _, err := os.Stat(expectedGlobalPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist in global path", expectedGlobalPath)
	}
}




func TestResolveMapping_TargetOverride(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-target-override-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	globalDir := filepath.Join(tmpDir, "global")
	_ = os.MkdirAll(globalDir, 0755)

	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"test-agent": {
				Target: ".local/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {
						core.MappingConfig{
							Target: globalDir,
							Path:   "global-cmds",
							Ext:    ".md",
						},
					},
					"skills": {
						core.MappingConfig{
							Path: "local-skills",
							Ext:  ".md",
						},
					},
				},
			},
		},
	}

	bp := &core.Blueprint{ID: "test"}

	// Test case 1: Command uses override target
	mappings, err := resolveMappings(kitConfig, "commands/test.yaml", bp, "test-agent")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping := mappings[0]
	expectedPath := filepath.Clean(filepath.Join(globalDir, "global-cmds"))
	if mapping.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, mapping.Path)
	}

	// Test case 2: Skill uses default tool target
	mappings, err = resolveMappings(kitConfig, "skills/test.yaml", bp, "test-agent")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping = mappings[0]
	expectedPath = filepath.Clean(filepath.Join(".local/", "local-skills"))
	if mapping.Path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, mapping.Path)
	}
}

func TestResolveMapping_TildeExpansionFailure(t *testing.T) {
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"test-agent": {
				Target: "~/hidden-path",
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "cmds", Ext: ".md"}},
				},
			},
		},
	}

	oldHome := os.Getenv("HOME")
	_ = os.Unsetenv("HOME")
	_ = os.Unsetenv("USERPROFILE") // For windows compatibility in tests if ever run there
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	// macOS and some other platforms might have fallback mechanisms for UserHomeDir
	// that don't rely on environment variables. If it still succeeds, we skip the test.
	if _, err := os.UserHomeDir(); err == nil {
		t.Skip("os.UserHomeDir() still succeeds after unsetting HOME; platform has fallback")
	}

	bp := &core.Blueprint{ID: "test"}
	_, err := resolveMappings(kitConfig, "commands/test.yaml", bp, "test-agent")
	if err == nil {
		t.Fatal("expected error for failed tilde expansion, got nil")
	}
	if !strings.Contains(err.Error(), "failed to resolve home directory") {
		t.Errorf("expected home directory resolution error, got: %v", err)
	}
}

func TestCodexExpansionHierarchy(t *testing.T) {
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"codex": {
				Mappings: map[string]core.MappingConfigs{
					"commands": {
						core.MappingConfig{
							Target: "${CODEX_PROMPTS_DIR:-${CODEX_HOME:-~/.codex}/prompts}",
							Path:   ".",
							Ext:    ".md",
						},
					},
				},
			},
		},
	}
	bp := &core.Blueprint{ID: "test"}

	// Cleanup env after test
	defer func() {
		_ = os.Unsetenv("CODEX_PROMPTS_DIR")
		_ = os.Unsetenv("CODEX_HOME")
	}()

	// Case 1: CODEX_PROMPTS_DIR
	_ = os.Setenv("CODEX_PROMPTS_DIR", "/path/to/prompts")
	mappings, err := resolveMappings(kitConfig, "commands/test.yaml", bp, "codex")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping := mappings[0]
	if mapping.Path != "/path/to/prompts" {
		t.Errorf("Priority 1 failed: expected /path/to/prompts, got %s", mapping.Path)
	}

	// Case 2: CODEX_HOME
	_ = os.Unsetenv("CODEX_PROMPTS_DIR")
	_ = os.Setenv("CODEX_HOME", "/home/user/.codex")
	mappings, err = resolveMappings(kitConfig, "commands/test.yaml", bp, "codex")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping = mappings[0]
	if mapping.Path != "/home/user/.codex/prompts" {
		t.Errorf("Priority 2 failed: expected /home/user/.codex/prompts, got %s", mapping.Path)
	}

	// Case 3: Default ~/.codex/prompts
	_ = os.Unsetenv("CODEX_HOME")
	home, _ := os.UserHomeDir()
	mappings, err = resolveMappings(kitConfig, "commands/test.yaml", bp, "codex")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping = mappings[0]
	expected := filepath.Clean(filepath.Join(home, ".codex/prompts"))
	if mapping.Path != expected {
		t.Errorf("Priority 3 failed: expected %s, got %s", expected, mapping.Path)
	}
}

func TestCodexGlobalLocalIsolation(t *testing.T) {
	globalDir := "/tmp/global-codex"

	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"codex": {
				Target: ".codex/", // Local target
				Mappings: map[string]core.MappingConfigs{
					"commands": {
						core.MappingConfig{
							Target: globalDir, // Global override
							Path:   ".",
							Ext:    ".md",
						},
					},
					"skills": {
						core.MappingConfig{
							Path: "skills",
							Ext:  ".md",
						},
					},
				},
			},
		},
	}
	bp := &core.Blueprint{ID: "test"}

	// 1. Verify Command is Global (Absolute)
	mappings, err := resolveMappings(kitConfig, "commands/test.yaml", bp, "codex")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping := mappings[0]
	if mapping.Path != globalDir {
		t.Errorf("expected global path %s, got %s", globalDir, mapping.Path)
	}

	// 2. Verify Skill is Local (Relative)
	mappings, err = resolveMappings(kitConfig, "skills/test.yaml", bp, "codex")
	if err != nil {
		t.Fatalf("resolveMappings failed: %v", err)
	}
	mapping = mappings[0]
	expectedLocal := filepath.Clean(filepath.Join(".codex/", "skills"))
	if mapping.Path != expectedLocal {
		t.Errorf("expected local path %s, got %s", expectedLocal, mapping.Path)
	}
}
