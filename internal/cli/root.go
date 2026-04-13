package cli

import "github.com/spf13/cobra"

// NewRootCmd builds and returns the root cobra command with all sub-commands
// registered.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gorate",
		Short: "Apply and inspect rate-limit policies on HTTP endpoints",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(NewRunCmd())
	root.AddCommand(NewInspectCmd())
	root.AddCommand(NewValidateCmd())
	root.AddCommand(NewStatsCmd())
	root.AddCommand(NewExportCmd())
	root.AddCommand(NewTagCmd())
	root.AddCommand(NewDedupCmd())
	root.AddCommand(NewMergeCmd())
	root.AddCommand(NewRenameCmd())
	root.AddCommand(NewCloneCmd())

	return root
}
