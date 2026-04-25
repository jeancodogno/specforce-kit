package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var initCmd = &cobra.Command{
	Use:   "init [agents...]",
	Short: "Initialize a new Specforce project",
	Long: `Initialize a new Specforce project in the current directory.
If agents are provided, they will be initialized immediately.
Otherwise, an interactive TUI will allow you to select agents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleInit(cmd.Context(), appUI, args...)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
