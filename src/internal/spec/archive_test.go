package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveSpec(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "specforce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	slug := "0001-test"
	specsPath := filepath.Join(tempDir, ".specforce", "specs", slug)
	archivePath := filepath.Join(tempDir, ".specforce", "archive", slug)

	if err := os.MkdirAll(specsPath, 0755); err != nil {
		t.Fatalf("failed to create specs path: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsPath, "requirements.md"), []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Happy Path
	if err := ArchiveSpec(context.Background(), tempDir, slug); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Errorf("expected archive path to exist")
	}
	if _, err := os.Stat(specsPath); !os.IsNotExist(err) {
		t.Errorf("expected source path to be removed")
	}

	// Error Case: Spec not found
	err = ArchiveSpec(context.Background(), tempDir, "9999-missing")
	if err == nil || err.Error() != "spec not found: 9999-missing" {
		t.Errorf("expected 'spec not found' error, got %v", err)
	}

	// Error Case: Archive already exists
	// Re-create source
	if err := os.MkdirAll(specsPath, 0755); err != nil {
		t.Fatalf("failed to recreate specs path: %v", err)
	}
	err = ArchiveSpec(context.Background(), tempDir, slug)
	if err == nil || err.Error() != "archive already exists: 0001-test" {
		t.Errorf("expected 'archive already exists' error, got %v", err)
	}
}
