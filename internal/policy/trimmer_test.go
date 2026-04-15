package policy

import (
	"strings"
	"testing"
)

func makeTrimmerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/a", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/b", Method: "*", Limit: 50, Window: 60},
		{Name: "gamma", Endpoint: "/api/c", Method: "POST", Limit: 500, Window: 60},
		{Name: "delta", Endpoint: "/api/d", Method: "PUT", Limit: 20, Window: 60},
	}
}

func TestTrim_NoOptions_ReturnsAll(t *testing.T) {
	policies := makeTrimmerPolicies()
	kept, results := Trim(policies, TrimOptions{})
	if len(kept) != len(policies) {
		t.Errorf("expected %d kept, got %d", len(policies), len(kept))
	}
	for _, r := range results {
		if r.Removed {
			t.Errorf("expected no removals, but %q was removed", r.Name)
		}
	}
}

func TestTrim_RemoveWildcard(t *testing.T) {
	policies := makeTrimmerPolicies()
	kept, results := Trim(policies, TrimOptions{RemoveWildcard: true})
	if len(kept) != 3 {
		t.Errorf("expected 3 kept, got %d", len(kept))
	}
	for _, r := range results {
		if r.Name == "beta" && !r.Removed {
			t.Error("expected beta to be removed")
		}
		if r.Name == "beta" && !strings.Contains(r.Reason, "wildcard") {
			t.Errorf("unexpected reason: %s", r.Reason)
		}
	}
}

func TestTrim_MaxLimit(t *testing.T) {
	policies := makeTrimmerPolicies()
	kept, _ := Trim(policies, TrimOptions{MaxLimit: 100})
	for _, p := range kept {
		if p.Limit > 100 {
			t.Errorf("policy %q with limit %d should have been removed", p.Name, p.Limit)
		}
	}
	if len(kept) != 3 {
		t.Errorf("expected 3 kept, got %d", len(kept))
	}
}

func TestTrim_RequireTag_RemovesUntagged(t *testing.T) {
	policies := makeTrimmerPolicies()
	policies[0] = AddTag(policies[0], "approved")
	kept, _ := Trim(policies, TrimOptions{RequireTag: "approved"})
	if len(kept) != 1 {
		t.Errorf("expected 1 kept, got %d", len(kept))
	}
	if kept[0].Name != "alpha" {
		t.Errorf("expected alpha to be kept, got %s", kept[0].Name)
	}
}

func TestFormatTrimResults_NoneRemoved(t *testing.T) {
	results := []TrimResult{
		{Name: "alpha", Removed: false},
	}
	out := FormatTrimResults(results)
	if !strings.Contains(out, "no policies trimmed") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatTrimResults_SomeRemoved(t *testing.T) {
	results := []TrimResult{
		{Name: "alpha", Removed: true, Reason: "wildcard method"},
		{Name: "beta", Removed: false},
	}
	out := FormatTrimResults(results)
	if !strings.Contains(out, "trimmed 1") {
		t.Errorf("expected trimmed count in output: %s", out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output: %s", out)
	}
}
