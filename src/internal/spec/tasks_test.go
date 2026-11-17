package spec

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func setupTasksTest(t *testing.T, projectRoot, slug, content string) string {
	tasksDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	tasksPath := filepath.Join(tasksDir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}
	return tasksPath
}

func TestUpdateTaskStatusFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-tasks-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "test-slug"
	content := `
## 2. Tasks

### T1.1: Task 1
**State:** [ ]
**Target:** target1

### T1.2: Task 2
**State:** [PENDING]
**Target:** target2
`
	tasksPath := setupTasksTest(t, projectRoot, slug, content)

	// Test updating T1.1 to finished
	if err := updateTaskStatusFile(projectRoot, slug, "T1.1", "finished"); err != nil {
		t.Fatalf("updateTaskStatusFile failed: %v", err)
	}

	updatedContent, _ := os.ReadFile(tasksPath)
	expected := `
## 2. Tasks

### T1.1: Task 1
**State:** [FINISHED]
**Target:** target1

### T1.2: Task 2
**State:** [PENDING]
**Target:** target2
`
	if string(updatedContent) != expected {
		t.Errorf("Unexpected content after update:\n%s", string(updatedContent))
	}

	// Test updating T1.2 to in-progress
	if err := updateTaskStatusFile(projectRoot, slug, "T1.2", "in-progress"); err != nil {
		t.Fatalf("updateTaskStatusFile failed: %v", err)
	}

	updatedContent, _ = os.ReadFile(tasksPath)
	expected2 := `
## 2. Tasks

### T1.1: Task 1
**State:** [FINISHED]
**Target:** target1

### T1.2: Task 2
**State:** [IN-PROGRESS]
**Target:** target2
`
	if string(updatedContent) != expected2 {
		t.Errorf("Unexpected content after update 2:\n%s", string(updatedContent))
	}
}

func TestUpdateTaskStatusFile_WithPhases(t *testing.T) {
	projectRoot, _ := os.MkdirTemp("", "specforce-update-tasks-*")
	defer func() { _ = os.RemoveAll(projectRoot) }()

	slug := "test-slug"
	tasksDir := filepath.Join(projectRoot, ".specforce", "specs", slug)
	_ = os.MkdirAll(tasksDir, 0755)
	tasksPath := filepath.Join(tasksDir, "tasks.md")

	content := `
## 2. Tasks

### Phase 1: Setup
#### T1.1: Task One
**State:** [PENDING]
**Target:** target/one

### Phase 2: Implementation
#### T2.1: Task Two
**State:** [PENDING]
**Target:** target/two
`
	_ = os.WriteFile(tasksPath, []byte(content), 0644)

	// Update T1.1
	if err := updateTaskStatusFile(projectRoot, slug, "T1.1", "finished"); err != nil {
		t.Fatalf("updateTaskStatusFile failed for T1.1: %v", err)
	}

	// Update T2.1
	if err := updateTaskStatusFile(projectRoot, slug, "T2.1", "in-progress"); err != nil {
		t.Fatalf("updateTaskStatusFile failed for T2.1: %v", err)
	}

	updatedContent, _ := os.ReadFile(tasksPath)
	expected := `
## 2. Tasks

### Phase 1: Setup
#### T1.1: Task One
**State:** [FINISHED]
**Target:** target/one

### Phase 2: Implementation
#### T2.1: Task Two
**State:** [IN-PROGRESS]
**Target:** target/two
`
	if string(updatedContent) != expected {
		t.Errorf("Unexpected content after update:\n%s", string(updatedContent))
	}
}

func TestUpdateTaskStatusFile_WithChecklists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-tasks-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "checklist-feature"
	content := `
## 2. Tasks

### Phase 1: Mixed
- [ ] T1.1: Modern Checklist
**Target:** src/one.go

#### T1.2: Classic with Checkbox
- [ ] T1.2: Classic with Checkbox
**Target:** src/two.go

- [/] T1.3: Working Checklist
**Target:** src/three.go
`
	tasksPath := setupTasksTest(t, projectRoot, slug, content)

	// Update T1.1 to finished -> expect [x]
	if err := updateTaskStatusFile(projectRoot, slug, "T1.1", "finished"); err != nil {
		t.Fatalf("updateTaskStatusFile failed for T1.1: %v", err)
	}

	// Update T1.2 to in-progress -> expect [/]
	if err := updateTaskStatusFile(projectRoot, slug, "T1.2", "in-progress"); err != nil {
		t.Fatalf("updateTaskStatusFile failed for T1.2: %v", err)
	}

	// Update T1.3 to finished -> expect [x]
	if err := updateTaskStatusFile(projectRoot, slug, "T1.3", "finished"); err != nil {
		t.Fatalf("updateTaskStatusFile failed for T1.3: %v", err)
	}

	updatedContent, _ := os.ReadFile(tasksPath)
	expected := `
## 2. Tasks

### Phase 1: Mixed
- [x] T1.1: Modern Checklist
**Target:** src/one.go

#### T1.2: Classic with Checkbox
- [/] T1.2: Classic with Checkbox
**Target:** src/two.go

- [x] T1.3: Working Checklist
**Target:** src/three.go
`
	if string(updatedContent) != expected {
		t.Errorf("Unexpected content after checklist update:\n%s", string(updatedContent))
	}
}

