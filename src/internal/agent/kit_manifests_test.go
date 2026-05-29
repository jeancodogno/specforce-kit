package agent

import (
	"os"
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
		if antigravity.DirName != ".agents/" {
			t.Errorf("expected antigravity DirName '.agents/', got '%s'", antigravity.DirName)
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

	// Check cursor
	cursor, ok := registry.GetAgent("cursor")
	if !ok {
		t.Error("expected to find 'cursor' agent in kit")
	} else {
		if cursor.DirName != ".cursor/" {
			t.Errorf("expected cursor DirName '.cursor/', got '%s'", cursor.DirName)
		}
	}
}

func TestCursorMapping(t *testing.T) {
	// We use the local FS to test the recent kit.yaml changes
	kitFS := os.DirFS("kit")
	registry := &Registry{}
	err := registry.Initialize(kitFS, "")
	if err != nil {
		t.Fatalf("failed to initialize registry with local kit FS: %v", err)
	}

	cursor, ok := registry.GetAgent("cursor")
	if !ok {
		t.Fatal("expected to find 'cursor' agent in local kit")
	}

	if cursor.Name != "Cursor" {
		t.Errorf("expected Name 'Cursor', got '%s'", cursor.Name)
	}
	if cursor.DirName != ".cursor/" {
		t.Errorf("expected DirName '.cursor/', got '%s'", cursor.DirName)
	}
}
