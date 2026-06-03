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
