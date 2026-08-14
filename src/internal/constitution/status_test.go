package constitution

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func getTestRegistry(t *testing.T) *Registry {
	// The artifacts are in src/internal/agent/artifacts/constitution/
	// We are in src/internal/constitution/
	artifactsRoot := "../../internal/agent/artifacts/constitution"
	artifactsFS := os.DirFS(artifactsRoot)

	registry, err := NewRegistry(artifactsFS)
	if err != nil {
		t.Fatalf("Failed to create registry for testing: %v", err)
	}
	return registry
}

func TestGetStatus(t *testing.T) {
	registry := getTestRegistry(t)

	// Create a temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	t.Run("Empty directory", func(t *testing.T) {
		runEmptyDirTest(t, tmpDir, registry)
	})

	t.Run("All core files present", func(t *testing.T) {
		runAllFilesPresentTest(t, tmpDir, registry)
	})

	t.Run("Partially present", func(t *testing.T) {
		runPartiallyPresentTest(t, registry)
	})

	t.Run("Missing docs directory", func(t *testing.T) {
		runMissingDocsDirTest(t, registry)
	})
}

func runEmptyDirTest(t *testing.T, tmpDir string, registry *Registry) {
	status, err := GetStatus(context.Background(), tmpDir, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	// (principles, architecture, ui-ux, security, engineering, governance)
	expectedTotal := 6
	if status.Total != expectedTotal {
		t.Errorf("expected %d artifacts, got %d", expectedTotal, status.Total)
	}
	if status.Found != 0 {
		t.Errorf("expected 0 found, got %d", status.Found)
	}
	if status.Progress != 0 {
		t.Errorf("expected 0%% progress, got %d%%", status.Progress)
	}
	for _, a := range status.Artifacts {
		if a.Exists {
			t.Errorf("artifact %s should not exist", a.Name)
		}
	}
}

func runAllFilesPresentTest(t *testing.T, tmpDir string, registry *Registry) {
	docsDir := filepath.Join(tmpDir, ".specforce", "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	// Create a dummy module
	modulesDir := filepath.Join(docsDir, "modules")
	if err := os.MkdirAll(modulesDir, 0755); err != nil {
		t.Fatalf("failed to create modules dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modulesDir, "auth.md"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	artifacts := registry.List()
	coreCount := 0
	for _, art := range artifacts {
		if art.Slug == "module" {
			continue
		}
		coreCount++
		filePath := filepath.Join(tmpDir, art.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("failed to create dir for %s: %v", art.Name, err)
		}
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file %s: %v", art.Name, err)
		}
	}

	status, err := GetStatus(context.Background(), tmpDir, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.Found != coreCount {
		t.Errorf("expected %d found, got %d", coreCount, status.Found)
	}
	if status.Progress != 100 {
		t.Errorf("expected 100%% progress, got %d%%", status.Progress)
	}
	if len(status.Modules) != 1 || status.Modules[0] != "auth" {
		t.Errorf("expected module 'auth', got %v", status.Modules)
	}
}

func runPartiallyPresentTest(t *testing.T, registry *Registry) {
	partialTmpDir, err := os.MkdirTemp("", "specforce-test-partial-*")
	if err != nil {
		t.Fatalf("failed to create partial temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(partialTmpDir)
	}()

	docsDir := filepath.Join(partialTmpDir, ".specforce", "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	artifacts := registry.List()
	coreArtifacts := make([]Artifact, 0)
	for _, art := range artifacts {
		if art.Slug != "module" {
			coreArtifacts = append(coreArtifacts, art)
		}
	}

	// Create only the first 2 core artifacts
	for i := 0; i < 2; i++ {
		filePath := filepath.Join(partialTmpDir, coreArtifacts[i].Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("failed to create dir for partial test: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write partial test file: %v", err)
		}
	}

	status, err := GetStatus(context.Background(), partialTmpDir, registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.Found != 2 {
		t.Errorf("expected 2 found, got %d", status.Found)
	}
	// (2/7) * 100 = 28.57... -> 28
	expectedProgress := (2 * 100) / len(coreArtifacts)
	if status.Progress != expectedProgress {
		t.Errorf("expected %d%% progress, got %d%%", expectedProgress, status.Progress)
	}
	if len(status.Modules) != 0 {
		t.Errorf("expected 0 modules, got %d", len(status.Modules))
	}
}

func runMissingDocsDirTest(t *testing.T, registry *Registry) {
	status, err := GetStatus(context.Background(), "/non-existent-dir-12345", registry)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Found != 0 {
		t.Errorf("expected 0 found, got %d", status.Found)
	}
}
