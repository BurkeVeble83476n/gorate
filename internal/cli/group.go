package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewGroupCmd returns the CLI command for grouping policies.
func NewGroupCmd() *cobra.Command {
	var file string
	var by string

	cmd := &cobra.Command{
		Use:   "group",
		Short: "Group policies by a field (method, endpoint, window)",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			groups, err := policy.Group(policies, policy.GroupBy(by))
			if err != nil {
				return fmt.Errorf("grouping policies: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			for _, g := range groups {
				fmt.Fprintf(w, "[%s] (%d policies)\n", g.Key, len(g.Policies))
				for _, p := range g.Policies {
					fmt.Fprintf(w, "  %-20s\t%s\t%s\t%d req\n",
						p.Name, p.Method, p.Endpoint, p.Limit)
				}
			}
			w.Flush()
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().StringVarP(&by, "by", "b", "method", "Field to group by: method, endpoint, window")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
