package spec

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var timestampRegex = regexp.MustCompile(`^\d{8}-\d{4}-`)

// PrepareSlug prepends a YYYYMMDD-HHMM timestamp to the final segment of the slug
// if it doesn't already have one.
func PrepareSlug(rawSlug string) string {
	if rawSlug == "" {
		return ""
	}

	// Use ToSlash to handle cross-platform path separators consistently
	path := filepath.ToSlash(rawSlug)
	segments := strings.Split(path, "/")
	
	finalSegment := segments[len(segments)-1]
	
	// Check if already has timestamp
	if timestampRegex.MatchString(finalSegment) {
		return rawSlug
	}

	// Generate timestamp
	ts := time.Now().Format("20060102-1504")
	
	// Sanitize final segment to prevent double hyphens
	sanitized := strings.TrimPrefix(finalSegment, "-")
	newSegment := ts + "-" + sanitized
	
	segments[len(segments)-1] = newSegment
	
	return filepath.FromSlash(strings.Join(segments, "/"))
}

// ResolveSlug attempts to find the actual directory name for a given slug.
// It handles exact matches and fuzzy matches for timestamped slugs (YYYYMMDD-HHMM-slug).
// It also handles sub-paths (e.g., team-a/my-feature).
func ResolveSlug(projectRoot string, slug string) string {
	if slug == "" {
		return ""
	}

	path := filepath.ToSlash(slug)
	segments := strings.Split(path, "/")
	parentPath := strings.Join(segments[:len(segments)-1], "/")
	finalSegment := segments[len(segments)-1]

	specsDir := filepath.Join(projectRoot, ".specforce", "specs", parentPath)

	// 1. Exact match
	if _, err := os.Stat(filepath.Join(specsDir, finalSegment)); err == nil {
		return slug
	}

	// 2. Fuzzy match in specsDir
	entries, err := os.ReadDir(specsDir)
	if err == nil {
		// Sort entries in reverse order to prioritize newest timestamps
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() > entries[j].Name()
		})

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if timestampRegex.MatchString(name) {
				prefix := timestampRegex.FindString(name)
				if name == prefix+finalSegment {
					return filepath.Join(parentPath, name)
				}
			}
		}
	}

	// 3. Repeat for archive if not found in specs
	archiveDir := filepath.Join(projectRoot, ".specforce", "archive", parentPath)
	if _, err := os.Stat(filepath.Join(archiveDir, finalSegment)); err == nil {
		return slug
	}

	entries, err = os.ReadDir(archiveDir)
	if err == nil {
		// Sort entries in reverse order to prioritize newest timestamps
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() > entries[j].Name()
		})

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if timestampRegex.MatchString(name) {
				prefix := timestampRegex.FindString(name)
				if name == prefix+finalSegment {
					return filepath.Join(parentPath, name)
				}
			}
		}
	}

	return slug
}

// GetSpecDir returns the absolute or relative path to the specification directory,
// searching first in .specforce/specs and then in .specforce/archive.
func GetSpecDir(projectRoot string, slug string) (string, bool) {
	resolved := ResolveSlug(projectRoot, slug)

	activePath := filepath.Join(projectRoot, ".specforce", "specs", resolved)
	if _, err := os.Stat(activePath); err == nil {
		return activePath, true
	}

	archivePath := filepath.Join(projectRoot, ".specforce", "archive", resolved)
	if _, err := os.Stat(archivePath); err == nil {
		return archivePath, true
	}

	return "", false
}
