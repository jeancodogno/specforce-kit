package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// LegacyAgents lists all deprecated agent blueprints.
var LegacyAgents = []string{
	"specforce-architect",
	"specforce-developer",
	"specforce-planner",
	"specforce-product-analyst",
	"specforce-qa",
	"specforce-spec-reviewer",
}

// LegacySkills lists all deprecated built-in skill blueprints.
var LegacySkills = []string{
	"opportunity-framing",
	"pragmatic-product-owner",
	"task-atomic-decomposition",
}

// KnownToolDirs lists all standard tool integration directories.
var KnownToolDirs = []string{
	".agents",
	".cursor",
	".claude",
	".qwen",
	".opencode",
	".kilocode",
	".codex",
	".kimi",
}

// DetectLegacyAssets scans tool integration directories for deprecated agents, skills, workflows, commands, and gemini.
func DetectLegacyAssets(root string) ([]string, error) {
	var found []string

	// Check root-level .gemini directory
	geminiDir := filepath.Join(root, ".gemini")
	if _, err := os.Stat(geminiDir); err == nil {
		found = append(found, geminiDir)
	}

	for _, toolDir := range KnownToolDirs {
		baseDir := filepath.Join(root, toolDir)
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			continue
		}

		// 1. Check legacy skills (e.g. .agents/skills/pragmatic-product-owner)
		skillsDir := filepath.Join(baseDir, "skills")
		if info, err := os.Stat(skillsDir); err == nil && info.IsDir() {
			for _, legacySkill := range LegacySkills {
				matches, err := findMatches(skillsDir, legacySkill)
				if err != nil {
					return nil, err
				}
				found = append(found, matches...)
			}
		}

		// 2. Check legacy agents (e.g. .agents/agents/specforce-architect)
		agentsDir := filepath.Join(baseDir, "agents")
		if info, err := os.Stat(agentsDir); err == nil && info.IsDir() {
			for _, legacyAgent := range LegacyAgents {
				matches, err := findMatches(agentsDir, legacyAgent)
				if err != nil {
					return nil, err
				}
				found = append(found, matches...)
			}
		}

		// 3. Check legacy workflows directory (e.g. .agents/workflows)
		workflowsDir := filepath.Join(baseDir, "workflows")
		if _, err := os.Stat(workflowsDir); err == nil {
			found = append(found, workflowsDir)
		}

		// 4. Check legacy commands directory (e.g. .claude/commands, .cursor/commands)
		commandsDir := filepath.Join(baseDir, "commands")
		if _, err := os.Stat(commandsDir); err == nil {
			found = append(found, commandsDir)
		}
	}

	return deduplicatePaths(found), nil
}

func findMatches(parentDir, namePattern string) ([]string, error) {
	var matches []string

	entries, err := os.ReadDir(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		// Exact match on folder/file or prefix match (e.g. namePattern, namePattern.md, namePattern.json)
		baseWithoutExt := strings.TrimSuffix(entryName, filepath.Ext(entryName))
		if entryName == namePattern || baseWithoutExt == namePattern {
			matches = append(matches, filepath.Join(parentDir, entryName))
		}
	}

	return matches, nil
}

func deduplicatePaths(paths []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, p := range paths {
		cleaned := filepath.Clean(p)
		if _, exists := seen[cleaned]; !exists {
			seen[cleaned] = struct{}{}
			result = append(result, cleaned)
		}
	}
	return result
}

// CleanupLegacyAssets removes the specified legacy asset paths safely.
func CleanupLegacyAssets(root string, assetPaths []string, ui core.UI) error {
	if len(assetPaths) == 0 {
		return nil
	}

	for _, path := range assetPaths {
		secPath, err := core.SecurePath(root, path)
		if err != nil {
			// In case path is already absolute and within root
			if filepath.IsAbs(path) {
				secPath = path
			} else {
				secPath = filepath.Join(root, path)
			}
		}

		if _, err := os.Stat(secPath); err == nil {
			if err := os.RemoveAll(secPath); err != nil {
				if ui != nil {
					ui.Warn(fmt.Sprintf("Failed to remove legacy asset %s: %v", path, err))
				}
			} else if ui != nil {
				rel, err := filepath.Rel(root, secPath)
				if err != nil {
					rel = secPath
				}
				dotCount := 40 - len(rel)
				if dotCount < 1 {
					dotCount = 1
				}
				dots := strings.Repeat(".", dotCount)
				ui.LogSubTask(fmt.Sprintf("%s %s DELETED", rel, dots))
			}
		}
	}

	return nil
}

// PromptAndCleanupLegacyAssets checks for legacy assets and prompts the user before cleaning them up.
func PromptAndCleanupLegacyAssets(root string, ui core.UI) error {
	legacyAssets, err := DetectLegacyAssets(root)
	if err != nil || len(legacyAssets) == 0 {
		return err
	}

	shouldDelete := false
	if ui != nil {
		shouldDelete = ui.Confirm("Legacy Specforce agents/skills detected. Do you want to remove them?")
	}

	if shouldDelete {
		return CleanupLegacyAssets(root, legacyAssets, ui)
	}

	return nil
}
