package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/gorate/internal/policy"
)

func writeDeprecatePolicyFile(t *testing.T, policies []policy.Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := writeFile(path, policies); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestDeprecateSetCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "alpha", Endpoint: "/alpha", Method: "GET", Limit: 100, Window: 60},
	}
	file := writeDeprecatePolicyFile(t, policies)

	cmd := NewDeprecateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"set", "--file", file, "--name", "alpha", "--reason", "old"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alpha") {
		t.Error("expected output to mention policy name")
	}
}

func TestDeprecateSetCmd_MissingNameFlag(t *testing.T) {
	file := writeDeprecatePolicyFile(t, []policy.Policy{
		{Name: "alpha", Endpoint: "/alpha", Method: "GET", Limit: 10, Window: 60},
	})
	cmd := NewDeprecateCmd()
	cmd.SetArgs([]string{"set", "--file", file})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing --name flag")
	}
}

func TestDeprecateSetCmd_FileNotFound(t *testing.T) {
	cmd := NewDeprecateCmd()
	cmd.SetArgs([]string{"set", "--file", "/nonexistent/path.yaml", "--name", "alpha"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDeprecateUnsetCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "beta", Endpoint: "/beta", Method: "POST", Limit: 50, Window: 30,
			Annotations: map[string]string{"deprecated": "true", "deprecated-since": "2024-01-01T00:00:00Z"}},
	}
	file := writeDeprecatePolicyFile(t, policies)

	cmd := NewDeprecateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"unset", "--file", file, "--name", "beta"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "beta") {
		t.Error("expected output to mention policy name")
	}
}

func TestDeprecateListCmd_ShowsDeprecated(t *testing.T) {
	policies := []policy.Policy{
		{Name: "old", Endpoint: "/old", Method: "GET", Limit: 10, Window: 60,
			Annotations: map[string]string{"deprecated": "true", "deprecated-since": "2024-01-01T00:00:00Z", "deprecated-reason": "sunset"}},
		{Name: "new", Endpoint: "/new", Method: "GET", Limit: 100, Window: 60},
	}
	file := writeDeprecatePolicyFile(t, policies)

	cmd := NewDeprecateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list", "--file", file})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "old") {
		t.Error("expected output to contain deprecated policy name")
	}
	if strings.Contains(buf.String(), "new") {
		t.Error("expected output not to contain non-deprecated policy")
	}
}

func TestDeprecateListCmd_NoneDeprecated(t *testing.T) {
	policies := []policy.Policy{
		{Name: "active", Endpoint: "/active", Method: "GET", Limit: 100, Window: 60},
	}
	file := writeDeprecatePolicyFile(t, policies)

	cmd := NewDeprecateCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list", "--file", file})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no deprecated") {
		t.Error("expected 'no deprecated policies' message")
	}
}

func TestDeprecateListCmd_FileNotFound(t *testing.T) {
	cmd := NewDeprecateCmd()
	cmd.SetArgs([]string{"list", "--file", filepath.Join(os.TempDir(), "ghost.yaml")})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
