package spec

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

	// 4. Test refinement persistence
	loaded.Refinement.IterationCount = 2
	loaded.Refinement.IsValid = true
	loaded.Refinement.Errors = []string{"Error 1", "Error 2"}
	if err := SaveMetadata(tmpDir, slug, loaded); err != nil {
		t.Fatalf("failed to save refinement metadata: %v", err)
	}

	refined, err := LoadMetadata(tmpDir, slug)
	if err != nil {
		t.Fatalf("failed to load refined metadata: %v", err)
	}
	if refined.Refinement.IterationCount != 2 {
		t.Errorf("expected iteration 2, got %d", refined.Refinement.IterationCount)
	}
	if !refined.Refinement.IsValid {
		t.Error("expected IsValid to be true")
	}
	if len(refined.Refinement.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(refined.Refinement.Errors))
	}
}

func TestMetadataSessionManagement(t *testing.T) {
	meta := &Metadata{
		Slug: "test",
		Name: "Test",
	}

	taskID := "T1.1"

	// 1. Start session
	meta.StartSession(taskID)
	if len(meta.TimeLogs[taskID].Sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(meta.TimeLogs[taskID].Sessions))
	}
	if meta.TimeLogs[taskID].Sessions[0].CompletedAt != nil {
		t.Error("expected session to be open")
	}

	// 2. End session
	time.Sleep(10 * time.Millisecond)
	meta.EndSession(taskID)
	if meta.TimeLogs[taskID].Sessions[0].CompletedAt == nil {
		t.Fatal("expected session to be closed")
	}

	duration := meta.GetTaskDuration(taskID)
	if duration < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %v", duration)
	}

	// 3. Resume session
	meta.StartSession(taskID)
	if len(meta.TimeLogs[taskID].Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(meta.TimeLogs[taskID].Sessions))
	}
	
	time.Sleep(10 * time.Millisecond)
	meta.EndSession(taskID)
	
	totalDuration := meta.GetTaskDuration(taskID)
	if totalDuration < 20*time.Millisecond {
		t.Errorf("expected total duration >= 20ms, got %v", totalDuration)
	}
}

func TestValidateSize(t *testing.T) {
	tests := []struct {
		size  SpecSize
		valid bool
	}{
		{SpecSizeSmall, true},
		{SpecSizeMedium, true},
		{SpecSizeLarge, true},
		{SpecSizeComplex, true},
		{"", false},
		{"unknown", false},
		{"huge", false},
	}

	for _, tt := range tests {
		if got := ValidateSize(tt.size); got != tt.valid {
			t.Errorf("ValidateSize(%q) = %v, want %v", tt.size, got, tt.valid)
		}
	}
}

func TestMetadataSize(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spec-metadata-size-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "size-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0750); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	// 1. Default size on missing file
	missingMeta, err := LoadMetadata(tmpDir, "non-existent")
	if err != nil {
		t.Fatalf("expected no error loading missing metadata, got %v", err)
	}
	if missingMeta.Size != SpecSizeMedium {
		t.Errorf("expected default size %q, got %q", SpecSizeMedium, missingMeta.Size)
	}

	// 2. Default size on spec.yaml without size field
	yamlWithoutSize := []byte("slug: size-spec\nname: Size Spec\ntype: feature\n")
	if err := os.WriteFile(filepath.Join(specDir, "spec.yaml"), yamlWithoutSize, 0600); err != nil {
		t.Fatalf("failed to write spec.yaml: %v", err)
	}
	loaded, err := LoadMetadata(tmpDir, slug)
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}
	if loaded.Size != SpecSizeMedium {
		t.Errorf("expected empty size to default to %q, got %q", SpecSizeMedium, loaded.Size)
	}

	// 3. Round-trip persistence for each size
	sizes := []SpecSize{SpecSizeSmall, SpecSizeMedium, SpecSizeLarge, SpecSizeComplex}
	for _, sz := range sizes {
		m := &Metadata{
			Slug: slug,
			Name: "Size Spec",
			Type: "feature",
			Size: sz,
		}
		if err := SaveMetadata(tmpDir, slug, m); err != nil {
			t.Fatalf("failed to save metadata for size %q: %v", sz, err)
		}

		reloaded, err := LoadMetadata(tmpDir, slug)
		if err != nil {
			t.Fatalf("failed to reload metadata for size %q: %v", sz, err)
		}
		if reloaded.Size != sz {
			t.Errorf("expected size %q after save/load, got %q", sz, reloaded.Size)
		}
	}
}

