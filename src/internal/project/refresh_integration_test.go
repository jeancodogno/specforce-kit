package project_test

import (
	"context"
	"fmt"
	"io/fs"
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
  claude:
    name: "Claude"
    target: ".claude/"
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

// listToolFiles describes every file and directory under the tool folder, so
// a failure shows what the installation actually created.
func listToolFiles(root string) string {
	var found []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			found = append(found, fmt.Sprintf("%s (walk error: %v)", p, err))
			return nil
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			rel = p
		}
		if d.IsDir() {
			rel += string(filepath.Separator)
		}
		found = append(found, rel)
		return nil
	})
	if err != nil {
		found = append(found, fmt.Sprintf("(walk failed: %v)", err))
	}
	if len(found) == 0 {
		return fmt.Sprintf("nothing under %s", root)
	}
	return fmt.Sprintf("files under %s: [%s]", root, strings.Join(found, ", "))
}

func TestService_UpdateTools_RefreshesAllCategories(t *testing.T) {
	kitFS := createMockKitFS()
	artifactsFS := fstest.MapFS{}
	tmpDir := t.TempDir()

	svc := project.NewService(kitFS, artifactsFS, tmpDir)
	ui := &mockUI{}

	// 1. Initial Installation
	err := svc.UpdateTools(context.Background(), ui, []string{"claude"})
	if err != nil {
		t.Fatalf("Initial UpdateTools failed: %v", err)
	}

	// Paths
	agentPath := filepath.Join(tmpDir, ".claude/agents/coder.md")
	skillPath := filepath.Join(tmpDir, ".claude/skills/tdd.md")
	commandPath := filepath.Join(tmpDir, ".claude/commands/run.md")

	// 2. Modify files manually (simulate user customization or old version)
	filesToModify := []string{agentPath, skillPath, commandPath}
	for _, p := range filesToModify {
		if err := os.WriteFile(p, []byte("user-modified-content"), 0644); err != nil {
			t.Fatalf("failed to modify file %s: %v; %s", p, err, listToolFiles(filepath.Join(tmpDir, ".claude")))
		}
	}

	// 3. Run Update
	err = svc.UpdateTools(context.Background(), ui, []string{"claude"})
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
			t.Errorf("failed to read file %s: %v; %s", tt.path, err, listToolFiles(filepath.Join(tmpDir, ".claude")))
			continue
		}
		// Note: AdaptArtifacts might inject headers, so we check if it contains the kit content
		if !strings.Contains(string(content), tt.expected) {
			t.Errorf("file %s was not refreshed. Got: %s, Expected to contain: %s", tt.path, string(content), tt.expected)
		}
	}
}
