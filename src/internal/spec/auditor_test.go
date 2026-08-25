package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDeterministicCheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spec-auditor-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0750); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	// 1. Create requirements.md with US-1 and US-2
	reqContent := `---
slug: test-spec
---
# Feature: Test
## 2. Out of Scope
- [US-3] should be ignored.
## 3. Acceptance Criteria
### [US-1] Requirement 1
### [US-2] Requirement 2
`
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte(reqContent), 0600); err != nil {
		t.Fatalf("failed to write requirements.md: %v", err)
	}

	// 2. Create tasks.md with only US-1
	tasksContent := `---
slug: test-spec
---
# Roadmap
### Phase 1
- [ ] T1.1: Task 1
**Context:** [US-1]
`
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	auditor := NewAuditor(tmpDir)
	errors := auditor.DeterministicCheck(context.Background(), slug)

	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	if errors[0].Code != "MISSING_TASK" {
		t.Errorf("expected MISSING_TASK code, got %q", errors[0].Code)
	}

	if errors[0].Context != "## [US-2]" {
		t.Errorf("expected context '## [US-2]', got %q", errors[0].Context)
	}
}

func TestParseErrors(t *testing.T) {
	auditor := &SpecAuditor{}
	content := `
Random preamble.
[COHERENCE_ERROR]
Artifact: tasks
Code: MISSING_TASK
Message: US-3 is missing
Context: ## [US-3]
[/COHERENCE_ERROR]

[COHERENCE_ERROR]
Artifact: design
Code: TECH_DRIFT
Message: DB mismatch
Context: Use Postgres
[/COHERENCE_ERROR]
`
	errors := auditor.ParseErrors(content)

	if len(errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errors))
	}

	if errors[0].Artifact != "tasks" || errors[0].Code != "MISSING_TASK" {
		t.Errorf("error 0 mismatch: %+v", errors[0])
	}

	if errors[1].Artifact != "design" || errors[1].Message != "DB mismatch" {
		t.Errorf("error 1 mismatch: %+v", errors[1])
	}
}

func TestDeterministicCheck_SmallSpec(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spec-auditor-small-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	auditor := NewAuditor(tmpDir)

	t.Run("Small spec without requirements.md produces no error", func(t *testing.T) {
		slug := "small-no-reqs"
		specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
		_ = os.MkdirAll(specDir, 0750)
		_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: small\n"), 0600)
		_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte("### Phase 1\n- [ ] T1.1: Task\n"), 0600)

		errors := auditor.DeterministicCheck(context.Background(), slug)
		if len(errors) != 0 {
			t.Errorf("expected no errors for small spec missing requirements.md, got: %v", errors)
		}
	})

	t.Run("Medium spec without requirements.md produces FILE_MISSING error", func(t *testing.T) {
		slug := "medium-no-reqs"
		specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
		_ = os.MkdirAll(specDir, 0750)
		_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: medium\n"), 0600)
		_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte("### Phase 1\n- [ ] T1.1: Task\n"), 0600)

		errors := auditor.DeterministicCheck(context.Background(), slug)
		if len(errors) != 1 || errors[0].Code != "FILE_MISSING" || errors[0].Artifact != "requirements" {
			t.Errorf("expected FILE_MISSING error for requirements, got: %v", errors)
		}
	})

	t.Run("Small spec with requirements.md still cross-references US tags", func(t *testing.T) {
		slug := "small-with-reqs"
		specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
		_ = os.MkdirAll(specDir, 0750)
		_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: small\n"), 0600)
		reqContent := "# Requirements\n### [US-1] Req 1\n### [US-2] Req 2\n"
		_ = os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte(reqContent), 0600)
		tasksContent := "### Phase 1\n- [ ] T1.1: Task\n**Context:** [US-1]\n"
		_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(tasksContent), 0600)

		errors := auditor.DeterministicCheck(context.Background(), slug)
		if len(errors) != 1 || errors[0].Code != "MISSING_TASK" {
			t.Errorf("expected 1 MISSING_TASK error, got: %v", errors)
		}
	})
}
