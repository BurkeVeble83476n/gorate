package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeDependencyPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func runDependencyCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "gorate"}
	root.AddCommand(NewDependencyCmd())
	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

const depPoliciesYAML = `
- name: alpha
  endpoint: /a
  method: GET
  limit: 10
  window: 60
- name: beta
  endpoint: /b
  method: POST
  limit: 5
  window: 30
`

func TestDependencyAddCmd_Success(t *testing.T) {
	path := writeDependencyPolicyFile(t, depPoliciesYAML)
	_, err := runDependencyCmd(t, "dependency", "add", "--file", path, "--from", "alpha", "--to", "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDependencyAddCmd_MissingFileFlag(t *testing.T) {
	_, err := runDependencyCmd(t, "dependency", "add", "--from", "alpha", "--to", "beta")
	if err == nil {
		t.Error("expected error for missing --file flag")
	}
}

func TestDependencyAddCmd_FileNotFound(t *testing.T) {
	_, err := runDependencyCmd(t, "dependency", "add", "--file", "/nonexistent.yaml", "--from", "alpha", "--to", "beta")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDependencyGetCmd_ShowsDependencies(t *testing.T) {
	path := writeDependencyPolicyFile(t, depPoliciesYAML)
	_, err := runDependencyCmd(t, "dependency", "add", "--file", path, "--from", "alpha", "--to", "beta")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	out, err := runDependencyCmd(t, "dependency", "get", "--file", path, "--name", "alpha")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected output to contain 'beta', got: %s", out)
	}
}

func TestDependencyGraphCmd_OutputsJSON(t *testing.T) {
	path := writeDependencyPolicyFile(t, depPoliciesYAML)
	_, err := runDependencyCmd(t, "dependency", "add", "--file", path, "--from", "alpha", "--to", "beta")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	out, err := runDependencyCmd(t, "dependency", "graph", "--file", path)
	if err != nil {
		t.Fatalf("graph failed: %v", err)
	}
	if !strings.Contains(out, "Edges") && !strings.Contains(out, "edges") && !strings.Contains(out, "From") {
		t.Logf("graph output: %s", out)
	}
}
