package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewReviewCmd returns the root review command with subcommands.
func NewReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Manage review status of rate-limit policies",
	}
	cmd.AddCommand(newReviewSetCmd())
	cmd.AddCommand(newReviewGetCmd())
	cmd.AddCommand(newReviewListCmd())
	return cmd
}

func newReviewSetCmd() *cobra.Command {
	var file, name, status, reviewer, comment string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set the review status of a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			updated, err := policy.SetReview(policies, name, policy.ReviewStatus(status), reviewer, comment)
			if err != nil {
				return err
			}
			if err := writeFile(file, updated); err != nil {
				return fmt.Errorf("saving policies: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Review status for %q set to %s\n", name, status)
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Policy name (required)")
	cmd.Flags().StringVarP(&status, "status", "s", "", "Review status: approved|rejected|pending (required)")
	cmd.Flags().StringVar(&reviewer, "reviewer", "", "Reviewer name")
	cmd.Flags().StringVar(&comment, "comment", "", "Review comment")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("status")
	return cmd
}

func newReviewGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the review status of a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			entry, err := policy.GetReview(policies, name)
			if err != nil {
				return err
			}
			fmt.Print(policy.FormatReview(entry))
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Policy name (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newReviewListCmd() *cobra.Command {
	var file, status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List policies by review status",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			s := policy.ReviewStatus(status)
			if s == "" {
				s = policy.ReviewPending
			}
			entries := policy.ListByReviewStatus(policies, s)
			if len(entries) == 0 {
				fmt.Fprintf(os.Stdout, "No policies with status %q\n", s)
				return nil
			}
			for _, e := range entries {
				fmt.Print(policy.FormatReview(e))
				fmt.Println("---")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status: approved|rejected|pending")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
