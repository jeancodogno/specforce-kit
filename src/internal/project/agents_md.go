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
- **Specs First:** Never write implementation code until a corresponding Specification (requirements.md, design.md, tasks.md) exists and is approved. You MUST consult spec artifacts using the Specforce CLI (e.g., specforce spec list, specforce spec status <slug> --json, and specforce spec artifact <name> --json) rather than reading the raw markdown files directly.
- **Atomic Tasks:** Follow the exact sequence of the tasks.md roadmap. Mark tasks as [DONE] or [FINISHED] sequentially.
- **Verification:** Execute the exact verification/TDD steps defined in the tasks.md file before marking a task as complete.

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
  on_task_finish:
    command: "make lint && make test"
  on_spec_finish:
    command: "make build"
` + "```" + `
If a hook fails, the state transition will be blocked.

*Note: The content above is managed by Specforce. Do not edit inside these markers.*
<!-- SPECFORCE_AGENTS_END -->
`

// EnsureAgentsMD creates or updates the AGENTS.md file in the project root.
func EnsureAgentsMD(root string, ui core.UI) error {
	path := filepath.Join(root, "AGENTS.md")
	replacement := generateAgentsContent()

	var existing string
	if _, err := os.Stat(path); err == nil {
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

	if err := os.WriteFile(path, []byte(merged), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	return nil
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
