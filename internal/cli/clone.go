package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewCloneCmd returns a cobra command that clones an existing policy under a
// new name within a policies file.
func NewCloneCmd() *cobra.Command {
	var (
		file     string
		newName  string
		override bool
		output   string
	)

	cmd := &cobra.Command{
		Use:   "clone <source-name>",
		Short: "Clone an existing rate-limit policy under a new name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceName := args[0]

			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			result, err := policy.Clone(policies, sourceName, policy.CloneOptions{
				NewName:  newName,
				Override: override,
			})
			if err != nil {
				return err
			}

			dest := file
			if output != "" {
				dest = output
			}

			data, err := policy.Export(result, "yaml")
			if err != nil {
				return fmt.Errorf("serialising policies: %w", err)
			}

			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Cloned %q → %q in %s\n", sourceName, newName, dest)
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "policies.yaml", "path to policies file")
	cmd.Flags().StringVarP(&newName, "name", "n", "", "new policy name (required)")
	cmd.Flags().BoolVar(&override, "override", false, "overwrite if a policy with the new name already exists")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to this file instead of the source file")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
