package spec

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func TestCheckTriadArtifacts(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-feature"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	err = os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Case 1: All missing
	ok, missing := CheckTriadArtifacts(tmpDir, slug)
	if ok {
		t.Error("Expected false, got true")
	}
	if len(missing) != 3 {
		t.Errorf("Expected 3 missing, got %d", len(missing))
	}

	// Case 2: Some missing
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	ok, missing = CheckTriadArtifacts(tmpDir, slug)
	if ok {
		t.Error("Expected false, got true")
	}
	if len(missing) != 2 {
		t.Errorf("Expected 2 missing, got %d", len(missing))
	}

	// Case 3: None missing
	if err := os.WriteFile(filepath.Join(specDir, "design.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	ok, missing = CheckTriadArtifacts(tmpDir, slug)
	if !ok {
		t.Error("Expected true, got false")
	}
	if len(missing) != 0 {
		t.Errorf("Expected 0 missing, got %d", len(missing))
	}
}

func TestGetContextFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-feature"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	docsDir := filepath.Join(tmpDir, ".specforce", "docs")
	err = os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "architecture.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := GetContextFiles(tmpDir, slug)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file (from spec dir), got %d. Global docs from .specforce/docs should be excluded.", len(files))
	}

	for _, f := range files {
		if !filepath.IsAbs(f) {
			t.Errorf("Expected absolute path, got %s", f)
		}
	}
}

func setupTasksFile(t *testing.T, tmpDir, slug, content string) {
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	err := os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func getSampleTasksMD() string {
	return `
# Implementation Roadmap

## 1. Execution Strategy
Strategy details here.

## 2. Tasks

### Phase 1: Setup
#### T1.1: Task One
**State:** [FINISHED]
**Target:** target/one
**Context:** context/one

**Action Steps:**
- step 1
- step 2

**Acceptance Check:**
run verify 1

### Phase 2: Implementation
#### T1.2: Task Two
**State:** [IN-PROGRESS]
**Target:** target/two
**Context:** context/two

**Action Steps:**
- step A

**Acceptance Check:**
run verify 2

## 3. Pre-emptive Mitigations
Mitigation details here.
`
}

func TestParseTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-feature"
	setupTasksFile(t, tmpDir, slug, getSampleTasksMD())

	report, err := ParseTasks(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(report.Phases) != 2 {
		t.Errorf("Expected 2 phases, got %d", len(report.Phases))
	}

	if report.Phases[0].Title != "Setup" || report.Phases[1].Title != "Implementation" {
		t.Errorf("Unexpected phase titles: %s, %s", report.Phases[0].Title, report.Phases[1].Title)
	}

	if len(report.Tasks()) != 2 {
		t.Errorf("Expected 2 tasks total, got %d", len(report.Tasks()))
	}

	t1 := report.Phases[0].Tasks[0]
	if t1.ID != "T1.1" || t1.Title != "Task One" {
		t.Errorf("Task 1.1 mismatch: %+v", t1)
	}
	
	if report.ExecutionStrategy != "Strategy details here." {
		t.Errorf("Expected strategy, got %v", report.ExecutionStrategy)
	}
	
	if report.PreemptiveMitigations != "Mitigation details here." {
		t.Errorf("Expected mitigations, got %v", report.PreemptiveMitigations)
	}
}

func TestParseTasks_EdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "edge-feature"
	content := `
# No sections
#### T1.1: Standalone Task
**State:** [PENDING]
`
	setupTasksFile(t, tmpDir, slug, content)

	report, err := ParseTasks(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatal(err)
	}
	
	if report.ExecutionStrategy != "" {
		t.Errorf("Expected empty strategy, got %v", report.ExecutionStrategy)
	}
	
	if len(report.Phases) != 1 || report.Phases[0].Title != "Initial Tasks" {
		t.Errorf("Expected 1 default phase, got %v", len(report.Phases))
	}
}

