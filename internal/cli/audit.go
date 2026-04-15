package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewAuditCmd returns the CLI command for viewing the in-memory audit log.
func NewAuditCmd(log *policy.AuditLog) *cobra.Command {
	var filterAction string
	var filterName string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Display the policy audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries := log.Filter(filterAction, filterName)
			fmt.Print(policy.FormatAuditLog(entries))
			return nil
		},
	}

	cmd.Flags().StringVar(&filterAction, "action", "", "Filter by action (e.g. CREATE, UPDATE, DELETE)")
	cmd.Flags().StringVar(&filterName, "name", "", "Filter by policy name")

	return cmd
}
