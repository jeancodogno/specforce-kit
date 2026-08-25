package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestGetStatus(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	registry := setupMockRegistry(t)
	slug := "test-slug"
	specDir := setupSpecDir(t, tmpDir, slug)

	t.Run("No files exist", func(t *testing.T) {
		testGetStatus_NoFiles(t, tmpDir, slug, registry)
	})

	t.Run("One file exists", func(t *testing.T) {
		testGetStatus_OneFile(t, tmpDir, slug, specDir, registry)
	})

	t.Run("Both files exist", func(t *testing.T) {
		testGetStatus_BothFiles(t, tmpDir, slug, specDir, registry)
	})
}

func testGetStatus_NoFiles(t *testing.T, tmpDir, slug string, registry *Registry) {
	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Found != 0 {
		t.Errorf("expected 0 found artifacts, got %d", status.Found)
	}
	if status.Total != 2 {
		t.Errorf("expected 2 total artifacts, got %d", status.Total)
	}
	if status.Progress != 0 {
		t.Errorf("expected 0%% progress, got %d%%", status.Progress)
	}

	// Verify blocking logic
	for _, art := range status.Artifacts {
		if art.Name == "requirements" {
			if art.Blocked {
				t.Errorf("requirements should not be blocked")
			}
			if art.Dependency != "" {
				t.Errorf("requirements should not have dependency")
			}
		}
		if art.Name == "design" {
			if !art.Blocked {
				t.Errorf("design should be blocked by requirements")
			}
			if art.Dependency != "requirements" {
				t.Errorf("expected dependency 'requirements', got '%s'", art.Dependency)
			}
		}
	}
}

func testGetStatus_OneFile(t *testing.T, tmpDir, slug, specDir string, registry *Registry) {
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Found != 1 {
		t.Errorf("expected 1 found artifact, got %d", status.Found)
	}
	if status.Progress != 50 {
		t.Errorf("expected 50%% progress, got %d%%", status.Progress)
	}

	// Verify blocking logic
	for _, art := range status.Artifacts {
		if art.Name == "design" {
			if art.Blocked {
				t.Errorf("design should NOT be blocked when requirements exists")
			}
		}
	}
}

func testGetStatus_BothFiles(t *testing.T, tmpDir, slug, specDir string, registry *Registry) {
	if err := os.WriteFile(filepath.Join(specDir, "design.md"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Found != 2 {
		t.Errorf("expected 2 found artifacts, got %d", status.Found)
	}
	if status.Progress != 100 {
		t.Errorf("expected 100%% progress, got %d%%", status.Progress)
	}
}

func TestGetStatus_ConditionalValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-test-conditional-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	registry := setupConditionalValidationRegistry(t)
	slug := "conditional-spec"
	specDir := setupConditionalValidationFiles(t, tmpDir, slug)

	t.Run("Validation skipped when artifacts are missing", func(t *testing.T) {
		status, err := GetStatus(context.Background(), tmpDir, slug, registry)
		if err != nil {
			t.Fatalf("GetStatus failed: %v", err)
		}

		// Found should be 1 (tasks.md exists, requirements.md missing)
		if status.Found != 1 {
			t.Errorf("expected 1 found artifact, got %d", status.Found)
		}

		for _, art := range status.Artifacts {
			if art.Name == "feature-tasks" {
				if len(art.ValidationErrors) > 0 {
					t.Errorf("expected no validation errors when missing artifacts, got %v", art.ValidationErrors)
				}
			}
		}
	})

	t.Run("Validation performed when all artifacts are present", func(t *testing.T) {
		_ = os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("reqs"), 0644)
		status, err := GetStatus(context.Background(), tmpDir, slug, registry)
		if err != nil {
			t.Fatalf("GetStatus failed: %v", err)
		}

		if status.Found != 2 {
			t.Errorf("expected 2 found artifacts, got %d", status.Found)
		}

		verifyTasksValidationErrors(t, status)
	})
}

func setupConditionalValidationRegistry(t *testing.T) *Registry {
	mockFS := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("\ndescription: \"Req\"\ninstruction: \"instr\"\ntemplate: \"tpl\"\n"),
		},
		"tasks.yaml": {
			Data: []byte("\ndescription: \"Tasks\"\ninstruction: \"instr\"\ntemplate: \"tpl\"\n"),
		},
	}
	registry, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}
	return registry
}

