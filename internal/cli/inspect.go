package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

func NewInspectCmd() *cobra.Command {
	var policiesFile string

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect and display loaded rate-limit policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(policiesFile)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			if len(policies) == 0 {
				fmt.Println("No policies found.")
				return nil
			}

			fmt.Printf("%-20s %-30s %-10s %-10s %s\n", "NAME", "ENDPOINT", "METHOD", "LIMIT", "WINDOW")
			fmt.Println("--------------------------------------------------------------------------------")
			for _, p := range policies {
				fmt.Printf("%-20s %-30s %-10s %-10d %s\n",
					p.Name,
					p.Endpoint,
					p.Method,
					p.Limit,
					p.Window,
				)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&policiesFile, "policies", "p", "policies.yaml", "Path to policies YAML file")

	return cmd
}
