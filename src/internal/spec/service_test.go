package spec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

type mockConfigProvider struct {
	config *core.ProjectConfig
}

func (m *mockConfigProvider) GetConfig(ctx context.Context) (*core.ProjectConfig, error) {
	return m.config, nil
}

func TestGetArtifactInstructions(t *testing.T) {
	artifactsFS := fstest.MapFS{
		"requirements.yaml": &fstest.MapFile{Data: []byte(`
description: Requirements Template
instruction: Base Requirements Instruction
template: Requirements Template Content
`)},
	}
	reg, _ := NewRegistry(artifactsFS)

	config := &core.ProjectConfig{
		Instructions: map[string][]string{
			"requirements": {"Project Rule 1"},
		},
	}
	svc := NewService(reg, &mockConfigProvider{config: config})

	art, err := svc.GetArtifact(context.Background(), "requirements")
	if err != nil {
		t.Fatalf("GetArtifact failed: %v", err)
	}

	expectedPrefix := "## Project Specific Instructions\n- Project Rule 1\n\n"
	if !strings.HasPrefix(art.Instruction, expectedPrefix) {
		t.Errorf("expected instruction to start with project rules, got:\n%s", art.Instruction)
	}

	if !strings.Contains(art.Instruction, "Base Requirements Instruction") {
		t.Errorf("base instructions missing")
	}
}

func TestGetArtifact(t *testing.T) {
	artifactsFS := fstest.MapFS{
		"requirements.yaml": &fstest.MapFile{Data: []byte(`
description: Requirements Template
instruction: Requirements Instruction
template: Requirements Template Content
`)},
	}
	reg, _ := NewRegistry(artifactsFS)
	
	t.Run("success without custom instructions", func(t *testing.T) {
		svc := NewService(reg, nil)
		art, err := svc.GetArtifact(context.Background(), "requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if art.Name != "requirements" {
			t.Errorf("expected 'requirements', got %v", art.Name)
		}
	})

	t.Run("success with custom instructions", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"requirements": {"Custom 1", "Custom 2"},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		art, err := svc.GetArtifact(context.Background(), "requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if !strings.Contains(art.Instruction, "Custom 1") || !strings.Contains(art.Instruction, "Custom 2") {
			t.Errorf("custom instructions not injected")
		}
	})
	
	t.Run("success with empty custom instructions", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"requirements": {},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		art, err := svc.GetArtifact(context.Background(), "requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if strings.Contains(art.Instruction, "Project Specific Instructions") {
			t.Errorf("expected no project specific instructions, got %v", art.Instruction)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewService(reg, nil)
		_, err := svc.GetArtifact(context.Background(), "non-existent")
		if err == nil {
			t.Fatal("expected error for non-existent artifact")
		}
	})
}

func TestGetArtifact_Layered(t *testing.T) {
	artifactsFS := fstest.MapFS{
		"requirements.yaml": &fstest.MapFile{Data: []byte(`
description: Requirements Template
instruction: Base Requirements Instruction
template: Requirements Template Content
`)},
		"design.yaml": &fstest.MapFile{Data: []byte(`
description: Design Template
instruction: Base Design Instruction
template: Design Template Content
`)},
	}
	reg, _ := NewRegistry(artifactsFS)

	t.Run("generic and specific combination", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"requirements":         {"Generic 1"},
				"feature-requirements": {"Specific 1"},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		art, err := svc.GetArtifact(context.Background(), "feature-requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if !strings.Contains(art.Instruction, "Generic 1") {
			t.Errorf("generic instructions missing")
		}
		if !strings.Contains(art.Instruction, "Specific 1") {
			t.Errorf("specific instructions missing")
		}
		// Order check: Generic should come before Specific
		genericIdx := strings.Index(art.Instruction, "Generic 1")
		specificIdx := strings.Index(art.Instruction, "Specific 1")
		if genericIdx > specificIdx {
			t.Errorf("expected generic instructions before specific")
		}
	})

	t.Run("generic only mapping (prefix resolution)", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"requirements": {"Generic Only"},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		art, err := svc.GetArtifact(context.Background(), "bug-requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if !strings.Contains(art.Instruction, "Generic Only") {
			t.Errorf("generic instructions not applied to prefixed name")
		}
	})

	testInferenceAndDeduplication(t, reg)
}

