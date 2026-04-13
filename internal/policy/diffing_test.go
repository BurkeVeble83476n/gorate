package policy

import (
	"strings"
	"testing"
)

func makeDiffPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: "1m"},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 5, Window: "30s"},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 20, Window: "2m"},
	}
}

func TestDiff_NoChanges(t *testing.T) {
	base := makeDiffPolicies()
	updated := makeDiffPolicies()
	result := Diff(base, updated)
	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("expected no diff, got added=%d removed=%d changed=%d",
			len(result.Added), len(result.Removed), len(result.Changed))
	}
}

func TestDiff_Added(t *testing.T) {
	base := makeDiffPolicies()
	updated := append(makeDiffPolicies(), Policy{Name: "delta", Endpoint: "/d", Method: "GET", Limit: 1, Window: "1s"})
	result := Diff(base, updated)
	if len(result.Added) != 1 || result.Added[0].Name != "delta" {
		t.Errorf("expected 1 added policy named delta, got %+v", result.Added)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := makeDiffPolicies()
	updated := makeDiffPolicies()[:2]
	result := Diff(base, updated)
	if len(result.Removed) != 1 || result.Removed[0].Name != "gamma" {
		t.Errorf("expected 1 removed policy named gamma, got %+v", result.Removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	base := makeDiffPolicies()
	updated := makeDiffPolicies()
	updated[0].Limit = 99
	result := Diff(base, updated)
	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed policy, got %d", len(result.Changed))
	}
	if result.Changed[0].Before.Limit != 10 || result.Changed[0].After.Limit != 99 {
		t.Errorf("unexpected change values: %+v", result.Changed[0])
	}
}

func TestFormatDiff_NoDiff(t *testing.T) {
	result := Diff(makeDiffPolicies(), makeDiffPolicies())
	out := FormatDiff(result)
	if out != "No differences found." {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestFormatDiff_ContainsSymbols(t *testing.T) {
	base := makeDiffPolicies()
	updated := makeDiffPolicies()[1:]
	updated = append(updated, Policy{Name: "new", Endpoint: "/n", Method: "DELETE", Limit: 3, Window: "10s"})
	updated[0].Limit = 999
	result := Diff(base, updated)
	out := FormatDiff(result)
	if !strings.Contains(out, "+ [added]") {
		t.Error("expected added marker in output")
	}
	if !strings.Contains(out, "- [removed]") {
		t.Error("expected removed marker in output")
	}
	if !strings.Contains(out, "~ [changed]") {
		t.Error("expected changed marker in output")
	}
}
