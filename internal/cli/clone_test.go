package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/gorate/internal/policy"
)

func writeClonePolicyFile(t *testing.T, policies []policy.Policy) string {
	t.Helper()
	data, err := policy.Export(policies, "yaml")
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "policies-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCloneCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
	}
	file := writeClonePolicyFile(t, policies)
	out := filepath.Join(t.TempDir(), "out.yaml")

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"clone", "alpha", "--file", file, "--name", "alpha-copy", "--output", out})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alpha-copy") {
		t.Errorf("expected output to mention new policy name")
	}
}

func TestCloneCmd_MissingNameFlag(t *testing.T) {
	policies := []policy.Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
	}
	file := writeClonePolicyFile(t, policies)

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"clone", "alpha", "--file", file})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --name flag is missing")
	}
}

func TestCloneCmd_FileNotFound(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"clone", "alpha", "--file", "/no/such/file.yaml", "--name", "copy"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