func testInferenceAndDeduplication(t *testing.T, reg *Registry) {
	t.Run("right-to-left inference", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"design": {"Design Rules"},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		// tasks-for-design should resolve to 'design' as it is the right-most keyword
		art, err := svc.GetArtifact(context.Background(), "tasks-for-design")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		if !strings.Contains(art.Instruction, "Design Rules") {
			t.Errorf("failed to infer base type 'design' from 'tasks-for-design'")
		}
	})

	t.Run("deduplication", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"requirements":         {"Rule A", "Rule B"},
				"feature-requirements": {"Rule B", "Rule C"},
			},
		}
		svc := NewService(reg, &mockConfigProvider{config: config})
		art, err := svc.GetArtifact(context.Background(), "feature-requirements")
		if err != nil {
			t.Fatalf("GetArtifact failed: %v", err)
		}
		// Count occurrences of Rule B
		count := strings.Count(art.Instruction, "Rule B")
		if count != 1 {
			t.Errorf("expected 'Rule B' to be deduplicated, got count %d", count)
		}
	})
}

func TestGetImplementationStatus(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "test-slug"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `
# Implementation Tasks
### Phase 1: Core
#### T1.1: Task 1
**State:** [PENDING]
`
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create requirements.md and design.md to avoid "blocked" status
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("# Requirements"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "design.md"), []byte("# Design"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("success with instructions", func(t *testing.T) {
		config := &core.ProjectConfig{
			Instructions: map[string][]string{
				"implementation": {"Do it right"},
			},
		}
		svc := NewService(nil, &mockConfigProvider{config: config})
		report, err := svc.GetImplementationStatus(context.Background(), tmpDir, slug)
		if err != nil {
			t.Fatalf("GetImplementationStatus failed: %v", err)
		}
		if report.Status != "ready" {
			t.Errorf("expected status 'ready', got %v", report.Status)
		}
		if len(report.Instructions) == 0 || report.Instructions[0] != "Do it right" {
			t.Errorf("instructions not injected")
		}
	})
	
	t.Run("success without configProvider", func(t *testing.T) {
		svc := NewService(nil, nil)
		report, err := svc.GetImplementationStatus(context.Background(), tmpDir, slug)
		if err != nil {
			t.Fatalf("GetImplementationStatus failed: %v", err)
		}
		if len(report.Instructions) != 0 {
			t.Errorf("expected no instructions")
		}
	})
}

func TestGetImplementationStatus_TieredSizing(t *testing.T) {
	mockFS := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("description: Requirements\ninstruction: Requirements Instruction\ntemplate: Requirements Template Content\n"),
		},
		"design.yaml": {
			Data: []byte("description: Design\ninstruction: Design Instruction\ntemplate: Design Template Content\ndependency: requirements\n"),
		},
		"tasks.yaml": {
			Data: []byte("description: Tasks\ninstruction: Tasks Instruction\ntemplate: Tasks Template Content\ndependency: design\n"),
		},
	}
	registry, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	svc := NewService(registry, nil)

	t.Run("Small spec with only tasks.md is ready", func(t *testing.T) {
		testSmallSpecImplementationStatus(t, svc)
	})

	t.Run("Medium spec with requirements.md and tasks.md is ready", func(t *testing.T) {
		testMediumSpecImplementationStatus(t, svc)
	})

	t.Run("Large spec missing design.md is blocked", func(t *testing.T) {
		testLargeSpecImplementationStatus(t, svc)
	})
}

const testTieredTasksContent = `
# Implementation Tasks
### Phase 1: Core
#### T1.1: Task 1
**State:** [PENDING]
`

func testSmallSpecImplementationStatus(t *testing.T, svc *Service) {
	t.Helper()
	tmpDir := t.TempDir()
	slug := "small-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: small\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(testTieredTasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := svc.GetImplementationStatus(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("GetImplementationStatus failed: %v", err)
	}
	if report.Status != "ready" {
		t.Errorf("expected status 'ready', got %q", report.Status)
	}
	if len(report.MissingArtifacts) != 0 {
		t.Errorf("expected no missing artifacts, got %v", report.MissingArtifacts)
	}
}

func testMediumSpecImplementationStatus(t *testing.T, svc *Service) {
	t.Helper()
	tmpDir := t.TempDir()
	slug := "medium-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: medium\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("# Requirements"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(testTieredTasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := svc.GetImplementationStatus(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("GetImplementationStatus failed: %v", err)
	}
	if report.Status != "ready" {
		t.Errorf("expected status 'ready', got %q", report.Status)
	}
	if len(report.MissingArtifacts) != 0 {
		t.Errorf("expected no missing artifacts, got %v", report.MissingArtifacts)
	}
}

