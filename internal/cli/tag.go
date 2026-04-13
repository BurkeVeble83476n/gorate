package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewTagCmd returns a cobra command for filtering policies by tag.
func NewTagCmd() *cobra.Command {
	var filePath string
	var tag string

	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Filter and display policies by tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			if tag == "" {
				return fmt.Errorf("--tag flag is required")
			}

			matched := policy.FilterByTag(policies, tag)

			if len(matched) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No policies found with tag %q\n", tag)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Policies with tag %q:\n", tag)
			policy.PrintTable(matched, cmd.OutOrStdout())
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "policies.yaml", "Path to the policies file")
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag to filter policies by (required)")

	return cmd
}
