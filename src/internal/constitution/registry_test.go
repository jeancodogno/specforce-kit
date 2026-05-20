package constitution

import (
	"os"
	"testing"
	"testing/fstest"
)

func TestNewRegistry(t *testing.T) {
	// The artifacts are in src/internal/agent/artifacts/constitution/
	// We are in src/internal/constitution/
	artifactsRoot := "../../internal/agent/artifacts/constitution"
	artifactsFS := os.DirFS(artifactsRoot)

	registry, err := NewRegistry(artifactsFS)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	artifacts := registry.List()
	if len(artifacts) < 7 {
		t.Errorf("Expected at least 7 artifacts, got %d", len(artifacts))
	}

	// Verify Architecture artifact
	arch, ok := registry.Get("architecture")
	if !ok {
		t.Fatal("Architecture artifact not found")
	}

	if arch.Description == "" {
		t.Error("Architecture description should not be empty")
	}
	if arch.Instruction == "" {
		t.Error("Architecture instruction should not be empty")
	}
	if arch.Template == "" {
		t.Error("Architecture template should not be empty")
	}
	if arch.Name != "architecture" {
		t.Errorf("Expected architecture, got %s", arch.Name)
	}
	if arch.Path != ".specforce/docs/architecture.md" {
		t.Errorf("Expected .specforce/docs/architecture.md, got %s", arch.Path)
	}

	// Verify Memorial artifact (custom path)
	mem, ok := registry.Get("memorial")
	if !ok {
		t.Fatal("Memorial artifact not found")
	}
	if mem.Path != ".specforce/memorial/ROUTING.md" {
		t.Errorf("Expected .specforce/memorial/ROUTING.md, got %s", mem.Path)
	}
}

func TestNewRegistry_IndexSpecialCase(t *testing.T) {
	mockFS := fstest.MapFS{
		"index.yaml": &fstest.MapFile{
			Data: []byte("description: Index Description\ninstruction: Index Instruction\ntemplate: Index Template"),
		},
	}

	reg, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("Failed to create registry with mock FS: %v", err)
	}

	idx, ok := reg.Get("index")
	if !ok {
		t.Fatal("Index artifact not found")
	}

	if idx.Name != "index" {
		t.Errorf("Expected index, got %s", idx.Name)
	}
	if idx.Path != ".specforce/docs/_index.md" {
		t.Errorf("Expected .specforce/docs/_index.md, got %s", idx.Path)
	}
}

func TestNewRegistry_LeadingUnderscoreSpecialCase(t *testing.T) {
	mockFS := fstest.MapFS{
		"_index.yaml": &fstest.MapFile{
			Data: []byte("description: Index Description\ninstruction: Index Instruction\ntemplate: Index Template"),
		},
	}

	reg, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("Failed to create registry with mock FS: %v", err)
	}

	idx, ok := reg.Get("index")
	if !ok {
		// If it's stored as _index, this will fail
		t.Fatal("Index artifact not found by 'index' slug")
	}

	if idx.Name != "index" {
		t.Errorf("Expected name 'index', got %s", idx.Name)
	}
	if idx.Path != ".specforce/docs/_index.md" {
		t.Errorf("Expected path .specforce/docs/_index.md, got %s", idx.Path)
	}
}

func TestNewRegistry_PathOverride(t *testing.T) {
	mockFS := fstest.MapFS{
		"custom.yaml": &fstest.MapFile{
			Data: []byte("description: Custom Description\ninstruction: Custom Instruction\ntemplate: Custom Template\npath: .specforce/custom/PATH.md"),
		},
	}

	reg, err := NewRegistry(mockFS)
	if err != nil {
		t.Fatalf("Failed to create registry with mock FS: %v", err)
	}

	custom, ok := reg.Get("custom")
	if !ok {
		t.Fatal("Custom artifact not found")
	}

	if custom.Path != ".specforce/custom/PATH.md" {
		t.Errorf("Expected path .specforce/custom/PATH.md, got %s", custom.Path)
	}
}
