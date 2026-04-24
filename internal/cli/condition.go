package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewConditionCmd returns the root condition command with subcommands.
func NewConditionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "condition",
		Short: "Manage conditions on rate-limit policies",
	}
	cmd.AddCommand(newConditionSetCmd())
	cmd.AddCommand(newConditionGetCmd())
	return cmd
}

func newConditionSetCmd() *cobra.Command {
	var file, name, field, operator, value string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Attach a condition to a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			updated, err := policy.SetCondition(policies, name, field, operator, value)
			if err != nil {
				return err
			}
			if err := writeFile(file, updated); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "condition set on policy %q\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "path to policy file")
	cmd.Flags().StringVar(&name, "name", "", "policy name")
	cmd.Flags().StringVar(&field, "field", "", "attribute field to evaluate")
	cmd.Flags().StringVar(&operator, "operator", "eq", "operator: eq, neq, contains, prefix")
	cmd.Flags().StringVar(&value, "value", "", "expected value")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newConditionGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show the condition attached to a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			cond, err := policy.GetCondition(policies, name)
			if err != nil {
				return err
			}
			if cond == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "no condition set on policy %q\n", name)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "field=%s operator=%s value=%s\n", cond.Field, cond.Operator, cond.Value)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "path to policy file")
	cmd.Flags().StringVar(&name, "name", "", "policy name")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
