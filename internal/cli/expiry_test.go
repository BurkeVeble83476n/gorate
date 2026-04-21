package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/gorate/internal/cli"
	"github.com/yourusername/gorate/internal/policy"
)

func writeExpiryPolicyFile(t *testing.T, policies []policy.Policy) string {
	t.Helper()
	data, _ := json.Marshal(policies)
	f, err := os.CreateTemp("", "expiry-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), data, 0644)
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestExpirySetCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60},
	}
	file := writeExpiryPolicyFile(t, policies)
	future := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	root := cli.NewRootCmd()
	root.AddCommand(cli.NewExpiryCmd())
	root.SetArgs([]string{"expiry", "set", "-f", file, "-n", "alpha", "--at", future})
	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "expiry set on") {
		t.Errorf("expected confirmation message, got: %s", buf.String())
	}
}

func TestExpirySetCmd_MissingNameFlag(t *testing.T) {
	root := cli.NewRootCmd()
	root.AddCommand(cli.NewExpiryCmd())
	root.SetArgs([]string{"expiry", "set", "-f", "some.json", "--at", "2099-01-01T00:00:00Z"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing name flag")
	}
}

func TestExpirySetCmd_FileNotFound(t *testing.T) {
	root := cli.NewRootCmd()
	root.AddCommand(cli.NewExpiryCmd())
	root.SetArgs([]string{"expiry", "set", "-f", "nonexistent.json", "-n", "alpha", "--at", "2099-01-01T00:00:00Z"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExpiryRemoveCmd_Success(t *testing.T) {
	policies := []policy.Policy{
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 20, Window: 30},
	}
	file := writeExpiryPolicyFile(t, policies)

	root := cli.NewRootCmd()
	root.AddCommand(cli.NewExpiryCmd())
	root.SetArgs([]string{"expiry", "remove", "-f", file, "-n", "beta"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "expiry removed") {
		t.Errorf("expected removal message, got: %s", buf.String())
	}
}

func TestExpiryStatusCmd_ShowsExpired(t *testing.T) {
	policies := []policy.Policy{
		{Name: "gamma", Endpoint: "/api/gamma", Method: "*", Limit: 5, Window: 120,
			Annotations: map[string]string{"expiry": time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)}},
	}
	file := writeExpiryPolicyFile(t, policies)

	root := cli.NewRootCmd()
	root.AddCommand(cli.NewExpiryCmd())
	root.SetArgs([]string{"expiry", "status", "-f", file})
	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "EXPIRED") {
		t.Errorf("expected EXPIRED in output, got: %s", buf.String())
	}
}
