package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UpgradeProgressModel is a Bubble Tea model for showing upgrade progress.
type UpgradeProgressModel struct {
	Percent  int
	Status   string
	Finished bool
	Err      error
}

// UpgradeProgressMsg is sent to update the progress bar.
type UpgradeProgressMsg struct {
	Percent int
	Status  string
}

// UpgradeFinishedMsg is sent when the upgrade is complete.
type UpgradeFinishedMsg struct {
	Err error
}

func (m UpgradeProgressModel) Init() tea.Cmd {
	return nil
}

func (m UpgradeProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpgradeProgressMsg:
		m.Percent = msg.Percent
		m.Status = msg.Status
		return m, nil
	case UpgradeFinishedMsg:
		m.Finished = true
		m.Err = msg.Err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m UpgradeProgressModel) View() string {
	if m.Finished {
		if m.Err != nil {
			return fmt.Sprintf("\n  %s %v\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("✘"), m.Err)
		}
		return fmt.Sprintf("\n  %s Upgrade complete!\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✔"))
	}

	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	bar := RenderProgressBar(m.Percent, 40)

	return fmt.Sprintf("\n  %s %s\n  %s\n", statusStyle.Render("›"), m.Status, bar)
}

// NewUpgradeProgressProgram creates a new tea.Program for showing upgrade progress.
func NewUpgradeProgressProgram() *tea.Program {
	return tea.NewProgram(UpgradeProgressModel{Status: "Starting upgrade..."})
}
