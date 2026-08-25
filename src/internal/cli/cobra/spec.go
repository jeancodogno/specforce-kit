package cobra

import (
	"fmt"

	"github.com/jeancodogno/specforce-kit/src/internal/tui"
	"github.com/spf13/cobra"
)

var jsonMode bool
var forceArchive bool
var specType string
var specSize string

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Manage feature specification artifacts",
	Long:  `Manage feature specification artifacts, including initialization, status tracking, and archiving.`,
	Annotations: map[string]string{
		"IsAgentCommand": "true",
	},
}

var specInitCmd = &cobra.Command{
	Use:   "init [slug]",
	Short: "Initialize a new feature specification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleSpecInit(cmd.Context(), appUI, args[0], jsonMode, specType, specSize)
	},
}

var specResizeCmd = &cobra.Command{
	Use:   "resize [slug]",
	Short: "Resize an existing specification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleSpecResize(cmd.Context(), appUI, args[0], specSize, jsonMode)
	},
}

var specListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active feature specifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleSpecList(cmd.Context(), appUI, jsonMode)
	},
}

var specStatusCmd = &cobra.Command{
	Use:   "status [slug]",
	Short: "Show the completeness status of a feature specification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		// We need to pass the json flag somehow. 
		// For now, I'll keep the existing handler's internal flag parsing if possible,
		// or refactor the handler.
		return executor.HandleSpecStatus(cmd.Context(), appUI, args[0], jsonMode)
	},
}

var specArtifactCmd = &cobra.Command{
	Use:   "artifact [slug]",
	Short: "Show details of a specific spec artifact",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		slug := ""
		if len(args) > 0 {
			slug = args[0]
		}
		return executor.HandleSpecArtifact(cmd.Context(), appUI, slug, jsonMode)
	},
}

var specArchiveCmd = &cobra.Command{
	Use:   "archive [slug]",
	Short: "Archive a completed feature specification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleSpecArchive(cmd.Context(), appUI, args[0], forceArchive)
	},
}

var auditError string
var auditClear bool
var auditIteration int
var auditValid bool
var auditInvalid bool

var specAuditCmd = &cobra.Command{
	Use:   "audit [slug]",
	Short: "Update refinement audit state for a specification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		
		// Build manual args for HandleSpec if we want to reuse it, 
		// but HandleSpec calls handleSpecAuditCmd which parses flags from args.
		// It's cleaner to just call HandleSpec with reconstructed args.
		cliArgs := []string{"audit", args[0]}
		if auditError != "" {
			cliArgs = append(cliArgs, "--error", auditError)
		}
		if auditClear {
			cliArgs = append(cliArgs, "--clear")
		}
		if auditIteration > 0 {
			cliArgs = append(cliArgs, "--iteration", fmt.Sprintf("%d", auditIteration))
		}
		if auditValid {
			cliArgs = append(cliArgs, "--valid")
		}
		if auditInvalid {
			cliArgs = append(cliArgs, "--invalid")
		}
		if jsonMode {
			cliArgs = append(cliArgs, "--json")
		}

		return executor.HandleSpec(cmd.Context(), appUI, cliArgs...)
	},
}

func init() {
	specInitCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specInitCmd.Flags().StringVar(&specType, "type", "feature", "type of specification (feature, bug)")
	specInitCmd.Flags().StringVar(&specSize, "size", "medium", "Specification size (small, medium, large, complex)")
	specListCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specStatusCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specArtifactCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specArchiveCmd.Flags().BoolVar(&forceArchive, "force", false, "force archive even if pending tasks remain")

	specResizeCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specResizeCmd.Flags().StringVar(&specSize, "size", "medium", "Specification size (small, medium, large, complex)")

	specAuditCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specAuditCmd.Flags().StringVar(&auditError, "error", "", "append a coherence error message")
	specAuditCmd.Flags().BoolVar(&auditClear, "clear", false, "clear all coherence errors")
	specAuditCmd.Flags().IntVar(&auditIteration, "iteration", 0, "set current refinement iteration count")
	specAuditCmd.Flags().BoolVar(&auditValid, "valid", false, "mark spec as valid")
	specAuditCmd.Flags().BoolVar(&auditInvalid, "invalid", false, "mark spec as invalid")

	specCmd.AddCommand(specInitCmd)
	specCmd.AddCommand(specListCmd)
	specCmd.AddCommand(specStatusCmd)
	specCmd.AddCommand(specArtifactCmd)
	specCmd.AddCommand(specArchiveCmd)
	specCmd.AddCommand(specResizeCmd)
	specCmd.AddCommand(specAuditCmd)
	rootCmd.AddCommand(specCmd)
}
