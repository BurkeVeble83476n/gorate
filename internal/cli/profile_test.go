package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeProfilePolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return p
}

func runProfileCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "gorate"}
	root.AddCommand(NewProfileCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestProfileCmd_MissingFileFlag(t *testing.T) {
	_, err := runProfileCmd(t, "profile")
	if err == nil {
		t.Error("expected error when --file flag is missing")
	}
}

func TestProfileCmd_FileNotFound(t *testing.T) {
	_, err := runProfileCmd(t, "profile", "--file", "/nonexistent/path.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestProfileCmd_DisplaysResults(t *testing.T) {
	content := `
- name: api-limit
  endpoint: /api/data
  method: GET
  limit: 100
  window: 1m
`
	file := writeProfilePolicyFile(t, content)
	out, err := runProfileCmd(t, "profile", "--file", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "api-limit") {
		t.Errorf("expected output to contain policy name, got: %s", out)
	}
}

func TestProfileCmd_VerboseFlag(t *testing.T) {
	content := `
- name: verbose-policy
  endpoint: /verbose
  method: POST
  limit: 5000
  window: 1s
`
	file := writeProfilePolicyFile(t, content)
	out, err := runProfileCmd(t, "profile", "--file", file, "--verbose")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "verbose-policy") {
		t.Errorf("expected verbose output to contain policy name, got: %s", out)
	}
	if !strings.Contains(out, "Risk") {
		t.Errorf("expected verbose output to contain Risk field, got: %s", out)
	}
}
