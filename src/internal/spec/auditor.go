package spec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CoherenceError represents a logical gap or contradiction between spec artifacts.
type CoherenceError struct {
	Artifact string `json:"artifact"` // Target artifact to fix (requirements|design|tasks)
	Code     string `json:"code"`     // Error code (e.g., MISSING_TASK, TECH_DRIFT)
	Message  string `json:"message"`  // Human-readable description
	Context  string `json:"context"`  // Relevant snippet from the source of truth
}

// Auditor defines the contract for validating specification coherence.
type Auditor interface {
	// Audit performs a deep logical analysis, typically involving AI.
	Audit(ctx context.Context, slug string) ([]CoherenceError, error)

	// DeterministicCheck performs fast, rule-based validations (e.g., US-X tag presence).
	DeterministicCheck(ctx context.Context, slug string) []CoherenceError
}

// SpecAuditor is the default implementation of the Auditor interface.
type SpecAuditor struct {
	projectRoot string
}

// NewAuditor creates a new SpecAuditor.
func NewAuditor(projectRoot string) *SpecAuditor {
	return &SpecAuditor{projectRoot: projectRoot}
}

// DeterministicCheck performs fast, rule-based validations.
func (a *SpecAuditor) DeterministicCheck(ctx context.Context, slug string) []CoherenceError {
	specDir := filepath.Join(a.projectRoot, ".specforce", "specs", slug)
	reqPath := filepath.Join(specDir, "requirements.md")
	tasksPath := filepath.Join(specDir, "tasks.md")

	meta, _ := LoadMetadata(a.projectRoot, slug)

	// 1. Extract US-X from requirements
	// #nosec G304 - internal spec file
	reqData, err := os.ReadFile(reqPath)
	if err != nil {
		if meta != nil && meta.Size == SpecSizeSmall {
			return nil
		}
		return []CoherenceError{{
			Artifact: "requirements",
			Code:     "FILE_MISSING",
			Message:  "requirements.md is missing",
		}}
	}

	requirements := extractRequirements(reqData)
	return crossReferenceTasks(tasksPath, requirements)
}

func extractRequirements(reqData []byte) map[string]bool {
	usRegex := regexp.MustCompile(`\[(US-\d+)\]`)
	matches := usRegex.FindAllStringSubmatch(string(reqData), -1)
	requirements := make(map[string]bool)
	for _, m := range matches {
		requirements[m[1]] = true
	}

	outOfScope := false
	lines := strings.Split(string(reqData), "\n")
	for _, line := range lines {
		if strings.Contains(line, "## 2. Out of Scope") {
			outOfScope = true
		} else if strings.HasPrefix(line, "## ") && outOfScope {
			outOfScope = false
		}
		if outOfScope {
			m := usRegex.FindStringSubmatch(line)
			if len(m) > 1 {
				delete(requirements, m[1])
			}
		}
	}
	return requirements
}

func crossReferenceTasks(tasksPath string, requirements map[string]bool) []CoherenceError {
	// #nosec G304 - internal spec file
	tasksData, err := os.ReadFile(tasksPath)
	if err != nil {
		return []CoherenceError{{
			Artifact: "tasks",
			Code:     "FILE_MISSING",
			Message:  "tasks.md is missing",
		}}
	}

	var errors []CoherenceError
	tasksContent := string(tasksData)
	for us := range requirements {
		if !strings.Contains(tasksContent, fmt.Sprintf("**Context:** [%s]", us)) {
			errors = append(errors, CoherenceError{
				Artifact: "tasks",
				Code:     "MISSING_TASK",
				Message:  fmt.Sprintf("Requirement [%s] has no execution steps in tasks.md", us),
				Context:  fmt.Sprintf("## [%s]", us),
			})
		}
	}
	return errors
}

// Audit runs the full suite of checks (Deterministic + AI-driven parsing).
func (a *SpecAuditor) Audit(ctx context.Context, slug string) ([]CoherenceError, error) {
	// 1. Run deterministic checks first
	errors := a.DeterministicCheck(ctx, slug)

	// 2. The AI-driven part is handled by the orchestrator calling the agent.
	// This function serves as the repository for parsing that agent's output
	// if it's already been persisted to spec.yaml or passed as raw text.
	return errors, nil
}

// ParseErrors extracts [COHERENCE_ERROR] blocks from raw AI output.
func (a *SpecAuditor) ParseErrors(content string) []CoherenceError {
	var errors []CoherenceError
	
	// Regex to match [COHERENCE_ERROR] ... [/COHERENCE_ERROR]
	// Using (?s) for dot-all mode
	blockRegex := regexp.MustCompile(`(?s)\[COHERENCE_ERROR\](.*?)\[/COHERENCE_ERROR\]`)
	blocks := blockRegex.FindAllStringSubmatch(content, -1)

	for _, b := range blocks {
		body := b[1]
		err := CoherenceError{}
		
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Artifact:") {
				err.Artifact = strings.TrimSpace(strings.TrimPrefix(line, "Artifact:"))
			} else if strings.HasPrefix(line, "Code:") {
				err.Code = strings.TrimSpace(strings.TrimPrefix(line, "Code:"))
			} else if strings.HasPrefix(line, "Message:") {
				err.Message = strings.TrimSpace(strings.TrimPrefix(line, "Message:"))
			} else if strings.HasPrefix(line, "Context:") {
				err.Context = strings.TrimSpace(strings.TrimPrefix(line, "Context:"))
			}
		}
		
		if err.Artifact != "" && err.Message != "" {
			errors = append(errors, err)
		}
	}

	return errors
}
