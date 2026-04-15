package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewTransformCmd returns the cobra command for the transform subcommand.
func NewTransformCmd() *cobra.Command {
	var (
		file       string
		capLimit   int
		defWindow  string
		upMethod   bool
		outputFile string
	)

	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply bulk transformations to policies in a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			var fns []policy.TransformFunc
			if upMethod {
				fns = append(fns, policy.UppercaseMethod())
			}
			if capLimit > 0 {
				fns = append(fns, policy.CapLimit(capLimit))
			}
			if defWindow != "" {
				fns = append(fns, policy.SetDefaultWindow(defWindow))
			}

			if len(fns) == 0 {
				return fmt.Errorf("no transformations specified; use --uppercase-method, --cap-limit, or --default-window")
			}

			out, results, err := policy.Transform(policies, fns...)
			if err != nil {
				return fmt.Errorf("transform: %w", err)
			}

			for _, r := range results {
				status := "unchanged"
				if r.Changed {
					status = "changed"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %-20s %s (%s)\n", r.Name, status, r.Note)
			}

			dest := file
			if outputFile != "" {
				dest = outputFile
			}
			data, err := policy.Export(out, "yaml")
			if err != nil {
				return fmt.Errorf("exporting: %w", err)
			}
			if err := writeFile(dest, data); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved %s policies to %s\n", strconv.Itoa(len(out)), dest)
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file to transform (required)")
	cmd.Flags().IntVar(&capLimit, "cap-limit", 0, "Cap all policy limits to this value")
	cmd.Flags().StringVar(&defWindow, "default-window", "", "Set window for policies missing one")
	cmd.Flags().BoolVar(&upMethod, "uppercase-method", false, "Uppercase all method fields")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (defaults to input file)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
