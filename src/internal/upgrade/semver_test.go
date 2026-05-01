package upgrade

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"v1.0.0", "v1.0.1", -1},
		{"1.0.0", "1.0.1", -1},
		{"v1.1.0", "v1.0.9", 1},
		{"1.1.0", "1.1.0", 0},
		{"v1.1.0", "1.1.0", 0},
		{"v1.1.0-alpha", "v1.1.0", -1},
	}

	for _, tt := range tests {
		result := CompareVersions(tt.v1, tt.v2)
		if result != tt.expected {
			t.Errorf("CompareVersions(%s, %s): expected %d, got %d", tt.v1, tt.v2, tt.expected, result)
		}
	}
}

func TestIsNewer(t *testing.T) {
	if !IsNewer("v1.0.0", "v1.0.1") {
		t.Error("expected v1.0.1 to be newer than v1.0.0")
	}
	if IsNewer("v1.0.1", "v1.0.0") {
		t.Error("expected v1.0.0 NOT to be newer than v1.0.1")
	}
	if IsNewer("v1.0.0", "v1.0.0") {
		t.Error("expected v1.0.0 NOT to be newer than v1.0.0")
	}
}
