package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/gorate/internal/policy"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFromFile_ValidYAML(t *testing.T) {
	content := `
policies:
  - name: test-policy
    endpoint: /api/test
    method: GET
    limit: 10
    window: 60
`
	path := writeTempFile(t, content)
	policies, err := policy.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	if policies[0].Name != "test-policy" {
		t.Errorf("expected name 'test-policy', got '%s'", policies[0].Name)
	}
}

func TestLoadFromFile_FileNotFound(t *testing.T) {
	_, err := policy.LoadFromFile("/nonexistent/path/policies.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	path := writeTempFile(t, ":::invalid yaml:::")
	_, err := policy.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoadFromFile_InvalidPolicy(t *testing.T) {
	content := `
policies:
  - name: ""
    endpoint: /api/test
    limit: 10
    window: 60
`
	path := writeTempFile(t, content)
	_, err := policy.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for missing name, got nil")
	}
}

func TestLoadFromFile_MultiplePolicies(t *testing.T) {
	content := `
policies:
  - name: policy-one
    endpoint: /api/one
    limit: 5
    window: 30
  - name: policy-two
    endpoint: /api/two
    method: POST
    limit: 20
    window: 60
`
	path := writeTempFile(t, content)
	policies, err := policy.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(policies))
	}
}
