package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Global framework installation",
	Long:  `Perform the initial setup and install the Specforce framework globally.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleInstall(cmd.Context(), appUI)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
