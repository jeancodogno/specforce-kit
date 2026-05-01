package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UpgradePromptModel is a Bubble Tea model for the upgrade confirmation.
type UpgradePromptModel struct {
	Version  string
	Choice   bool
	Quitting bool
}

func (m UpgradePromptModel) Init() tea.Cmd {
	return nil
}

func (m UpgradePromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.Choice = true
			m.Quitting = true
			return m, tea.Quit
		case "n", "N", "esc", "enter":
			m.Choice = false
			m.Quitting = true
			return m, tea.Quit
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m UpgradePromptModel) View() string {
	if m.Quitting {
		return ""
	}

	promptStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF"))

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Bold(true)

	return fmt.Sprintf("\n %s Do you want to upgrade to version %s now? [y/N] ", 
		promptStyle.Render("?"), 
		versionStyle.Render(m.Version))
}

// PromptForUpgrade displays an interactive prompt and returns true if the user confirmed.
func PromptForUpgrade(version string) (bool, error) {
	p := tea.NewProgram(UpgradePromptModel{Version: version})
	m, err := p.Run()
	if err != nil {
		return false, err
	}
	return m.(UpgradePromptModel).Choice, nil
}
