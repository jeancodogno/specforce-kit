package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

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

	taskBlock := contentStr[start:end]
	newContent := contentStr

	// 1. Update checklist header if present in the block
	checklistHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^- \[[ xX/]?\] %s: .*$`, regexp.QuoteMeta(taskID)))
	checklistMatch := checklistHeaderRegex.FindStringIndex(taskBlock)
	if checklistMatch != nil {
		char := mapStatusToCheckbox(newStatus)
		startInContent := start + checklistMatch[0]
		// The checkbox is at index 3 in "- [ ]"
		newContent = newContent[:startInContent+3] + char + newContent[startInContent+4:]
	}

	// 2. Find and replace the state line in the block (if it exists)
	stateRegex := regexp.MustCompile(`\*\*State:\*\* (\x60?)\[(.*?)\](\x60?)`)
	stateMatch := stateRegex.FindStringSubmatchIndex(taskBlock)

	if stateMatch != nil {
		preTick := taskBlock[stateMatch[2]:stateMatch[3]]
		postTick := taskBlock[stateMatch[6]:stateMatch[7]]
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
