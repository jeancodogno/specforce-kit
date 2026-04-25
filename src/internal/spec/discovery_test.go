package spec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDiscovery(t *testing.T) {
	// Setup temp project root
	if err := os.MkdirAll(filepath.Join("testdata", "discovery"), 0755); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("testdata") }()

	projectRoot := filepath.Join("testdata", "discovery")
	specsDir := filepath.Join(projectRoot, ".specforce", "specs")
	archiveDir := filepath.Join(projectRoot, ".specforce", "archive")

	// Create active spec
	if err := os.MkdirAll(filepath.Join(specsDir, "active-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	// Create archived spec
	if err := os.MkdirAll(filepath.Join(archiveDir, "archived-spec"), 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("SpecExists active", func(t *testing.T) {
		exists, status := SpecExists(projectRoot, "active-spec")
		if !exists || status != "active" {
			t.Errorf("expected active-spec to exist and be active, got %v, %s", exists, status)
		}
	})

	t.Run("SpecExists archived", func(t *testing.T) {
		exists, status := SpecExists(projectRoot, "archived-spec")
		if !exists || status != "archived" {
			t.Errorf("expected archived-spec to exist and be archived, got %v, %s", exists, status)
		}
	})

	t.Run("SpecExists none", func(t *testing.T) {
		exists, status := SpecExists(projectRoot, "none")
		if exists {
			t.Errorf("expected none to not exist, got %v, %s", exists, status)
		}
	})

	t.Run("ListActiveSpecs", func(t *testing.T) {
		specs, err := ListActiveSpecs(context.Background(), projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		if len(specs) != 1 || specs[0].Slug != "active-spec" {
			t.Errorf("expected 1 active spec (active-spec), got %d: %+v", len(specs), specs)
		}
	})
}
