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

func setupLegacyWfCmdGemini(t *testing.T, tmpDir string) (string, string, string) {
	t.Helper()
	wfDir := filepath.Join(tmpDir, ".agents", "workflows")
	if err := os.MkdirAll(wfDir, 0755); err != nil {
		t.Fatalf("failed to create workflows dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wfDir, "spf-discovery.md"), []byte("# Discovery"), 0600); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}

	cmdDir := filepath.Join(tmpDir, ".claude", "commands")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatalf("failed to create commands dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cmdDir, "spf-spec.md"), []byte("# Spec"), 0600); err != nil {
		t.Fatalf("failed to write command file: %v", err)
	}

	geminiDir := filepath.Join(tmpDir, ".gemini")
	if err := os.MkdirAll(geminiDir, 0755); err != nil {
		t.Fatalf("failed to create gemini dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), []byte("{}"), 0600); err != nil {
		t.Fatalf("failed to write gemini settings: %v", err)
	}
	return wfDir, cmdDir, geminiDir
}

func TestDetectLegacyAssets_WorkflowsCommandsGemini(t *testing.T) {
	tmpDir := t.TempDir()
	wfDir, cmdDir, geminiDir := setupLegacyWfCmdGemini(t, tmpDir)

	assets, err := DetectLegacyAssets(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 3 {
		t.Fatalf("expected 3 legacy assets, got %d: %v", len(assets), assets)
	}

	expectedPaths := map[string]bool{
		filepath.Clean(wfDir):     false,
		filepath.Clean(cmdDir):    false,
		filepath.Clean(geminiDir): false,
	}
	for _, a := range assets {
		if _, ok := expectedPaths[a]; ok {
			expectedPaths[a] = true
		}
	}
	for p, found := range expectedPaths {
		if !found {
			t.Errorf("expected legacy asset %s to be detected", p)
		}
	}

	ui := &legacyMockUI{confirmResult: true}
	if err := PromptAndCleanupLegacyAssets(tmpDir, ui); err != nil {
		t.Fatalf("PromptAndCleanupLegacyAssets failed: %v", err)
	}
	if !ui.confirmCalled {
		t.Errorf("expected Confirm to be called")
	}

	for _, dir := range []string{wfDir, cmdDir, geminiDir} {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("expected %s to be deleted, got err: %v", dir, err)
		}
	}
}

