package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root cobra command for gorate.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gorate",
		Short: "A lightweight CLI tool for applying and inspecting rate-limit policies",
		Long: `gorate helps you apply, inspect, and manage HTTP rate-limit policies
during local development. Use subcommands to run a proxy, inspect policies,
validate configuration, export, diff, merge, and more.`,
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
	root.AddCommand(NewPatchCmd())
	root.AddCommand(NewDiffCmd())
	root.AddCommand(NewSnapshotCmd())

	return root
}
