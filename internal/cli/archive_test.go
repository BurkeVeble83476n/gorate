package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/gorate/internal/policy"
)

func writeArchivePolicyFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "policies-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func writeArchiveFile(t *testing.T, dir, label string, policies []policy.Policy) string {
	t.Helper()
	archive := policy.Archive{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Policies:  policies,
	}
	data, _ := json.MarshalIndent(archive, "", "  ")
	path := filepath.Join(dir, "test_archive.json")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestArchiveSaveCmd_MissingFileFlag(t *testing.T) {
	cmd := NewArchiveCmd()
	cmd.SetArgs([]string{"save"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when --file flag is missing")
	}
}

func TestArchiveSaveCmd_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	cmd := NewArchiveCmd()
	cmd.SetArgs([]string{"save", "--file", "/nonexistent/policies.yaml", "--dir", dir})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing policy file")
	}
}

func TestArchiveListCmd_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	cmd := NewArchiveCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list", "--dir", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected output for empty archive list")
	}
}

func TestArchiveLoadCmd_Success(t *testing.T) {
	dir := t.TempDir()
	policies := []policy.Policy{
		{Name: "p1", Endpoint: "/test", Method: "GET", Limit: 10, Window: 60},
	}
	archivePath := writeArchiveFile(t, dir, "testlabel", policies)

	cmd := NewArchiveCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"load", "--path", archivePath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestArchiveLoadCmd_FileNotFound(t *testing.T) {
	cmd := NewArchiveCmd()
	cmd.SetArgs([]string{"load", "--path", "/nonexistent/archive.json"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing archive file")
	}
}