func TestImplementationReport_Tasks(t *testing.T) {
	report := &ImplementationReport{
		Phases: []Phase{
			{
				ID:    "1",
				Title: "Phase 1",
				Tasks: []ImplementationTask{
					{ID: "T1.1", Title: "Task 1.1"},
					{ID: "T1.2", Title: "Task 1.2"},
				},
			},
			{
				ID:    "2",
				Title: "Phase 2",
				Tasks: []ImplementationTask{
					{ID: "T2.1", Title: "Task 2.1"},
				},
			},
		},
	}

	tasks := report.Tasks()
	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != "T1.1" || tasks[1].ID != "T1.2" || tasks[2].ID != "T2.1" {
		t.Errorf("Unexpected task order or IDs: %v", tasks)
	}
}

func TestParseTasks_HybridFormat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "hybrid-feature"
	content := `
# Implementation Roadmap

### Phase 1: Mixed
#### T1.1: Classic Header
**State:** [FINISHED]
**Target:** src/one.go

- [ ] T1.2: Modern Checklist
**Target:** src/two.go

- [/] T1.3: Working Checklist
**Target:** src/three.go

- [x] T1.4: Done Checklist
**Target:** src/four.go

#### T1.5: Classic with Checkbox
- [x] T1.5: Classic with Checkbox
**Target:** src/five.go
`
	setupTasksFile(t, tmpDir, slug, content)

	report, err := ParseTasks(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	tasks := report.Tasks()
	if len(tasks) != 5 {
		t.Fatalf("Expected 5 tasks, got %d", len(tasks))
	}

	// Verify T1.1 (Classic)
	if tasks[0].ID != "T1.1" || tasks[0].State != "FINISHED" {
		t.Errorf("T1.1 mismatch: ID=%s, State=%s", tasks[0].ID, tasks[0].State)
	}

	// Verify T1.2 (Modern - [ ])
	if tasks[1].ID != "T1.2" || tasks[1].State != "READY" && tasks[1].State != "PENDING" { // It should be READY now because preceding T1.1 is FINISHED
		t.Errorf("T1.2 mismatch: ID=%s, State=%s", tasks[1].ID, tasks[1].State)
	}

	// Verify T1.3 (Modern - [/])
	if tasks[2].ID != "T1.3" || tasks[2].State != "IN-PROGRESS" {
		t.Errorf("T1.3 mismatch: ID=%s, State=%s", tasks[2].ID, tasks[2].State)
	}

	// Verify T1.4 (Modern - [x])
	if tasks[3].ID != "T1.4" || tasks[3].State != "FINISHED" {
		t.Errorf("T1.4 mismatch: ID=%s, State=%s", tasks[3].ID, tasks[3].State)
	}
}

