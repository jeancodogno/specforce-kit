package installer

import (
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// Options defines configuration for the installation/extraction process.
type Options struct {
	// ToolsOnly, when true, restricts installation to tool-specific directories
	// and explicitly excludes the .specforce directory.
	ToolsOnly bool
}

// ShouldInstall determines if a given relative path should be installed based on the provided options.
func ShouldInstall(path string, opts Options) bool {
	if !opts.ToolsOnly {
		return true
	}

	// Always allow absolute paths in ToolsOnly mode, as they represent global tool exports
	if filepath.IsAbs(path) {
		return true
	}

	// Ensure the path is normalized for prefix checking
	normalized := strings.TrimPrefix(path, "./")

	// Explicitly exclude .specforce/ directory in ToolsOnly mode
	if strings.HasPrefix(normalized, ".specforce/") {
		return false
	}

	// Whitelist tool-specific directories
	for _, prefix := range core.ToolPrefixes {
		if strings.HasPrefix(normalized, prefix) {
			return true
		}
	}

	// Default to false in ToolsOnly mode if not matched by whitelist
	return false
}
