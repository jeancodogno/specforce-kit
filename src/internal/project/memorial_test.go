package project

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMemorialService_Initialize(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memorial-test-init")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	svc := NewMemorialService(tmpDir)
	ctx := context.Background()

	err = svc.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	routingPath := filepath.Join(tmpDir, ".specforce", "memorial", "ROUTING.md")
	if _, err := os.Stat(routingPath); os.IsNotExist(err) {
		t.Error("ROUTING.md not created")
	}
}

func TestMemorialService_RecordAndConsolidate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memorial-test-record")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	svc := NewMemorialService(tmpDir)
	ctx := context.Background()

	_ = svc.Initialize(ctx)

	f1 := Fragment{
		Date:    time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
		Scope:   "Security",
		Author:  "AgentX",
		Type:    FragmentAction,
		Title:   "Hardened Permissions",
		Content: "Applied 0600 to all state files.",
	}
	if err := svc.Record(ctx, f1); err != nil {
		t.Fatalf("Record f1 failed: %v", err)
	}

	f2 := Fragment{
		Date:    time.Date(2026, 5, 8, 11, 0, 0, 0, time.UTC),
		Scope:   "SDD",
		Author:  "AgentY",
		Type:    FragmentLesson,
		Title:   "Trigger Optimization",
		Content: "Reduced token usage by 30%.",
	}
	if err := svc.Record(ctx, f2); err != nil {
		t.Fatalf("Record f2 failed: %v", err)
	}

	consolidated, err := svc.Consolidate(ctx, 10)
	if err != nil {
		t.Fatalf("Consolidate failed: %v", err)
	}

	expectedStrings := []string{
		"FOR AI AGENTS: RULES OF ENGAGEMENT",
		"Hardened Permissions",
		"Trigger Optimization",
	}
	for _, s := range expectedStrings {
		if !strings.Contains(consolidated, s) {
			t.Errorf("Consolidated output missing expected string: %s", s)
		}
	}

	idx1 := strings.Index(consolidated, "Hardened Permissions")
	idx2 := strings.Index(consolidated, "Trigger Optimization")
	if idx2 > idx1 {
		t.Error("Fragments not ordered newest first")
	}
}

func TestMemorialService_Migration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memorial-test-migration")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create legacy memorial.md
	legacyDir := filepath.Join(tmpDir, ".specforce", "docs")
	_ = os.MkdirAll(legacyDir, 0750)
	legacyPath := filepath.Join(legacyDir, "memorial.md")
	legacyContent := "# Legacy Memorial\nSome old context."
	_ = os.WriteFile(legacyPath, []byte(legacyContent), 0600)

	svc := NewMemorialService(tmpDir)
	ctx := context.Background()

	err = svc.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize with migration failed: %v", err)
	}

	// Check if legacy.md exists in new location
	migratedPath := filepath.Join(tmpDir, ".specforce", "memorial", "legacy.md")
	data, err := os.ReadFile(migratedPath)
	if err != nil {
		t.Fatalf("Migrated file not found: %v", err)
	}
	if string(data) != legacyContent {
		t.Errorf("Migrated content mismatch: expected %q, got %q", legacyContent, string(data))
	}

	// Check if old file was renamed
	if _, err := os.Stat(legacyPath + ".deprecated"); os.IsNotExist(err) {
		t.Error("Legacy file not renamed to .deprecated")
	}
}
