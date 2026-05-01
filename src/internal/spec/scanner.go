package spec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// StateCategory defines the top-level grouping for project state items.
type StateCategory string

const (
	CategoryConstitution    StateCategory = "Constitution"
	CategoryActiveSpecs      StateCategory = "Active Specs"
	CategoryImplementations StateCategory = "Active Implementations"
	CategoryArchived        StateCategory = "Archived"
)

// StateItem represents a single navigatable entry in the Specforce Console.
type StateItem struct {
	Slug          string        `json:"slug"`
	Name          string        `json:"name"`
	Path          string        `json:"path"`
	Category      StateCategory `json:"category"`
	Status        string        `json:"status"` // PENDING | IN-PROGRESS | FINISHED
	Progress      int           `json:"progress"`
	Description   string        `json:"description"`
	ArtifactCount int           `json:"artifact_count"`
	ArtifactTotal int           `json:"artifact_total"`
	TaskCount     int           `json:"task_count"`
	TaskTotal     int           `json:"task_total"`
	AnyTaskWorking bool          `json:"any_task_working"`
	CurrentTaskID string        `json:"current_task_id"`
	CurrentTask   string        `json:"current_task"`
	ArchivedDate  string        `json:"archived_date"`
	Worktree      string        `json:"worktree"`
}

// StateTree represents the full hierarchical state of a Specforce project.
type StateTree struct {
	Categories map[StateCategory][]StateItem `json:"categories"`
}

// NewStateTree initializes a new, empty state tree.
func NewStateTree() *StateTree {
	return &StateTree{
		Categories: make(map[StateCategory][]StateItem),
	}
}

// ScanProject synchronously parses the .specforce directory to populate a StateTree.
func ScanProject(ctx context.Context, projectRoot string, registry *Registry) (*StateTree, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tree := NewStateTree()

	worktrees, err := discoverWorktrees(ctx, projectRoot)
	if err != nil {
		// Log error but continue with current root
		worktrees = []worktreeInfo{{Path: projectRoot, Branch: ""}}
	}

	for _, wt := range worktrees {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Robust main root check
		isMainRoot := evalPath(wt.Path) == evalPath(projectRoot)

		// 1. Scan Constitution documents (.specforce/docs/) - ONLY MAIN ROOT
		if isMainRoot {
			if err := scanConstitution(ctx, wt.Path, tree); err != nil {
				return nil, err
			}
		}

		// 2. Scan Active Specs and Implementations (.specforce/specs/) - ALL WORKTREES
		branchLabel := ""
		if !isMainRoot {
			branchLabel = wt.Branch
		}
		if err := scanActiveSpecs(ctx, wt.Path, tree, registry, branchLabel); err != nil {
			return nil, err
		}

		// 3. Scan Archived specs (.specforce/archive/) - ALL WORKTREES
		if err := scanArchivedSpecs(ctx, wt.Path, tree, branchLabel); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func evalPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	eval, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return abs
	}
	return eval
}

func scanConstitution(ctx context.Context, projectRoot string, tree *StateTree) error {
	docsDir := filepath.Join(projectRoot, ".specforce", "docs")
	entries, err := os.ReadDir(docsDir)
	if err != nil {
		return nil // No docs dir is fine
	}

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "current-state.md" || entry.Name() == "memorial.md" {
			continue
		}
		slug := strings.TrimSuffix(entry.Name(), ".md")
		tree.Categories[CategoryConstitution] = append(tree.Categories[CategoryConstitution], StateItem{
			Slug:     slug,
			Name:     cases.Title(language.Und).String(strings.ReplaceAll(slug, "-", " ")),
			Path:     filepath.Join(projectRoot, ".specforce", "docs", entry.Name()),
			Category: CategoryConstitution,
			Status:   "FINISHED",
			Progress: 100,
		})
	}
	return nil
}

func scanActiveSpecs(ctx context.Context, projectRoot string, tree *StateTree, registry *Registry, worktree string) error {
	specsDir := filepath.Join(projectRoot, ".specforce", "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		item := scanSingleActiveSpec(ctx, projectRoot, entry.Name(), registry, worktree)
		tree.Categories[item.Category] = append(tree.Categories[item.Category], item)
	}
	return nil
}

