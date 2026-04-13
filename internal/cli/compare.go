package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewCompareCmd returns a cobra command that compares two policy files.
func NewCompareCmd() *cobra.Command {
	var fileA, fileB string

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare two policy files and report differences",
		RunE: func(cmd *cobra.Command, args []string) error {
			policiesA, err := policy.LoadFromFile(fileA)
	 {
				return fmt.Errorf("loading file-a: %w", err)
			}

			policiesB, err := policy.LoadFromFile(fileB)
			if err != nil {
				return fmt.Errorf("loading file-b: %w", err)
			}

			result := policy.Compare(policiesA, policiesB)
			out := policy.FormatCompare(result)
			fmt.Fprint(os.Stdout, out)

			if len(result.Conflicts) > 0 || len(result.OnlyInA) > 0 || len(result.OnlyInB) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&fileA, "file-a", "", "first policy file (required)")
	cmd.Flags().StringVar(&fileB, "file-b", "", "second policy file (required)")
	_ = cmd.MarkFlagRequired("file-a")
	_ = cmd.MarkFlagRequired("file-b")

	return cmd
}
