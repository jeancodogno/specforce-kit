package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Launch the Specforce Console TUI",
	Long:  `Launch an interactive console to manage specifications, constitution, and implementation tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleConsole(cmd.Context(), appUI)
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
