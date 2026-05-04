package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

type ImplementationTask struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	State        string   `json:"state"`
	Target       string   `json:"target"`
	Context      string   `json:"context"`
	ActionSteps  []string `json:"action_steps"`
	Verification string   `json:"verification"`
}

type Phase struct {
	ID    string               `json:"id"`
	Title string               `json:"title"`
	Tasks []ImplementationTask `json:"tasks"`
}

type ImplementationReport struct {
	Name                  string               `json:"name"`
	Status                string               `json:"status"` // ready | blocked
	MissingArtifacts      []string             `json:"missing_artifacts,omitempty"`
	ContextFiles          []string             `json:"context_files"`
	Instructions          []string             `json:"instructions,omitempty"`
	Phases                []Phase              `json:"phases"`
	ExecutionStrategy     string               `json:"execution_strategy"`
	PreemptiveMitigations string               `json:"preemptive_mitigations"`
}

// Tasks returns a flat list of tasks across all phases for backward compatibility.
func (r *ImplementationReport) Tasks() []ImplementationTask {
	var tasks []ImplementationTask
	for _, p := range r.Phases {
		tasks = append(tasks, p.Tasks...)
	}
	return tasks
}

// FindProjectRoot looks for the .specforce directory in the current or parent directories.
func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, ".specforce")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf(".specforce directory not found")
}

// CheckTriadArtifacts verifies if requirements.md, design.md, and tasks.md exist.
func CheckTriadArtifacts(projectRoot, slug string) (bool, []string) {
	required := []string{"requirements.md", "design.md", "tasks.md"}
	missing := []string{}

	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)

	for _, art := range required {
		if _, err := os.Stat(filepath.Join(specDir, art)); os.IsNotExist(err) {
			missing = append(missing, art)
		}
	}

	return len(missing) == 0, missing
}

// GetContextFiles returns absolute paths for all files in the spec directory and global docs.
func GetContextFiles(projectRoot, slug string) ([]string, error) {
	files := []string{}

	specDir := filepath.Join(projectRoot, ".specforce", "specs", slug)

	// Helper to add files from a directory
	addFiles := func(dir string) error {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil
		}

		return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				abs, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				files = append(files, abs)
			}
			return nil
		})
	}

	if err := addFiles(specDir); err != nil {
		return nil, err
	}

	return files, nil
}

// ParseTasks extracts implementation details from tasks.md.
func ParseTasks(ctx context.Context, projectRoot, slug string) (*ImplementationReport, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tasksPath, err := core.SecurePath(projectRoot, filepath.Join(".specforce", "specs", slug, "tasks.md"))
	if err != nil {
		return nil, fmt.Errorf("security: %w", err)
	}
	// #nosec G304 - Path is secured by SecurePath
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file %s: %w", tasksPath, errors.Join(core.ErrInvalidSpecFile, err))
	}

	report := &ImplementationReport{
		Name:   slug,
		Phases: []Phase{},
	}

	extractStrategyAndMitigations(content, report)

	if err := extractTasksFromContent(ctx, content, report); err != nil {
		return nil, err
	}

	report.Status = calculateReportStatus(report.Tasks())

	return report, nil
}

func extractStrategyAndMitigations(content []byte, report *ImplementationReport) {
	strategyRegex := regexp.MustCompile(`(?s)## 1\. Execution Strategy\n(.*?)(?:\n##|$)`)
	strategyMatch := strategyRegex.FindSubmatch(content)
	if len(strategyMatch) > 1 {
		report.ExecutionStrategy = strings.TrimSpace(string(strategyMatch[1]))
	}

	mitigationsRegex := regexp.MustCompile(`(?s)## 3\. Pre-emptive Mitigations\n(.*?)(?:\n##|$)`)
	mitigationsMatch := mitigationsRegex.FindSubmatch(content)
	if len(mitigationsMatch) > 1 {
		report.PreemptiveMitigations = strings.TrimSpace(string(mitigationsMatch[1]))
	}
}

func extractTasksFromContent(ctx context.Context, content []byte, report *ImplementationReport) error {
	phaseHeaderRegex := regexp.MustCompile(`(?m)^### Phase (\d+): (.*)$`)
	taskHeaderRegex := regexp.MustCompile(`(?m)^(#{3,4}|- \[[ xX/]?\]) (T[\d.]+): (.*)$`)

	phaseMatches := phaseHeaderRegex.FindAllSubmatchIndex(content, -1)
	taskMatches := taskHeaderRegex.FindAllSubmatchIndex(content, -1)

	seenTasks := make(map[string]bool)
	for _, tm := range taskMatches {
		if err := ctx.Err(); err != nil {
			return err
		}

		task := parseTaskBlock(content, tm, taskMatches, phaseMatches)
		if seenTasks[task.ID] {
			continue
		}
		seenTasks[task.ID] = true

		// Find or create phase
		currentPhaseIdx := findPhaseIdxForTask(content, tm[0], phaseMatches, report)
		report.Phases[currentPhaseIdx].Tasks = append(report.Phases[currentPhaseIdx].Tasks, task)
	}

	return nil
}

