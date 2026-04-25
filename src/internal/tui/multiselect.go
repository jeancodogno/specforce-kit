package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
		agents = append(agents, m.choices[i].Name)
	}
	s := WarningStyle.Render("⚠️ WARNING: You have unselected existing agents.\n")
	s += BodyStyle.Render(fmt.Sprintf("The following agent directories will be decommissioned: %v\n\n", agents))
	s += BodyStyle.Render("Are you sure you want to proceed? (y/n)")
	return s
}

func (m model) viewSelection() string {
	s := HeaderStyle.Render("Select AI agents to initialize in this project:") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ActiveArrowStyle.Render(ArrowGlyph)
		}

		checked := UnselectedBulletStyle.Render(EmptyBulletGlyph)
		if _, ok := m.selected[i]; ok {
			checked = SelectedBulletStyle.Render(BulletGlyph)
		}

		existsInfo := ""
		if choice.Exists {
			existsInfo = DimmedStyle.Render(" (already exists)")
		}

		s += fmt.Sprintf("%s %s %s%s\n", cursor, checked, BodyStyle.Render(choice.Name), existsInfo)
	}

	s += m.viewFooter()
	return s
}

func (m model) viewFooter() string {
	footer := "\n" + DimmedStyle.Render("Press space/enter to toggle, 'y' to confirm, 'q' to quit.") + "\n"
	if len(m.toRemove) > 0 {
		footer += "\n" + WarningStyle.Render("⚠️ Some existing agents will be removed upon confirmation.") + "\n"
	}
	return footer
}

// SelectAgents launches the TUI and returns the selected agent IDs.
func SelectAgents(available []agent.AgentMetadata, existing []string) ([]string, error) {
	p := tea.NewProgram(initialModel(available, existing))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	finalModel := m.(model)
	if finalModel.aborted {
		return nil, fmt.Errorf("aborted")
	}

	var selected []string
	for i := range finalModel.selected {
		selected = append(selected, finalModel.choices[i].ID)
	}

	return selected, nil
}
