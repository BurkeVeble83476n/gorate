package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/gorate/internal/policy"
)

func NewExpiryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expiry",
		Short: "Manage policy expiry times",
	}
	cmd.AddCommand(newExpirySetCmd())
	cmd.AddCommand(newExpiryRemoveCmd())
	cmd.AddCommand(newExpiryStatusCmd())
	return cmd
}

func newExpirySetCmd() *cobra.Command {
	var file, name, expiresAt string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set an expiry time on a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			at, err := time.Parse(time.RFC3339, expiresAt)
			if err != nil {
				return fmt.Errorf("invalid time format (use RFC3339): %w", err)
			}
			updated, err := policy.SetExpiry(policies, name, at)
			if err != nil {
				return err
			}
			if err := writeFile(file, updated); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "expiry set on %q until %s\n", name, at.Format(time.RFC3339))
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "policy name (required)")
	cmd.Flags().StringVar(&expiresAt, "at", "", "expiry time in RFC3339 format (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("at")
	return cmd
}

func newExpiryRemoveCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove expiry from a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			updated, err := policy.RemoveExpiry(policies, name)
			if err != nil {
				return err
			}
			if err := writeFile(file, updated); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "expiry removed from %q\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "policy name (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newExpiryStatusCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show expiry status for all policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			results := policy.EvaluateExpiry(policies)
			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no policies with expiry set")
				return nil
			}
			for _, r := range results {
				status := "active"
				if r.Expired {
					status = "EXPIRED"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-8s %s\n", r.Name, status, r.ExpiresAt.Format(time.RFC3339))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "policy file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
