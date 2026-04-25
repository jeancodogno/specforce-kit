package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
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

**Verification (TDD):**
run verify 1

### Phase 2: Implementation
#### T1.2: Task Two
**State:** [IN-PROGRESS]
**Target:** target/two
**Context:** context/two

**Action Steps:**
- step A

**Verification (TDD):**
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
