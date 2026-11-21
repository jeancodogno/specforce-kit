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

func TestCheckQAGate(t *testing.T) {
	write := func(t *testing.T, yamlBody string) (string, string) {
		t.Helper()
		root := t.TempDir()
		slug := "gate-slug"
		dir := filepath.Join(root, ".specforce", "specs", slug)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "spec.yaml"), []byte(yamlBody), 0600); err != nil {
			t.Fatal(err)
		}
		return root, slug
	}

	t.Run("Missing", func(t *testing.T) {
		root, slug := write(t, "slug: gate-slug\nname: n\ntype: feature\nsize: medium\n")
		res, err := CheckQAGate(root, slug)
		if err != nil || res.State != QAGateMissing {
			t.Fatalf("expected missing, got %+v err=%v", res, err)
		}
	})
	t.Run("Passed", func(t *testing.T) {
		root, slug := write(t, "slug: gate-slug\nqa:\n  status: PASSED\n  run: 1\n  criteria:\n    - {id: US-1, scenario: s, evidence: e, verdict: PASSED}\n")
		res, err := CheckQAGate(root, slug)
		if err != nil || res.State != QAGatePassed || len(res.Failed) != 0 {
			t.Fatalf("expected passed, got %+v err=%v", res, err)
		}
	})
	t.Run("Failed", func(t *testing.T) {
		root, slug := write(t, "slug: gate-slug\nqa:\n  status: FAILED\n  run: 1\n  criteria:\n    - {id: US-1, scenario: s1, evidence: e, verdict: FAILED}\n    - {id: US-2, scenario: s2, evidence: e, verdict: PASSED}\n")
		res, err := CheckQAGate(root, slug)
		if err != nil || res.State != QAGateFailed {
			t.Fatalf("expected failed, got %+v err=%v", res, err)
		}
		if len(res.Failed) != 1 || res.Failed[0].ID != "US-1" || res.Failed[0].Scenario != "s1" {
			t.Errorf("bad failed list: %+v", res.Failed)
		}
	})
	t.Run("Unreadable", func(t *testing.T) {
		root, slug := write(t, "slug: gate-slug\nqa: [not, a, mapping]\n")
		if _, err := CheckQAGate(root, slug); err == nil {
			t.Fatal("expected error for malformed qa block")
		}
	})
}
