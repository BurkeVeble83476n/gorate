package policy

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeSnapshotPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 5, Window: 30},
	}
}

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	err := SaveSnapshot(makeSnapshotPolicies(), "test-label", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("snapshot file was not created")
	}
}

func TestLoadSnapshot_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	policies := makeSnapshotPolicies()
	_ = SaveSnapshot(policies, "my-snap", path)

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Label != "my-snap" {
		t.Errorf("expected label 'my-snap', got %q", snap.Label)
	}
	if len(snap.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(snap.Policies))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDescribeSnapshot_Format(t *testing.T) {
	snap := &PolicySnapshot{
		Label:     "release-v1",
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Policies:  makeSnapshotPolicies(),
	}
	desc := DescribeSnapshot(snap)
	if desc == "" {
		t.Fatal("expected non-empty description")
	}
	for _, substr := range []string{"release-v1", "2", "2024-06-01"} {
		if !containsStr(desc, substr) {
			t.Errorf("expected description to contain %q, got: %s", substr, desc)
		}
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
