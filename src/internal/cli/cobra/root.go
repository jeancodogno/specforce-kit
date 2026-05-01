package cobra

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/jeancodogno/specforce-kit/src/internal/cli"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
)

var (
	devMode        bool
	upgradeService *upgrade.Service
	appVersion     string
)

var rootCmd = &cobra.Command{
	Use:   "specforce",
	Short: "Specforce: AI-Native Software Design and Delivery Framework",
	Long: `Specforce is an AI-native framework designed to streamline the software 
development lifecycle through Spec-Driven Development (SDD). It coordinates 
AI agents and human developers to build high-quality software from 
well-defined specifications.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Upgrade Service
		mgr, err := upgrade.NewStateManager()
		if err != nil {
			return nil // Don't block CLI if state fails
		}

		// TODO: Detect if installed via NPM to choose provider
		provider := upgrade.NewGitHubProvider()
		upgradeService = upgrade.NewService(mgr, provider, appVersion)

		// Check for update if TTY and not an agent command
		isTTY := tui.IsTTY()
		isAgentCmd := false
		for c := cmd; c != nil; c = c.Parent() {
			if c.Annotations["IsAgentCommand"] == "true" {
				isAgentCmd = true
				break
			}
		}

		if os.Getenv("DEBUG") == "1" {
			_, _ = fmt.Fprintf(os.Stderr, "[DEBUG] UpgradeCheck: TTY=%v, AgentCmd=%v\n", isTTY, isAgentCmd)
		}

		if isTTY && !isAgentCmd {
			upgradeService.CheckForUpdate(cmd.Context())
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if upgradeService == nil {
			return
		}

		// Only show notification if TTY and not an agent command
		isTTY := tui.IsTTY()
		isAgentCmd := false
		for c := cmd; c != nil; c = c.Parent() {
			if c.Annotations["IsAgentCommand"] == "true" {
				isAgentCmd = true
				break
			}
		}

		if isTTY && !isAgentCmd {
			if available, latest := upgradeService.IsUpdateAvailable(); available {
				fmt.Print(tui.RenderUpdateNotification(appVersion, latest))
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, v string) error {
	appVersion = v
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
	e := cli.NewExecutor(appVersion)
	e.DevMode = devMode
	return e
}
