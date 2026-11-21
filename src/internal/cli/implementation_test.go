package cli

import (
	"context"
	"fmt"
	"io"
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

	err := executor.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.1"}, "finished")
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

	err := executor.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.1"}, "finished")
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

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = orig

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

const batchCLITasksMD = `
### Phase 1: Core
#### T1.1: Task 1
**State:** [PENDING]

#### T1.2: Task 2
**State:** [PENDING]

#### T1.3: Task 3
**State:** [PENDING]
`

func TestImplementationUpdate_Batch(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	slug := "batch-cli-test"
	specDir := filepath.Join(".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(batchCLITasksMD), 0644)

	executor := &Executor{Version: "1.0.0"}
	ui := tui.NewUI()

	t.Run("batch update multiple tasks and stdout feedback", func(t *testing.T) {
		output := captureStdout(t, func() {
			err := executor.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.1", "T1.2"}, "finished")
			if err != nil {
				t.Fatalf("HandleImplementationUpdate failed: %v", err)
			}
		})

		expectedMsg := fmt.Sprintf("[OK] Tasks T1.1, T1.2 status updated to finished in %s\n", slug)
		if output != expectedMsg {
			t.Errorf("expected stdout %q, got %q", expectedMsg, output)
		}

		content, _ := os.ReadFile(filepath.Join(specDir, "tasks.md"))
		if !strings.Contains(string(content), "#### T1.1: Task 1\n**State:** [FINISHED]") {
			t.Errorf("T1.1 not finished")
		}
		if !strings.Contains(string(content), "#### T1.2: Task 2\n**State:** [FINISHED]") {
			t.Errorf("T1.2 not finished")
		}
	})

	t.Run("single task stdout feedback", func(t *testing.T) {
		output := captureStdout(t, func() {
			err := executor.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.3"}, "finished")
			if err != nil {
				t.Fatalf("HandleImplementationUpdate failed: %v", err)
			}
		})

		expectedMsg := fmt.Sprintf("[OK] Task T1.3 status updated to finished in %s\n", slug)
		if output != expectedMsg {
			t.Errorf("expected stdout %q, got %q", expectedMsg, output)
		}
	})
}

const qaValidJSON = `{"criteria":[{"id":"US-1.1","scenario":"s","evidence":"e","verdict":"PASSED"},{"id":"US-1.2","scenario":"s2","evidence":"e2","verdict":"NOT-VERIFIABLE"}],"findings":[{"title":"f"}]}`

func qaProject(t *testing.T, slug string) string {
	t.Helper()
	root := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.MkdirAll(filepath.Join(root, ".specforce", "specs", slug), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestHandleImplementationQARecord(t *testing.T) {
	slug := "qa-demo"
	e := &Executor{Version: "1.0.0"}

	t.Run("records and prints success line", func(t *testing.T) {
		root := qaProject(t, slug)
		var err error
		out := captureStdout(t, func() {
			err = e.HandleImplementationQARecord(context.Background(), tui.NewUI(), slug, strings.NewReader(qaValidJSON))
		})
		if err != nil {
			t.Fatalf("record failed: %v", err)
		}
		want := "[OK] QA report recorded for qa-demo: PASSED (1 passed, 0 failed, 1 not-verifiable, 1 findings, run 1)"
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in %q", want, out)
		}
		data, _ := os.ReadFile(filepath.Join(root, ".specforce", "specs", slug, "spec.yaml"))
		if !strings.Contains(string(data), "qa:") || !strings.Contains(string(data), "US-1.1") {
			t.Errorf("spec.yaml missing qa block: %s", data)
		}
	})

	t.Run("invalid inputs are refused and stored report survives", func(t *testing.T) {
		root := qaProject(t, slug)
		if err := e.HandleImplementationQARecord(context.Background(), tui.NewUI(), slug, strings.NewReader(qaValidJSON)); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, ".specforce", "specs", slug, "spec.yaml")
		before, _ := os.ReadFile(path)
		bad := map[string]string{
			"empty":   "",
			"badjson": "{",
			"verdict": `{"criteria":[{"id":"US-9","scenario":"s","evidence":"e","verdict":"MAYBE"}]}`,
			"zero":    `{"criteria":[]}`,
			"unknown": `{"criteria":[{"id":"a","scenario":"s","evidence":"e","verdict":"PASSED"}],"x":1}`,
		}
		for name, in := range bad {
			err := e.HandleImplementationQARecord(context.Background(), tui.NewUI(), slug, strings.NewReader(in))
			if err == nil {
				t.Errorf("%s: expected error", name)
				continue
			}
			if name == "verdict" && !strings.Contains(err.Error(), "US-9") {
				t.Errorf("verdict error must name the criterion: %v", err)
			}
		}
		after, _ := os.ReadFile(path)
		if string(before) != string(after) {
			t.Error("invalid call changed the stored report")
		}
	})

	t.Run("unknown slug", func(t *testing.T) {
		qaProject(t, slug)
		if err := e.HandleImplementationQARecord(context.Background(), tui.NewUI(), "nope", strings.NewReader(qaValidJSON)); err == nil {
			t.Error("expected error for unknown slug")
		}
	})
}

func TestHandleImplementationUpdate_RedNotice(t *testing.T) {
	slug := "red-notice"
	setup := func(t *testing.T) {
		root := qaProject(t, slug)
		tasks := "### T1.1: [RED] tests\n**State:** [PENDING]\n\n### T1.2: [GREEN] code\n**State:** [PENDING]\n"
		writeCLIFile(t, filepath.Join(root, ".specforce", "specs", slug, "tasks.md"), tasks)
		writeCLIFile(t, filepath.Join(root, ".specforce", "config.yaml"), "hooks:\n  on_task_finished:\n    - \"true\"\n")
	}
	e := &Executor{Version: "1.0.0"}

	t.Run("RED only emits warning before OK", func(t *testing.T) {
		setup(t)
		ui := &archiveMockUI{}
		var err error
		out := captureStdout(t, func() {
			err = e.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.1"}, "finished")
		})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(ui.warns, ""), "RED task: skipped per-task hooks: true") {
			t.Errorf("expected RED warning, got %v", ui.warns)
		}
		if !strings.Contains(out, "[OK] Task T1.1") {
			t.Errorf("expected OK line, got %q", out)
		}
	})

	t.Run("GREEN emits no warning", func(t *testing.T) {
		setup(t)
		ui := &archiveMockUI{}
		if err := e.HandleImplementationUpdate(context.Background(), ui, slug, []string{"T1.2"}, "finished"); err != nil {
			t.Fatal(err)
		}
		if len(ui.warns) != 0 {
			t.Errorf("unexpected warnings: %v", ui.warns)
		}
	})
}
