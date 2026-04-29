package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/installer"
)

func TestSkillNameBug(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-bug-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	commandsDir := filepath.Join(kitDir, "commands")
	_ = os.MkdirAll(commandsDir, 0755)

	// Command with empty metadata name
	commandYAML := "description: Test Command Description\ncontent: |\n  # Command Content\n"
	_ = os.WriteFile(filepath.Join(commandsDir, "archive.yaml"), []byte(commandYAML), 0644)
	
	// kit.yaml
	kitYAML := "tools:\n  kimi-code:\n    target: .kimi/\n    mappings:\n      commands:\n        - path: skills/spf-*\n          name: SKILL\n          ext: .md\n"
	_ = os.WriteFile(filepath.Join(kitDir, "kit.yaml"), []byte(kitYAML), 0644)

	kitFS := os.DirFS(kitDir)
	projectRoot := tmpDir

	err = AdaptArtifacts(context.Background(), projectRoot, kitFS, "kimi-code", nil, installer.Options{})
	if err != nil {
		t.Fatalf("AdaptArtifacts failed: %v", err)
	}

	// Filename is SKILL.md as expected
	cmdMD := filepath.Join(projectRoot, ".kimi/skills/spf-archive/SKILL.md")
	
	data, err := os.ReadFile(cmdMD)
	if err != nil {
		t.Fatalf("failed to read generated file %s: %v", cmdMD, err)
	}
	content := string(data)
	
	if strings.Contains(content, "name: SKILL") {
		t.Errorf("BUG STILL PRESENT: Header name is 'SKILL', expected a unique name (e.g., spf.archive). Content: %q", content)
	}
	if !strings.Contains(content, "name: spf.archive") {
		t.Errorf("Expected name 'spf.archive' in header. Content: %q", content)
	}
}

func TestHeaderNameUniqueness(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-unique-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kitDir := filepath.Join(tmpDir, "kit")
	commandsDir := filepath.Join(kitDir, "commands")
	_ = os.MkdirAll(commandsDir, 0755)

	// Two different blueprints
	_ = os.WriteFile(filepath.Join(commandsDir, "cmd1.yaml"), []byte("name: duplicate\ncontent: 1"), 0644)
	_ = os.WriteFile(filepath.Join(commandsDir, "cmd2.yaml"), []byte("name: duplicate\ncontent: 2"), 0644)
	
	// kit.yaml
	kitYAML := "tools:\n  test-agent:\n    target: .test/\n    mappings:\n      commands:\n        - path: cmds\n          ext: .md\n"
	_ = os.WriteFile(filepath.Join(kitDir, "kit.yaml"), []byte(kitYAML), 0644)

	kitFS := os.DirFS(kitDir)
	projectRoot := tmpDir

	// This should fail because both blueprints results in 'name: duplicate'
	err = AdaptArtifacts(context.Background(), projectRoot, kitFS, "test-agent", nil, installer.Options{})
	if err == nil {
		t.Error("Expected error due to duplicate header names, but got nil")
	} else if !strings.Contains(err.Error(), "duplicate header name") {
		t.Errorf("Expected 'duplicate header name' error, got: %v", err)
	}
}
