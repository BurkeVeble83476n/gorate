package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/gorate/internal/policy"
)

func writeReviewPolicyFile(t *testing.T, policies []policy.Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.yaml")
	if err := writeFile(path, policies); err != nil {
		t.Fatalf("writing policy file: %v", err)
	}
	return path
}

func runReviewCmd(args ...string) (string, error) {
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestReviewSetCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "api-read", Endpoint: "/api", Method: "GET", Limit: 100, Window: 60},
	}
	path := writeReviewPolicyFile(t, policies)
	_, err := runReviewCmd("review", "set", "-f", path, "-n", "api-read", "-s", "approved", "--reviewer", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated, _ := policy.LoadFromFile(path)
	entry, _ := policy.GetReview(updated, "api-read")
	if entry.Status != policy.ReviewApproved {
		t.Errorf("expected approved, got %s", entry.Status)
	}
	if entry.Reviewer != "alice" {
		t.Errorf("expected reviewer alice, got %s", entry.Reviewer)
	}
}

func TestReviewSetCmd_MissingNameFlag(t *testing.T) {
	path := writeReviewPolicyFile(t, []policy.Policy{})
	_, err := runReviewCmd("review", "set", "-f", path, "-s", "approved")
	if err == nil {
		t.Fatal("expected error for missing --name flag")
	}
}

func TestReviewSetCmd_FileNotFound(t *testing.T) {
	_, err := runReviewCmd("review", "set", "-f", "/no/such/file.yaml", "-n", "x", "-s", "approved")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReviewGetCmd_ShowsStatus(t *testing.T) {
	policies := []policy.Policy{
		{Name: "svc", Endpoint: "/svc", Method: "*", Limit: 50, Window: 30,
			Annotations: map[string]string{"review/status": "rejected", "review/reviewer": "bob"}},
	}
	path := writeReviewPolicyFile(t, policies)
	out, err := runReviewCmd("review", "get", "-f", path, "-n", "svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "rejected") {
		t.Errorf("expected 'rejected' in output, got: %s", out)
	}
	if !strings.Contains(out, "bob") {
		t.Errorf("expected reviewer 'bob' in output, got: %s", out)
	}
}

func TestReviewListCmd_FiltersByStatus(t *testing.T) {
	policies := []policy.Policy{
		{Name: "p1", Endpoint: "/p1", Method: "GET", Limit: 10, Window: 60,
			Annotations: map[string]string{"review/status": "approved"}},
		{Name: "p2", Endpoint: "/p2", Method: "GET", Limit: 20, Window: 60},
	}
	path := writeReviewPolicyFile(t, policies)
	out, err := runReviewCmd("review", "list", "-f", path, "-s", "approved")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "p1") {
		t.Errorf("expected p1 in output, got: %s", out)
	}
	if strings.Contains(out, "p2") {
		t.Errorf("did not expect p2 in approved list, got: %s", out)
	}
	_ = os.Remove(path)
}
