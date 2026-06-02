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
	tests = append(tests, getDensityErrorCases()...)
	tests = append(tests, getParallelTaskErrorCases()...)
	return tests
}

func getParallelTaskErrorCases() []validateTasksTestCase {
	return []validateTasksTestCase{
		getUnknownParallelTaskErrorCase(),
		getCrossPhaseParallelTaskErrorCase(),
		getAsymmetricCrossPhaseParallelTaskErrorCase(),
		getTargetIsolationConflictErrorCase(),
	}
}

func getUnknownParallelTaskErrorCase() validateTasksTestCase {
	return validateTasksTestCase{
		name: "Unknown Parallel Task",
		content: `### Phase 1: Phase One
- [ ] T1.1: First Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T9.9
**Acceptance Check:**
Check it`,
		expected: []string{
			"Task T1.1 (line 2) references unknown parallel task T9.9",
		},
	}
}

func getCrossPhaseParallelTaskErrorCase() validateTasksTestCase {
	return validateTasksTestCase{
		name: "Cross Phase Parallel Task",
		content: `### Phase 1: Phase One
- [ ] T1.1: First Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T2.1
**Acceptance Check:**
Check it
### Phase 2: Phase Two
- [ ] T2.1: Second Task
**Target:** API
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T1.1
**Acceptance Check:**
Check it`,
		expected: []string{
			"Task T1.1 (line 2) references parallel task T2.1 from a different phase",
			"Task T2.1 (line 12) references parallel task T1.1 from a different phase",
		},
	}
}

func getAsymmetricCrossPhaseParallelTaskErrorCase() validateTasksTestCase {
	return validateTasksTestCase{
		name: "Asymmetric Cross Phase Parallel Task",
		content: `### Phase 1: Phase One
- [ ] T1.1: First Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T2.1
**Acceptance Check:**
Check it
### Phase 2: Phase Two
- [ ] T2.1: Second Task
**Target:** API
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Acceptance Check:**
Check it`,
		expected: []string{
			"Task T1.1 (line 2) references parallel task T2.1 from a different phase",
			"Task T2.1 (line 12) references parallel task T1.1 from a different phase",
		},
	}
}

func getTargetIsolationConflictErrorCase() validateTasksTestCase {
	return validateTasksTestCase{
		name: "Target Isolation Conflict",
		content: `### Phase 1: Phase One
- [ ] T1.1: First Task
**Target:** src/common.go
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T1.2
**Acceptance Check:**
Check it
- [ ] T1.2: Second Task
**Target:** src/common.go
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Parallel With:** T1.1
**Acceptance Check:**
Check it`,
		expected: []string{
			"Parallel conflict: T1.1 and T1.2 both target \"src/common.go\"",
		},
	}
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
- Verify output
**Acceptance Check:**
Check files`,
			expected: nil,
		},
		{
			name: "Implicit Parallel Symmetry",
			content: `### Phase 1: Setup
- [ ] T1.1: Task 1
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do this
- Do that
**Parallel With:** T1.2
**Acceptance Check:**
Check files
- [ ] T1.2: Task 2
**Target:** API
**Context:** US-1
**Action Steps:**
- Do this
- Do that
**Acceptance Check:**
Check files`,
			expected: nil,
		},
	}
}

func getDensityErrorCases() []validateTasksTestCase {
	return []validateTasksTestCase{
		{
			name: "Low Action Density",
			content: `### Phase 1: Phase One
- [ ] T1.1: Low Density Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Only one step
**Acceptance Check:**
Verify`,
			expected: []string{
				"Task T1.1 (line 2) has insufficient action density (found 1, expected at least 2)",
			},
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
- Step 2
**Acceptance Check:**
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
- Step 2
**Acceptance Check:**
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
- Step 2
**Acceptance Check:**
Check it
- [ ] T1.3: Gapped Task
**Target:** CLI
**Context:** US-1
**Action Steps:**
- Do something
- Step 2
**Acceptance Check:**
Check it`,
			expected: []string{"Task sequence gap at line 10: expected T1.2, found T1.3"},
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
				"Task T1.1 (line 2) is missing mandatory **Acceptance Check:** section",
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
- Step 2
**Acceptance Check:**
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
**Acceptance Check:**
- Step that looks like action step but is under verification`,
			expected: []string{
				"Task T1.1 (line 2) is missing mandatory items under **Action Steps:**",
			},
		},
	}
}
