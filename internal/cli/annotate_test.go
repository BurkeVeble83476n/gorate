package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/gorate/internal/policy"
)

func writeAnnotatePolicyFile(t *testing.T, policies []policy.Policy) string {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "policies.yaml")
	data := "policies:\n"
	for _, p := range policies {
		data += "  - name: " + p.Name + "\n"
		data += "    endpoint: " + p.Endpoint + "\n"
		data += "    method: " + p.Method + "\n"
		data += "    limit: 100\n"
		data += "    window: 60\n"
	}
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("writing policy file: %v", err)
	}
	return path
}

func TestAnnotateSetCmd_Success(t *testing.T) {
	path := writeAnnotatePolicyFile(t, []policy.Policy{
		{Name: "api-get", Endpoint: "/api", Method: "GET"},
	})
	cmd := NewAnnotateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"set", "-f", path, "-n", "api-get", "-k", "owner", "-v", "team-a"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "owner") {
		t.Errorf("expected output to mention key 'owner', got: %s", buf.String())
	}
}

func TestAnnotateSetCmd_MissingNameFlag(t *testing.T) {
	path := writeAnnotatePolicyFile(t, []policy.Policy{
		{Name: "api-get", Endpoint: "/api", Method: "GET"},
	})
	cmd := NewAnnotateCmd()
	cmd.SetArgs([]string{"set", "-f", path, "-k", "owner", "-v", "team-a"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing --name flag")
	}
}

func TestAnnotateSetCmd_FileNotFound(t *testing.T) {
	cmd := NewAnnotateCmd()
	cmd.SetArgs([]string{"set", "-f", "/nonexistent.yaml", "-n", "api-get", "-k", "owner", "-v", "team-a"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestAnnotateGetCmd_NoAnnotations(t *testing.T) {
	path := writeAnnotatePolicyFile(t, []policy.Policy{
		{Name: "api-get", Endpoint: "/api", Method: "GET"},
	})
	cmd := NewAnnotateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"get", "-f", path, "-n", "api-get"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No annotations") {
		t.Errorf("expected 'No annotations' message, got: %s", buf.String())
	}
}

func TestAnnotateRemoveCmd_KeyNotFound(t *testing.T) {
	path := writeAnnotatePolicyFile(t, []policy.Policy{
		{Name: "api-get", Endpoint: "/api", Method: "GET"},
	})
	cmd := NewAnnotateCmd()
	cmd.SetArgs([]string{"remove", "-f", path, "-n", "api-get", "-k", "missing"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when removing non-existent annotation key")
	}
}
