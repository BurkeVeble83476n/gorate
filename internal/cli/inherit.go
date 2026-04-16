package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewInheritCmd returns a cobra command that resolves policy inheritance
// from a YAML file and prints the resulting policies.
func NewInheritCmd() *cobra.Command {
	var file string
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "inherit",
		Short: "Resolve policy inheritance and display merged results",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			resolved, notes, err := policy.ApplyInheritance(policies)
			if err != nil {
				return fmt.Errorf("inheritance resolution failed: %w", err)
			}

			for _, note := range notes {
				fmt.Fprintln(cmd.OutOrStdout(), "  "+note)
			}

			if outputJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(resolved)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nResolved %d policies:\n", len(resolved))
			for _, p := range resolved {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s  %s %s  limit=%d window=%s\n",
					p.Name, p.Method, p.Endpoint, p.Limit, p.Window)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output resolved policies as JSON")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func init() {
	_ = os.Stderr // ensure os import is used
}
