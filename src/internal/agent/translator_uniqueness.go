package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// nameTracker tracks header names per session to detect duplicates.
type nameTracker struct {
	names map[string]string // name -> source blueprint path
}

func newNameTracker() *nameTracker {
	return &nameTracker{names: make(map[string]string)}
}

func (t *nameTracker) validate(name, path string) error {
	if existingPath, ok := t.names[name]; ok && existingPath != path {
		return fmt.Errorf("security: duplicate header name %q detected in %s (conflicts with %s)", name, path, existingPath)
	}
	t.names[name] = path
	return nil
}

func resolveHeaderName(bp *core.Blueprint, mapping core.MappingConfig, category string) string {
	name := bp.Metadata.Name
	if name == "" {
		name = mapping.Name
		if name == "SKILL" {
			name = strings.TrimSuffix(filepath.Base(bp.ID), filepath.Ext(bp.ID))
			if category == "commands" {
				name = "spf." + name
			}
		}
	}
	return name
}
