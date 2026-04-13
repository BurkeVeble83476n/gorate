package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewDiffCmd returns a cobra command that diffs two policy files.
func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare two policy files and show differences",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseFile, _ := cmd.Flags().GetString("base")
			updatedFile, _ := cmd.Flags().GetString("updated")

			if baseFile == "" || updatedFile == "" {
				return fmt.Errorf("both --base and --updated flags are required")
			}

			base, err := policy.LoadFromFile(baseFile)
			if err != nil {
				return fmt.Errorf("loading base file: %w", err)
			}

			updated, err := policy.LoadFromFile(updatedFile)
			if err != nil {
				return fmt.Errorf("loading updated file: %w", err)
			}

			result := policy.Diff(base, updated)
			out := policy.FormatDiff(result)
			fmt.Fprint(os.Stdout, out)
			return nil
		},
	}

	cmd.Flags().String("base", "", "Path to the base policy file")
	cmd.Flags().String("updated", "", "Path to the updated policy file")
	return cmd
}
