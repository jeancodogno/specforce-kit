package agent

import (
	"testing"
)

func TestKitManifests(t *testing.T) {
	kitFS, err := GetKitFS()
	if err != nil {
		t.Fatalf("failed to get kit FS: %v", err)
	}

	registry := &Registry{}
	err = registry.Initialize(kitFS, "")
	if err != nil {
		t.Fatalf("failed to initialize registry with kit FS: %v", err)
	}

	// Check open-code
	opencode, ok := registry.GetAgent("open-code")
	if !ok {
		t.Error("expected to find 'open-code' agent in kit")
	} else {
		if opencode.DirName != ".opencode/" {
			t.Errorf("expected open-code DirName '.opencode/', got '%s'", opencode.DirName)
		}
	}

	// Check kilo-code
	kilocode, ok := registry.GetAgent("kilo-code")
	if !ok {
		t.Error("expected to find 'kilo-code' agent in kit")
	} else {
		if kilocode.DirName != ".kilocode/" {
			t.Errorf("expected kilo-code DirName '.kilocode/', got '%s'", kilocode.DirName)
		}
	}

	// Check antigravity (to be sure)
	antigravity, ok := registry.GetAgent("antigravity")
	if !ok {
		t.Error("expected to find 'antigravity' agent in kit")
	} else {
		if antigravity.DirName != ".agent/" {
			t.Errorf("expected antigravity DirName '.agent/', got '%s'", antigravity.DirName)
		}
	}

	// Check kimi-code
	kimicode, ok := registry.GetAgent("kimi-code")
	if !ok {
		t.Error("expected to find 'kimi-code' agent in kit")
	} else {
		if kimicode.DirName != ".kimi/" {
			t.Errorf("expected kimi-code DirName '.kimi/', got '%s'", kimicode.DirName)
		}
	}
}
