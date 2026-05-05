package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// taskValidationState tracks the context of the tasks.md parser.
type taskValidationState struct {
	currentPhase     int
	nextTaskIdx      int
	lastPhaseLine    int
	taskCountInPhase int
	errors           []string
}

// taskBlock tracks mandatory fields within a task block.
type taskBlock struct {
	id              string
	line            int
	hasTarget       bool
	hasContext      bool
	hasActionHeader bool
	hasActionItems  bool
	hasVerify       bool
	inActionSteps   bool
}

// ValidateTasks performs an exhaustive structural and content validation of tasks.md.
func ValidateTasks(ctx context.Context, projectRoot, slug string) ([]string, error) {
	tasksPath, err := core.SecurePath(projectRoot, filepath.Join(".specforce", "specs", slug, "tasks.md"))
	if err != nil {
		return nil, fmt.Errorf("security: %w", err)
	}

	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks.md: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	state := &taskValidationState{nextTaskIdx: 1}
	var currentTask *taskBlock

	phaseRegex := regexp.MustCompile(`^### Phase (\d+): (.*)$`)
	taskRegex := regexp.MustCompile(`^(#{4}|- \[[ xX/]?\]) (T(\d+)\.(\d+)): (.*)$`)

	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if pm := phaseRegex.FindStringSubmatch(trimmed); pm != nil {
			currentTask = state.handlePhaseHeader(pm, lineNum, currentTask)
			continue
		}

		if tm := taskRegex.FindStringSubmatch(trimmed); tm != nil {
			currentTask = state.handleTaskHeader(tm, lineNum, currentTask)
			continue
		}

		if currentTask != nil {
			updateTaskBlockState(currentTask, trimmed)
		}
	}

	state.validateLastTask(currentTask)
	state.checkEmptyPhase()

	return state.errors, nil
}

func (s *taskValidationState) handlePhaseHeader(match []string, lineNum int, currentTask *taskBlock) *taskBlock {
	s.validateLastTask(currentTask)
	s.checkEmptyPhase()

	phaseID, _ := strconv.Atoi(match[1])
	if phaseID != s.currentPhase+1 {
		s.errors = append(s.errors, fmt.Sprintf("Phase ID %d (line %d) is out of sequence, expected %d", phaseID, lineNum, s.currentPhase+1))
	}
	s.currentPhase = phaseID
	s.nextTaskIdx = 1
	s.taskCountInPhase = 0
	s.lastPhaseLine = lineNum
	return nil
}

func (s *taskValidationState) handleTaskHeader(match []string, lineNum int, currentTask *taskBlock) *taskBlock {
	s.validateLastTask(currentTask)

	if s.currentPhase == 0 {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) found before any Phase definition", match[2], lineNum))
	}

	taskID := match[2]
	pID, _ := strconv.Atoi(match[3])
	tIdx, _ := strconv.Atoi(match[4])

	if pID != s.currentPhase {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) does not match the parent Phase %d", taskID, lineNum, s.currentPhase))
	}
	if tIdx != s.nextTaskIdx {
		s.errors = append(s.errors, fmt.Sprintf("Task sequence gap at line %d: expected T%d.%d, found %s", lineNum, s.currentPhase, s.nextTaskIdx, taskID))
	}

	s.nextTaskIdx = tIdx + 1
	s.taskCountInPhase++
	return &taskBlock{id: taskID, line: lineNum}
}

func updateTaskBlockState(task *taskBlock, trimmed string) {
	if strings.HasPrefix(trimmed, "**Target:**") {
		task.hasTarget = true
		task.inActionSteps = false
	} else if strings.HasPrefix(trimmed, "**Context:**") {
		task.hasContext = true
		task.inActionSteps = false
	} else if strings.HasPrefix(trimmed, "**Action Steps:**") {
		task.hasActionHeader = true
		task.inActionSteps = true
	} else if task.inActionSteps && strings.HasPrefix(trimmed, "- ") && !strings.HasPrefix(trimmed, "- [") {
		task.hasActionItems = true
	} else if strings.HasPrefix(trimmed, "**Verification (TDD):**") {
		task.hasVerify = true
		task.inActionSteps = false
	}
}

func (s *taskValidationState) validateLastTask(task *taskBlock) {
	if task == nil {
		return
	}
	if !task.hasTarget {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Target:** field", task.id, task.line))
	}
	if !task.hasContext {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Context:** field", task.id, task.line))
	}
	if !task.hasActionHeader {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Action Steps:** header", task.id, task.line))
	} else if !task.hasActionItems {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory items under **Action Steps:**", task.id, task.line))
	}
	if !task.hasVerify {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Verification (TDD):** section", task.id, task.line))
	}
}

func (s *taskValidationState) checkEmptyPhase() {
	if s.taskCountInPhase == 0 && s.currentPhase > 0 {
		s.errors = append(s.errors, fmt.Sprintf("Phase %d (line %d) has no tasks", s.currentPhase, s.lastPhaseLine))
	}
}

