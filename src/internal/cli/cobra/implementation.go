package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var implementationCmd = &cobra.Command{
	Use:   "implementation",
	Short: "Task tracking and implementation status",
	Annotations: map[string]string{
		"IsAgentCommand": "true",
	},
}

var implementationStatusCmd = &cobra.Command{
	Use:   "status [slug]",
	Short: "Show implementation status for a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleImplementationStatus(cmd.Context(), appUI, args[0], jsonMode)
	},
}

var implementationUpdateCmd = &cobra.Command{
	Use:   "update [slug]",
	Short: "Update task status for a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		// We'll need to pass the other flags too.
		return executor.HandleImplementationUpdate(cmd.Context(), appUI, args[0], taskId, taskStatus)
	},
}

var taskId string
var taskStatus string

func init() {
	implementationStatusCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	
	implementationUpdateCmd.Flags().StringVar(&taskId, "task", "", "task ID to update")
	implementationUpdateCmd.Flags().StringVar(&taskStatus, "status", "", "new status for the task")
	_ = implementationUpdateCmd.MarkFlagRequired("task")
	_ = implementationUpdateCmd.MarkFlagRequired("status")

	implementationCmd.AddCommand(implementationStatusCmd)
	implementationCmd.AddCommand(implementationUpdateCmd)
	rootCmd.AddCommand(implementationCmd)
}
