package core

import (
	"testing"
)

func TestToolPrefixes(t *testing.T) {
	expected := ".cursor/"
	found := false
	for _, p := range ToolPrefixes {
		if p == expected {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected ToolPrefixes to contain %q", expected)
	}

	for _, p := range ToolPrefixes {
		if p == ".gemini/" {
			t.Errorf("ToolPrefixes should not contain .gemini/")
		}
	}
}
