package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ExportFormat represents the supported export formats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatYAML ExportFormat = "yaml"
)

// Export writes policies to the given file path in the specified format.
func Export(policies []Policy, path string, format ExportFormat) error {
	var data []byte
	var err error

	switch format {
	case FormatJSON:
		data, err = json.MarshalIndent(policies, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal policies to JSON: %w", err)
		}
	case FormatYAML:
		data, err = yaml.Marshal(policies)
		if err != nil {
			return fmt.Errorf("failed to marshal policies to YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// ParseFormat resolves a format string to an ExportFormat, defaulting to YAML.
func ParseFormat(s string) (ExportFormat, error) {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON, nil
	case "yaml", "yml", "":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be json or yaml", s)
	}
}
