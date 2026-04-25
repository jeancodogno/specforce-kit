package tui

import (
	"fmt"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/spec"
)

// RenderImplementationStatus renders the implementation status report.
func RenderImplementationStatus(report *spec.ImplementationReport) error {
	if IsTTY() {
		PrintBranding()
	}

	renderStatusHeader(report)
	renderArtifactsCheck(report)
	renderImplementationContext(report)
	renderExecutionStrategy(report)
	renderAtomicTasks(report)
	renderProjectInstructions(report)
	renderPreemptiveMitigations(report)

	return nil
}

func renderStatusHeader(report *spec.ImplementationReport) {
	fmt.Println("\n" + HeaderStyle.Render("IMPLEMENTATION READINESS: "+strings.ToUpper(report.Name)))
	
	var statusBadge string
	switch strings.ToLower(report.Status) {
	case "ready":
		statusBadge = SuccessBadgeStyle.Render(" READY ")
	case "blocked":
		statusBadge = ErrorBadgeStyle.Render(" BLOCKED ")
	case "in-progress":
		statusBadge = WarningBadgeStyle.Render(" IN-PROGRESS ")
	case "finished":
		statusBadge = SuccessBadgeStyle.Render(" FINISHED ")
	default:
		statusBadge = SubtitleStyle.Render(strings.ToUpper(report.Status))
	}
	fmt.Printf("%s\n\n", statusBadge)
}

func renderArtifactsCheck(report *spec.ImplementationReport) {
	fmt.Println(HeaderStyle.Render("ARTIFACTS CHECK"))
	if len(report.MissingArtifacts) == 0 {
		fmt.Println(SuccessStyle.Render("  [v] All triad artifacts present (requirements, design, tasks)"))
	} else {
		for _, art := range report.MissingArtifacts {
			fmt.Println(ErrorStyle.Render("  [x] Missing: " + art))
		}
	}
	fmt.Println()
}

func renderImplementationContext(report *spec.ImplementationReport) {
	fmt.Println(HeaderStyle.Render("IMPLEMENTATION CONTEXT"))
	if len(report.ContextFiles) == 0 {
		fmt.Println(SubtitleStyle.Render("  No context files found."))
	} else {
		for _, file := range report.ContextFiles {
			fmt.Println(SubtitleStyle.Render("  - " + file))
		}
	}
	fmt.Println()
}

func renderExecutionStrategy(report *spec.ImplementationReport) {
	if report.ExecutionStrategy != "" {
		fmt.Println(HeaderStyle.Render("EXECUTION STRATEGY"))
		fmt.Println(BodyStyle.Render("  " + report.ExecutionStrategy))
		fmt.Println()
	}
}

func renderAtomicTasks(report *spec.ImplementationReport) {
	fmt.Println(HeaderStyle.Render("ATOMIC TASKS"))
	if len(report.Phases) == 0 || (len(report.Phases) == 1 && len(report.Phases[0].Tasks) == 0) {
		fmt.Println(SubtitleStyle.Render("  No tasks found in tasks.md."))
	} else {
		for _, phase := range report.Phases {
			renderPhase(phase, len(report.Phases))
		}
	}
	fmt.Println()
}

func renderPhase(phase spec.Phase, totalPhases int) {
	allFinished := len(phase.Tasks) > 0
	for _, t := range phase.Tasks {
		if strings.ToUpper(t.State) != "FINISHED" {
			allFinished = false
			break
		}
	}

	phaseTitle := phase.Title
	if allFinished {
		phaseTitle = "✓ " + phaseTitle
	}
	
	if totalPhases > 1 || (phase.ID != "0" && phase.ID != "") {
		phaseStyle := HeaderStyle
		if allFinished {
			phaseStyle = SuccessStyle
		}
		fmt.Printf("\n  %s\n", phaseStyle.Render(strings.ToUpper(phaseTitle)))
	}

	for _, task := range phase.Tasks {
		renderTask(task)
	}
}

func renderTask(task spec.ImplementationTask) {
	stateIcon := "○"
	stateStyle := BodyStyle
	switch strings.ToUpper(task.State) {
	case "FINISHED":
		stateIcon = "◉"
		stateStyle = SuccessStyle
	case "IN-PROGRESS":
		stateStyle = WarningStyle
	case "FAILED":
		stateStyle = ErrorStyle
	}
	fmt.Printf("    %s %s: %s\n", stateStyle.Render(stateIcon), HeaderStyle.Render(task.ID), BodyStyle.Render(task.Title))
	if task.Target != "" {
		fmt.Printf("        %s %s\n", SubtitleStyle.Render("Target:"), BodyStyle.Render(task.Target))
	}
}

func renderProjectInstructions(report *spec.ImplementationReport) {
	if len(report.Instructions) > 0 {
		fmt.Println(HeaderStyle.Render("PROJECT SPECIFIC INSTRUCTIONS"))
		for _, inst := range report.Instructions {
			fmt.Println(BodyStyle.Render("  - " + inst))
		}
		fmt.Println()
	}
}

func renderPreemptiveMitigations(report *spec.ImplementationReport) {
	if report.PreemptiveMitigations != "" {
		fmt.Println(HeaderStyle.Render("PRE-EMPTIVE MITIGATIONS"))
		fmt.Println(BodyStyle.Render("  " + report.PreemptiveMitigations))
		fmt.Println()
	}
}
