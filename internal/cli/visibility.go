package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/gorate/internal/policy"
)

func NewVisibilityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "visibility",
		Short: "Manage policy visibility (public, internal, private)",
	}
	cmd.AddCommand(newVisibilitySetCmd(), newVisibilityGetCmd(), newVisibilityFilterCmd())
	return cmd
}

func newVisibilitySetCmd() *cobra.Command {
	var file, name, level string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set visibility on a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			updated, err := policy.SetVisibility(policies, name, level)
			if err != nil {
				return err
			}
			return writeFile(file, updated)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Policy file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	cmd.Flags().StringVar(&level, "level", "", "Visibility level (public, internal, private)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("level")
	return cmd
}

func newVisibilityGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get visibility of a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			v, err := policy.GetVisibility(policies, name)
			if err != nil {
				return err
			}
			if v == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(v)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Policy file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newVisibilityFilterCmd() *cobra.Command {
	var file, level string
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "List policies by visibility",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			results := policy.FilterByVisibility(policies, level)
			policy.PrintTable(results)
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Policy file")
	cmd.Flags().StringVar(&level, "level", "", "Visibility level to filter by")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("level")
	return cmd
}