func scanSingleActiveSpec(ctx context.Context, projectRoot, slug string, registry *Registry, worktree string) StateItem {
	item := StateItem{
		Slug:     slug,
		Name:     slug,
		Path:     filepath.Join(projectRoot, ".specforce", "specs", slug),
		Category: CategoryActiveSpecs,
		Status:   "PENDING",
		Worktree: worktree,
	}

	if status, err := GetStatus(ctx, projectRoot, slug, registry); err == nil {
		item.Progress = status.Progress
		item.ArtifactCount = status.Found
		item.ArtifactTotal = status.Total
	}

	tasksPath := filepath.Join(projectRoot, ".specforce", "specs", slug, "tasks.md")
	if _, err := os.Stat(tasksPath); err == nil {
		item.Category = CategoryImplementations
		if report, err := ParseTasks(ctx, projectRoot, slug); err == nil {
			item.Status = strings.ToUpper(report.Status)
			item.TaskTotal = len(report.Tasks())
			foundFinished := 0
			for _, t := range report.Tasks() {
				state := strings.ToUpper(t.State)
				if state == "FINISHED" || state == "X" {
					foundFinished++
				} else {
					if state == "IN-PROGRESS" {
						item.AnyTaskWorking = true
					}
					if item.CurrentTaskID == "" {
						item.CurrentTaskID = t.ID
						item.CurrentTask = t.Title
					}
				}
			}
			item.TaskCount = foundFinished
			if item.TaskTotal > 0 {
				item.Progress = (foundFinished * 100) / item.TaskTotal
			}
			if item.CurrentTask == "" && foundFinished == item.TaskTotal {
				item.CurrentTask = "All tasks complete"
			}
		}
	}
	return item
}

func scanArchivedSpecs(ctx context.Context, projectRoot string, tree *StateTree, worktree string) error {
	archiveDir := filepath.Join(projectRoot, ".specforce", "archive")
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil
	}

	archivedItems := []StateItem{}
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		date := "2026-04-16"
		if info, err := entry.Info(); err == nil {
			date = info.ModTime().Format("2006-01-02")
		}

		archivedItems = append(archivedItems, StateItem{
			Slug:         entry.Name(),
			Name:         entry.Name(),
			Path:         filepath.Join(projectRoot, ".specforce", "archive", entry.Name()),
			Category:     CategoryArchived,
			Status:       "FINISHED",
			Progress:     100,
			ArchivedDate: date,
			Worktree:     worktree,
		})
	}

	sort.Slice(archivedItems, func(i, j int) bool {
		return archivedItems[i].ArchivedDate > archivedItems[j].ArchivedDate
	})

	if len(archivedItems) > 10 {
		archivedItems = archivedItems[:10]
	}

	// We append instead of overwriting because we might have multiple worktrees
	tree.Categories[CategoryArchived] = append(tree.Categories[CategoryArchived], archivedItems...)
	return nil
}

type worktreeInfo struct {
	Path   string
	Branch string
}

func discoverWorktrees(ctx context.Context, projectRoot string) ([]worktreeInfo, error) {
	// git worktree list --porcelain
	cmd := exec.CommandContext(ctx, "git", "worktree", "list", "--porcelain")
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		// Not a git repo or git not installed - return current root as a fallback
		return []worktreeInfo{{Path: projectRoot, Branch: ""}}, nil
	}

	return parseWorktreePorcelain(string(out)), nil
}

func parseWorktreePorcelain(output string) []worktreeInfo {
	var worktrees []worktreeInfo
	var current worktreeInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if current.Path != "" {
				if _, err := os.Stat(current.Path); err == nil {
					worktrees = append(worktrees, current)
				}
			}
			current = worktreeInfo{Path: strings.TrimPrefix(line, "worktree ")}
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
			// Simplified branch name
			current.Branch = strings.TrimPrefix(current.Branch, "refs/heads/")
		}
	}
	// Add the last one
	if current.Path != "" {
		if _, err := os.Stat(current.Path); err == nil {
			worktrees = append(worktrees, current)
		}
	}

	return worktrees
}
