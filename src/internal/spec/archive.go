package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ArchiveSpec moves a specification from the specs directory to the archive directory.
func ArchiveSpec(ctx context.Context, basePath, slug string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	specsDir := filepath.Join(basePath, ".specforce", "specs", slug)
	archiveDir := filepath.Join(basePath, ".specforce", "archive", slug)

	// Check if source exists
	srcInfo, err := os.Stat(specsDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("spec not found: %s", slug)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("spec path is not a directory: %s", specsDir)
	}

	// Check if destination already exists
	_, err = os.Stat(archiveDir)
	if err == nil {
		return fmt.Errorf("archive already exists: %s", slug)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check archive destination: %w", err)
	}

	// Ensure archive parent directory exists
	if err := os.MkdirAll(filepath.Dir(archiveDir), 0750); err != nil {
		return fmt.Errorf("failed to create archive base directory: %w", err)
	}

	// Perform atomic move
	if err := os.Rename(specsDir, archiveDir); err != nil {
		// Fallback for cross-partition move or other rename failures could be implemented here
		return fmt.Errorf("failed to archive spec: %w", err)
	}

	return nil
}
