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

func writeSnapshotPolicyFile(t *testing.T, dir string) string {
	t.Helper()
	content := `policies:
  - name: snap-policy
    endpoint: /api/snap
    method: GET
    limit: 20
    window: 60
`
	path := filepath.Join(dir, "policies.yaml")
	_ = os.WriteFile(path, []byte(content), 0644)
	return path
}

func writeSnapshotFile(t *testing.T, dir string) string {
	t.Helper()
	snap := policy.PolicySnapshot{
		Timestamp: time.Now().UTC(),
		Label:     "test-snap",
		Policies: []policy.Policy{
			{Name: "snap-policy", Endpoint: "/api/snap", Method: "GET", Limit: 20, Window: 60},
		},
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	path := filepath.Join(dir, "snap.json")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestSnapshotSaveCmd_Success(t *testing.T) {
	dir := t.TempDir()
	policyFile := writeSnapshotPolicyFile(t, dir)
	outFile := filepath.Join(dir, "out.json")

	root := NewRootCmd()
	root.AddCommand(NewSnapshotCmd())
	root.SetArgs([]string{"snapshot", "save", "-f", policyFile, "-o", outFile, "-l", "v1"})

	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to be created")
	}
}

func TestSnapshotSaveCmd_MissingFileFlag(t *testing.T) {
	root := NewRootCmd()
	root.AddCommand(NewSnapshotCmd())
	root.SetArgs([]string{"snapshot", "save"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error when --file flag is missing")
	}
}

func TestSnapshotLoadCmd_Success(t *testing.T) {
	dir := t.TempDir()
	snapFile := writeSnapshotFile(t, dir)

	root := NewRootCmd()
	root.AddCommand(NewSnapshotCmd())
	root.SetArgs([]string{"snapshot", "load", "-s", snapFile})

	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected output from snapshot load")
	}
}

func TestSnapshotLoadCmd_FileNotFound(t *testing.T) {
	root := NewRootCmd()
	root.AddCommand(NewSnapshotCmd())
	root.SetArgs([]string{"snapshot", "load", "-s", "/nonexistent/snap.json"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot file")
	}
}