func TestImplementationTaskJSON(t *testing.T) {
	task := ImplementationTask{
		ID:           "T1.1",
		Title:        "Test Task",
		ParallelWith: []string{"T1.2", "T1.3"},
		IsParallel:   true,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	expectedParallelWith := `"parallel_with":["T1.2","T1.3"]`
	expectedIsParallel := `"is_parallel":true`

	strData := string(data)
	if !strings.Contains(strData, expectedParallelWith) {
		t.Errorf("JSON output missing or incorrect parallel_with field. Got: %s", strData)
	}
	if !strings.Contains(strData, expectedIsParallel) {
		t.Errorf("JSON output missing or incorrect is_parallel field. Got: %s", strData)
	}

	var unmarshaled ImplementationTask
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaled.IsParallel != true {
		t.Errorf("Expected IsParallel to be true, got %v", unmarshaled.IsParallel)
	}
	if len(unmarshaled.ParallelWith) != 2 || unmarshaled.ParallelWith[0] != "T1.2" || unmarshaled.ParallelWith[1] != "T1.3" {
		t.Errorf("Expected ParallelWith [T1.2, T1.3], got %v", unmarshaled.ParallelWith)
	}
}

func TestParseTaskBlock_Parallel(t *testing.T) {
	content := []byte(`### T1.2 - [PARSER] Implement Parallel Metadata Extraction
**State:** [PENDING]
**Target:** src/internal/spec/implementation.go
**Context:** [REQ-1](requirements.md)
**Parallel With:** T1.1, T1.2
**Acceptance Check:**
- Ensure tests pass.

- Extract ParallelWith.
`)
	tm := []int{
		0, len(content),
		0, 3,
		4, 8,
		11, strings.Index(string(content), "\n"),
	}
	
	task := parseTaskBlock(content, tm, [][]int{}, [][]int{}, nil)

	if task.ID != "T1.2" {
		t.Errorf("Expected ID T1.2, got %s", task.ID)
	}
	if !task.IsParallel {
		t.Errorf("Expected IsParallel to be true")
	}
	if len(task.ParallelWith) != 2 {
		t.Fatalf("Expected 2 parallel tasks, got %d: %v", len(task.ParallelWith), task.ParallelWith)
	}
	if task.ParallelWith[0] != "T1.1" {
		t.Errorf("Expected T1.1, got %s", task.ParallelWith[0])
	}
	if task.ParallelWith[1] != "T1.2" {
		t.Errorf("Expected T1.2, got %s", task.ParallelWith[1])
	}
}

func TestParseTasks_ParallelTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "parallel-feature"
	content := `
# Implementation Roadmap

### Phase 1: Setup
#### T1.1: Task One
**State:** [FINISHED]
**Target:** target/one

#### T1.2: Task Two
**State:** [PENDING]
**Target:** target/two
**Parallel With:** T1.1, T1.3

#### T1.3: Task Three
**State:** [PENDING]
**Target:** target/three
**Parallel With:** T1.1, T1.2

#### T1.4: Task Four
**State:** [PENDING]
**Target:** target/four
`
	setupTasksFile(t, tmpDir, slug, content)

	report, err := ParseTasks(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	tasks := report.Tasks()
	if len(tasks) != 4 {
		t.Fatalf("Expected 4 tasks, got %d", len(tasks))
	}

	if tasks[1].State != "READY" {
		t.Errorf("Expected T1.2 to be READY, got %s", tasks[1].State)
	}
	if tasks[2].State != "READY" {
		t.Errorf("Expected T1.3 to be READY, got %s", tasks[2].State)
	}
	if tasks[3].State != "PENDING" {
		t.Errorf("Expected T1.4 to be PENDING because T1.2/T1.3 are not finished, got %s", tasks[3].State)
	}
}

func createHookScript(t *testing.T, tmpDir, name, logPath string) string {
	t.Helper()
	scriptPath := filepath.Join(tmpDir, name)
	scriptContent := fmt.Sprintf("#!/bin/sh\necho %s >> %s\n", name, logPath)
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write hook script: %v", err)
	}
	return scriptPath
}

func countLogLines(t *testing.T, logPath, match string) int {
	t.Helper()
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("failed to read log: %v", err)
	}
	count := 0
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == match {
			count++
		}
	}
	return count
}

func TestService_UpdateTaskStatus_BatchHooks(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "batch-hooks-slug"
	tasksDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(tasksDir, 0755)

	content := "### Phase 1: Core\n#### T1.1: Task 1\n**State:** [PENDING]\n#### T1.2: Task 2\n**State:** [PENDING]\n#### T1.3: Task 3\n**State:** [PENDING]"
	_ = os.WriteFile(filepath.Join(tasksDir, "tasks.md"), []byte(content), 0644)

	logPath := filepath.Join(tmpDir, "hook.log")
	scriptPath := createHookScript(t, tmpDir, "task_finished.sh", logPath)

	config := &core.ProjectConfig{
		Hooks: core.HooksConfig{
			OnTaskFinished: []string{scriptPath},
		},
	}
	svc := NewService(nil, &mockConfigProvider{config: config})

	err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1", "T1.2"}, "finished")
	if err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}

	if count := countLogLines(t, logPath, "task_finished.sh"); count != 1 {
		t.Errorf("expected OnTaskFinished to execute 1 time, got %d", count)
	}

	updatedContent, _ := os.ReadFile(filepath.Join(tasksDir, "tasks.md"))
	sContent := string(updatedContent)
	if !strings.Contains(sContent, "#### T1.1: Task 1\n**State:** [FINISHED]") {
		t.Errorf("T1.1 not marked FINISHED")
	}
	if !strings.Contains(sContent, "#### T1.2: Task 2\n**State:** [FINISHED]") {
		t.Errorf("T1.2 not marked FINISHED")
	}
	if !strings.Contains(sContent, "#### T1.3: Task 3\n**State:** [PENDING]") {
		t.Errorf("T1.3 should remain PENDING")
	}
}

