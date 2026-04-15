package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTransformPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write policy file: %v", err)
	}
	return path
}

func runTransformCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "gorate"}
	root.AddCommand(NewTransformCmd())
	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"transform"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestTransformCmd_MissingFileFlag(t *testing.T) {
	_, err := runTransformCmd(t, "--uppercase-method")
	if err == nil {
		t.Error("expected error when --file flag is missing")
	}
}

func TestTransformCmd_FileNotFound(t *testing.T) {
	_, err := runTransformCmd(t, "--file", "/nonexistent/path.yaml", "--uppercase-method")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTransformCmd_NoTransformFlags(t *testing.T) {
	content := `
- name: alpha
  endpoint: /a
  method: get
  limit: 10
  window: 1m
`
	path := writeTransformPolicyFile(t, content)
	_, err := runTransformCmd(t, "--file", path)
	if err == nil {
		t.Error("expected error when no transform flags provided")
	}
}

func TestTransformCmd_UppercaseMethod(t *testing.T) {
	content := `
- name: alpha
  endpoint: /a
  method: get
  limit: 10
  window: 1m
`
	path := writeTransformPolicyFile(t, content)
	out, err := runTransformCmd(t, "--file", path, "--uppercase-method")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected output to mention alpha, got: %s", out)
	}
}

func TestTransformCmd_CapLimit(t *testing.T) {
	content := `
- name: beta
  endpoint: /b
  method: POST
  limit: 500
  window: 1m
`
	path := writeTransformPolicyFile(t, content)
	out, err := runTransformCmd(t, "--file", path, "--cap-limit", "100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected output to mention beta, got: %s", out)
	}
}
