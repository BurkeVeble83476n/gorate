package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewSnapshotCmd creates the snapshot command with save and load subcommands.
func NewSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save or load a snapshot of policies",
	}
	cmd.AddCommand(newSnapshotSaveCmd())
	cmd.AddCommand(newSnapshotLoadCmd())
	return cmd
}

func newSnapshotSaveCmd() *cobra.Command {
	var file, output, label string
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save current policies to a snapshot file",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}
			if err := policy.SaveSnapshot(policies, label, output); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s\n", output)
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file to snapshot (required)")
	cmd.Flags().StringVarP(&output, "output", "o", "snapshot.json", "Output snapshot file path")
	cmd.Flags().StringVarP(&label, "label", "l", "", "Label for this snapshot")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newSnapshotLoadCmd() *cobra.Command {
	var snapshotFile string
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load and describe a snapshot file",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := policy.LoadSnapshot(snapshotFile)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), policy.DescribeSnapshot(snap))
			policy.PrintTable(snap.Policies, cmd.OutOrStdout())
			return nil
		},
	}
	cmd.Flags().StringVarP(&snapshotFile, "snapshot", "s", "", "Snapshot file to load (required)")
	_ = cmd.MarkFlagRequired("snapshot")
	return cmd
}
