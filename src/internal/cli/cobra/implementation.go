package cobra

import (
	"os"

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
		return executor.HandleImplementationUpdate(cmd.Context(), appUI, args[0], taskIds, taskStatus)
	},
}

var implementationQACmd = &cobra.Command{
	Use:   "qa",
	Short: "Persist the independent QA report",
}

var implementationQARecordCmd = &cobra.Command{
	Use:   "record [slug]",
	Short: "Record a QA report (JSON on stdin) for a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleImplementationQARecord(cmd.Context(), appUI, args[0], os.Stdin)
	},
}

var taskIds []string
var taskStatus string

func init() {
	implementationStatusCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	
	implementationUpdateCmd.Flags().StringSliceVar(&taskIds, "task", nil, "task ID(s) to update (comma-separated or repeated)")
	implementationUpdateCmd.Flags().StringVar(&taskStatus, "status", "", "new status for the task")
	_ = implementationUpdateCmd.MarkFlagRequired("task")
	_ = implementationUpdateCmd.MarkFlagRequired("status")

	implementationCmd.AddCommand(implementationStatusCmd)
	implementationCmd.AddCommand(implementationUpdateCmd)
	implementationQACmd.AddCommand(implementationQARecordCmd)
	implementationCmd.AddCommand(implementationQACmd)
	rootCmd.AddCommand(implementationCmd)
}
