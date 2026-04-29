package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
	"gopkg.in/yaml.v3"
)

func TestDualInstallation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-dual-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	_ = os.MkdirAll(filepath.Join(kitDir, "commands"), 0755)

	blueprintYAML := "description: Test Command\ncontent: |\n  # Hello\n"
	_ = os.WriteFile(filepath.Join(kitDir, "commands/test-cmd.yaml"), []byte(blueprintYAML), 0644)

	kitFS := os.DirFS(kitDir)

	// kit.yaml with dual mapping for commands
	kitYAML := `
tools:
  test-agent:
    target: ".test/"
    mappings:
      commands:
        - path: "cmds"
          ext: ".md"
        - path: "skills/spf-*"
          name: "SKILL"
          ext: ".md"
`
	var kitConfig core.KitConfig
	if err := yaml.Unmarshal([]byte(kitYAML), &kitConfig); err != nil {
		t.Fatalf("failed to unmarshal kit.yaml: %v", err)
	}

	projectRoot := filepath.Join(tmpDir, "project")

	err = processBlueprint(context.Background(), projectRoot, kitFS, &kitConfig, "commands/test-cmd.yaml", "test-agent", installer.Options{})
	if err != nil {
		t.Fatalf("processBlueprint failed: %v", err)
	}

	// Verify both files exist
	expectedCmd := filepath.Join(projectRoot, ".test/cmds/test-cmd.md")
	expectedSkill := filepath.Join(projectRoot, ".test/skills/spf-test-cmd/SKILL.md")

	if _, err := os.Stat(expectedCmd); os.IsNotExist(err) {
		t.Errorf("expected command file %s to exist", expectedCmd)
	}
	if _, err := os.Stat(expectedSkill); os.IsNotExist(err) {
		t.Errorf("expected skill file %s to exist", expectedSkill)
	}
}