func testLargeSpecImplementationStatus(t *testing.T, svc *Service) {
	t.Helper()
	tmpDir := t.TempDir()
	slug := "large-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "spec.yaml"), []byte("type: feature\nsize: large\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("# Requirements"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(testTieredTasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := svc.GetImplementationStatus(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("GetImplementationStatus failed: %v", err)
	}
	if report.Status != "blocked" {
		t.Errorf("expected status 'blocked', got %q", report.Status)
	}
	foundDesign := false
	for _, missing := range report.MissingArtifacts {
		if missing == "design.md" {
			foundDesign = true
			break
		}
	}
	if !foundDesign {
		t.Errorf("expected missing artifacts to contain 'design.md', got %v", report.MissingArtifacts)
	}
}

func TestService_GetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "test-slug"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("# Req"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reg, _ := NewRegistry(fstest.MapFS{
		"requirements.yaml": &fstest.MapFile{Data: []byte("description: R\ninstruction: I\ntemplate: T")},
	})
	svc := NewService(reg, nil)
	
	status, err := svc.GetStatus(context.Background(), tmpDir, slug)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Slug != slug {
		t.Errorf("expected slug %v, got %v", slug, status.Slug)
	}
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "test-slug"
	tasksDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := "### Phase 1\n#### T1.1: Init\n**State:** [PENDING]\n#### T1.2: Next\n**State:** [PENDING]"
	if err := os.WriteFile(filepath.Join(tasksDir, "tasks.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	config := &core.ProjectConfig{
		Hooks: core.HooksConfig{OnTaskFinished: []string{"echo task"}},
	}
	svc := NewService(nil, &mockConfigProvider{config: config})

	t.Run("normal task finished", func(t *testing.T) {
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1"}, "finished")
		if err != nil {
			t.Fatalf("UpdateTaskStatus failed: %v", err)
		}
		updatedContent, _ := os.ReadFile(filepath.Join(tasksDir, "tasks.md"))
		if !strings.Contains(string(updatedContent), "#### T1.1: Init\n**State:** [FINISHED]") {
			t.Errorf("T1.1 not updated")
		}
	})
	
	t.Run("update non-finished status", func(t *testing.T) {
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.2"}, "in-progress")
		if err != nil {
			t.Fatal(err)
		}
		updatedContent, _ := os.ReadFile(filepath.Join(tasksDir, "tasks.md"))
		if !strings.Contains(string(updatedContent), "#### T1.2: Next\n**State:** [IN-PROGRESS]") {
			t.Errorf("T1.2 not updated")
		}
	})
}

func TestUpdateTaskStatus_Hooks(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "hooks-slug"
	tasksDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "### Phase 1\n#### T1.1: Task\n**State:** [PENDING]"
	if err := os.WriteFile(filepath.Join(tasksDir, "tasks.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	config := &core.ProjectConfig{
		Hooks: core.HooksConfig{OnTaskFinished: []string{"false"}},
	}
	svc := NewService(nil, &mockConfigProvider{config: config})

	t.Run("failing hook blocks update", func(t *testing.T) {
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1"}, "finished")
		if err == nil {
			t.Fatal("expected error from failing hook")
		}
	})
	
	t.Run("nil configProvider", func(t *testing.T) {
		svcNil := NewService(nil, nil)
		err := svcNil.UpdateTaskStatus(context.Background(), tmpDir, slug, []string{"T1.1"}, "finished")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestResizeSpec(t *testing.T) {
	tmpDir := t.TempDir()
	slug := "resize-test"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	meta := &Metadata{
		Slug: slug,
		Name: slug,
		Type: "feature",
		Size: SpecSizeSmall,
	}
	if err := SaveMetadata(tmpDir, slug, meta); err != nil {
		t.Fatal(err)
	}

	svc := NewService(nil, nil)
	svc.SetProjectRoot(tmpDir)

	t.Run("valid resize", func(t *testing.T) {
		updated, err := svc.ResizeSpec(context.Background(), slug, SpecSizeLarge)
		if err != nil {
			t.Fatalf("ResizeSpec failed: %v", err)
		}
		if updated.Size != SpecSizeLarge {
			t.Errorf("expected size %v, got %v", SpecSizeLarge, updated.Size)
		}
	})

	t.Run("invalid size", func(t *testing.T) {
		_, err := svc.ResizeSpec(context.Background(), slug, SpecSize("invalid"))
		if err == nil {
			t.Fatal("expected error for invalid size, got nil")
		}
	})

	t.Run("non-existent spec", func(t *testing.T) {
		_, err := svc.ResizeSpec(context.Background(), "non-existent", SpecSizeMedium)
		if err == nil {
			t.Fatal("expected error for non-existent spec, got nil")
		}
	})
}
