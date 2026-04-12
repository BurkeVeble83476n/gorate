package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"gorate/internal/cli"
)

func writePolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestRunCmd_MissingTargetFlag(t *testing.T) {
	cmd := cli.NewRootCmd("test")
	cmd.SetArgs([]string{"run", "--policies", "some.yaml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --target is missing")
	}
}

func TestRunCmd_InvalidPoliciesFile(t *testing.T) {
	cmd := cli.NewRootCmd("test")
	cmd.SetArgs([]string{"run", "--target", "http://localhost:9090", "--policies", "/nonexistent/path.yaml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing policies file")
	}
}

func TestRunCmd_InvalidYAML(t *testing.T) {
	path := writePolicyFile(t, ":::invalid yaml:::")
	cmd := cli.NewRootCmd("test")
	cmd.SetArgs([]string{"run", "--target", "http://localhost:9090", "--policies", path})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
