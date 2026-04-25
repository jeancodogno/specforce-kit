package cobra

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/cli"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

var (
	devMode bool
)

var rootCmd = &cobra.Command{
	Use:   "specforce",
	Short: "Specforce: AI-Native Software Design and Delivery Framework",
	Long: `Specforce is an AI-native framework designed to streamline the software 
development lifecycle through Spec-Driven Development (SDD). It coordinates 
AI agents and human developers to build high-quality software from 
well-defined specifications.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, v string) error {
	rootCmd.Version = v
	rootCmd.SetVersionTemplate("specforce version {{.Version}}\n")
	tui.AppVersion = v
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "use local kit/ for development instead of embedded FS")
}

// GetDevMode returns whether development mode is active.
func GetDevMode() bool {
	return devMode
}

// GetExecutor returns a new CLI executor with current settings.
func GetExecutor() *cli.Executor {
	e := cli.NewExecutor(rootCmd.Version)
	e.DevMode = devMode
	return e
}
