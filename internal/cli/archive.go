package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewArchiveCmd returns the root archive command with save and list subcommands.
func NewArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive and restore policy snapshots",
	}

	cmd.AddCommand(newArchiveSaveCmd())
	cmd.AddCommand(newArchiveListCmd())
	cmd.AddCommand(newArchiveLoadCmd())

	return cmd
}

func newArchiveSaveCmd() *cobra.Command {
	var file, dir, label string

	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save current policies to an archive",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("failed to load policies: %w", err)
			}

			path, err := policy.SaveArchive(dir, label, policies)
			if err != nil {
				return fmt.Errorf("failed to save archive: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Archive saved to: %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file to archive (required)")
	cmd.Flags().StringVarP(&dir, "dir", "d", "./archives", "Directory to store archives")
	cmd.Flags().StringVarP(&label, "label", "l", "backup", "Label for the archive")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newArchiveListCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved archives",
		RunE: func(cmd *cobra.Command, args []string) error {
			archives, err := policy.ListArchives(dir)
			if err != nil {
				return fmt.Errorf("failed to list archives: %w", err)
			}

			if len(archives) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No archives found.")
				return nil
			}

			for _, a := range archives {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s — %d policies\n",
					a.CreatedAt.Format("2006-01-02 15:04:05"), a.Label, len(a.Policies))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "./archives", "Directory to read archives from")
	return cmd
}

func newArchiveLoadCmd() *cobra.Command {
	var archivePath string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Inspect policies from an archive file",
		RunE: func(cmd *cobra.Command, args []string) error {
			archive, err := policy.LoadArchive(archivePath)
			if err != nil {
				return fmt.Errorf("failed to load archive: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Archive: %s (created %s)\n",
				archive.Label, archive.CreatedAt.Format("2006-01-02 15:04:05"))
			policy.PrintTable(cmd.OutOrStdout(), archive.Policies)
			return nil
		},
	}

	cmd.Flags().StringVarP(&archivePath, "path", "p", "", "Path to archive file (required)")
	_ = cmd.MarkFlagRequired("path")

	return cmd
}
