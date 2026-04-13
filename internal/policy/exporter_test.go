package policy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/yourorg/gorate/internal/policy"
)

func makeExportPolicies() []policy.Policy {
	return []policy.Policy{
		{Name: "api-read", Endpoint: "/api/read", Method: "GET", Limit: 100, Window: 60},
		{Name: "api-write", Endpoint: "/api/write", Method: "POST", Limit: 20, Window: 60},
	}
}

func TestExport_JSON(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "policies.json")

	policies := makeExportPolicies()
	if err := policy.Export(policies, out, policy.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	var result []policy.Policy
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 policies, got %d", len(result))
	}
}

func TestExport_YAML(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "policies.yaml")

	policies := makeExportPolicies()
	if err := policy.Export(policies, out, policy.FormatYAML); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	var result []policy.Policy
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("output is not valid YAML: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 policies, got %d", len(result))
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "policies.csv")

	if err := policy.Export(makeExportPolicies(), out, "csv"); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected policy.ExportFormat
	}{
		{"json", policy.FormatJSON},
		{"JSON", policy.FormatJSON},
		{"yaml", policy.FormatYAML},
		{"yml", policy.FormatYAML},
		{"", policy.FormatYAML},
	}
	for _, tc := range cases {
		f, err := policy.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, f, tc.expected)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := policy.ParseFormat("toml"); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
