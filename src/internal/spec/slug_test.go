package spec

import (
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
