package cobra

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jeancodogno/specforce-kit/src/internal/cli"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
	"github.com/jeancodogno/specforce-kit/src/internal/upgrade"
	"github.com/spf13/cobra"
)

var (
	devMode        bool
	upgradeService *upgrade.Service
	appVersion     string

	// isTTY reports whether stdout is a terminal; injectable for tests.
	isTTY = tui.IsTTY
)

var rootCmd = &cobra.Command{
	Use:   "specforce",
	Short: tui.Tagline,
	Long: tui.Tagline + `

Specforce is an AI-native framework designed to streamline the software 
development lifecycle through Spec-Driven Development (SDD). It coordinates 
AI agents and human developers to build high-quality software from 
well-defined specifications.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Upgrade Service
		mgr, err := upgrade.NewStateManager()
		if err != nil {
			return nil // Don't block CLI if state fails
		}

		// Version source follows the installation method (npm channel or GitHub)
		exe, _ := os.Executable()
		upgradeService = upgrade.NewService(mgr, upgrade.SelectProvider(exe, appVersion), appVersion)

		// Check for update if TTY and not an agent command
		interactive := isTTY()

		// Cleanup .old file from previous upgrade if it exists
		if exe, err := os.Executable(); err == nil {
			if _, err := os.Stat(exe + ".old"); err == nil {
				_ = os.Remove(exe + ".old")
			}
		}

		isAgentCmd := isAgentCommand(cmd)

		if os.Getenv("DEBUG") == "1" {
			_, _ = fmt.Fprintf(os.Stderr, "[DEBUG] UpgradeCheck: TTY=%v, AgentCmd=%v\n", interactive, isAgentCmd)
		}

		if interactive && !isAgentCmd {
			upgradeService.CheckForUpdate()
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		showUpdateNotice(cmd)
	},
}

// isAgentCommand reports whether cmd or one of its parents is an agent command.
func isAgentCommand(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations["IsAgentCommand"] == "true" {
			return true
		}
	}
	return false
}

// stableSwitchCommand is the command that moves an npm installation to the
// stable channel; switching channels is always done by hand.
const stableSwitchCommand = "npm install -g @jeancodogno/specforce-kit@latest"

// showUpdateNotice writes the pending updater notice to the command's error
// stream, only for interactive (stdout is a terminal), non-agent commands.
// Notices are marked as shown only after they were written.
func showUpdateNotice(cmd *cobra.Command) {
	if upgradeService == nil || !isTTY() || isAgentCommand(cmd) {
		return
	}

	notice := upgradeService.PendingNotice()
	var text string
	switch notice.Kind {
	case upgrade.NoticeApplied:
		text = tui.RenderUpdateAppliedNotification(notice.Version)
	case upgrade.NoticeManual:
		text = tui.RenderManualUpdateNotification(appVersion, notice.Version, upgradeService.ManualUpdateCommand(notice.Version))
	case upgrade.NoticeStable:
		text = tui.RenderStableNotification(notice.Version, upgrade.ChannelOf(appVersion), stableSwitchCommand)
	case upgrade.NoticeAvailable:
		text = tui.RenderUpdateNotification(appVersion, notice.Version)
	default:
		return
	}

	if _, err := fmt.Fprint(cmd.ErrOrStderr(), text); err != nil {
		return
	}
	_ = upgradeService.MarkNoticeShown(notice)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, v string) error {
	tui.ConfigureColor(os.Getenv, tui.IsTTY(), lipgloss.HasDarkBackground)
	appVersion = v
	rootCmd.Version = v
	rootCmd.SetVersionTemplate("specforce version {{.Version}}\n")
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