func parseTaskBlock(content []byte, tm []int, taskMatches, phaseMatches [][]int) ImplementationTask {
	prefix := string(content[tm[2]:tm[3]])
	task := ImplementationTask{
		ID:    string(content[tm[4]:tm[5]]),
		Title: strings.TrimSpace(string(content[tm[6]:tm[7]])),
	}

	nextHeaderPos := findNextHeaderPos(content, tm[0], tm[1], taskMatches, phaseMatches)
	taskBlock := content[tm[0]:nextHeaderPos]

	task.State = extractField(taskBlock, `\*\*State:\*\* \x60?\[([^\]\x60\n]*)\]\x60?`)
	isChecklist := strings.HasPrefix(prefix, "- [")
	if task.State == "" && isChecklist {
		task.State = mapCheckboxToState(prefix)
	}

	task.Target = extractField(taskBlock, `\*\*Target:\*\* \x60?([^\x60\n]*)\x60?`)
	task.Context = extractField(taskBlock, `\*\*Context:\*\* \x60?([^\x60\n]*)\x60?`)
	task.Verification = extractField(taskBlock, `(?s)\*\*Verification \(TDD\):\*\*\n(.*?)(?:\n\n|\n###|\n##|$)`)

	actionStepsRegex := regexp.MustCompile(`(?m)^- (.*)$`)
	actionMatches := actionStepsRegex.FindAllSubmatch(taskBlock, -1)
	for i, am := range actionMatches {
		if i == 0 && isChecklist {
			continue
		}
		step := strings.TrimSpace(string(am[1]))
		if step != "" {
			task.ActionSteps = append(task.ActionSteps, step)
		}
	}

	return task
}

func findNextHeaderPos(content []byte, taskStart, taskEnd int, taskMatches, phaseMatches [][]int) int {
	nextHeaderPos := len(content)
	for _, nextTm := range taskMatches {
		if nextTm[0] > taskStart && nextTm[0] < nextHeaderPos {
			nextHeaderPos = nextTm[0]
			break
		}
	}
	for _, pm := range phaseMatches {
		if pm[0] > taskStart && pm[0] < nextHeaderPos {
			nextHeaderPos = pm[0]
			break
		}
	}
	nextSectionRegex := regexp.MustCompile(`(?m)^## `)
	sectionMatch := nextSectionRegex.FindIndex(content[taskEnd:])
	if sectionMatch != nil {
		absSectionPos := taskEnd + sectionMatch[0]
		if absSectionPos < nextHeaderPos {
			nextHeaderPos = absSectionPos
		}
	}
	return nextHeaderPos
}

func mapCheckboxToState(prefix string) string {
	char := prefix[3:4]
	switch strings.ToLower(char) {
	case "x":
		return "FINISHED"
	case "/":
		return "IN-PROGRESS"
	default:
		return "PENDING"
	}
}

func findPhaseIdxForTask(content []byte, taskStart int, phaseMatches [][]int, report *ImplementationReport) int {
	matchIdx := -1
	for i := len(phaseMatches) - 1; i >= 0; i-- {
		if phaseMatches[i][0] < taskStart {
			matchIdx = i
			break
		}
	}

	if matchIdx == -1 {
		if len(report.Phases) == 0 || report.Phases[0].ID != "0" {
			report.Phases = append([]Phase{{ID: "0", Title: "Initial Tasks", Tasks: []ImplementationTask{}}}, report.Phases...)
		}
		return 0
	}

	phaseID := string(content[phaseMatches[matchIdx][2]:phaseMatches[matchIdx][3]])
	phaseTitle := strings.TrimSpace(string(content[phaseMatches[matchIdx][4]:phaseMatches[matchIdx][5]]))

	for idx, p := range report.Phases {
		if p.ID == phaseID && p.Title == phaseTitle {
			return idx
		}
	}

	report.Phases = append(report.Phases, Phase{ID: phaseID, Title: phaseTitle, Tasks: []ImplementationTask{}})
	return len(report.Phases) - 1
}

func calculateReportStatus(tasks []ImplementationTask) string {
	if len(tasks) == 0 {
		return "ready"
	}
	allFinished := true
	allPending := true
	for _, t := range tasks {
		st := strings.ToUpper(t.State)
		if st != "FINISHED" {
			allFinished = false
		}
		if st != "PENDING" && st != "" {
			allPending = false
		}
	}

	if allFinished {
		return "finished"
	} else if !allPending {
		return "in-progress"
	}
	return "ready"
}

func extractField(block []byte, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindSubmatch(block)
	if len(match) > 1 {
		return strings.TrimSpace(string(match[1]))
	}
	return ""
}
