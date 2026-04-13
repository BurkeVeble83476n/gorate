package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeComparePolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write policy file: %v", err)
	}
	return path
}

const comparePoliciesA = `
policies:
  - name: alpha
    endpoint: /api/users
    method: GET
    limit: 100
    window: 1m
  - name: beta
    endpoint: /api/orders
    method: POST
    limit: 50
    window: 30s
`

const comparePoliciesB = `
policies:
  - name: alpha
    endpoint: /api/users
    method: GET
    limit: 200
    window: 1m
  - name: gamma
    endpoint: /api/items
    method: GET
    limit: 10
    window: 10s
`

func runCompareCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "gorate"}
	root.AddCommand(NewCompareCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestCompareCmd_MissingFileAFlag(t *testing.T) {
	_, err := runCompareCmd(t, "compare", "--file-b", "some.yaml")
	if err == nil {
		t.Error("expected error for missing --file-a flag")
	}
}

func TestCompareCmd_FileNotFound(t *testing.T) {
	_, err := runCompareCmd(t, "compare", "--file-a", "missing.yaml", "--file-b", "also_missing.yaml")
	if err == nil {
		t.Error("expected error for missing files")
	}
}

func TestCompareCmd_DetectsConflicts(t *testing.T) {
	fileA := writeComparePolicyFile(t, comparePoliciesA)
	fileB := writeComparePolicyFile(t, comparePoliciesB)
	out, _ := runCompareCmd(t, "compare", "--file-a", fileA, "--file-b", fileB)
	if !strings.Contains(out, "conflict") && !strings.Contains(out, "alpha") {
		t.Errorf("expected conflict output, got: %s", out)
	}
}
