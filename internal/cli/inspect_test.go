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

func TestInspectCmd_FileNotFound(t *testing.T) {
	cmd := cli.NewRootCmd("test")
	cmd.SetArgs([]string{"inspect", "--policies", "/no/such/file.yaml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
