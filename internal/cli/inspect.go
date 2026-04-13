package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/gorate/internal/policy"
)

// NewInspectCmd returns a cobra command that loads and displays rate-limit policies.
func NewInspectCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect rate-limit policies from a configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			if len(policies) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No policies defined.")
				return nil
			}

			policy.PrintTable(cmd.OutOrStdout(), policies)
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "policies.yaml", "Path to the policies YAML file")

	_ = os.Stderr // ensure os import used
	return cmd
}
