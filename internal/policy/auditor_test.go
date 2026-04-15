package policy

import (
	"strings"
	"testing"
)

func TestRecord_AddsEntry(t *testing.T) {
	log := NewAuditLog()
	log.Record("CREATE", "api-limit", "created with limit 100")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Action != "CREATE" || e.PolicyName != "api-limit" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	log := NewAuditLog()
	log.Record("CREATE", "p1", "")
	log.Record("UPDATE", "p1", "limit changed")
	log.Record("DELETE", "p2", "removed")
	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}
}

func TestFilter_ByAction(t *testing.T) {
	log := NewAuditLog()
	log.Record("CREATE", "p1", "")
	log.Record("UPDATE", "p1", "")
	log.Record("CREATE", "p2", "")
	results := log.Filter("CREATE", "")
	if len(results) != 2 {
		t.Errorf("expected 2 CREATE entries, got %d", len(results))
	}
}

func TestFilter_ByPolicyName(t *testing.T) {
	log := NewAuditLog()
	log.Record("CREATE", "alpha", "")
	log.Record("UPDATE", "beta", "")
	log.Record("DELETE", "alpha", "")
	results := log.Filter("", "alpha")
	if len(results) != 2 {
		t.Errorf("expected 2 entries for alpha, got %d", len(results))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	log := NewAuditLog()
	log.Record("CREATE", "p1", "")
	results := log.Filter("DELETE", "p1")
	if len(results) != 0 {
		t.Errorf("expected 0 entries, got %d", len(results))
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(nil)
	if out != "No audit entries found." {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatAuditLog_ContainsFields(t *testing.T) {
	log := NewAuditLog()
	log.Record("UPDATE", "my-policy", "window changed to 60s")
	out := FormatAuditLog(log.Entries)
	if !strings.Contains(out, "UPDATE") || !strings.Contains(out, "my-policy") {
		t.Errorf("output missing expected fields: %s", out)
	}
}