// updateTaskStatusFile updates the status of a task in tasks.md.
func updateTaskStatusFile(projectRoot, slug, taskID, newStatus string) error {
	tasksPath, err := core.SecurePath(projectRoot, filepath.Join(".specforce", "specs", slug, "tasks.md"))
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	// #nosec G304 - Path is secured by SecurePath
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read tasks.md: %w", err)
	}

	contentStr := string(content)
	start, end, err := findTaskBlock(contentStr, taskID)
	if err != nil {
		return err
	}

	taskBlockContent := contentStr[start:end]
	newContent := contentStr

	// 1. Update checklist header if present in the block
	checklistHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^- \[[ xX/]?\] %s: .*$`, regexp.QuoteMeta(taskID)))
	checklistMatch := checklistHeaderRegex.FindStringIndex(taskBlockContent)
	if checklistMatch != nil {
		char := mapStatusToCheckbox(newStatus)
		startInContent := start + checklistMatch[0]
		// The checkbox is at index 3 in "- [ ]"
		newContent = newContent[:startInContent+3] + char + newContent[startInContent+4:]
	}

	// 2. Find and replace the state line in the block (if it exists)
	stateRegex := regexp.MustCompile(`\*\*State:\*\* (\x60?)\[(.*?)\](\x60?)`)
	stateMatch := stateRegex.FindStringSubmatchIndex(taskBlockContent)

	if stateMatch != nil {
		preTick := taskBlockContent[stateMatch[2]:stateMatch[3]]
		postTick := taskBlockContent[stateMatch[6]:stateMatch[7]]
		mappedStatus := mapStatus(newStatus)
		newStateLine := fmt.Sprintf("**State:** %s[%s]%s", preTick, mappedStatus, postTick)

		startInContent := start + stateMatch[0]
		endInContent := start + stateMatch[1]
		newContent = newContent[:startInContent] + newStateLine + newContent[endInContent:]
	} else if checklistMatch == nil {
		return fmt.Errorf("state field or checkbox not found for task %s", taskID)
	}

	// #nosec G306, G703 - Path is secured by SecurePath
	return os.WriteFile(tasksPath, []byte(newContent), 0600)
}

func findTaskBlock(content, taskID string) (int, int, error) {
	// 1. Find the task header (Classic or Checklist)
	classicHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^#{3,4} %s: .*$`, regexp.QuoteMeta(taskID)))
	classicMatches := classicHeaderRegex.FindAllStringIndex(content, -1)

	checklistHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^- \[[ xX/]?\] %s: .*$`, regexp.QuoteMeta(taskID)))
	checklistMatches := checklistHeaderRegex.FindAllStringIndex(content, -1)

	var headerMatches [][]int
	if len(classicMatches) > 0 {
		headerMatches = classicMatches
	} else {
		headerMatches = checklistMatches
	}

	if len(headerMatches) == 0 {
		return 0, 0, fmt.Errorf("task %s not found in tasks.md", taskID)
	}
	if len(headerMatches) > 1 {
		return 0, 0, fmt.Errorf("ambiguous task ID %s: found %d matches", taskID, len(headerMatches))
	}

	headerEnd := headerMatches[0][1]

	// 2. Find the next task header or end of tasks section to define the block
	// We must avoid stopping at a redundant header for the SAME task ID (e.g. checklist under classic)
	nextBoundaryRegex := regexp.MustCompile(`(?m)^(?:#{3,4}|- \[[ xX/]?\]) (?:T[\d.]+|Phase [\d.]+): `)
	allBoundaries := nextBoundaryRegex.FindAllStringIndex(content[headerEnd:], -1)

	blockEnd := len(content)
	for _, b := range allBoundaries {
		boundaryText := content[headerEnd+b[0] : headerEnd+b[1]]
		// If the boundary contains our current taskID, it's likely a hybrid format redundant header, skip it.
		if !strings.Contains(boundaryText, " "+taskID+":") {
			blockEnd = headerEnd + b[0]
			break
		}
	}

	// Also check for end of section ##
	nextSectionRegex := regexp.MustCompile(`(?m)^## `)
	nextSectionMatch := nextSectionRegex.FindStringIndex(content[headerEnd:])
	if nextSectionMatch != nil && (headerEnd+nextSectionMatch[0] < blockEnd) {
		blockEnd = headerEnd + nextSectionMatch[0]
	}

	return headerMatches[0][0], blockEnd, nil
}

func mapStatusToCheckbox(status string) string {
	switch strings.ToLower(status) {
	case "finished":
		return "x"
	case "in-progress":
		return "/"
	case "pending":
		return " "
	default:
		return " "
	}
}

func mapStatus(status string) string {
	statusMap := map[string]string{
		"pending":     "PENDING",
		"in-progress": "IN-PROGRESS",
		"finished":    "FINISHED",
		"failed":      "FAILED",
	}

	mapped, ok := statusMap[strings.ToLower(status)]
	if !ok {
		return strings.ToUpper(status)
	}
	return mapped
}