func TestService_UpdateTaskStatus_PhaseAndSpecHooks(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "phase-spec-hooks-slug"
	tasksDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(tasksDir, 0755)

	content := "### Phase 1: Core\n#### T1.1: Task 1\n**State:** [PENDING]\n#### T1.2: Task 2\n**State:** [PENDING]\n### Phase 2: Polish\n#### T2.1: Task 3\n**State:** [PENDING]\n#### T2.2: Task 4\n**State:** [PENDING]"
	_ = os.WriteFile(filepath.Join(tasksDir, "tasks.md"), []byte(content), 0644)

	logPath := filepath.Join(tmpDir, "hook.log")
	taskHook := createHookScript(t, tmpDir, "task.sh", logPath)
	phaseHook := createHookScript(t, tmpDir, "phase.sh", logPath)
	specHook := createHookScript(t, tmpDir, "spec.sh", logPath)

	config := &core.ProjectConfig{
		Hooks: core.HooksConfig{
			OnTaskFinished:     []string{taskHook},
			OnPhaseFinished:    []string{phaseHook},
			OnAllTasksFinished: []string{specHook},
		},
	}
	svc := NewService(nil, &mockConfigProvider{config: config})

	// 1. Batch without last task in phase (T1.1)
	if err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1"}, "finished"); err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}
	if count := countLogLines(t, logPath, "task.sh"); count != 1 {
		t.Errorf("expected 1 task.sh, got %d", count)
	}
	if count := countLogLines(t, logPath, "phase.sh"); count != 0 {
		t.Errorf("expected 0 phase.sh, got %d", count)
	}

	// 2. Batch including last task in phase 1 (T1.2)
	if err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.2"}, "finished"); err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}
	if count := countLogLines(t, logPath, "phase.sh"); count != 1 {
		t.Errorf("expected 1 phase.sh, got %d", count)
	}
	if count := countLogLines(t, logPath, "spec.sh"); count != 0 {
		t.Errorf("expected 0 spec.sh, got %d", count)
	}

	// 3. Batch including last task in phase 2 and last in spec (T2.1, T2.2)
	if err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T2.1", "T2.2"}, "finished"); err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}
	if count := countLogLines(t, logPath, "spec.sh"); count != 1 {
		t.Errorf("expected 1 spec.sh, got %d", count)
	}
}

func TestService_UpdateTaskStatus_HookFailureRollback(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "rollback-slug"
	tasksDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(tasksDir, 0755)

	content := "### Phase 1: Core\n#### T1.1: Task 1\n**State:** [PENDING]\n#### T1.2: Task 2\n**State:** [PENDING]"
	_ = os.WriteFile(filepath.Join(tasksDir, "tasks.md"), []byte(content), 0644)

	config := &core.ProjectConfig{
		Hooks: core.HooksConfig{
			OnTaskFinished: []string{"false"},
		},
	}
	svc := NewService(nil, &mockConfigProvider{config: config})

	err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1", "T1.2"}, "finished")
	if err == nil {
		t.Fatal("expected error from failing hook, got nil")
	}

	// Verify tasks.md was not mutated
	updatedContent, _ := os.ReadFile(filepath.Join(tasksDir, "tasks.md"))
	if string(updatedContent) != content {
		t.Errorf("tasks.md was modified despite hook failure:\n%s", string(updatedContent))
	}
}
