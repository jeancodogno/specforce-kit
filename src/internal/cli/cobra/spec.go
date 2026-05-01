package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var jsonMode bool
var forceArchive bool

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
		// We'll pass the slug as the first argument in a slice
		return executor.HandleSpecInit(cmd.Context(), appUI, args[0], jsonMode)
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

func init() {
	specInitCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specListCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specStatusCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specArtifactCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	specArchiveCmd.Flags().BoolVar(&forceArchive, "force", false, "force archive even if pending tasks remain")

	specCmd.AddCommand(specInitCmd)
	specCmd.AddCommand(specListCmd)
	specCmd.AddCommand(specStatusCmd)
	specCmd.AddCommand(specArtifactCmd)
	specCmd.AddCommand(specArchiveCmd)
	rootCmd.AddCommand(specCmd)
}
