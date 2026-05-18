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

var (
	memorialType    string
	memorialTitle   string
	memorialContent string
	memorialAuthor  string
)

var (
	distillSlugs   []string
	distillSummary string
	distillAuthor  string
)

var archiveMemorialCmd = &cobra.Command{
	Use:   "memorial <slug>",
	Short: "Record a memorial fragment for a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		slug := args[0]
		return executor.HandleArchiveMemorial(cmd.Context(), appUI, slug, memorialType, memorialTitle, memorialContent, memorialAuthor)
	},
}

var archiveDistillCmd = &cobra.Command{
	Use:   "distill",
	Short: "Consolidate old memory fragments into DISTILLED.md",
	RunE: func(cmd *cobra.Command, args []string) error {
		executor := GetExecutor()
		appUI := tui.NewUI()
		return executor.HandleArchiveDistill(cmd.Context(), appUI, distillSlugs, distillSummary, distillAuthor)
	},
}

func init() {
	archiveMemorialCmd.Flags().StringVar(&memorialType, "type", "", "Type of memorial (lesson, decision, context)")
	archiveMemorialCmd.Flags().StringVar(&memorialTitle, "title", "", "Title of the memorial")
	archiveMemorialCmd.Flags().StringVar(&memorialContent, "content", "", "Content of the memorial")
	archiveMemorialCmd.Flags().StringVar(&memorialAuthor, "author", "agent", "Author of the memorial")
	_ = archiveMemorialCmd.MarkFlagRequired("type")
	_ = archiveMemorialCmd.MarkFlagRequired("title")
	_ = archiveMemorialCmd.MarkFlagRequired("content")

	archiveDistillCmd.Flags().StringSliceVar(&distillSlugs, "slug", []string{}, "Slugs of the fragments to distill")
	archiveDistillCmd.Flags().StringVar(&distillSummary, "summary", "", "Consolidated summary of the fragments")
	archiveDistillCmd.Flags().StringVar(&distillAuthor, "author", "agent", "Author of the distillation")
	_ = archiveDistillCmd.MarkFlagRequired("slug")
	_ = archiveDistillCmd.MarkFlagRequired("summary")

	archiveCmd.AddCommand(archiveInstructionsCmd)
	archiveCmd.AddCommand(archiveMemorialCmd)
	archiveCmd.AddCommand(archiveDistillCmd)
	rootCmd.AddCommand(archiveCmd)
}
