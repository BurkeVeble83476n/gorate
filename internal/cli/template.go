package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/yourorg/gorate/internal/policy"
)

// NewTemplateCmd returns a cobra command that applies a policy template
// with caller-supplied variable substitutions and prints the result.
func NewTemplateCmd() *cobra.Command {
	var (
		name     string
		endpoint string
		method   string
		limit    int
		window   int
		varsRaw  []string
		format   string
	)

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Apply a policy template with variable substitutions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name flag is required")
			}
			if endpoint == "" {
				return fmt.Errorf("--endpoint flag is required")
			}

			vars := parseVarFlags(varsRaw)

			tmpl := policy.Template{
				Name:     name,
				Endpoint: endpoint,
				Method:   method,
				Limit:    limit,
				Window:   window,
			}

			p, err := policy.ApplyTemplate(tmpl, vars)
			if err != nil {
				return fmt.Errorf("failed to apply template: %w", err)
			}

			return printTemplateResult(p, format)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Template name (supports {{VAR}} placeholders)")
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "Endpoint pattern (supports {{VAR}} placeholders)")
	cmd.Flags().StringVar(&method, "method", "", "HTTP method")
	cmd.Flags().IntVar(&limit, "limit", 60, "Request limit")
	cmd.Flags().IntVar(&window, "window", 60, "Window in seconds")
	cmd.Flags().StringArrayVar(&varsRaw, "var", nil, "Variable substitution in KEY=VALUE format")
	cmd.Flags().StringVar(&format, "format", "yaml", "Output format: yaml or json")

	return cmd
}

func parseVarFlags(raw []string) map[string]string {
	vars := make(map[string]string)
	for _, entry := range raw {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}

func printTemplateResult(p policy.Policy, format string) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(p)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(p)
	default:
		return fmt.Errorf("unsupported format %q: use yaml or json", format)
	}
}
