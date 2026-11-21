package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	isSmall          bool
	errors           []string
	warnings         []string
	allTasks         []*taskBlock
}

// maxTaskPaths is the number of distinct file paths a task may cite before a warning is raised.
const maxTaskPaths = 3

// TaskValidation holds the outcome of validating tasks.md: blocking errors and non-blocking warnings.
type TaskValidation struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

var taskPathTokenRegex = regexp.MustCompile("`([^`\\s]+)`")

// taskBlock tracks mandatory fields within a task block.
type taskBlock struct {
	id            string
	line          int
	phaseID       int
	hasTarget     bool
	targetPath    string
	hasContext    bool
	hasVerify     bool
	inConstraints bool
	inFence       bool
	hasFence      bool
	citedPaths    map[string]bool
	parallelWith  []string
	isParallel    bool
}

// ValidateTasks performs an exhaustive structural and content validation of tasks.md and returns only the blocking errors.
func ValidateTasks(ctx context.Context, projectRoot, slug string) ([]string, error) {
	report, err := ValidateTasksReport(ctx, projectRoot, slug)
	if err != nil {
		return nil, err
	}
	return report.Errors, nil
}

// ValidateTasksReport validates tasks.md and returns blocking errors and non-blocking warnings.
func ValidateTasksReport(ctx context.Context, projectRoot, slug string) (TaskValidation, error) {
	slug = ResolveSlug(projectRoot, slug)

	tasksPath, err := core.SecurePath(projectRoot, filepath.Join(".specforce", "specs", slug, "tasks.md"))
	if err != nil {
		return TaskValidation{}, fmt.Errorf("security: %w", err)
	}

	// #nosec G304 - Path is secured by core.SecurePath
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return TaskValidation{}, fmt.Errorf("failed to read tasks.md: %w", err)
	}

	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return TaskValidation{}, fmt.Errorf("failed to load metadata: %w", err)
	}

	isSmall := meta != nil && meta.Size == SpecSizeSmall

	lines := strings.Split(string(content), "\n")
	state := &taskValidationState{nextTaskIdx: 1, isSmall: isSmall}
	var currentTask *taskBlock

	phaseRegex := regexp.MustCompile(`^### Phase (\d+): (.*)$`)
	taskRegex := regexp.MustCompile(`^(#{4}|- \[[ xX/]?\]) (T(\d+)\.(\d+)): (.*)$`)

	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return TaskValidation{}, err
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

	if state.currentPhase == 0 {
		state.errors = append(state.errors, "No valid Phase (### Phase N: Name) found in tasks.md")
	}

	state.validateParallelReferences()
	state.validateParallelTargets()

	return TaskValidation{Errors: state.errors, Warnings: state.warnings}, nil
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

	newTask := &taskBlock{id: taskID, line: lineNum, phaseID: s.currentPhase}
	s.allTasks = append(s.allTasks, newTask)
	return newTask
}

func updateTaskBlockState(task *taskBlock, trimmed string) {
	if strings.HasPrefix(trimmed, "```") {
		task.inFence = !task.inFence
		task.hasFence = true
		return
	}
	if task.inFence {
		return
	}

	if strings.HasPrefix(trimmed, "**") {
		// Any bold field heading closes the constraints block (matches extractConstraints).
		task.inConstraints = false
	}

	switch {
	case strings.HasPrefix(trimmed, "**Target:**"):
		task.hasTarget = true
		task.inConstraints = false
		raw := strings.TrimSpace(strings.TrimPrefix(trimmed, "**Target:**"))
		task.targetPath = strings.Trim(raw, "`")
		task.collectPaths(raw)
	case strings.HasPrefix(trimmed, "**Context:**"):
		task.hasContext = true
		task.inConstraints = false
	case strings.HasPrefix(trimmed, "**Constraints:**"), strings.HasPrefix(trimmed, "**Action Steps:**"):
		task.inConstraints = true
	case task.inConstraints && strings.HasPrefix(trimmed, "- ") && !strings.HasPrefix(trimmed, "- ["):
		task.collectPaths(trimmed)
	case strings.HasPrefix(trimmed, "**Acceptance Check:**"), strings.HasPrefix(trimmed, "**Verification (TDD):**"):
		task.hasVerify = true
		task.inConstraints = false
	case strings.HasPrefix(trimmed, "**Parallel With:**"):
		task.isParallel = true
		peersStr := strings.TrimSpace(strings.TrimPrefix(trimmed, "**Parallel With:**"))
		for _, p := range strings.Split(peersStr, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				task.parallelWith = append(task.parallelWith, p)
			}
		}
	}
}

