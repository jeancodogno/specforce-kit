package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates new config", func(t *testing.T) {
		err := EnsureConfigExists(tmpDir)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		configPath := filepath.Join(tmpDir, ".specforce", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("expected config file to be created")
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != DefaultConfigContent {
			t.Errorf("expected default config content, got %s", string(data))
		}
	})

	t.Run("does not overwrite existing config", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, ".specforce", "config.yaml")
		existingContent := "instructions: {}"
		if err := os.WriteFile(configPath, []byte(existingContent), 0644); err != nil {
			t.Fatal(err)
		}

		err := EnsureConfigExists(tmpDir)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != existingContent {
			t.Errorf("expected existing config content to be preserved, got %s", string(data))
		}
	})

	t.Run("returns error on security failure", func(t *testing.T) {
		err := EnsureConfigExists("/nonexistent/root")
		if err == nil {
			t.Fatal("expected error for nonexistent root")
		}
	})
}

func TestLoadConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	specforceDir := filepath.Join(tmpDir, ".specforce")
	if err := os.MkdirAll(specforceDir, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("load hooks from config", func(t *testing.T) {
		configContent := `
hooks:
  on_task_finished:
    - "echo task"
  on_phase_finished:
    - "echo phase"
  on_all_tasks_finished:
    - "echo all"
`
		configPath := filepath.Join(specforceDir, "config.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		config := LoadConfig(tmpDir)
		if len(config.Hooks.OnTaskFinished) == 0 || config.Hooks.OnTaskFinished[0] != "echo task" {
			t.Errorf("expected 'echo task', got %v", config.Hooks.OnTaskFinished)
		}
		if len(config.Hooks.OnPhaseFinished) == 0 || config.Hooks.OnPhaseFinished[0] != "echo phase" {
			t.Errorf("expected 'echo phase', got %v", config.Hooks.OnPhaseFinished)
		}
		if len(config.Hooks.OnAllTasksFinished) == 0 || config.Hooks.OnAllTasksFinished[0] != "echo all" {
			t.Errorf("expected 'echo all', got %v", config.Hooks.OnAllTasksFinished)
		}
	})

	t.Run("load instructions from config", func(t *testing.T) {
		configContent := `
instructions:
  requirements:
    - "Instruction 1"
`
		configPath := filepath.Join(specforceDir, "config.yaml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		config := LoadConfig(tmpDir)
		if len(config.Instructions["requirements"]) == 0 || config.Instructions["requirements"][0] != "Instruction 1" {
			t.Errorf("expected 'Instruction 1', got %v", config.Instructions["requirements"])
		}
	})
}

func TestLoadConfig_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	specforceDir := filepath.Join(tmpDir, ".specforce")
	if err := os.MkdirAll(specforceDir, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("returns empty config if file missing", func(t *testing.T) {
		emptyDir := t.TempDir()
		config := LoadConfig(emptyDir)
		if config == nil {
			t.Fatal("expected non-nil config")
		}
		if len(config.Instructions) != 0 {
			t.Errorf("expected empty instructions, got %v", config.Instructions)
		}
	})

	t.Run("returns empty config if malformed", func(t *testing.T) {
		configPath := filepath.Join(specforceDir, "config.yaml")
		if err := os.WriteFile(configPath, []byte("invalid: yaml: :"), 0644); err != nil {
			t.Fatal(err)
		}

		config := LoadConfig(tmpDir)
		if config == nil {
			t.Fatal("expected non-nil config")
		}
		if len(config.Instructions) != 0 {
			t.Errorf("expected empty instructions, got %v", config.Instructions)
		}
	})
}
