package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func setupTestProject(t *testing.T, tmpDir, slug string) {
	specDir := filepath.Join(".specforce", "specs", slug)
	err := os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	tasksMD := `
# Implementation Roadmap

## 1. Execution Strategy
Strategy

## 2. Tasks

### T1.1: [SCAFFOLD] Task 1
**State:** [PENDING]
**Target:** target/1
**Context:** context/1

**Action Steps:**
- step 1

**Verification (TDD):**
verify 1

## 3. Pre-emptive Mitigations
Mitigation
`
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(tasksMD), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("req"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "design.md"), []byte("design"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestHandleImplementationStatus(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	slug := "0015-test"
	setupTestProject(t, tmpDir, slug)

	executor := NewExecutor("1.0.0")
	executor.DevMode = false
	ui := tui.NewUI()

	t.Run("TUI Mode", func(t *testing.T) {
		err := executor.HandleImplementation(context.Background(), ui, "status", slug)
		if err != nil {
			t.Errorf("HandleImplementation failed: %v", err)
		}
	})

	t.Run("JSON Mode", func(t *testing.T) {
		err := executor.HandleImplementation(context.Background(), ui, "status", slug, "--json")
		if err != nil {
			t.Errorf("HandleImplementation failed: %v", err)
		}
	})
}

func TestHandleImplementationUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	t.Run("Update status to finished", testUpdateStatusToFinished)
	t.Run("Hook failure blocks update and returns error", testHookFailureBlocksUpdate)
}

func testUpdateStatusToFinished(t *testing.T) {
	slug := "0015-test-update"
	specDir := filepath.Join(".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)

	tasksMD := `
### T1.1: Task 1
**State:** [PENDING]
`
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(tasksMD), 0644)

	executor := &Executor{Version: "1.0.0"}
	ui := tui.NewUI()

	err := executor.HandleImplementationUpdate(context.Background(), ui, slug, "T1.1", "finished")
	if err != nil {
		t.Errorf("HandleImplementationUpdate failed: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(specDir, "tasks.md"))
	if !strings.Contains(string(content), "**State:** [FINISHED]") {
		t.Errorf("Expected status to be [FINISHED], got:\n%s", string(content))
	}
}

func testHookFailureBlocksUpdate(t *testing.T) {
	slug := "hook-fail-test"
	specDir := filepath.Join(".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)

	tasksMD := `
### T1.1: Task 1
**State:** [PENDING]
`
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(tasksMD), 0644)

	configContent := `
hooks:
  on_task_finished:
    - "false"
`
	_ = os.MkdirAll(".specforce", 0755)
	_ = os.WriteFile(filepath.Join(".specforce", "config.yaml"), []byte(configContent), 0644)

	executor := &Executor{Version: "1.0.0"}
	ui := tui.NewUI()

	err := executor.HandleImplementationUpdate(context.Background(), ui, slug, "T1.1", "finished")
	if err == nil {
		t.Errorf("Expected error from HandleImplementationUpdate due to hook failure, got nil")
	}

	if !strings.Contains(err.Error(), "hook failures") {
		t.Errorf("Expected error message to contain 'hook failures', got: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(specDir, "tasks.md"))
	if !strings.Contains(string(content), "**State:** [PENDING]") {
		t.Errorf("Expected status to remain [PENDING], got:\n%s", string(content))
	}
}