// collectPaths records the backtick-quoted tokens without spaces that contain a slash.
func (t *taskBlock) collectPaths(text string) {
	for _, m := range taskPathTokenRegex.FindAllStringSubmatch(text, -1) {
		if strings.Contains(m[1], "/") {
			if t.citedPaths == nil {
				t.citedPaths = make(map[string]bool)
			}
			t.citedPaths[m[1]] = true
		}
	}
}

func (s *taskValidationState) validateLastTask(task *taskBlock) {
	if task == nil {
		return
	}
	if !task.hasTarget {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Target:** field", task.id, task.line))
	}
	if !s.isSmall && !task.hasContext {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Context:** field", task.id, task.line))
	}
	if !task.hasVerify {
		s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) is missing mandatory **Acceptance Check:** section", task.id, task.line))
	}
	if task.hasFence {
		s.warnings = append(s.warnings, fmt.Sprintf("Task %s (line %d) contains a fenced code block", task.id, task.line))
	}
	if n := len(task.citedPaths); n > maxTaskPaths {
		s.warnings = append(s.warnings, fmt.Sprintf("Task %s (line %d) cites %d distinct file paths (limit %d)", task.id, task.line, n, maxTaskPaths))
	}
}

func (s *taskValidationState) checkEmptyPhase() {
	if s.taskCountInPhase == 0 && s.currentPhase > 0 {
		s.errors = append(s.errors, fmt.Sprintf("Phase %d (line %d) has no tasks", s.currentPhase, s.lastPhaseLine))
	}
}

func (s *taskValidationState) validateParallelReferences() {
	taskMap := s.buildTaskMap()
	adjacencyMap := s.buildAdjacencyMap()

	for _, t := range s.allTasks {
		peers := adjacencyMap[t.id]
		for _, peerID := range peers {
			peer, exists := taskMap[peerID]
			if !exists {
				s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) references unknown parallel task %s", t.id, t.line, peerID))
			} else if peer.phaseID != t.phaseID {
				s.errors = append(s.errors, fmt.Sprintf("Task %s (line %d) references parallel task %s from a different phase", t.id, t.line, peerID))
			}
		}
	}
}

func (s *taskValidationState) buildTaskMap() map[string]*taskBlock {
	taskMap := make(map[string]*taskBlock)
	for _, t := range s.allTasks {
		taskMap[t.id] = t
	}
	return taskMap
}

func (s *taskValidationState) buildAdjacencyMap() map[string][]string {
	adjacencyMap := make(map[string][]string)

	for _, t := range s.allTasks {
		if !t.isParallel {
			continue
		}
		for _, peerID := range t.parallelWith {
			addUniquePeer(adjacencyMap, t.id, peerID)
			addUniquePeer(adjacencyMap, peerID, t.id)
		}
	}
	return adjacencyMap
}

func addUniquePeer(adjacencyMap map[string][]string, key, peer string) {
	for _, p := range adjacencyMap[key] {
		if p == peer {
			return
		}
	}
	adjacencyMap[key] = append(adjacencyMap[key], peer)
}

func (s *taskValidationState) validateParallelTargets() {
	taskMap := s.buildTaskMap()
	adjacencyMap := s.buildAdjacencyMap()
	visited := make(map[string]bool)

	for _, t := range s.allTasks {
		if visited[t.id] || len(adjacencyMap[t.id]) == 0 {
			continue
		}
		group := getParallelGroup(t.id, adjacencyMap, visited)
		s.checkGroupTargetConflicts(group, taskMap)
	}
}

func getParallelGroup(startNode string, adjacencyMap map[string][]string, visited map[string]bool) []string {
	var group []string
	queue := []string{startNode}
	visited[startNode] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		group = append(group, curr)

		for _, peer := range adjacencyMap[curr] {
			if !visited[peer] {
				visited[peer] = true
				queue = append(queue, peer)
			}
		}
	}
	return group
}

func (s *taskValidationState) checkGroupTargetConflicts(group []string, taskMap map[string]*taskBlock) {
	for i := 0; i < len(group); i++ {
		for j := i + 1; j < len(group); j++ {
			t1, ok1 := taskMap[group[i]]
			t2, ok2 := taskMap[group[j]]
			if !ok1 || !ok2 || t1.targetPath == "" || t2.targetPath == "" {
				continue
			}
			if t1.targetPath == t2.targetPath {
				id1, id2 := t1.id, t2.id
				if id1 > id2 {
					id1, id2 = id2, id1
				}
				s.errors = append(s.errors, fmt.Sprintf("Parallel conflict: %s and %s both target %q", id1, id2, t1.targetPath))
			}
		}
	}
}

