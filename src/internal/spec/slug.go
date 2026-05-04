package spec

import (
	"path/filepath"
	"regexp"
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
