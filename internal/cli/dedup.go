package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewDedupCmd returns a cobra command that checks for and optionally removes
// duplicate endpoint+method policies from a policy file.
func NewDedupCmd() *cobra.Command {
	var fix bool

	cmd := &cobra.Command{
		Use:   "dedup --file <policies.yaml>",
		Short: "Find or remove duplicate endpoint+method policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("--file flag is required")
			}

			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			duplicates := policy.FindDuplicates(policies)
			if len(duplicates) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No duplicate policies found.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Found %d duplicate(s):\n", len(duplicates))
			for _, e := range duplicates {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", e.Error())
			}

			if !fix {
				fmt.Fprintln(cmd.OutOrStdout(), "\nRun with --fix to remove duplicates (keeps first occurrence).")
				return nil
			}

			deduped := policy.DeduplicatePolicies(policies)
			out, err := policy.Export(deduped, "yaml")
			if err != nil {
				return fmt.Errorf("failed to export deduplicated policies: %w", err)
			}

			if err := os.WriteFile(file, []byte(out), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Removed %d duplicate(s). File updated: %s\n",
				len(duplicates), file)
			return nil
		},
	}

	cmd.Flags().String("file", "", "Path to the policy YAML file")
	cmd.Flags().BoolVar(&fix, "fix", false, "Remove duplicates and overwrite the file")
	return cmd
}
