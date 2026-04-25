package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NeonSpinner represents a styled progress indicator with a mandatory label.
type NeonSpinner struct {
	spinner spinner.Model
	label   string
	done    bool
	err     error
}

// NewNeonSpinner creates a new spinner with the specified label.
func NewNeonSpinner(label string) NeonSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(brandCyan)
	return NeonSpinner{
		spinner: s,
		label:   label,
	}
}

// Init initializes the spinner animation.
func (m NeonSpinner) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles spinner animation ticks.
func (m NeonSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

// View renders the spinner and its label.
func (m NeonSpinner) View() string {
	if m.done {
		return fmt.Sprintf("  %s %s\n", SuccessStyle.Render("✓"), m.label)
	}
	if m.err != nil {
		return fmt.Sprintf("  %s %s\n", ErrorStyle.Render("✖"), m.label)
	}
	return fmt.Sprintf("  %s %s\n", m.spinner.View(), m.label)
}

// SetDone marks the operation as successfully completed.
func (m *NeonSpinner) SetDone() {
	m.done = true
}

// SetError marks the operation as failed with an error.
func (m *NeonSpinner) SetError(err error) {
	m.err = err
}
