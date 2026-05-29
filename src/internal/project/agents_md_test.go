package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func TestEnsureAgentsMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Run("Create new file", func(t *testing.T) {
		err := EnsureAgentsMD(tempDir, nil, []string{})
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

		err = EnsureAgentsMD(tempDir, nil, []string{})
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

func TestEnsurePlatformConfigs_NoSelection(t *testing.T) {
	tempDir := t.TempDir()

	err := EnsureAgentsMD(tempDir, nil, []string{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify no directories were created
	for _, dir := range []string{".gemini", ".claude"} {
		path := filepath.Join(tempDir, dir)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("expected directory %s NOT to exist", dir)
		}
	}
}

func TestEnsurePlatformConfigs_Gemini(t *testing.T) {
	t.Run("Creates configs if selected", func(t *testing.T) {
		subTempDir := t.TempDir()
		err := EnsureAgentsMD(subTempDir, nil, []string{"gemini-cli"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		geminiPath := filepath.Join(subTempDir, ".gemini", "settings.json")
		if _, err := os.Stat(geminiPath); err != nil {
			t.Errorf("Gemini settings.json not created: %v", err)
		}
	})

	t.Run("Creates configs if dir exists even if not selected", func(t *testing.T) {
		subTempDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(subTempDir, ".gemini"), 0755); err != nil {
			t.Fatal(err)
		}

		err := EnsureAgentsMD(subTempDir, nil, []string{})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		geminiPath := filepath.Join(subTempDir, ".gemini", "settings.json")
		if _, err := os.Stat(geminiPath); err != nil {
			t.Errorf("Gemini settings.json should be created because directory exists")
		}
	})
}

func TestEnsurePlatformConfigs_Claude(t *testing.T) {
	t.Run("Creates configs if selected", func(t *testing.T) {
		subTempDir := t.TempDir()
		err := EnsureAgentsMD(subTempDir, nil, []string{"claude"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		claudeLink := filepath.Join(subTempDir, ".claude", "rules", "AGENTS.md")
		if _, err := os.Lstat(claudeLink); err != nil {
			t.Errorf("Claude Code symlink not created: %v", err)
		}
	})
}

func TestEnsurePlatformConfigs_Cursor(t *testing.T) {
	t.Run("Does NOT create symlinks for Cursor", func(t *testing.T) {
		subTempDir := t.TempDir()
		err := EnsureAgentsMD(subTempDir, nil, []string{"cursor"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Cursor symlink (should NOT exist)
		cursorLink := filepath.Join(subTempDir, ".cursor", "rules", "AGENTS.md")
		if _, err := os.Lstat(cursorLink); err == nil {
			t.Errorf("Cursor symlink should NOT be created")
		}

		cursorDir := filepath.Join(subTempDir, ".cursor")
		if _, err := os.Stat(cursorDir); err == nil {
			t.Errorf("Cursor directory should NOT be created by EnsureAgentsMD")
		}
	})

	t.Run("Does NOT create symlinks for Cursor even if directory exists", func(t *testing.T) {
		subTempDir := t.TempDir()
		cursorDir := filepath.Join(subTempDir, ".cursor")
		if err := os.MkdirAll(cursorDir, 0755); err != nil {
			t.Fatal(err)
		}

		err := EnsureAgentsMD(subTempDir, nil, []string{})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Cursor symlink (should NOT exist)
		cursorLink := filepath.Join(subTempDir, ".cursor", "rules", "AGENTS.md")
		if _, err := os.Lstat(cursorLink); err == nil {
			t.Errorf("Cursor symlink should NOT be created even if .cursor directory exists")
		}
	})
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
	if !strings.Contains(content, "## 2. CLI Execution & Environment") {
		t.Errorf("content does not contain CLI Execution section")
	}
	if !strings.Contains(content, "## 3. Environment Recovery") {
		t.Errorf("content does not contain Environment Recovery section")
	}

	expectedRecoveryCmd := "npm i -g @jeancodogno/specforce-kit@" + core.Version
	if !strings.Contains(content, expectedRecoveryCmd) {
		t.Errorf("content does not contain the correct recovery command with version: %s", expectedRecoveryCmd)
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
