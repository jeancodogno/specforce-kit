package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestGetStatus_BugTasksValidation_Deep(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-bug-status-deep-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	mockFS := fstest.MapFS{
		"bug-requirements.yaml": {
			Data: []byte("description: Bug Req\ninstruction: I\ntemplate: T\n"),
		},
		"tasks.yaml": {
			Data: []byte("description: Tasks\ninstruction: I\ntemplate: T\ndependency: design\n"),
		},
	}
	registry, _ := NewRegistry(mockFS)

	slug := "test-bug"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)

	// Create spec.yaml with type: bug
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("slug: test-bug\ntype: bug\n"), 0644)

	// Create a malformed tasks.md (no phases)
	content := `
# Implementation Roadmap: Broken Bug Tasks
Just some text without phases or tasks.
`
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(content), 0644)

	// Create requirements.md to ensure found == total and validation is triggered
	_ = os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("bug requirements"), 0644)

	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	foundTasks := false
	for _, art := range status.Artifacts {
		if art.Name == "bug-tasks" {
			foundTasks = true
			if len(art.ValidationErrors) == 0 {
				t.Errorf("expected validation errors for malformed tasks.md in bug spec, got 0. Name: %s", art.Name)
			} else {
				t.Logf("Bug tasks validation errors: %v", art.ValidationErrors)
			}
		}
	}

	if !foundTasks {
		t.Errorf("tasks artifact not found (as bug-tasks) in status artifacts list")
	}

	if status.IsValid {
		t.Errorf("expected status.IsValid to be false for malformed tasks.md")
	}
}
