package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"gorate/internal/cli"
)

const validPoliciesYAML = `
- name: test-policy
  endpoint: /api/v1/resource
  method: GET
  limit: 10
  window: 1m
`

// writePolicyFile is defined in a shared test helper file.

// TestInspectCmd_DisplaysPolicies verifies that the inspect command prints
// policy name and endpoint when given a valid policies file.
func TestInspectCmd_DisplaysPolicies(t *testing.T) {
	path := writePolicyFile(t, validPoliciesYAML)

	buf := &bytes.Buffer{}
	cmd := cli.NewRootCmd("test")
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"inspect", "--policies", path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test-policy") {
		t.Errorf("expected output to contain policy name, got: %s", output)
	}
	if !strings.Contains(output, "/api/v1/resource") {
		t.Errorf("expected output to contain endpoint, got: %s", output)
	}
}

// TestInspectCmd_DisplaysMethod verifies that the inspect command includes
// the HTTP method in its output.
func TestInspectCmd_DisplaysMethod(t *testing.T) {
	path := writePolicyFile(t, validPoliciesYAML)

	buf := &bytes.Buffer{}
	cmd := cli.NewRootCmd("test")
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"inspect", "--policies", path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "GET") {
		t.Errorf("expected output to contain HTTP method GET, got: %s", buf.String())
	}
}

// TestInspectCmd_EmptyPolicies verifies that the inspect command handles an
// empty policy list without returning an error.
func TestInspectCmd_EmptyPolicies(t *testing.T) {
	path := writePolicyFile(t, "[]")

	buf := &bytes.Buffer{}
	cmd := cli.NewRootCmd("test")
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"inspect", "--policies", path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestInspectCmd_FileNotFound verifies that the inspect command returns an
// error when the specified policies file does not exist.
func TestInspectCmd_FileNotFound(t *testing.T) {
	cmd := cli.NewRootCmd("test")
	cmd.SetArgs([]string{"inspect", "--policies", "/no/such/file.yaml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
