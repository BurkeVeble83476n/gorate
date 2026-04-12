package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gorate",
		Short:   "A lightweight CLI tool for rate-limiting HTTP endpoints during local development",
		Version: version,
	}

	cmd.AddCommand(NewRunCmd())
	cmd.AddCommand(NewInspectCmd())

	return cmd
}
