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

func setupSpecDir(t *testing.T, tmpDir, slug string) string {
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
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
