package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spec-metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0750); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	// 1. Test default loading (backward compatibility)
	meta, err := LoadMetadata(tmpDir, slug)
	if err != nil {
		t.Errorf("expected no error loading missing metadata, got %v", err)
	}
	if meta.Type != "feature" {
		t.Errorf("expected default type 'feature', got %q", meta.Type)
	}

	// 2. Test saving metadata
	original := &Metadata{
		Slug: slug,
		Name: "Test Spec",
		Type: "bug",
	}
	if err := SaveMetadata(tmpDir, slug, original); err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	// 3. Test loading saved metadata
	loaded, err := LoadMetadata(tmpDir, slug)
	if err != nil {
		t.Fatalf("failed to load saved metadata: %v", err)
	}
	if loaded.Type != "bug" {
		t.Errorf("expected type 'bug', got %q", loaded.Type)
	}
	if loaded.Name != "Test Spec" {
		t.Errorf("expected name 'Test Spec', got %q", loaded.Name)
	}
}
