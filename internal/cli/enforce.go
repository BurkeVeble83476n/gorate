package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewEnforceCmd returns a cobra command that evaluates policies against an enforcement mode.
func NewEnforceCmd() *cobra.Command {
	var file string
	var mode string
	var failOnViolation bool

	cmd := &cobra.Command{
		Use:   "enforce",
		Short: "Evaluate policies against an enforcement mode",
		Long: `Evaluate loaded policies using a specified enforcement mode.

Modes:
  strict   — block policies with violations (exits non-zero)
  warn     — report violations but allow all policies
  disable  — skip enforcement entirely`,
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			em, err := policy.ParseEnforcementMode(mode)
			if err != nil {
				return err
			}

			results := policy.Enforce(policies, em)
			fmt.Print(policy.FormatEnforcement(results))

			if failOnViolation {
				for _, r := range results {
					if !r.Enforced {
						fmt.Fprintf(os.Stderr, "enforcement failed: %d polic(ies) blocked\n", countBlocked(results))
						os.Exit(1)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to policy YAML file (required)")
	cmd.Flags().StringVarP(&mode, "mode", "m", "warn", "Enforcement mode: strict, warn, or disable")
	cmd.Flags().BoolVar(&failOnViolation, "fail", false, "Exit with non-zero status if any policy is blocked")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func countBlocked(results []policy.EnforcementResult) int {
	n := 0
	for _, r := range results {
		if !r.Enforced {
			n++
		}
	}
	return n
}
