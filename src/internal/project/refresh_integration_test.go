package project_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jeancodogno/specforce-kit/src/internal/project"
)

func createMockKitFS() fstest.MapFS {
	return fstest.MapFS{
		"kit.yaml": &fstest.MapFile{Data: []byte(`
tools:
  gemini:
    name: "Gemini"
    target: ".gemini/"
    mappings:
      agents: { path: "agents", ext: ".md" }
      skills: { path: "skills", ext: ".md" }
      commands: { path: "commands", ext: ".md" }
`)},
		"agents/coder.yaml": &fstest.MapFile{Data: []byte("content: kit-agent-content")},
		"skills/tdd.yaml":   &fstest.MapFile{Data: []byte("content: kit-skill-content")},
		"commands/run.yaml": &fstest.MapFile{Data: []byte("content: kit-command-content")},
	}
}

func TestService_UpdateTools_RefreshesAllCategories(t *testing.T) {
	kitFS := createMockKitFS()
	artifactsFS := fstest.MapFS{}
	tmpDir := t.TempDir()

	svc := project.NewService(kitFS, artifactsFS, tmpDir)
	ui := &mockUI{}

	// 1. Initial Installation
	err := svc.UpdateTools(context.Background(), ui, []string{"gemini"})
	if err != nil {
		t.Fatalf("Initial UpdateTools failed: %v", err)
	}

	// Paths
	agentPath := filepath.Join(tmpDir, ".gemini/agents/coder.md")
	skillPath := filepath.Join(tmpDir, ".gemini/skills/tdd.md")
	commandPath := filepath.Join(tmpDir, ".gemini/commands/run.md")

	// 2. Modify files manually (simulate user customization or old version)
	filesToModify := []string{agentPath, skillPath, commandPath}
	for _, p := range filesToModify {
		if err := os.WriteFile(p, []byte("user-modified-content"), 0644); err != nil {
			t.Fatalf("failed to modify file %s: %v", p, err)
		}
	}

	// 3. Run Update
	err = svc.UpdateTools(context.Background(), ui, []string{"gemini"})
	if err != nil {
		t.Fatalf("Second UpdateTools failed: %v", err)
	}

	// 4. Verify Refresh
	tests := []struct {
		path     string
		expected string
	}{
		{agentPath, "kit-agent-content"},
		{skillPath, "kit-skill-content"},
		{commandPath, "kit-command-content"},
	}

	for _, tt := range tests {
		content, err := os.ReadFile(tt.path)
		if err != nil {
			t.Errorf("failed to read file %s: %v", tt.path, err)
			continue
		}
		// Note: AdaptArtifacts might inject headers, so we check if it contains the kit content
		if !strings.Contains(string(content), tt.expected) {
			t.Errorf("file %s was not refreshed. Got: %s, Expected to contain: %s", tt.path, string(content), tt.expected)
		}
	}
}
