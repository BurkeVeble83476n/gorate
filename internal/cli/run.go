package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
	"gorate/internal/proxy"
)

func NewRunCmd() *cobra.Command {
	var (
		policiesFile string
		targetURL    string
		addr         string
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the rate-limiting proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(policiesFile)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			srv, err := proxy.NewServer(targetURL, addr, policies)
			if err != nil {
				return fmt.Errorf("creating server: %w", err)
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			log.Printf("gorate proxy listening on %s -> %s", addr, targetURL)

			if err := srv.Start(ctx); err != nil {
				return fmt.Errorf("server error: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&policiesFile, "policies", "p", "policies.yaml", "Path to policies YAML file")
	cmd.Flags().StringVarP(&targetURL, "target", "t", "", "Target upstream URL (required)")
	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Address to listen on")
	_ = cmd.MarkFlagRequired("target")

	return cmd
}
