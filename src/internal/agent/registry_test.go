package agent

import (
	"testing"
	"testing/fstest"
)

func TestRegistry_Initialize(t *testing.T) {
	mockFS := fstest.MapFS{
		"kit.yaml": {
			Data: []byte(`
tools:
  claude:
    name: "Claude"
    description: "Claude agent"
    target: ".claude"
  gemini:
    name: "Gemini"
    description: "Gemini agent"
    target: ".gemini"
`),
		},
		"some-other-file.txt": {
			Data: []byte("ignored"),
		},
	}

	registry := &Registry{}
	err := registry.Initialize(mockFS)
	if err != nil {
		t.Fatalf("failed to initialize registry: %v", err)
	}

	agents := registry.GetAgents()
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}

	claude, ok := registry.GetAgent("claude")
	if !ok {
		t.Error("expected to find 'claude' agent")
	}
	if claude.Name != "Claude" {
		t.Errorf("expected Name 'Claude', got '%s'", claude.Name)
	}

	_, ok = registry.GetAgent("non-existent")
	if ok {
		t.Error("did not expect to find 'non-existent' agent")
	}
}
