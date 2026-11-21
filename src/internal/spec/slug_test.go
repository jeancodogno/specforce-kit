package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareSlug(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		contains string
		wantsTS  bool
	}{
		{
			name:     "simple slug",
			raw:      "my-feature",
			contains: "-my-feature",
			wantsTS:  true,
		},
		{
			name:     "nested path",
			raw:      "docs/specs/my-feature",
			contains: "docs/specs/",
			wantsTS:  true,
		},
		{
			name:     "already has timestamp",
			raw:      "20230101-1200-already-done",
			contains: "20230101-1200-already-done",
			wantsTS:  false,
		},
		{
			name:     "handles leading hyphen",
			raw:      "-feature",
			contains: "-feature",
			wantsTS:  true,
		},
		{
			name:     "handles double slashes",
			raw:      "team-a//feature",
			contains: "team-a/",
			wantsTS:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrepareSlug(tt.raw)
			base := filepath.Base(got)
			if tt.wantsTS {
				if !timestampRegex.MatchString(base) {
					t.Errorf("PrepareSlug(%q) = %v, expected timestamp prefix in base", tt.raw, got)
				}
			} else {
				if got != tt.raw {
					t.Errorf("PrepareSlug(%q) = %v, expected no change", tt.raw, got)
				}
			}
			// For nested paths, verify the directory structure remains
			if !strings.Contains(filepath.ToSlash(got), tt.contains) {
				t.Errorf("PrepareSlug(%q) = %v, should contain %v", tt.raw, got, tt.contains)
			}
		})
	}
}

func TestResolveSlug(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-resolve-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	specsDir := filepath.Join(tmpDir, ".specforce", "specs")
	archiveDir := filepath.Join(tmpDir, ".specforce", "archive")

	// 1. Exact match in specs
	_ = os.MkdirAll(filepath.Join(specsDir, "exact"), 0755)
	
	// 2. Timestamped match in specs
	_ = os.MkdirAll(filepath.Join(specsDir, "20260520-1234-fuzzy"), 0755)
	
	// 3. Sub-path match
	_ = os.MkdirAll(filepath.Join(specsDir, "team-a", "20260520-1111-api"), 0755)

	// 4. Match in archive
	_ = os.MkdirAll(filepath.Join(archiveDir, "20260521-0000-old"), 0755)

	// 5. Multiple timestamps for same slug
	_ = os.MkdirAll(filepath.Join(specsDir, "20260520-1000-conflict"), 0755)
	_ = os.MkdirAll(filepath.Join(specsDir, "20260521-1000-conflict"), 0755)

	tests := []struct {
		name     string
		slug     string
		expected string
	}{
		{"Exact match", "exact", "exact"},
		{"Fuzzy match", "fuzzy", "20260520-1234-fuzzy"},
		{"Sub-path fuzzy", "team-a/api", filepath.Join("team-a", "20260520-1111-api")},
		{"Archive fuzzy", "old", "20260521-0000-old"},
		{"Prioritize Newest", "conflict", "20260521-1000-conflict"},
		{"No match", "none", "none"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveSlug(tmpDir, tt.slug)
			if got != tt.expected {
				t.Errorf("ResolveSlug(%q) = %q, expected %q", tt.slug, got, tt.expected)
			}
		})
	}
}
