package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewMergeCmd returns a cobra command that merges two policy files.
func NewMergeCmd() *cobra.Command {
	var strategy string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "merge <file1> <file2>",
		Short: "Merge two policy files into one, resolving conflicts by strategy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			policiesA, err := policy.LoadFromFile(args[0])
			if err != nil {
				return fmt.Errorf("loading %s: %w", args[0], err)
			}

			policiesB, err := policy.LoadFromFile(args[1])
			if err != nil {
				return fmt.Errorf("loading %s: %w", args[1], err)
			}

			merged, err := policy.Merge(policy.MergeStrategy(strategy), policiesA, policiesB)
			if err != nil {
				return fmt.Errorf("merging policies: %w", err)
			}

			format := "yaml"
			if outputFile != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "writing %d policies to %s\n", len(merged), outputFile)
			}

			out, err := policy.Export(merged, policy.ParseFormat(format))
			if err != nil {
				return fmt.Errorf("exporting merged policies: %w", err)
			}

			fmt.Fprint(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "keep-first",
		"Conflict resolution strategy: keep-first, keep-last, highest-limit, lowest-limit")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"Optional output file path (stdout if omitted)")

	return cmd
}
