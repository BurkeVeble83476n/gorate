package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGroupPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return p
}

const groupPolicyYAML = `
- name: list-users
  endpoint: /api/users
  method: GET
  limit: 100
  window: 1m
- name: create-user
  endpoint: /api/users
  method: POST
  limit: 20
  window: 1m
- name: list-orders
  endpoint: /api/orders
  method: GET
  limit: 50
  window: 1m
`

func TestGroupCmd_ByMethod(t *testing.T) {
	file := writeGroupPolicyFile(t, groupPolicyYAML)
	cmd := NewGroupCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--file", file, "--by", "method"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "GET") {
		t.Errorf("expected GET group in output, got: %s", out)
	}
	if !strings.Contains(out, "POST") {
		t.Errorf("expected POST group in output, got: %s", out)
	}
}

func TestGroupCmd_FileNotFound(t *testing.T) {
	cmd := NewGroupCmd()
	cmd.SetArgs([]string{"--file", "/nonexistent/path.yaml"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestGroupCmd_InvalidBy(t *testing.T) {
	file := writeGroupPolicyFile(t, groupPolicyYAML)
	cmd := NewGroupCmd()
	cmd.SetArgs([]string{"--file", file, "--by", "badfield"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for invalid group-by field")
	}
}
