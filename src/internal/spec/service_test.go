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
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, "T1.1", "finished")
		if err != nil {
			t.Fatalf("UpdateTaskStatus failed: %v", err)
		}
		updatedContent, _ := os.ReadFile(filepath.Join(tasksDir, "tasks.md"))
		if !strings.Contains(string(updatedContent), "#### T1.1: Init\n**State:** [FINISHED]") {
			t.Errorf("T1.1 not updated")
		}
	})
	
	t.Run("update non-finished status", func(t *testing.T) {
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, "T1.2", "in-progress")
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
		err := svc.UpdateTaskStatus(context.Background(), tmpDir, slug, "T1.1", "finished")
		if err == nil {
			t.Fatal("expected error from failing hook")
		}
	})
	
	t.Run("nil configProvider", func(t *testing.T) {
		svcNil := NewService(nil, nil)
		err := svcNil.UpdateTaskStatus(context.Background(), tmpDir, slug, "T1.1", "finished")
		if err != nil {
			t.Fatal(err)
		}
	})
}