func TestUpdateTaskStatusFile_TimeTracking(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-time-tracking-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "time-tracked-spec"
	content := `
## 2. Tasks
- [ ] T1.1: Tracked Task
**Target:** src/main.go
`
	setupTasksTest(t, projectRoot, slug, content)

	// 1. Update to in-progress
	if err := updateTaskStatusFile(projectRoot, slug, "T1.1", "in-progress"); err != nil {
		t.Fatalf("updateTaskStatusFile failed: %v", err)
	}

	// Check spec.yaml
	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}
	if len(meta.TimeLogs["T1.1"].Sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(meta.TimeLogs["T1.1"].Sessions))
	}
	if meta.TimeLogs["T1.1"].Sessions[0].CompletedAt != nil {
		t.Error("expected session to be open")
	}

	// 2. Update to finished
	if err := updateTaskStatusFile(projectRoot, slug, "T1.1", "finished"); err != nil {
		t.Fatalf("updateTaskStatusFile failed: %v", err)
	}

	meta, _ = LoadMetadata(projectRoot, slug)
	if meta.TimeLogs["T1.1"].Sessions[0].CompletedAt == nil {
		t.Error("expected session to be closed")
	}
}

func TestUpdateTaskBlockState_ParallelWith(t *testing.T) {
	task := &taskBlock{}
	
	// Test basic parallel with
	updateTaskBlockState(task, "**Parallel With:** T1.2, T1.3")
	if !task.isParallel {
		t.Error("expected isParallel to be true")
	}
	if len(task.parallelWith) != 2 || task.parallelWith[0] != "T1.2" || task.parallelWith[1] != "T1.3" {
		t.Errorf("unexpected parallelWith: %v", task.parallelWith)
	}

	// Test spacing and empty peer logic
	task2 := &taskBlock{}
	updateTaskBlockState(task2, "**Parallel With:**   T2.1  ,  , T2.2   ")
	if !task2.isParallel {
		t.Error("expected isParallel to be true")
	}
	if len(task2.parallelWith) != 2 || task2.parallelWith[0] != "T2.1" || task2.parallelWith[1] != "T2.2" {
		t.Errorf("unexpected parallelWith: %v", task2.parallelWith)
	}

	// Test empty parallel with
	task3 := &taskBlock{}
	updateTaskBlockState(task3, "**Parallel With:**")
	if !task3.isParallel {
		t.Error("expected isParallel to be true")
	}
	if len(task3.parallelWith) != 0 {
		t.Errorf("unexpected parallelWith: %v", task3.parallelWith)
	}
}

func assertTaskFileContent(t *testing.T, tasksPath, expected string) {
	t.Helper()
	updatedContent, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks.md: %v", err)
	}
	if string(updatedContent) != expected {
		t.Errorf("Unexpected content after update:\n%s", string(updatedContent))
	}
}

func assertTaskSessions(t *testing.T, meta *Metadata, taskIDs []string, expectedCount int, isOpen bool) {
	t.Helper()
	for _, taskID := range taskIDs {
		log, exists := meta.TimeLogs[taskID]
		if !exists {
			t.Fatalf("expected time log for %s", taskID)
		}
		if len(log.Sessions) != expectedCount {
			t.Fatalf("expected %d session(s) for %s, got %d", expectedCount, taskID, len(log.Sessions))
		}
		if isOpen && log.Sessions[expectedCount-1].CompletedAt != nil {
			t.Errorf("expected session for %s to be open", taskID)
		} else if !isOpen && log.Sessions[expectedCount-1].CompletedAt == nil {
			t.Errorf("expected session for %s to be closed", taskID)
		}
	}
}

const (
	batchTasksInitialContent = `
## 2. Tasks

### Phase 1: Batch Updates
- [ ] T1.1: First Task
**Target:** src/first.go
**State:** [PENDING]

- [ ] T1.2: Second Task
**Target:** src/second.go
**State:** [PENDING]

- [ ] T1.3: Third Task
**Target:** src/third.go
**State:** [PENDING]
`
	batchTasksExpectedInProgress = `
## 2. Tasks

### Phase 1: Batch Updates
- [/] T1.1: First Task
**Target:** src/first.go
**State:** [IN-PROGRESS]

- [/] T1.2: Second Task
**Target:** src/second.go
**State:** [IN-PROGRESS]

- [ ] T1.3: Third Task
**Target:** src/third.go
**State:** [PENDING]
`
	batchTasksExpectedFinished = `
## 2. Tasks

### Phase 1: Batch Updates
- [x] T1.1: First Task
**Target:** src/first.go
**State:** [FINISHED]

- [x] T1.2: Second Task
**Target:** src/second.go
**State:** [FINISHED]

- [ ] T1.3: Third Task
**Target:** src/third.go
**State:** [PENDING]
`
)

