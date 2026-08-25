package agent

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestImplementBlueprintGuardrails(t *testing.T) {
	kitFS, err := GetKitFS()
	if err != nil {
		t.Fatalf("failed to get kit FS: %v", err)
	}

	// 1. Check implement.yaml
	implData, err := os.ReadFile("kit/commands/implement.yaml")
	if err != nil {
		t.Fatalf("failed to read implement.yaml: %v", err)
	}
	implContent := string(implData)

	requiredWorkerDirectives := []string{
		"NON-NEGOTIABLE WORKER GUARDRAILS",
		"Uncertainty & Ambiguity",
		"Multiple Interpretations",
		"Flawed Approach",
		"Scope Containment",
		"Do NOT implement features that were not requested",
		"Do NOT refactor code that was not requested",
		"Do NOT implement tests for impossible scenarios",
		"Do NOT remove pre-existing dead code unless explicitly requested",
		"Test Integrity & Invariance (Tests = Specification)",
		"NEVER remove tests to reduce the failure count",
		"NEVER weaken existing tests to make them pass",
		"NEVER use mechanisms to skip, ignore, or bypass tests",
		"If a test is genuinely wrong/broken, STOP and confirm with the user",
		"Mid/Post-Implementation Spec Gate",
	}

	for _, directive := range requiredWorkerDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("implement.yaml missing mandatory directive: %q", directive)
		}
	}

	// 2. Check engineering.yaml
	engData, err := os.ReadFile("artifacts/constitution/engineering.yaml")
	if err != nil {
		t.Fatalf("failed to read engineering.yaml: %v", err)
	}
	engContent := string(engData)

	requiredEngDirectives := []string{
		"AI Coding Constraints & Test Invariance",
		"Pre-Coding Protocol",
		"During Coding Protocol",
		"Post-Implementation & Change Protocol",
		"Tests are the specification",
	}

	for _, directive := range requiredEngDirectives {
		if !strings.Contains(engContent, directive) {
			t.Errorf("engineering.yaml missing mandatory directive: %q", directive)
		}
	}
	_ = kitFS
}

func TestSpecBlueprintTieredSizing(t *testing.T) {
	specData, err := os.ReadFile("kit/commands/spec.yaml")
	if err != nil {
		t.Fatalf("failed to read spec.yaml: %v", err)
	}

	var parsed map[string]any
	if err := yaml.Unmarshal(specData, &parsed); err != nil {
		t.Fatalf("spec.yaml is not valid YAML: %v", err)
	}

	specContent := string(specData)

	requiredDirectives := []string{
		"Scope & Complexity Assessment (Spec Sizing)",
		"`small`: Trivial / single component / quick fix. Skips exhaustive 5-dimension grill, conducts single-turn confirmation, generates only `tasks.md`.",
		"`medium`: Standard scoped feature or bug fix. Focuses on business rules and edge cases, generates `requirements.md` and `tasks.md`.",
		"`large`: Multi-component or persistent change. Standard consultative grill (3-5 dimensions), generates full triad (`requirements.md`, `design.md`, `tasks.md`).",
		"`complex`: High architectural risk or cross-cutting redesign. Adversarial grill across all 5 dimensions + gray areas exploration, generates full triad + deep notes.",
		"specforce spec init <slug> --type <feature|bug> --size <small|medium|large|complex>",
		"Mid-Flight Scope Changes & Resizing Directive",
		"specforce spec resize <slug> --size <small|medium|large|complex>",
		"If Promoted",
		"If Demoted",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(specContent, directive) {
			t.Errorf("spec.yaml missing mandatory directive: %q", directive)
		}
	}
}