// updateTaskStatusesFile updates the status of multiple tasks in tasks.md atomically.
func updateTaskStatusesFile(projectRoot, slug string, taskIDs []string, newStatus string) error {
	if len(taskIDs) == 0 {
		return nil
	}

	tasksPath, err := core.SecurePath(projectRoot, filepath.Join(".specforce", "specs", slug, "tasks.md"))
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	// #nosec G304 - Path is secured by SecurePath
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read tasks.md: %w", err)
	}

	uniqueTaskIDs := deduplicateStrings(taskIDs)
	newContent, err := applyTaskReplacements(string(content), uniqueTaskIDs, newStatus)
	if err != nil {
		return err
	}

	// #nosec G306, G703 - Path is secured by SecurePath
	if err := os.WriteFile(tasksPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write tasks.md: %w", err)
	}

	updateTaskMetadataLogs(projectRoot, slug, newStatus, uniqueTaskIDs)
	return nil
}

type taskReplacement struct {
	start      int
	end        int
	newContent string
}

func applyTaskReplacements(contentStr string, uniqueTaskIDs []string, newStatus string) (string, error) {
	replacements := make([]taskReplacement, 0, len(uniqueTaskIDs))

	for _, taskID := range uniqueTaskIDs {
		start, end, err := findTaskBlock(contentStr, taskID)
		if err != nil {
			return "", err
		}

		updatedBlock, err := updateTaskBlockContent(contentStr[start:end], taskID, newStatus)
		if err != nil {
			return "", err
		}

		replacements = append(replacements, taskReplacement{
			start:      start,
			end:        end,
			newContent: updatedBlock,
		})
	}

	sort.Slice(replacements, func(i, j int) bool {
		return replacements[i].start > replacements[j].start
	})

	newContent := contentStr
	for _, rep := range replacements {
		newContent = newContent[:rep.start] + rep.newContent + newContent[rep.end:]
	}

	return newContent, nil
}

func deduplicateStrings(items []string) []string {
	seen := make(map[string]bool, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func updateTaskMetadataLogs(projectRoot, slug, newStatus string, uniqueTaskIDs []string) {
	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return
	}

	switch strings.ToLower(newStatus) {
	case "in-progress":
		for _, taskID := range uniqueTaskIDs {
			meta.StartSession(taskID)
		}
	case "finished", "todo", "pending":
		for _, taskID := range uniqueTaskIDs {
			meta.EndSession(taskID)
		}
	}
	_ = SaveMetadata(projectRoot, slug, meta)
}

func updateTaskBlockContent(taskBlockContent, taskID, newStatus string) (string, error) {
	// 1. Check for checklist header in the block
	checklistHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^- \[[ xX/]?\] %s: .*$`, regexp.QuoteMeta(taskID)))
	checklistMatch := checklistHeaderRegex.FindStringIndex(taskBlockContent)

	// 2. Check for state line in the block
	stateRegex := regexp.MustCompile(`\*\*State:\*\* (\x60?)\[(.*?)\](\x60?)`)
	stateMatch := stateRegex.FindStringSubmatchIndex(taskBlockContent)

	if checklistMatch == nil && stateMatch == nil {
		return "", fmt.Errorf("state field or checkbox not found for task %s", taskID)
	}

	res := taskBlockContent
	if checklistMatch != nil {
		char := mapStatusToCheckbox(newStatus)
		chkIdx := checklistMatch[0] + 3
		res = res[:chkIdx] + char + res[chkIdx+1:]
	}

	stateMatch = stateRegex.FindStringSubmatchIndex(res)
	if stateMatch != nil {
		preTick := res[stateMatch[2]:stateMatch[3]]
		postTick := res[stateMatch[6]:stateMatch[7]]
		mappedStatus := mapStatus(newStatus)
		newStateLine := fmt.Sprintf("**State:** %s[%s]%s", preTick, mappedStatus, postTick)

		res = res[:stateMatch[0]] + newStateLine + res[stateMatch[1]:]
	}

	return res, nil
}

// updateTaskStatusFile updates the status of a task in tasks.md.
func updateTaskStatusFile(projectRoot, slug, taskID, newStatus string) error {
	return updateTaskStatusesFile(projectRoot, slug, []string{taskID}, newStatus)
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
