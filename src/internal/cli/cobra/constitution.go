package cobra

import (
	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var constitutionCmd = &cobra.Command{
	Use:   "constitution",
	Short: "Manage project constitution docs",
	Annotations: map[string]string{
		"IsAgentCommand": "true",
	},
}

var constitutionStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the completeness status of the project constitution",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleConstitutionStatus(cmd.Context(), appUI, jsonMode)
	},
}

var constitutionArtifactCmd = &cobra.Command{
	Use:   "artifact [slug]",
	Short: "Show details of a specific constitution artifact",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		slug := ""
		if len(args) > 0 {
			slug = args[0]
		}
		return executor.HandleConstitutionArtifact(cmd.Context(), appUI, slug, jsonMode)
	},
}

func init() {
	constitutionStatusCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")
	constitutionArtifactCmd.Flags().BoolVar(&jsonMode, "json", false, "output in machine-readable JSON format")

	constitutionCmd.AddCommand(constitutionStatusCmd)
	constitutionCmd.AddCommand(constitutionArtifactCmd)
	rootCmd.AddCommand(constitutionCmd)
}
