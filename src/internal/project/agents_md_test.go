package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureAgentsMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Run("Create new file", func(t *testing.T) {
		err := EnsureAgentsMD(tempDir, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		path := filepath.Join(tempDir, "AGENTS.md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("AGENTS.md was not created")
		}

		content, _ := os.ReadFile(path)
		if !strings.Contains(string(content), "<!-- SPECFORCE_AGENTS_START -->") {
			t.Errorf("created file missing markers")
		}
	})

	t.Run("Update existing file", func(t *testing.T) {
		path := filepath.Join(tempDir, "AGENTS.md")
		customContent := "CUSTOM START\n<!-- SPECFORCE_AGENTS_START -->\nOLD\n<!-- SPECFORCE_AGENTS_END -->\nCUSTOM END"
		err := os.WriteFile(path, []byte(customContent), 0644)
		if err != nil {
			t.Fatalf("failed to write existing file: %v", err)
		}

		err = EnsureAgentsMD(tempDir, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		content, _ := os.ReadFile(path)
		if !strings.HasPrefix(string(content), "CUSTOM START") {
			t.Errorf("custom content not preserved")
		}
		if !strings.Contains(string(content), "# AI Agent Collaboration Guide") {
			t.Errorf("managed content not updated")
		}
	})
}

func TestEnsurePlatformConfigs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	err = EnsureAgentsMD(tempDir, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Gemini
	geminiPath := filepath.Join(tempDir, ".gemini", "settings.json")
	if _, err := os.Stat(geminiPath); err != nil {
		t.Errorf("Gemini settings.json not created: %v", err)
	}
	data, _ := os.ReadFile(geminiPath)
	if !strings.Contains(string(data), "\"fileName\": [") {
		t.Errorf("Gemini settings.json content mismatch (expected array): %s", string(data))
	}
	if !strings.Contains(string(data), "\"AGENTS.md\"") {
		t.Errorf("Gemini settings.json missing AGENTS.md: %s", string(data))
	}

	// Antigravity symlink
	agentLink := filepath.Join(tempDir, ".agent", "rules", "AGENTS.md")
	if _, err := os.Lstat(agentLink); err != nil {
		t.Errorf("Antigravity symlink not created: %v", err)
	}
	target, _ := os.Readlink(agentLink)
	if target != "../../AGENTS.md" {
		t.Errorf("Antigravity symlink target mismatch: %s", target)
	}

	// Claude Code symlink
	claudeLink := filepath.Join(tempDir, ".claude", "rules", "AGENTS.md")
	if _, err := os.Lstat(claudeLink); err != nil {
		t.Errorf("Claude Code symlink not created: %v", err)
	}
	target, _ = os.Readlink(claudeLink)
	if target != "../../AGENTS.md" {
		t.Errorf("Claude Code symlink target mismatch: %s", target)
	}
}

func TestGenerateAgentsContent(t *testing.T) {
	content := generateAgentsContent()

	if !strings.Contains(content, "<!-- SPECFORCE_AGENTS_START -->") {
		t.Errorf("content does not contain start marker")
	}
	if !strings.Contains(content, "<!-- SPECFORCE_AGENTS_END -->") {
		t.Errorf("content does not contain end marker")
	}
	if !strings.Contains(content, "# AI Agent Collaboration Guide") {
		t.Errorf("content does not contain title")
	}
}

func TestMergeAgentsContent(t *testing.T) {
	replacement := "NEW CONTENT"

	t.Run("Empty existing", func(t *testing.T) {
		existing := ""
		result := mergeAgentsContent(existing, replacement)
		if result != replacement {
			t.Errorf("expected %q, got %q", replacement, result)
		}
	})

	t.Run("Existing without markers", func(t *testing.T) {
		existing := "CUSTOM CONTENT"
		result := mergeAgentsContent(existing, replacement)
		if !strings.HasPrefix(result, "CUSTOM CONTENT") {
			t.Errorf("custom content not preserved at start")
		}
		if !strings.Contains(result, replacement) {
			t.Errorf("replacement not found")
		}
	})

	t.Run("Existing with markers", func(t *testing.T) {
		existing := "CUSTOM START\n<!-- SPECFORCE_AGENTS_START -->\nOLD CONTENT\n<!-- SPECFORCE_AGENTS_END -->\nCUSTOM END"
		result := mergeAgentsContent(existing, replacement)
		expected := "CUSTOM START\n" + replacement + "\nCUSTOM END"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}
