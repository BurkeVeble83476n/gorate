package cli

import (
	"bytes"
	"strings"
	"testing"

	"gorate/internal/policy"
)

func buildAuditLog() *policy.AuditLog {
	log := policy.NewAuditLog()
	log.Record("CREATE", "api-limit", "initial creation")
	log.Record("UPDATE", "api-limit", "limit bumped to 200")
	log.Record("DELETE", "old-policy", "removed stale policy")
	return log
}

func runAuditCmd(log *policy.AuditLog, args ...string) (string, error) {
	cmd := NewAuditCmd(log)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestAuditCmd_ShowsAllEntries(t *testing.T) {
	log := buildAuditLog()
	out, err := runAuditCmd(log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "CREATE") || !strings.Contains(out, "DELETE") {
		t.Errorf("expected all entries in output, got: %s", out)
	}
}

func TestAuditCmd_FilterByAction(t *testing.T) {
	log := buildAuditLog()
	out, err := runAuditCmd(log, "--action", "UPDATE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "UPDATE") {
		t.Errorf("expected UPDATE in output, got: %s", out)
	}
	if strings.Contains(out, "DELETE") {
		t.Errorf("unexpected DELETE in filtered output")
	}
}

func TestAuditCmd_FilterByName(t *testing.T) {
	log := buildAuditLog()
	out, err := runAuditCmd(log, "--name", "old-policy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "old-policy") {
		t.Errorf("expected old-policy in output, got: %s", out)
	}
}

func TestAuditCmd_EmptyLog(t *testing.T) {
	log := policy.NewAuditLog()
	out, err := runAuditCmd(log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No audit entries found.") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
