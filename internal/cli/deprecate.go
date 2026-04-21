package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewDeprecateCmd returns the root deprecate command with subcommands.
func NewDeprecateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deprecate",
		Short: "Manage policy deprecation markers",
	}
	cmd.AddCommand(newDeprecateSetCmd())
	cmd.AddCommand(newDeprecateUnsetCmd())
	cmd.AddCommand(newDeprecateListCmd())
	return cmd
}

func newDeprecateSetCmd() *cobra.Command {
	var file, name, reason, replacement string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Mark a policy as deprecated",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			policies, err = policy.Deprecate(policies, name, reason, replacement)
			if err != nil {
				return err
			}
			if err := writeFile(file, policies); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "policy %q marked as deprecated\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "policy file (required)")
	cmd.Flags().StringVar(&name, "name", "", "policy name (required)")
	cmd.Flags().StringVar(&reason, "reason", "", "reason for deprecation")
	cmd.Flags().StringVar(&replacement, "replacement", "", "replacement policy name")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newDeprecateUnsetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Remove deprecation marker from a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			policies, err = policy.Undeprecate(policies, name)
			if err != nil {
				return err
			}
			if err := writeFile(file, policies); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deprecation removed from policy %q\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "policy file (required)")
	cmd.Flags().StringVar(&name, "name", "", "policy name (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newDeprecateListCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all deprecated policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			deprecated := policy.ListDeprecated(policies)
			if len(deprecated) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no deprecated policies")
				return nil
			}
			for _, p := range deprecated {
				info, _ := policy.GetDeprecationInfo(p)
				fmt.Fprintf(cmd.OutOrStdout(), "- %s (since: %s, reason: %s, replacement: %s)\n",
					info.Name, info.Since, info.Reason, info.Replacement)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "policy file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