func TestUpdateTaskStatusesFile_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-batch-tasks-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "batch-feature"
	tasksPath := setupTasksTest(t, projectRoot, slug, batchTasksInitialContent)

	t.Run("transition to in-progress", func(t *testing.T) {
		if err := updateTaskStatusesFile(projectRoot, slug, []string{"T1.1", "T1.2"}, "in-progress"); err != nil {
			t.Fatalf("updateTaskStatusesFile failed: %v", err)
		}
		assertTaskFileContent(t, tasksPath, batchTasksExpectedInProgress)
	})

	t.Run("transition to finished", func(t *testing.T) {
		if err := updateTaskStatusesFile(projectRoot, slug, []string{"T1.1", "T1.2"}, "finished"); err != nil {
			t.Fatalf("updateTaskStatusesFile failed: %v", err)
		}
		assertTaskFileContent(t, tasksPath, batchTasksExpectedFinished)
	})
}

func TestUpdateTaskStatusesFile_AtomicFailure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-batch-atomic-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "atomic-feature"
	initialContent := `
## 2. Tasks

### Phase 1: Atomic Test
- [ ] T1.1: First Task
**Target:** src/first.go
**State:** [PENDING]

- [ ] T1.2: Second Task
**Target:** src/second.go
**State:** [PENDING]
`
	tasksPath := setupTasksTest(t, projectRoot, slug, initialContent)

	// Attempt updating T1.1 (valid) and T99.99 (invalid)
	err = updateTaskStatusesFile(projectRoot, slug, []string{"T1.1", "T99.99"}, "finished")
	if err == nil {
		t.Fatal("expected error for nonexistent task ID, got nil")
	}

	// Verify file content remains completely unchanged
	currentContent, readErr := os.ReadFile(tasksPath)
	if readErr != nil {
		t.Fatalf("failed to read tasks.md: %v", readErr)
	}
	if string(currentContent) != initialContent {
		t.Errorf("File was modified despite atomic failure:\n%s", string(currentContent))
	}

	// Verify no session was created in metadata
	meta, err := LoadMetadata(projectRoot, slug)
	if err == nil && meta.TimeLogs != nil && len(meta.TimeLogs) > 0 {
		t.Errorf("expected no metadata time logs, got: %v", meta.TimeLogs)
	}
}

func TestUpdateTaskStatusesFile_MetadataTiming(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-batch-timing-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	projectRoot := tempDir
	slug := "timing-feature"
	content := `
## 2. Tasks

### Phase 1: Timing
- [ ] T1.1: Task 1
**Target:** src/t1.go
**State:** [PENDING]

- [ ] T1.2: Task 2
**Target:** src/t2.go
**State:** [PENDING]

- [ ] T1.3: Task 3
**Target:** src/t3.go
**State:** [PENDING]
`
	setupTasksTest(t, projectRoot, slug, content)

	// 1. Batch start sessions for T1.1 and T1.2
	if err := updateTaskStatusesFile(projectRoot, slug, []string{"T1.1", "T1.2"}, "in-progress"); err != nil {
		t.Fatalf("updateTaskStatusesFile failed: %v", err)
	}

	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}
	assertTaskSessions(t, meta, []string{"T1.1", "T1.2"}, 1, true)

	if _, exists := meta.TimeLogs["T1.3"]; exists {
		t.Errorf("expected no time log for T1.3")
	}

	// 2. Batch end sessions for T1.1 and T1.2
	if err := updateTaskStatusesFile(projectRoot, slug, []string{"T1.1", "T1.2"}, "finished"); err != nil {
		t.Fatalf("updateTaskStatusesFile failed: %v", err)
	}

	meta, err = LoadMetadata(projectRoot, slug)
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}
	assertTaskSessions(t, meta, []string{"T1.1", "T1.2"}, 1, false)
}


const scenarioFixtureTasks = `---
slug: scenario-guard
---

# Implementation Roadmap

## 1. Execution Strategy
- **Gravity Order:** one then two.

## 2. Tasks

### Phase 1: Core

- [x] T1.1: [CODE] First task
**Target:** ` + "`src/a.go`" + `
**Context:** [US-1]

**Constraints:**
- Keep it small.

**Acceptance Check:**
` + "`go test ./...`" + ` exits 0.

- [ ] T1.2: [CODE] Second task
**Target:** ` + "`src/b.go`" + `
**Context:** [US-2]

**Constraints:**
- Last real constraint.

**Acceptance Check:**
` + "`go test ./src/...`" + ` exits 0.
`

