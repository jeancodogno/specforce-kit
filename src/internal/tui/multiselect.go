package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeancodogno/specforce-kit/src/internal/agent"
)

type AgentOption struct {
	ID       string
	Name     string
	Selected bool
	Exists   bool
}

type model struct {
	choices    []AgentOption
	cursor     int
	selected   map[int]struct{}
	toRemove   map[int]struct{}
	confirming bool
	quitting   bool
	aborted    bool
}

func initialModel(available []agent.AgentMetadata, existing []string) model {
	choices := make([]AgentOption, len(available))
	for i, a := range available {
		choices[i] = AgentOption{
			ID:   a.ID,
			Name: a.Name,
		}
	}

	selected := make(map[int]struct{})
	for i, choice := range choices {
		for _, e := range existing {
			if choice.ID == e {
				choices[i].Exists = true
				selected[i] = struct{}{}
				break
			}
		}
	}

	return model{
		choices:  choices,
		selected: selected,
		toRemove: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.confirming {
		return m.handleConfirmation(keyMsg)
	}

	return m.handleNavigation(keyMsg)
}

func (m model) handleConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.quitting = true
		return m, tea.Quit
	case "n", "N", "esc":
		m.confirming = false
		return m, nil
	}
	return m, nil
}

func (m model) handleNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		m.aborted = true
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter", " ":
		m.toggleSelection()
	case "y":
		if len(m.toRemove) > 0 {
			m.confirming = true
		} else {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) toggleSelection() {
	_, ok := m.selected[m.cursor]
	if ok {
		delete(m.selected, m.cursor)
		if m.choices[m.cursor].Exists {
			m.toRemove[m.cursor] = struct{}{}
		}
	} else {
		m.selected[m.cursor] = struct{}{}
		delete(m.toRemove, m.cursor)
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	if m.confirming {
		return m.viewConfirmation()
	}

	return m.viewSelection()
}

func (m model) viewConfirmation() string {
	var agents []string
	for i := range m.toRemove {
		agents = append(agents, DimmedStyle.Render(m.choices[i].Name))
	}

	// High-fidelity Warning Header
	title := ErrorStyle.Bold(true).Render(" ⚠️   DECOMMISSION WARNING ")
	warningText := BodyStyle.Render("You have unselected agents that are already initialized.")
	listText := BodyStyle.Render("The following directories will be ") + ErrorStyle.Render("DELETED") + BodyStyle.Render(":")
	agentList := strings.Join(agents, ", ")
	prompt := BodyStyle.Render("Are you sure you want to proceed? (y/n)")

	content := lipgloss.JoinVertical(lipgloss.Left,
		warningText,
		"",
		listText,
		"↳ "+agentList,
		"",
		prompt,
	)

	// Ensure the title is treated as a separate header to avoid Lipgloss alignment issues with multi-line padding
	frameContent := title + "\n\n" + content

	frame := lipgloss.NewStyle().
		Border(CleanBorder).
		BorderForeground(errorRed).
		Padding(1, 2).
		Render(frameContent)

	return "\n" + frame
}

func (m model) viewSelection() string {
	var rows []string

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ActiveArrowStyle.Render(ArrowGlyph)
		}

		checked := UnselectedBulletStyle.Render(EmptyBulletGlyph)
		status := ReadyStatusStyle.Render("[ READY  ]")
		if _, ok := m.selected[i]; ok {
			checked = SelectedBulletStyle.Render(BulletGlyph)
			status = ActiveStatusStyle.Render("[ ACTIVE ]")
		}

		label := BodyStyle.Render(choice.Name)
		if choice.Exists && m.cursor != i {
			label = DimmedStyle.Render(choice.Name)
		}

		// Use dots for alignment to create a "surgical" look
		dotCount := 35 - len(choice.Name)
		if dotCount < 1 {
			dotCount = 1
		}
		dots := DimmedStyle.Render(strings.Repeat(".", dotCount))

		rows = append(rows, fmt.Sprintf("%s %s %s %s %s", cursor, checked, label, dots, status))
	}

	body := strings.Join(rows, "\n")
	
	// Wrap in CleanBorder frame
	frame := lipgloss.NewStyle().
		Border(CleanBorder).
		BorderForeground(mutedGrey).
		Padding(1, 2).
		Render(HeaderStyle.Render(" SELECT AI AGENTS ") + "\n\n" + body)

	return "\n" + frame + m.viewFooter()
}

func (m model) viewFooter() string {
	footer := "\n" + DimmedStyle.Render("Press space/enter to toggle, 'y' to confirm, 'q' to quit.") + "\n"
	if len(m.toRemove) > 0 {
		footer += "\n" + WarningStyle.Render("⚠️ Some existing agents will be removed upon confirmation.") + "\n"
	}
	return footer
}

// SelectAgents launches the TUI and returns the selected agent IDs and removed agent IDs.
func SelectAgents(available []agent.AgentMetadata, existing []string) ([]string, []string, error) {
	p := tea.NewProgram(initialModel(available, existing))
	m, err := p.Run()
	if err != nil {
		return nil, nil, err
	}

	finalModel := m.(model)
	if finalModel.aborted {
		return nil, nil, fmt.Errorf("aborted")
	}

	var selected []string
	for i := range finalModel.selected {
		selected = append(selected, finalModel.choices[i].ID)
	}

	var removed []string
	for i := range finalModel.toRemove {
		removed = append(removed, finalModel.choices[i].ID)
	}

	return selected, removed, nil
}
