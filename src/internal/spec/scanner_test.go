package spec

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func setupMockProject(t *testing.T, tempDir string) {
	// 1. Setup mock .specforce structure
	specforceDir := filepath.Join(tempDir, ".specforce")
	docsDir := filepath.Join(specforceDir, "docs")
	specsDir := filepath.Join(specforceDir, "specs")
	archiveDir := filepath.Join(specforceDir, "archive")

	dirs := []string{docsDir, specsDir, archiveDir}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("Failed to create mock dir %s: %v", d, err)
		}
	}

	// 2. Add mock files
	// Constitution doc
	if err := os.WriteFile(filepath.Join(docsDir, "architecture.md"), []byte("# Architecture"), 0644); err != nil {
		t.Fatalf("Failed to create mock doc: %v", err)
	}
	// Active Spec
	activeSpecDir := filepath.Join(specsDir, "0001-test-spec")
	if err := os.MkdirAll(activeSpecDir, 0755); err != nil {
		t.Fatalf("Failed to create mock spec dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(activeSpecDir, "requirements.md"), []byte("# Requirements"), 0644); err != nil {
		t.Fatalf("Failed to create mock requirement doc: %v", err)
	}
	// Active Implementation (has tasks.md)
	implSpecDir := filepath.Join(specsDir, "0002-test-impl")
	if err := os.MkdirAll(implSpecDir, 0755); err != nil {
		t.Fatalf("Failed to create mock impl dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(implSpecDir, "tasks.md"), []byte("## 2. Tasks\n### T1.1: [SCAFFOLD] Task\n**State:** `[PENDING]`"), 0644); err != nil {
		t.Fatalf("Failed to create mock task doc: %v", err)
	}
	// Archived Spec
	archivedSpecDir := filepath.Join(archiveDir, "0003-test-archive")
	if err := os.MkdirAll(archivedSpecDir, 0755); err != nil {
		t.Fatalf("Failed to create mock archive dir: %v", err)
	}
}

func TestScanProject(t *testing.T) {
	tempDir := t.TempDir()
	setupMockProject(t, tempDir)

	// 3. Mock registry
	registry := &Registry{
		artifacts: map[string]Artifact{
			"requirements": {Name: "requirements"},
			"design":       {Name: "design"},
			"tasks":        {Name: "tasks"},
		},
	}

	// 4. Execute scan
	tree, err := ScanProject(context.Background(), tempDir, registry)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// 5. Assertions
	if len(tree.Categories[CategoryConstitution]) != 1 {
		t.Errorf("Expected 1 constitution doc, got %d", len(tree.Categories[CategoryConstitution]))
	}
	if len(tree.Categories[CategoryActiveSpecs]) != 1 {
		t.Errorf("Expected 1 active spec, got %d", len(tree.Categories[CategoryActiveSpecs]))
	}
	if len(tree.Categories[CategoryImplementations]) != 1 {
		t.Errorf("Expected 1 active implementation, got %d", len(tree.Categories[CategoryImplementations]))
	}
	if len(tree.Categories[CategoryArchived]) != 1 {
		t.Errorf("Expected 1 archived spec, got %d", len(tree.Categories[CategoryArchived]))
	}

	// Check status
	if tree.Categories[CategoryImplementations][0].Status != "READY" {
		t.Errorf("Expected READY status for implementation, got %s", tree.Categories[CategoryImplementations][0].Status)
	}
}

func TestScanProject_Cancellation(t *testing.T) {
	tempDir := t.TempDir()

	// Setup many files to increase chance of catching it (though here it's synchronous)
	docsDir := filepath.Join(tempDir, ".specforce", "docs")
	_ = os.MkdirAll(docsDir, 0755)
	for i := 0; i < 100; i++ {
		_ = os.WriteFile(filepath.Join(docsDir, filepath.Join(string(rune(i)), ".md")), []byte("# Test"), 0644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := ScanProject(ctx, tempDir, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestParseTasks_InvalidFile_ReturnsErrInvalidSpecFile(t *testing.T) {
	tmpDir := t.TempDir()

	slug := "bad-spec"
	// Create the spec directory but do NOT create tasks.md
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	_, err := ParseTasks(context.Background(), tmpDir, slug)
	if err == nil {
		t.Fatal("expected an error for missing tasks.md, got nil")
	}

	if !errors.Is(err, core.ErrInvalidSpecFile) {
		t.Errorf("expected errors.Is(err, core.ErrInvalidSpecFile) to be true, got: %v", err)
	}
}

func TestScanProject_AnyTaskWorking(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Setup mock .specforce structure
	specsDir := filepath.Join(tempDir, ".specforce", "specs")
	implSpecDir := filepath.Join(specsDir, "working-impl")
	if err := os.MkdirAll(implSpecDir, 0755); err != nil {
		t.Fatalf("Failed to create mock impl dir: %v", err)
	}

	// 2. Add tasks.md with an in-progress task
	content := "## 2. Tasks\n### T1.1: [SCAFFOLD] Working Task\n**State:** `[IN-PROGRESS]`"
	if err := os.WriteFile(filepath.Join(implSpecDir, "tasks.md"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create mock tasks.md: %v", err)
	}

	// 3. Mock registry
	registry := &Registry{
		artifacts: map[string]Artifact{
			"tasks": {Name: "tasks"},
		},
	}

	// 4. Execute scan
	tree, err := ScanProject(context.Background(), tempDir, registry)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// 5. Assertions
	found := false
	for _, item := range tree.Categories[CategoryImplementations] {
		if item.Slug == "working-impl" {
			found = true
			if !item.AnyTaskWorking {
				t.Errorf("Expected AnyTaskWorking to be true for working-impl")
			}
		}
	}
	if !found {
		t.Errorf("Expected to find working-impl in CategoryImplementations")
	}

	t.Run("Task State X", func(t *testing.T) {
		implXDir := filepath.Join(specsDir, "x-impl")
		if err := os.MkdirAll(implXDir, 0755); err != nil {
			t.Fatal(err)
		}
		contentX := "## 2. Tasks\n### T1.1: Task X\n**State:** `[X]`"
		if err := os.WriteFile(filepath.Join(implXDir, "tasks.md"), []byte(contentX), 0644); err != nil {
			t.Fatal(err)
		}
		tree, _ := ScanProject(context.Background(), tempDir, registry)
		foundX := false
		for _, item := range tree.Categories[CategoryImplementations] {
			if item.Slug == "x-impl" {
				foundX = true
				if item.TaskCount != 1 {
					t.Errorf("Expected TaskCount 1 for X state, got %d", item.TaskCount)
				}
			}
		}
		if !foundX {
			t.Errorf("Expected to find x-impl")
		}
	})
}

func TestScanProject_ArchiveLimits(t *testing.T) {
	tempDir := t.TempDir()
	archiveDir := filepath.Join(tempDir, ".specforce", "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create 12 archived items
	for i := 0; i < 12; i++ {
		name := filepath.Join(archiveDir, filepath.Join(string(rune('a'+i)), "folder"))
		if err := os.MkdirAll(name, 0755); err != nil {
			t.Fatal(err)
		}
	}
	
	// Hidden file and hidden dir in archive
	if err := os.WriteFile(filepath.Join(archiveDir, ".hidden-file"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(archiveDir, ".hidden-dir"), 0755); err != nil {
		t.Fatal(err)
	}

	tree, err := ScanProject(context.Background(), tempDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	archived := tree.Categories[CategoryArchived]
	if len(archived) != 10 {
		t.Errorf("Expected archive limit of 10, got %d", len(archived))
	}
	
	for _, item := range archived {
		if strings.HasPrefix(item.Slug, ".") {
			t.Errorf("Archive scan should ignore hidden entries: %s", item.Slug)
		}
	}
}

func TestParseWorktreePorcelain(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a dummy dir to simulate another worktree
	otherDir := filepath.Join(tempDir, "other-wt")
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatalf("Failed to create dummy dir: %v", err)
	}

	input := "worktree " + tempDir + "\n" +
		"HEAD e46d9471c0c22760c6e7260ce1fafb3a5ace52d3\n" +
		"branch refs/heads/main\n" +
		"\n" +
		"worktree " + otherDir + "\n" +
		"HEAD 1234567890abcdef1234567890abcdef12345678\n" +
		"branch refs/heads/feature-x\n" +
		"\n" +
		"worktree /non/existent/path\n" +
		"HEAD abcdef1234567890abcdef1234567890abcdef\n" +
		"branch refs/heads/ghost\n"

	worktrees := parseWorktreePorcelain(input)

	if len(worktrees) != 2 {
		t.Errorf("Expected 2 worktrees, got %d", len(worktrees))
	}

	if worktrees[0].Path != tempDir || worktrees[0].Branch != "main" {
		t.Errorf("First worktree mismatch: %+v", worktrees[0])
	}

	if worktrees[1].Path != otherDir || worktrees[1].Branch != "feature-x" {
		t.Errorf("Second worktree mismatch: %+v", worktrees[1])
	}
}

func TestScanProject_WithWorktrees(t *testing.T) {
	tempDir := t.TempDir()
	
	// Initialize a real git repo in tempDir
	runCmd := func(dir string, name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run %s %v: %v", name, args, err)
		}
	}

	runCmd(tempDir, "git", "init")
	runCmd(tempDir, "git", "config", "user.email", "test@example.com")
	runCmd(tempDir, "git", "config", "user.name", "Test User")
	runCmd(tempDir, "git", "commit", "--allow-empty", "-m", "initial")

	// Create .specforce in main repo
	setupMockProject(t, tempDir)

	// Create a worktree
	wtDir := filepath.Join(t.TempDir(), "wt1")
	runCmd(tempDir, "git", "worktree", "add", "-b", "wt-branch", wtDir)

	// Create .specforce in worktree
	setupMockProject(t, wtDir)
	// Add a unique spec in worktree to distinguish
	wtSpecDir := filepath.Join(wtDir, ".specforce", "specs", "wt-spec")
	if err := os.MkdirAll(wtSpecDir, 0755); err != nil {
		t.Fatalf("Failed to create wt spec dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wtSpecDir, "requirements.md"), []byte("# WT Spec"), 0644); err != nil {
		t.Fatalf("Failed to write wt spec requirements: %v", err)
	}

	// Setup mock registry
	registry := &Registry{
		artifacts: map[string]Artifact{
			"requirements": {Name: "requirements"},
		},
	}

	tree, err := ScanProject(context.Background(), tempDir, registry)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// Should have items from both
	foundWTSpec := false
	for _, item := range tree.Categories[CategoryActiveSpecs] {
		if item.Slug == "wt-spec" {
			foundWTSpec = true
			if item.Worktree == "" {
				t.Error("wt-spec should have worktree tagged")
			}
		}
	}

	if !foundWTSpec {
		t.Error("Did not find spec from worktree")
	}

	// Constitution should ONLY be from main root
	if len(tree.Categories[CategoryConstitution]) != 1 {
		t.Errorf("Expected only 1 constitution doc (from main root), got %d", len(tree.Categories[CategoryConstitution]))
	}
}