func setupConditionalValidationFiles(t *testing.T, tmpDir, slug string) string {
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)

	// Create meta
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature"), 0644)

	// Create malformed tasks.md
	malformedTasks := "### Wrong Header"
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(malformedTasks), 0644)
	return specDir
}

func verifyTasksValidationErrors(t *testing.T, status SpecStatus) {
	foundErrors := false
	for _, art := range status.Artifacts {
		if art.Name == "feature-tasks" {
			if len(art.ValidationErrors) > 0 {
				foundErrors = true
			}
		}
	}
	if !foundErrors {
		t.Errorf("expected validation errors when all artifacts are present, but got none")
	}
}

func setupSpecDir(t *testing.T, tmpDir, slug string) string {
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: large\n"), 0644)
	return specDir
}

func setupMockRegistry(t *testing.T) *Registry {
	mockFS := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte(`
description: "Requirements"
instruction: "Write requirements"
template: "# Requirements"
`),
		},
		"design.yaml": {
			Data: []byte(`
description: "Design"
instruction: "Write design"
template: "# Design"
dependency: "requirements"
`),
		},
	}

	registry, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}
	return registry
}

func TestGetStatus_Sizing(t *testing.T) {
	mockFS := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("description: Req\ninstruction: instr\ntemplate: tpl\n"),
		},
		"design.yaml": {
			Data: []byte("description: Design\ninstruction: instr\ntemplate: tpl\ndependency: requirements\n"),
		},
		"tasks.yaml": {
			Data: []byte("description: Tasks\ninstruction: instr\ntemplate: tpl\ndependency: design\n"),
		},
	}
	registry, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "specforce-status-sizing-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("Small size spec returns only tasks artifact", func(t *testing.T) {
		testSmallSpecStatus(t, tmpDir, registry)
	})

	t.Run("Medium size spec returns requirements and tasks", func(t *testing.T) {
		testMediumSpecStatus(t, tmpDir, registry)
	})

	t.Run("Large size spec returns all artifacts", func(t *testing.T) {
		testLargeSpecStatus(t, tmpDir, registry)
	})
}

func testSmallSpecStatus(t *testing.T, tmpDir string, registry *Registry) {
	slug := "small-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: small\n"), 0644)

	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Total != 1 {
		t.Errorf("expected 1 total artifact for small spec, got %d", status.Total)
	}
	if len(status.Artifacts) != 1 || status.Artifacts[0].Name != "feature-tasks" {
		t.Fatalf("expected only feature-tasks artifact, got %v", status.Artifacts)
	}
	if status.Found != 0 {
		t.Errorf("expected 0 found artifacts, got %d", status.Found)
	}

	// Write tasks.md
	_ = os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte("### Phase 1: Test\n- [ ] T1.1: Step\n**Target:** file.go\n**Action Steps:**\n- item\n**Acceptance Check:**\n- verify\n"), 0644)
	status, err = GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Found != 1 {
		t.Errorf("expected 1 found artifact, got %d", status.Found)
	}
	if status.Progress != 100 {
		t.Errorf("expected 100%% progress, got %d%%", status.Progress)
	}
	if !status.IsValid {
		t.Errorf("expected status to be valid, got errors: %v", status.Artifacts[0].ValidationErrors)
	}
}

func testMediumSpecStatus(t *testing.T, tmpDir string, registry *Registry) {
	slug := "medium-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: medium\n"), 0644)

	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Total != 2 {
		t.Errorf("expected 2 total artifacts for medium spec, got %d", status.Total)
	}
	if len(status.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(status.Artifacts))
	}
}

func testLargeSpecStatus(t *testing.T, tmpDir string, registry *Registry) {
	slug := "large-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	_ = os.MkdirAll(specDir, 0755)
	_ = os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: large\n"), 0644)

	status, err := GetStatus(context.Background(), tmpDir, slug, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Total != 3 {
		t.Errorf("expected 3 total artifacts for large spec, got %d", status.Total)
	}
}
