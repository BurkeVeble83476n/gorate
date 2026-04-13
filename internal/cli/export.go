package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/gorate/internal/policy"
)

// NewExportCmd creates the `gorate export` subcommand.
func NewExportCmd() *cobra.Command {
	var (
		policiesFile string
		outputFile   string
		formatStr    string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export rate-limit policies to a file (JSON or YAML)",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(policiesFile)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			format, err := policy.ParseFormat(formatStr)
			if err != nil {
				return err
			}

			if outputFile == "" {
				switch format {
				case policy.FormatJSON:
					outputFile = "policies.json"
				default:
					outputFile = "policies.yaml"
				}
			}

			if err := policy.Export(policies, outputFile, format); err != nil {
				return fmt.Errorf("exporting policies: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Exported %d policies to %s\n", len(policies), outputFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&policiesFile, "policies", "p", "policies.yaml", "Path to the policies file")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: policies.<format>)")
	cmd.Flags().StringVarP(&formatStr, "format", "f", "yaml", "Export format: json or yaml")

	return cmd
}
