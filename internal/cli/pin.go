package cli

import (
	"fmt"

	"github.com/jdoe/gorate/internal/policy"
	"github.com/spf13/cobra"
)

// NewPinCmd returns the root pin command with pin/unpin/list subcommands.
func NewPinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin or unpin rate-limit policies to protect them from automated changes",
	}
	cmd.AddCommand(newPinSetCmd())
	cmd.AddCommand(newPinUnsetCmd())
	cmd.AddCommand(newPinListCmd())
	return cmd
}

func newPinSetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Pin a policy by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			updated, err := policy.Pin(policies, name)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Policy %q pinned.\n", name)
			_ = updated
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the policy to pin (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newPinUnsetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unpin a policy by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			updated, err := policy.Unpin(policies, name)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Policy %q unpinned.\n", name)
			_ = updated
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the policy to unpin (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newPinListCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			pinned := policy.ListPinned(policies)
			if len(pinned) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No pinned policies.")
				return nil
			}
			for _, p := range pinned {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s (%s %s)\n", p.Name, p.Method, p.Endpoint)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
