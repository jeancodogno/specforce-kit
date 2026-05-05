package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

const agentsMDTemplate = `<!-- SPECFORCE_AGENTS_START -->
# AI Agent Collaboration Guide

This project uses **Specforce** for Spec-Driven Development (SDD). As an AI agent, you MUST adhere to the following rules:

## 1. Spec-Driven Development (SDD) Protocol
You MUST operate exclusively through the Specforce workflow engines (commands/skills). They define your mindset and mandatory steps:

- **Discovery (` + "`/spf:discovery`" + `):** Activate for brainstorming, research, bug investigation, or root cause analysis. Purely read-only.
- **Planning (` + "`/spf:spec`" + `):** Activate for new features, structural changes, or to formalize a discovered fix strategy.
- **Governance (` + "`/spf:constitution`" + `):** Use to ensure proposals respect architecture, security, and principles.
- **Execution (` + "`/spf:implement`" + `):** Activate to perform the deterministic implementation cycle following an approved roadmap.
- **Archival (` + "`/spf:archive`" + `):** Activate once verified to harvest lessons, update Memorial, and clean up specs.

### Proactive Mandate
Do NOT wait for explicit slash commands. You MUST automatically activate the correct workflow based on the user's technical intent:
1. **Discovery Intent** (vague idea, "how to", bug report) --> Activate ` + "`spf.discovery`" + `.
2. **Planning Intent** (new feature, structural pivot, confirmed fix) --> Activate ` + "`spf.spec`" + `.
3. **Execution Intent** (approved roadmap exists, "go", "implement") --> Activate ` + "`spf.implement`" + `.

**CRITICAL: DIRECT EDIT PROHIBITION**
You are STRICTLY FORBIDDEN from using mutation tools (` + "`replace`" + `, ` + "`write_file`" + `) to modify the codebase unless a valid, approved specification exists. If the intent requires a change but no spec is active, you MUST pivot to ` + "`spf.discovery`" + ` or ` + "`spf.spec`" + ` first.

- **Specs First:** Never write implementation code until a fully approved Specification (requirements.md, design.md, tasks.md) exists.
- **Total Consistency:** If a change is required at any point (even mid-implementation), you MUST update ALL related artifacts. You are strictly forbidden from updating only tasks.md while leaving requirements.md or design.md inconsistent.
- **Atomic Execution:** Follow the exact sequence of the tasks.md roadmap. Mark tasks as [DONE] or [FINISHED] sequentially and ONLY after successful verification.

## 2. Project Constitution
Before proposing architectural changes or adding new patterns, you MUST review the relevant Constitution documents located in .specforce/docs/:
- principles.md: Core values, philosophy, and cultural/technical axioms.
- architecture.md: System boundaries, dependency direction, and persistence topology.
- ui-ux.md: Visual direction, interaction patterns, and aesthetic DNA.
- security.md: AuthZ, roles, permissions, and data protection rules.
- engineering.md: Coding standards, testing strategy, and refactoring guidelines.
- governance.md: Project lifecycle rules, ownership, and AI boundaries.
- memorial.md: Durable lessons learned and cross-session memory.

## 3. Custom Hooks Configuration
Specforce allows developers to gate state transitions (e.g., finishing a task) using custom hooks. You can configure these in the project root's config.yaml:

` + "```" + `yaml
# config.yaml example
hooks:
  on_task_finished:
    - "make lint"
    - "make test"
  on_phase_finished:
    - "go test ./src/internal/..."
  on_all_tasks_finished:
    - "go test ./..."
` + "```" + `
If a hook fails, the state transition will be blocked.

*Note: The content above is managed by Specforce. Do not edit inside these markers.*
<!-- SPECFORCE_AGENTS_END -->
`

// EnsureAgentsMD creates or updates the AGENTS.md file in the project root.
func EnsureAgentsMD(root string, ui core.UI, selectedAgents []string) error {
	root = filepath.Clean(root)
	path := filepath.Join(root, "AGENTS.md")
	replacement := generateAgentsContent()

	var existing string
	if _, err := os.Stat(path); err == nil {
		// #nosec G304
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read existing AGENTS.md: %w", err)
		}
		existing = string(data)
	}

	merged := mergeAgentsContent(existing, replacement)

	if existing == merged {
		return nil
	}

	if ui != nil {
		ui.SubTask("Updating AGENTS.md...")
	}

	// #nosec G306 G304
	if err := os.WriteFile(path, []byte(merged), 0600); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	return ensurePlatformConfigs(root, selectedAgents)
}

func ensurePlatformConfigs(root string, selectedAgents []string) error {
	// 1. Gemini
	geminiDir := filepath.Join(root, ".gemini")
	if shouldManageDir(root, ".gemini", []string{"gemini-cli"}, selectedAgents) {
		if err := os.MkdirAll(geminiDir, 0750); err != nil {
			return fmt.Errorf("failed to create .gemini directory: %w", err)
		}
		geminiSettings := filepath.Join(geminiDir, "settings.json")
		settingsContent := `{
  "context": {
    "fileName": [
      "AGENTS.md",
      "GEMINI.md"
    ]
  }
}`
		// Always write the file to ensure the configuration is correct and up to date
		if err := os.WriteFile(geminiSettings, []byte(settingsContent), 0600); err != nil {
			return fmt.Errorf("failed to write .gemini/settings.json: %w", err)
		}
	}

	// 2. Symlinks
	agentMappings := map[string][]string{
		".agent":  {"antigravity"},
		".claude": {"claude"},
	}

	for dir, requiredAgents := range agentMappings {
		if shouldManageDir(root, dir, requiredAgents, selectedAgents) {
			rulesDir := filepath.Join(root, dir, "rules")
			if err := os.MkdirAll(rulesDir, 0750); err != nil {
				return fmt.Errorf("failed to create %s/rules directory: %w", dir, err)
			}

			linkPath := filepath.Join(rulesDir, "AGENTS.md")
			if _, err := os.Lstat(linkPath); err == nil {
				// Remove existing link or file if it exists to ensure it's correct
				if err := os.Remove(linkPath); err != nil {
					return fmt.Errorf("failed to remove existing link at %s: %w", linkPath, err)
				}
			}

			if err := os.Symlink("../../AGENTS.md", linkPath); err != nil {
				return fmt.Errorf("failed to create symlink at %s: %w", linkPath, err)
			}
		}
	}

	return nil
}

func shouldManageDir(root, dir string, requiredAgents, selectedAgents []string) bool {
	// REQ-2 AC-3: If directory exists, we manage it regardless of selection
	path := filepath.Join(root, dir)
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}

	// REQ-1: Manage if tool is selected
	for _, ra := range requiredAgents {
		for _, sa := range selectedAgents {
			if ra == sa {
				return true
			}
		}
	}

	return false
}

func generateAgentsContent() string {
	return agentsMDTemplate
}

func mergeAgentsContent(existing, replacement string) string {
	startMarker := "<!-- SPECFORCE_AGENTS_START -->"
	endMarker := "<!-- SPECFORCE_AGENTS_END -->"

	startIdx := strings.Index(existing, startMarker)
	endIdx := strings.Index(existing, endMarker)

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace content between markers (including markers)
		return existing[:startIdx] + replacement + existing[endIdx+len(endMarker):]
	}

	if existing == "" {
		return replacement
	}

	// Append to end if markers not found
	if !strings.HasSuffix(existing, "\n") {
		existing += "\n"
	}
	return existing + replacement
}
