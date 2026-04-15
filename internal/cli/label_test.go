package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeLabelPolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policies.yaml")
	err := os.WriteFile(p, []byte(content), 0644)
	require.NoError(t, err)
	return p
}

const labelPolicyYAML = `
- name: alpha
  endpoint: /api/alpha
  method: GET
  limit: 10
  window: 60
  labels:
    - core
`

func TestLabelAddCmd_Success(t *testing.T) {
	file := writeLabelPolicyFile(t, labelPolicyYAML)
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"label", "add", "--file", file, "--name", "alpha", "--label", "public"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "public")
	assert.Contains(t, buf.String(), "alpha")
}

func TestLabelAddCmd_MissingNameFlag(t *testing.T) {
	file := writeLabelPolicyFile(t, labelPolicyYAML)
	root := NewRootCmd()
	root.SetArgs([]string{"label", "add", "--file", file, "--label", "public"})
	err := root.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--name")
}

func TestLabelAddCmd_FileNotFound(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"label", "add", "--file", "/nonexistent/path.yaml", "--name", "alpha", "--label", "x"})
	err := root.Execute()
	assert.Error(t, err)
}

func TestLabelGetCmd_ShowsLabels(t *testing.T) {
	file := writeLabelPolicyFile(t, labelPolicyYAML)
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"label", "get", "--file", file, "--name", "alpha"})
	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "core")
}

func TestLabelGetCmd_MissingNameFlag(t *testing.T) {
	file := writeLabelPolicyFile(t, labelPolicyYAML)
	root := NewRootCmd()
	root.SetArgs([]string{"label", "get", "--file", file})
	err := root.Execute()
	assert.Error(t, err)
}
