package spec

import (
	"os"
	"path/filepath"
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
