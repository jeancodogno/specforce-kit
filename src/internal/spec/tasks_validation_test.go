package spec

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type validateTasksTestCase struct {
	name     string
	content  string
	expected []string
}

func TestValidateTasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-feature"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	tasksPath := filepath.Join(specDir, "tasks.md")
	tests := getValidateTasksTestCases()
	runValidateTasksTests(t, tmpDir, slug, tasksPath, tests)
}

func runValidateTasksTests(t *testing.T, tmpDir, slug, tasksPath string, tests []validateTasksTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(tasksPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write tasks.md: %v", err)
			}

			errors, err := ValidateTasks(context.Background(), tmpDir, slug)
			if err != nil {
				t.Fatalf("ValidateTasks failed: %v", err)
			}
			if !reflect.DeepEqual(errors, tt.expected) {
				t.Errorf("got %q, want %q", errors, tt.expected)
			}
		})
	}
}

func getValidateTasksTestCases() []validateTasksTestCase {
	tests := []validateTasksTestCase{}
	tests = append(tests, getHappyPathCases()...)
	tests = append(tests, getHierarchyErrorCases()...)
	tests = append(tests, getFieldAndPhaseErrorCases()...)
	return tests
}

func getHappyPathCases() []validateTasksTestCase {
	return []validateTasksTestCase{
		{
			name: "Happy Path",
			content: `### Phase 1: Setup
- [ ] T1.1: Init
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Run init
**Verification (TDD):**
Check files`,
			expected: nil,
		},
	}
}

func getHierarchyErrorCases() []validateTasksTestCase {
	return []validateTasksTestCase{
		{
			name: "Task Before Phase",
			content: `- [ ] T1.1: Early Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
**Verification (TDD):**
Check it`,
			expected: []string{
				"Task T1.1 (line 1) found before any Phase definition",
				"Task T1.1 (line 1) does not match the parent Phase 0",
				"No valid Phase (### Phase N: Name) found in tasks.md",
			},
		},
		{
			name: "Phase ID Out of Sequence",
			content: `### Phase 2: Wrong Order`,
			expected: []string{
				"Phase ID 2 (line 1) is out of sequence, expected 1",
				"Phase 2 (line 1) has no tasks",
			},
		},
		{
			name: "Task Phase Mismatch",
			content: `### Phase 1: Phase One
- [ ] T2.1: Wrong Phase Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
**Verification (TDD):**
Check it`,
			expected: []string{"Task T2.1 (line 2) does not match the parent Phase 1"},
		},
		{
			name: "Task Sequence Gap",
			content: `### Phase 1: Phase One
- [ ] T1.1: First Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
**Verification (TDD):**
Check it
- [ ] T1.3: Gapped Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
**Verification (TDD):**
Check it`,
			expected: []string{"Task sequence gap at line 9: expected T1.2, found T1.3"},
		},
	}
}

func getFieldAndPhaseErrorCases() []validateTasksTestCase {
	return []validateTasksTestCase{
		{
			name: "Missing Mandatory Fields",
			content: `### Phase 1: Phase One
- [ ] T1.1: Minimal Task`,
			expected: []string{
				"Task T1.1 (line 2) is missing mandatory **Target:** field",
				"Task T1.1 (line 2) is missing mandatory **Context:** field",
				"Task T1.1 (line 2) is missing mandatory **Action Steps:** header",
				"Task T1.1 (line 2) is missing mandatory **Verification (TDD):** section",
			},
		},
		{
			name: "Empty Phase",
			content: `### Phase 1: Empty
### Phase 2: Next
- [ ] T2.1: Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Step
**Verification (TDD):**
Verify`,
			expected: []string{"Phase 1 (line 1) has no tasks"},
		},
		{
			name: "Action items under Verification header",
			content: `### Phase 1: Phase One
- [ ] T1.1: Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
**Verification (TDD):**
- Step that looks like action step but is under verification`,
			expected: []string{
				"Task T1.1 (line 2) is missing mandatory items under **Action Steps:**",
			},
		},
	}
}
