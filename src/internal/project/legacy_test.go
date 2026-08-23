package project

import (
	"os"
	"path/filepath"
	"testing"
)

type legacyMockUI struct {
	confirmResult bool
	confirmCalled bool
	subTasks      []string
}

func (m *legacyMockUI) Log(message string)        {}
func (m *legacyMockUI) Warn(message string)       {}
func (m *legacyMockUI) Error(message string)      {}
func (m *legacyMockUI) Success(message string)    {}
func (m *legacyMockUI) SubTask(message string)    {}
func (m *legacyMockUI) LogSubTask(message string) { m.subTasks = append(m.subTasks, message) }
func (m *legacyMockUI) StartSpinner(message string) {}
func (m *legacyMockUI) StopSpinner()               {}
func (m *legacyMockUI) Confirm(question string) bool {
	m.confirmCalled = true
	return m.confirmResult
}

func TestDetectLegacyAssets_Empty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-legacy-empty-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	assets, err := DetectLegacyAssets(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 0 {
		t.Errorf("expected 0 legacy assets, got %d: %v", len(assets), assets)
	}
}

func TestDetectLegacyAssets_Found(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-legacy-found-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create legacy skill in .agents/skills/pragmatic-product-owner
	skillDir := filepath.Join(tmpDir, ".agents", "skills", "pragmatic-product-owner")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skillDir: %v", err)
	}

	// Create legacy agent in .cursor/agents/specforce-architect.md
	cursorAgentsDir := filepath.Join(tmpDir, ".cursor", "agents")
	if err := os.MkdirAll(cursorAgentsDir, 0755); err != nil {
		t.Fatalf("failed to create cursorAgentsDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cursorAgentsDir, "specforce-architect.md"), []byte("# Specforce Architect"), 0600); err != nil {
		t.Fatalf("failed to write agent file: %v", err)
	}

	assets, err := DetectLegacyAssets(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 legacy assets, got %d: %v", len(assets), assets)
	}
}

func TestPromptAndCleanupLegacyAssets_Declined(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-legacy-decline-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	skillDir := filepath.Join(tmpDir, ".agents", "skills", "task-atomic-decomposition")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skillDir: %v", err)
	}

	ui := &legacyMockUI{confirmResult: false}
	err = PromptAndCleanupLegacyAssets(tmpDir, ui)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ui.confirmCalled {
		t.Errorf("expected Confirm to be called")
	}

	// Asset should still exist
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		t.Errorf("expected legacy asset to be preserved when user declines")
	}
}

func TestPromptAndCleanupLegacyAssets_Accepted(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-legacy-accept-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	skillDir := filepath.Join(tmpDir, ".agents", "skills", "opportunity-framing")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skillDir: %v", err)
	}

	ui := &legacyMockUI{confirmResult: true}
	err = PromptAndCleanupLegacyAssets(tmpDir, ui)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ui.confirmCalled {
		t.Errorf("expected Confirm to be called")
	}

	// Asset should be deleted
	if _, err := os.Stat(skillDir); !os.IsNotExist(err) {
		t.Errorf("expected legacy asset to be deleted when user accepts")
	}
}
