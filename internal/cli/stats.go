package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewStatsCmd returns a cobra command that prints live request stats
// recorded by a running gorate proxy session.
//
// For now the command operates on a shared in-process RequestStats instance
// passed via closure, making it easy to wire up from NewRootCmd.
func NewStatsCmd(stats *policy.RequestStats) *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Display rate-limit request statistics",
		Long:  "Print a summary of allowed and rejected requests per policy collected during the current session.",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap := stats.Snapshot()
			if len(snap) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No stats recorded yet.")
				return nil
			}

			if jsonOut {
				return printStatsJSON(cmd, snap)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%-24s %10s %10s\n", "Policy", "Allowed", "Rejected")
			fmt.Fprintln(cmd.OutOrStdout(), "----------------------------------------------")
			for _, ps := range snap {
				fmt.Fprintf(cmd.OutOrStdout(), "%-24s %10d %10d\n",
					ps.Name, ps.Allowed, ps.Rejected)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output stats as JSON")
	return cmd
}

func printStatsJSON(cmd *cobra.Command, snap []policy.PolicyStats) error {
	fmt.Fprintln(cmd.OutOrStdout(), "[")
	for i, ps := range snap {
		comma := ","
		if i == len(snap)-1 {
			comma = ""
		}
		fmt.Fprintf(cmd.OutOrStdout(),
			`  {"policy":%q,"allowed":%d,"rejected":%d}%s\n`,
			ps.Name, ps.Allowed, ps.Rejected, comma)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "]")
	return nil
}
