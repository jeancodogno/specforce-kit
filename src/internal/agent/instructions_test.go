package agent

import (
	"testing"
	"testing/fstest"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func TestInstructionManager_InjectVariables(t *testing.T) {
	config := &core.ProjectConfig{
		Context: map[string]string{
			"project_name": "Test Project",
		},
	}
	m := NewInstructionManager(fstest.MapFS{}, config)

	input := "Welcome to {{project_name}}"
	expected := "Welcome to Test Project"
	got := m.InjectVariables(input)

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInstructionManager_GetInstructions(t *testing.T) {
	mockFS := fstest.MapFS{
		"instructions/test.md": {Data: []byte("Base {{var}}")},
	}
	config := &core.ProjectConfig{
		Instructions: map[string][]string{
			"test": {"Custom"},
		},
		Context: map[string]string{
			"var": "Value",
		},
	}
	m := NewInstructionManager(mockFS, config)

	got, err := m.GetInstructions("test")
	if err != nil {
		t.Fatal(err)
	}

	expected := "Base Value\n\nCustom"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
