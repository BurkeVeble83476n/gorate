package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewRenameCmd creates the 'rename' subcommand which renames a policy within a file.
func NewRenameCmd() *cobra.Command {
	var file, oldName, newName string

	cmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename a rate-limit policy by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			updated, err := policy.Rename(policies, policy.RenameOptions{
				OldName: oldName,
				NewName: newName,
			})
			if err != nil {
				return fmt.Errorf("renaming policy: %w", err)
			}

			out, err := policy.Export(updated, policy.FormatYAML)
			if err != nil {
				return fmt.Errorf("exporting policies: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to the policy file (required)")
	cmd.Flags().StringVar(&oldName, "old", "", "Current name of the policy (required)")
	cmd.Flags().StringVar(&newName, "new", "", "New name for the policy (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("old")
	_ = cmd.MarkFlagRequired("new")

	return cmd
}
