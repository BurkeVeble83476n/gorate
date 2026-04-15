package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/gorate/internal/policy"
)

// NewWatchdogCmd returns a cobra command that evaluates watchdog rules against a stats snapshot.
func NewWatchdogCmd() *cobra.Command {
	var fileFlag string
	var rulesFlag []string

	cmd := &cobra.Command{
		Use:   "watchdog",
		Short: "Evaluate watchdog alert rules against current rate-limit stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(fileFlag)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}

			rs := policy.NewRequestStats(policies)
			snap := rs.Snapshot()

			rules, err := parseWatchdogRules(rulesFlag)
			if err != nil {
				return fmt.Errorf("parsing rules: %w", err)
			}

			alerts := policy.EvaluateWatchdog(snap, rules)
			fmt.Fprint(os.Stdout, policy.FormatAlerts(alerts))
			return nil
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "path to policy YAML file (required)")
	cmd.Flags().StringArrayVarP(&rulesFlag, "rule", "r", nil, "rule in format name:maxRejections:windowSeconds (repeatable)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

// parseWatchdogRules parses rule strings of the form "name:maxRejections:windowSeconds".
func parseWatchdogRules(raw []string) ([]policy.WatchdogRule, error) {
	var rules []policy.WatchdogRule
	for _, r := range raw {
		var name string
		var maxRej int
		var windowSec int
		_, err := fmt.Sscanf(r, "%s", &name)
		_ = err
		parts := splitN(r, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid rule format %q, expected name:maxRejections:windowSeconds", r)
		}
		name = parts[0]
		maxRej, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid maxRejections in rule %q: %w", r, err)
		}
		windowSec, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid windowSeconds in rule %q: %w", r, err)
		}
		rules = append(rules, policy.WatchdogRule{
			PolicyName:    name,
			MaxRejections: maxRej,
			Window:        time.Duration(windowSec) * time.Second,
		})
	}
	return rules, nil
}

func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := strings.Index(s, sep)
		if idx == -1 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}
