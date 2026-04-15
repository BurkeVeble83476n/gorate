package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdoe/gorate/internal/cli"
)

func writePinPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write policy file: %v", err)
	}
	return path
}

const pinPolicyYAML = `
- name: alpha
  endpoint: /api/alpha
  method: GET
  limit: 10
  window: 60
- name: beta
  endpoint: /api/beta
  method: POST
  limit: 5
  window: 30
`

func TestPinSetCmd_Success(t *testing.T) {
	file := writePinPolicyFile(t, pinPolicyYAML)
	root := NewRootCmd()
	root.AddCommand(NewPinCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"pin", "set", "--file", file, "--name", "alpha"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alpha") {
		t.Errorf("expected output to mention policy name, got: %s", buf.String())
	}
}

func TestPinSetCmd_MissingNameFlag(t *testing.T) {
	file := writePinPolicyFile(t, pinPolicyYAML)
	root := NewRootCmd()
	root.AddCommand(NewPinCmd())
	root.SetArgs([]string{"pin", "set", "--file", file})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing --name flag")
	}
}

func TestPinSetCmd_FileNotFound(t *testing.T) {
	root := NewRootCmd()
	root.AddCommand(NewPinCmd())
	root.SetArgs([]string{"pin", "set", "--file", "/nonexistent/path.yaml", "--name", "alpha"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestPinListCmd_NoPinnedPolicies(t *testing.T) {
	file := writePinPolicyFile(t, pinPolicyYAML)
	root := NewRootCmd()
	root.AddCommand(NewPinCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"pin", "list", "--file", file})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No pinned") {
		t.Errorf("expected no pinned message, got: %s", buf.String())
	}
}

func TestPinUnsetCmd_MissingFileFlag(t *testing.T) {
	root := NewRootCmd()
	root.AddCommand(NewPinCmd())
	root.SetArgs([]string{"pin", "unset", "--name", "alpha"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing --file flag")
	}
}
