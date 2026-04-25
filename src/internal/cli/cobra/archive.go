package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Manage feature archiving and lifecycle",
}

var archiveInstructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Show instructions for archiving a feature",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleArchiveInstructions(cmd.Context(), appUI)
	},
}

func init() {
	archiveCmd.AddCommand(archiveInstructionsCmd)
	rootCmd.AddCommand(archiveCmd)
}