const scenarioFixtureSection = `
## 3. Black-box Scenarios

- **S1** [US-1] First scenario
  - Command: ` + "`specforce spec list --json`" + `
  - Expected: exit 0 and JSON output.
- **S2** [US-2] Second scenario
  - Command: ` + "`specforce spec validate`" + `
  - Expected: no errors.
`

func parseScenarioFixture(t *testing.T, slug, content string) (*ImplementationReport, TaskValidation) {
	t.Helper()
	root := t.TempDir()
	setupTasksTest(t, root, slug, content)
	report, err := ParseTasks(context.Background(), root, slug)
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	validation, err := ValidateTasksReport(context.Background(), root, slug)
	if err != nil {
		t.Fatalf("ValidateTasksReport failed: %v", err)
	}
	return report, validation
}

func TestParseTasks_ScenarioSectionNotCountedAsTasks(t *testing.T) {
	base, baseVal := parseScenarioFixture(t, "scenario-base", scenarioFixtureTasks)
	with, withVal := parseScenarioFixture(t, "scenario-with", scenarioFixtureTasks+scenarioFixtureSection)

	baseTasks, withTasks := base.Tasks(), with.Tasks()
	if len(baseTasks) != 2 {
		t.Fatalf("baseline expected 2 tasks, got %d", len(baseTasks))
	}
	if len(withTasks) != len(baseTasks) {
		t.Fatalf("scenario section changed task count: got %d, want %d", len(withTasks), len(baseTasks))
	}
	if len(with.Phases) != len(base.Phases) {
		t.Errorf("scenario section changed phase count: got %d, want %d", len(with.Phases), len(base.Phases))
	}
	if with.Status != base.Status {
		t.Errorf("scenario section changed progress status: got %q, want %q", with.Status, base.Status)
	}
	for i := range baseTasks {
		if withTasks[i].ID != baseTasks[i].ID || withTasks[i].State != baseTasks[i].State {
			t.Errorf("task %d changed: got %s/%s, want %s/%s", i, withTasks[i].ID, withTasks[i].State, baseTasks[i].ID, baseTasks[i].State)
		}
	}

	last, lastBase := withTasks[len(withTasks)-1], baseTasks[len(baseTasks)-1]
	if !reflect.DeepEqual(last.Constraints, lastBase.Constraints) {
		t.Errorf("last task constraints absorbed scenario text: got %v, want %v", last.Constraints, lastBase.Constraints)
	}
	if last.Verification != lastBase.Verification {
		t.Errorf("last task verification absorbed scenario text: got %q, want %q", last.Verification, lastBase.Verification)
	}
	for _, c := range last.Constraints {
		if strings.Contains(c, "S1") || strings.Contains(c, "S2") || strings.Contains(c, "Command") {
			t.Errorf("scenario text leaked into constraints: %q", c)
		}
	}

	if len(withVal.Errors) != len(baseVal.Errors) || len(withVal.Warnings) != len(baseVal.Warnings) {
		t.Errorf("validation changed with scenario section: errors %v vs %v, warnings %v vs %v",
			withVal.Errors, baseVal.Errors, withVal.Warnings, baseVal.Warnings)
	}
	if len(withVal.Errors) != 0 {
		t.Errorf("expected valid roadmap, got errors: %v", withVal.Errors)
	}
}

func TestParseTasks_LegacyVerifyPhaseStaysValid(t *testing.T) {
	legacy := scenarioFixtureTasks + `
### Phase 2: Black-box Verification

- [ ] T2.1: [VERIFY] Scenario: legacy check
**Target:** ` + "`src/c.go`" + `
**Context:** [US-1]

**Constraints:**
- Run the legacy check.

**Acceptance Check:**
` + "`go test ./src/...`" + ` exits 0.
`
	report, validation := parseScenarioFixture(t, "scenario-legacy", legacy)

	if len(validation.Errors) != 0 {
		t.Fatalf("legacy [VERIFY] phase must stay valid, got errors: %v", validation.Errors)
	}
	tasks := report.Tasks()
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks including legacy [VERIFY], got %d", len(tasks))
	}
	if len(report.Phases) != 2 {
		t.Errorf("expected 2 phases, got %d", len(report.Phases))
	}
	last := tasks[len(tasks)-1]
	if last.ID != "T2.1" || !strings.Contains(last.Title, "[VERIFY]") {
		t.Errorf("legacy [VERIFY] task not counted as an ordinary task: %+v", last)
	}
}
