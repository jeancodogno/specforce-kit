package cli

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestIntegration_SpecInitTimestamp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-init-int-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Mkdir(".specforce", 0755); err != nil {
		t.Fatal(err)
	}

	executor := NewExecutor("1.0.0")
	ui := tui.NewUI()

	t.Run("Root slug", func(t *testing.T) { testRootSlug(t, executor, ui) })
	t.Run("Nested slug", func(t *testing.T) { testNestedSlug(t, executor, ui) })
	t.Run("Idempotency", func(t *testing.T) { testIdempotency(t, executor, ui) })
}

func testRootSlug(t *testing.T, executor *Executor, ui core.UI) {
	slug := "root-feature"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false)
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	matches, _ := filepath.Glob(".specforce/specs/*-root-feature")
	if len(matches) != 1 {
		t.Errorf("expected 1 directory matching *-root-feature, found %d", len(matches))
	}

	re := regexp.MustCompile(`^\.specforce/specs/\d{8}-\d{4}-root-feature$`)
	if !re.MatchString(filepath.ToSlash(matches[0])) {
		t.Errorf("directory name %s does not match expected pattern", matches[0])
	}
}

func testNestedSlug(t *testing.T, executor *Executor, ui core.UI) {
	slug := "team-x/api-v1"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false)
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	matches, _ := filepath.Glob(".specforce/specs/team-x/*-api-v1")
	if len(matches) != 1 {
		t.Errorf("expected 1 directory matching team-x/*-api-v1, found %d", len(matches))
	}

	re := regexp.MustCompile(`^\.specforce/specs/team-x/\d{8}-\d{4}-api-v1$`)
	if !re.MatchString(filepath.ToSlash(matches[0])) {
		t.Errorf("directory name %s does not match expected pattern", matches[0])
	}
}

func testIdempotency(t *testing.T, executor *Executor, ui core.UI) {
	slug := "20240101-1200-legacy-spec"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false)
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	expectedPath := filepath.Join(".specforce", "specs", slug)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist without changes", expectedPath)
	}
}
