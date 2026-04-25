package policy

import (
	"strings"
	"testing"
)

func TestRecord_AddsEntry(t *testing.T) {
	cl := NewChangelog()
	cl.Record("api-limit", "limit", "100", "200", "alice")
	if len(cl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(cl.Entries))
	}
	e := cl.Entries[0]
	if e.PolicyName != "api-limit" || e.Field != "limit" || e.OldValue != "100" || e.NewValue != "200" || e.Author != "alice" {
		t.Errorf("entry fields mismatch: %+v", e)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	cl := NewChangelog()
	cl.Record("p1", "method", "GET", "POST", "bob")
	cl.Record("p2", "window", "60s", "30s", "carol")
	if len(cl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cl.Entries))
	}
}

func TestFilterByPolicy_ReturnsMatchingEntries(t *testing.T) {
	cl := NewChangelog()
	cl.Record("alpha", "limit", "10", "20", "alice")
	cl.Record("beta", "limit", "5", "15", "bob")
	cl.Record("alpha", "window", "60s", "120s", "alice")

	results := cl.FilterByPolicy("alpha")
	if len(results) != 2 {
		t.Fatalf("expected 2 entries for 'alpha', got %d", len(results))
	}
	for _, e := range results {
		if e.PolicyName != "alpha" {
			t.Errorf("unexpected policy name: %s", e.PolicyName)
		}
	}
}

func TestFilterByPolicy_NoMatch(t *testing.T) {
	cl := NewChangelog()
	cl.Record("alpha", "limit", "10", "20", "alice")
	results := cl.FilterByPolicy("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 entries, got %d", len(results))
	}
}

func TestFilterByField_ReturnsMatchingEntries(t *testing.T) {
	cl := NewChangelog()
	cl.Record("p1", "limit", "10", "20", "alice")
	cl.Record("p2", "method", "GET", "POST", "bob")
	cl.Record("p3", "LIMIT", "5", "50", "carol")

	results := cl.FilterByField("limit")
	if len(results) != 2 {
		t.Fatalf("expected 2 entries for field 'limit', got %d", len(results))
	}
}

func TestFormatChangelog_ContainsFields(t *testing.T) {
	cl := NewChangelog()
	cl.Record("my-policy", "limit", "50", "100", "dev")
	out := FormatChangelog(cl.Entries)
	for _, want := range []string{"my-policy", "limit", "50", "100", "dev"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

func TestFormatChangelog_EmptyEntries(t *testing.T) {
	out := FormatChangelog(nil)
	if out != "no changelog entries" {
		t.Errorf("expected empty message, got: %s", out)
	}
}
