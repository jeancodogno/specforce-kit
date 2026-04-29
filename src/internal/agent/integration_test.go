package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

func TestIntegration_CodexGlobalExport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-int-codex-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	promptsDir := filepath.Join(tmpDir, "global-prompts")
	if err := os.Setenv("CODEX_PROMPTS_DIR", promptsDir); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() { _ = os.Unsetenv("CODEX_PROMPTS_DIR") }()

	projectRoot := filepath.Join(tmpDir, "project")
	_ = os.MkdirAll(projectRoot, 0755)

	// Mock kitFS with the actual kit.yaml we just modified
	kitFS := os.DirFS("kit")

	kitConfig, err := LoadKitConfig(kitFS, projectRoot)
	if err != nil {
		t.Fatalf("failed to load kit config: %v", err)
	}

	// Create a dummy blueprint
	kitDir := filepath.Join(tmpDir, "mock-kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "commands"), 0755)
	blueprintYAML := "description: Test Integration\ncontent: |\n  # Integrated\n"
	_ = os.WriteFile(filepath.Join(kitDir, "commands/test-int.yaml"), []byte(blueprintYAML), 0644)
	
	blueprintFS := os.DirFS(kitDir)

	err = processBlueprint(context.Background(), projectRoot, blueprintFS, kitConfig, "commands/test-int.yaml", "codex", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Expected: promptsDir/spf-test-int.md
	expectedPath := filepath.Join(promptsDir, "spf-test-int.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected integrated file %s to exist", expectedPath)
	}
}

func TestIntegration_AntigravityFlatExport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-int-antigravity-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectRoot := filepath.Join(tmpDir, "project")
	_ = os.MkdirAll(projectRoot, 0755)

	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"antigravity": {
				Target: ".agent/",
				Mappings: map[string]core.MappingConfigs{
					"commands": {
						core.MappingConfig{
							Path: ".",
							Name: "spf-*",
							Ext:  ".md",
						},
					},
				},
			},
		},
	}

	kitDir := filepath.Join(tmpDir, "mock-kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "commands"), 0755)
	_ = os.WriteFile(filepath.Join(kitDir, "commands/wf.yaml"), []byte("description: wf"), 0644)
	blueprintFS := os.DirFS(kitDir)

	err = processBlueprint(context.Background(), projectRoot, blueprintFS, kitConfig, "commands/wf.yaml", "antigravity", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Expected: .agent/spf-wf.md (Flat, no subfolder)
	expectedPath := filepath.Join(projectRoot, ".agent/spf-wf.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected flat file %s to exist", expectedPath)
	}
}

func TestIntegration_SecurityConstraint(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-int-security-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectRoot := filepath.Join(tmpDir, "project")
	_ = os.MkdirAll(projectRoot, 0755)

	globalTmp, _ := os.MkdirTemp("", "specforce-global-attack-*")
	defer func() { _ = os.RemoveAll(globalTmp) }()

	// Configure 'claude' (non-global) with an absolute path
	kitConfig := &core.KitConfig{
		Tools: map[string]core.ToolRoute{
			"claude": {
				Target: globalTmp,
				Mappings: map[string]core.MappingConfigs{
					"commands": {core.MappingConfig{Path: "spf", Ext: ".md"}},
				},
			},
		},
	}

	kitDir := filepath.Join(tmpDir, "mock-kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "commands"), 0755)
	_ = os.WriteFile(filepath.Join(kitDir, "commands/attack.yaml"), []byte("description: attack"), 0644)
	blueprintFS := os.DirFS(kitDir)

	err = processBlueprint(context.Background(), projectRoot, blueprintFS, kitConfig, "commands/attack.yaml", "claude", installer.Options{})
	
	// Should fail with security error
	if err == nil {
		t.Error("expected processBlueprint to fail for non-global agent with absolute path, but it succeeded")
	}
}



