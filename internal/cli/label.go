package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/gorate/internal/policy"
)

// NewLabelCmd returns the root label command with subcommands.
func NewLabelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels on rate-limit policies",
	}
	cmd.AddCommand(newLabelAddCmd())
	cmd.AddCommand(newLabelRemoveCmd())
	cmd.AddCommand(newLabelGetCmd())
	return cmd
}

func newLabelAddCmd() *cobra.Command {
	var file, name, label string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a label to a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name flag is required")
			}
			if label == "" {
				return fmt.Errorf("--label flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			updated, err := policy.AddLabel(policies, name, label)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "label %q added to policy %q\n", label, name)
			_ = updated
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&label, "label", "", "Label to add")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newLabelRemoveCmd() *cobra.Command {
	var file, name, label string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a label from a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			_, err = policy.RemoveLabel(policies, name, label)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "label %q removed from policy %q\n", label, name)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&label, "label", "", "Label to remove")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newLabelGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get labels for a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			labels, err := policy.GetLabels(policies, name)
			if err != nil {
				return err
			}
			if len(labels) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no labels)")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), strings.Join(labels, ", "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
