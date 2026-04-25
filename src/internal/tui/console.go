package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
)

// TickMsg is sent when the real-time refresh ticker fires.
type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// ConsoleModel represents the state of the Specforce Console TUI.
type ConsoleModel struct {
	ctx         context.Context
	StateTree   *spec.StateTree
	registry    *spec.Registry
	projectRoot string
	width       int
	height      int
	viewport    viewport.Model
	ready       bool
	err         error
	errorMsg    string
}

// NewConsoleModel initializes a new ConsoleModel with the project state.
func NewConsoleModel(ctx context.Context, tree *spec.StateTree, registry *spec.Registry, root string) *ConsoleModel {
	return &ConsoleModel{
		ctx:         ctx,
		StateTree:   tree,
		registry:    registry,
		projectRoot: root,
	}
}

// Init initializes the Bubbletea program.
func (m *ConsoleModel) Init() tea.Cmd {
	return tea.Batch(tick(), nil)
}

// ErrorMsg is sent when an external operation fails.
type ErrorMsg struct {
	Err error
}

// Update handles terminal events and updates the model state.
func (m *ConsoleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(m.width, m.height)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = m.height
		}

		m.viewport.SetContent(m.renderDashboard())

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.refreshState()
		default:
			// Handle viewport scrolling for keys like up/down/pgup/pgdown
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case TickMsg:
		m.refreshState()
		cmds = append(cmds, tick())

	case ErrorMsg:
		m.err = msg.Err
	}

	// Sync content after any change
	if m.ready {
		m.viewport.SetContent(m.renderDashboard())
	}

	return m, tea.Batch(cmds...)
}

func (m *ConsoleModel) refreshState() {
	tree, err := spec.ScanProject(m.ctx, m.projectRoot, m.registry)
	if err == nil {
		m.StateTree = tree
	}
}

// View renders the TUI to the terminal.
func (m *ConsoleModel) View() string {
	if m.err != nil {
		return ErrorStyle.Render("Console Error: " + m.err.Error())
	}

	if !m.ready {
		return "Initializing..."
	}

	return m.viewport.View()
}

func (m *ConsoleModel) renderDashboard() string {
	var s strings.Builder

	// 1. Logo
	s.WriteString(m.renderLogo())
	s.WriteString("\n\n")

	// Check if everything is empty
	isEmpty := true
	for _, items := range m.StateTree.Categories {
		if len(items) > 0 {
			isEmpty = false
			break
		}
	}

	if isEmpty {
		s.WriteString(m.renderEmptyStateContent())
		return s.String()
	}

	// 3. Sections
	categories := []struct {
		cat   spec.StateCategory
		label string
	}{
		{spec.CategoryConstitution, "🏛️  CONSTITUTION (Foundation)"},
		{spec.CategoryActiveSpecs, "📋 ACTIVE SPECIFICATIONS (Planning Phase)"},
		{spec.CategoryImplementations, "🚀 ACTIVE IMPLEMENTATIONS (Execution Phase)"},
		{spec.CategoryArchived, "📦 RECENTLY ARCHIVED (Completed Features)"},
	}

	for _, config := range categories {
		items, ok := m.StateTree.Categories[config.cat]
		if !ok || len(items) == 0 {
			continue
		}

		s.WriteString(HeaderStyle.Render(config.label))
		s.WriteString("\n")

		for _, item := range items {
			s.WriteString(m.renderItem(item))
		}
		s.WriteString("\n")
	}

	// 4. Error Message
	if m.errorMsg != "" {
		s.WriteString("\n")
		s.WriteString(RenderErrorBadge(m.errorMsg))
	}

	return lipgloss.NewStyle().Padding(1, 2).Width(m.width - 4).Render(s.String())
}

func (m *ConsoleModel) renderLogo() string {
	return GenerateLogo(true) + "\n"
}

func (m *ConsoleModel) renderItem(item spec.StateItem) string {
	var s strings.Builder

	marker := "  "
	style := BodyStyle

	// High-density info based on category
	switch item.Category {
	case spec.CategoryConstitution:
		icon := MutedStyle.Render("○")
		if item.Progress == 100 {
			icon = FinishedStyle.Render("◉")
		}
		name := style.Render(item.Name)
		metadata := MutedStyle.Render(fmt.Sprintf("Active (%s)", item.Path))
		fmt.Fprintf(&s, "%s%s %-15s : %s\n", marker, icon, name, metadata)
	case spec.CategoryActiveSpecs:
		icon := MutedStyle.Render("○")
		if item.Progress == 100 {
			icon = FinishedStyle.Render("◉")
		}
		artifactStatus := MutedStyle.Render(fmt.Sprintf("(%d/%d artifacts)", item.ArtifactCount, item.ArtifactTotal))
		fmt.Fprintf(&s, "%s%s %s - %s %s\n", marker, icon, style.Render(item.Slug), RenderProgressBar(item.Progress, 20), artifactStatus)
	case spec.CategoryImplementations:
		icon := MutedStyle.Render("○")
		if item.Progress == 100 {
			icon = FinishedStyle.Render("◉")
		} else if item.AnyTaskWorking {
			icon = InProgressStyle.Render("◉")
		}
		taskStatus := MutedStyle.Render(fmt.Sprintf("(%d/%d tasks)", item.TaskCount, item.TaskTotal))
		fmt.Fprintf(&s, "%s%s %s - %s %s\n", marker, icon, style.Render(item.Slug), RenderProgressBar(item.Progress, 20), taskStatus)
		// Contextual Task Detail
		var detail string
		if item.Progress == 100 {
			detail = "All tasks complete"
		} else if item.AnyTaskWorking {
			detail = fmt.Sprintf("Working on: %s %s", item.CurrentTaskID, item.CurrentTask)
		} else {
			detail = fmt.Sprintf("Next task: %s %s", item.CurrentTaskID, item.CurrentTask)
		}
		fmt.Fprintf(&s, "  %s %s\n", marker, MutedStyle.Render("↳ "+detail))
	case spec.CategoryArchived:
		icon := MutedStyle.Render("○")
		metadata := MutedStyle.Render(fmt.Sprintf("(Archived on: %s)", item.ArchivedDate))
		fmt.Fprintf(&s, "%s%s %-20s %s\n", marker, icon, style.Render(item.Slug), metadata)
	}

	return s.String()
}

func (m *ConsoleModel) renderEmptyStateContent() string {
	return `[ NO ACTIVE SPECIFICATIONS FOUND ]

Your .specforce/ directory is currently empty or uninitialized.

To get started:
1. Run 'specforce init' to set up the framework.
2. Run 'specforce spec init <slug>' to create your first spec.

Press 'q' or 'Esc' to exit.`
}
