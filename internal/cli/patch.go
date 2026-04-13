package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewPatchCmd creates a cobra command for patching a single policy field.
func NewPatchCmd() *cobra.Command {
	var file string
	var name string
	var limit int
	var window string
	var method string

	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Partially update a policy's fields by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			opts := policy.PatchOptions{}

			if cmd.Flags().Changed("limit") {
				v := limit
				opts.Limit = &v
			}
			if cmd.Flags().Changed("window") {
				opts.Window = &window
			}
			if cmd.Flags().Changed("method") {
				opts.Method = &method
			}

			updated, err := policy.Patch(policies, name, opts)
			if err != nil {
				return fmt.Errorf("patch failed: %w", err)
			}

			out, err := policy.Export(updated, "yaml")
			if err != nil {
				return fmt.Errorf("failed to export policies: %w", err)
			}

			if err := os.WriteFile(file, []byte(out), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Policy %q patched successfully.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "policies.yaml", "Path to the policies file")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the policy to patch (required)")
	cmd.Flags().IntVar(&limit, "limit", 0, "New request limit")
	cmd.Flags().StringVar(&window, "window", "", "New time window (e.g. 1m, 30s)")
	cmd.Flags().StringVar(&method, "method", "", "New HTTP method")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
