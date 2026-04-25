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

	// 3. Find and replace the state line in the block
	stateRegex := regexp.MustCompile(`\*\*State:\*\* (\x60?)\[(.*?)\](\x60?)`)
	stateMatch := stateRegex.FindStringSubmatchIndex(taskBlock)
	if stateMatch == nil {
		return fmt.Errorf("state field not found for task %s", taskID)
	}

	// Capture surrounding formatting
	preTick := taskBlock[stateMatch[2]:stateMatch[3]]
	postTick := taskBlock[stateMatch[6]:stateMatch[7]]

	mappedStatus := mapStatus(newStatus)
	newStateLine := fmt.Sprintf("**State:** %s[%s]%s", preTick, mappedStatus, postTick)

	// We need to replace it in the original contentStr
	startInContent := start + stateMatch[0]
	endInContent := start + stateMatch[1]

	newContent := contentStr[:startInContent] + newStateLine + contentStr[endInContent:]

	// #nosec G306, G703 - Path is secured by SecurePath
	return os.WriteFile(tasksPath, []byte(newContent), 0600)
}

func findTaskBlock(content, taskID string) (int, int, error) {
	// 1. Find the task header
	taskHeaderRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^#{3,4} %s: .*$`, regexp.QuoteMeta(taskID)))
	headerMatches := taskHeaderRegex.FindAllStringIndex(content, -1)

	if len(headerMatches) == 0 {
		return 0, 0, fmt.Errorf("task %s not found in tasks.md", taskID)
	}
	if len(headerMatches) > 1 {
		return 0, 0, fmt.Errorf("ambiguous task ID %s: found %d matches", taskID, len(headerMatches))
	}

	headerEnd := headerMatches[0][1]

	// 2. Find the next task header or end of tasks section to define the block
	nextBoundaryRegex := regexp.MustCompile(`(?m)^#{3,4} (T[\d.]+|Phase [\d.]+): `)
	nextBoundaryMatch := nextBoundaryRegex.FindStringIndex(content[headerEnd:])

	blockEnd := len(content)
	if nextBoundaryMatch != nil {
		blockEnd = headerEnd + nextBoundaryMatch[0]
	}

	// Also check for end of section ##
	nextSectionRegex := regexp.MustCompile(`(?m)^## `)
	nextSectionMatch := nextSectionRegex.FindStringIndex(content[headerEnd:])
	if nextSectionMatch != nil && (nextBoundaryMatch == nil || headerEnd+nextSectionMatch[0] < blockEnd) {
		blockEnd = headerEnd + nextSectionMatch[0]
	}

	return headerMatches[0][0], blockEnd, nil
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
