package cli

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewScoreCmd returns the cobra command for scoring policies.
func NewScoreCmd() *cobra.Command {
	var file string
	var sortByScore bool

	cmd := &cobra.Command{
		Use:   "score",
		Short: "Score policies by risk and permissiveness",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			scores := policy.Score(policies)

			if sortByScore {
				sort.Slice(scores, func(i, j int) bool {
					return scores[i].Score > scores[j].Score
				})
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSCORE\tBREAKDOWN")
			for _, s := range scores {
				breakdown := "-"
				if len(s.Breakdown) > 0 {
					breakdown = s.Breakdown[0]
					if len(s.Breakdown) > 1 {
						breakdown += fmt.Sprintf(" (+%d more)", len(s.Breakdown)-1)
					}
				}
				fmt.Fprintf(w, "%s\t%d\t%s\n", s.Name, s.Score, breakdown)
			}
			return w.Flush()
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().BoolVar(&sortByScore, "sort", false, "Sort output by score descending")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
