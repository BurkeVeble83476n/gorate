package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/gorate/internal/cli"
)

func writeExportPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp policy file: %v", err)
	}
	return path
}

const validExportYAML = `
- name: test-policy
  endpoint: /api/test
  method: GET
  limit: 50
  window: 30
`

func TestExportCmd_JSONFormat(t *testing.T) {
	policyFile := writeExportPolicyFile(t, validExportYAML)
	outFile := filepath.Join(t.TempDir(), "out.json")

	root := cli.NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"export", "--policies", policyFile, "--output", outFile, "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(outFile); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}

	if !strings.Contains(buf.String(), "Exported 1 policies") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestExportCmd_FileNotFound(t *testing.T) {
	root := cli.NewRootCmd()
	root.SetArgs([]string{"export", "--policies", "/nonexistent/path.yaml"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for missing policy file, got nil")
	}
}

func TestExportCmd_InvalidFormat(t *testing.T) {
	policyFile := writeExportPolicyFile(t, validExportYAML)

	root := cli.NewRootCmd()
	root.SetArgs([]string{"export", "--policies", policyFile, "--format", "toml"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}
