package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewScheduleCmd returns the cobra command for evaluating policy schedules.
func NewScheduleCmd() *cobra.Command {
	var file string
	var fromStr string
	var toStr string
	var policyName string

	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Evaluate time-based schedule for a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			found := false
			for _, p := range policies {
				if p.Name == policyName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("policy %q not found in file", policyName)
			}

			from, err := time.Parse(time.RFC3339, fromStr)
			if err != nil {
				return fmt.Errorf("invalid --from value: %w", err)
			}
			to, err := time.Parse(time.RFC3339, toStr)
			if err != nil {
				return fmt.Errorf("invalid --to value: %w", err)
			}

			s := policy.Schedule{
				PolicyName: policyName,
				ActiveFrom: from,
				ActiveTo:   to,
			}
			if err := policy.ValidateSchedule(s); err != nil {
				return fmt.Errorf("invalid schedule: %w", err)
			}

			results := policy.EvaluateSchedules([]policy.Schedule{s}, time.Now())
			for _, r := range results {
				status := "INACTIVE"
				if r.Active {
					status = "ACTIVE"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s — %s\n", status, r.PolicyName, r.Reason)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().StringVar(&fromStr, "from", "", "Schedule start time in RFC3339 format (required)")
	cmd.Flags().StringVar(&toStr, "to", "", "Schedule end time in RFC3339 format (required)")
	cmd.Flags().StringVar(&policyName, "name", "", "Name of the policy to schedule (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
