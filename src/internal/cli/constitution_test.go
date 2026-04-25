package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestHandleConstitutionArtifact(t *testing.T) {
	wd, _ := os.Getwd()
	projectRoot := filepath.Join(wd, "..", "..", "..")

	executor := &Executor{
		Version:       "1.0.0",
		DevMode:       true,
		ArtifactsRoot: filepath.Join(projectRoot, "src", "internal", "agent", "artifacts"),
	}
	ui := tui.NewUI()

	t.Run("Valid artifact", func(t *testing.T) {
		err := executor.HandleConstitutionArtifact(context.Background(), ui, "architecture", false)
		if err != nil {
			t.Errorf("HandleConstitutionArtifact failed: %v", err)
		}
	})

	t.Run("Invalid artifact", func(t *testing.T) {
		err := executor.HandleConstitutionArtifact(context.Background(), ui, "invalid", false)
		if err != nil {
			t.Errorf("HandleConstitutionArtifact failed: %v", err)
		}
	})

	t.Run("JSON mode", func(t *testing.T) {
		err := executor.HandleConstitutionArtifact(context.Background(), ui, "architecture", true)
		if err != nil {
			t.Errorf("HandleConstitutionArtifact failed: %v", err)
		}
	})

	t.Run("List mode", func(t *testing.T) {
		err := executor.HandleConstitutionArtifact(context.Background(), ui, "", false)
		if err != nil {
			t.Errorf("HandleConstitutionArtifact failed: %v", err)
		}
	})
}
