package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds and returns the root cobra command for gorate.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gorate",
		Short: "A lightweight CLI tool for rate-limit policy management",
		Long: `gorate helps you apply, inspect, and manage HTTP rate-limit policies
during local development. Use subcommands to load, validate, export, and
manipulate policy files.`,
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
	root.AddCommand(NewGroupCmd())
	root.AddCommand(NewTemplateCmd())
	root.AddCommand(NewScheduleCmd())
	root.AddCommand(NewCompareCmd())
	root.AddCommand(NewScoreCmd())
	root.AddCommand(NewProfileCmd())
	root.AddCommand(NewAnnotateCmd())
	root.AddCommand(NewLabelCmd())
	root.AddCommand(NewTransformCmd())
	root.AddCommand(NewWatchdogCmd())
	root.AddCommand(NewDependencyCmd())
	root.AddCommand(NewInheritCmd())
	root.AddCommand(NewVisibilityCmd())
	root.AddCommand(NewExpiryCmd())
	root.AddCommand(NewArchiveCmd())
	root.AddCommand(NewPinCmd())
	root.AddCommand(NewAuditCmd())
	root.AddCommand(NewDeprecateCmd())

	return root
}
