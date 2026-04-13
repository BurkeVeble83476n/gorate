package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewProfileCmd creates the `gorate profile` subcommand.
func NewProfileCmd() *cobra.Command {
	var file string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Show a risk and lint profile for each policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			if len(policies) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No policies found.")
				return nil
			}

			profiles := policy.Profile(policies)

			for i, p := range profiles {
				if i > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "---")
				}
				if verbose {
					fmt.Fprint(cmd.OutOrStdout(), policy.FormatProfile(p))
				} else {
					issue := ""
					if p.IssueCount > 0 {
						issue = fmt.Sprintf(" [%d issue(s)]", p.IssueCount)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%-20s risk=%-4d %s%s\n",
						p.Name, p.RiskScore, p.Endpoint, issue)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full profile details")
	_ = cmd.MarkFlagRequired("file")

	if err := cmd.MarkFlagRequired("file"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return cmd
}
