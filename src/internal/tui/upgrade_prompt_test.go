package tui

import (
	"bytes"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpgradePrompt(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"y", true},
		{"Y", true},
		{"n", false},
		{"N", false},
		{"\r", false}, // Enter
	}

	for _, tt := range tests {
		m := UpgradePromptModel{Version: "v1.0.0"}
		var in bytes.Buffer
		in.WriteString(tt.input)

		p := tea.NewProgram(m, tea.WithInput(&in))
		finalModel, err := p.Run()
		if err != nil {
			t.Fatalf("failed to run program: %v", err)
		}

		choice := finalModel.(UpgradePromptModel).Choice
		if choice != tt.expected {
			t.Errorf("for input %q, expected choice %v, got %v", tt.input, tt.expected, choice)
		}
	}
}
