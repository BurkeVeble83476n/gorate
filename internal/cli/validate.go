package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewValidateCmd returns a cobra command that validates a policies file
// and reports any errors found without starting a proxy server.
func NewValidateCmd() *cobra.Command {
	var policiesFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a rate-limit policies file",
		Long:  "Parse and validate a YAML policies file, reporting any structural or semantic errors.",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(policiesFile)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			var hasErrors bool
			for _, p := range policies {
				var errs []string

				if endpointErr := policy.ValidateEndpoint(p.Endpoint); endpointErr != nil {
					errs = append(errs, endpointErr.Error())
				}
				if methodErr := policy.ValidateMethod(p.Method); methodErr != nil {
					errs = append(errs, methodErr.Error())
				}
				errs = append(errs, policy.ValidateLimit(p.Limit, p.WindowSeconds)...)

				if len(errs) > 0 {
					ve := &policy.ValidationError{Name: p.Name, Errors: errs}
					fmt.Fprintln(cmd.ErrOrStderr(), ve.Error())
					hasErrors = true
				}
			}

			if hasErrors {
				return fmt.Errorf("one or more policies failed validation")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "All %d policy(ies) are valid.\n", len(policies))
			return nil
		},
	}

	cmd.Flags().StringVarP(&policiesFile, "policies", "p", "", "Path to the YAML policies file (required)")
	_ = cmd.MarkFlagRequired("policies")

	return cmd
}
