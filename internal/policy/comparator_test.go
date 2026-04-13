package policy

import (
	"strings"
	"testing"
)

func makeComparePolicies() ([]Policy, []Policy) {
	a := []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 100, Window: "1m"},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 50, Window: "30s"},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 10, Window: "10s"},
	}
	b := []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 200, Window: "1m"},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 50, Window: "30s"},
		{Name: "delta", Endpoint: "/d", Method: "DELETE", Limit: 5, Window: "5s"},
	}
	return a, b
}

func TestCompare_OnlyInA(t *testing.T) {
	a, b := makeComparePolicies()
	result := Compare(a, b)
	if len(result.OnlyInA) != 1 || result.OnlyInA[0].Name != "gamma" {
		t.Errorf("expected gamma only in A, got %v", result.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a, b := makeComparePolicies()
	result := Compare(a, b)
	if len(result.OnlyInB) != 1 || result.OnlyInB[0].Name != "delta" {
		t.Errorf("expected delta only in B, got %v", result.OnlyInB)
	}
}

func TestCompare_InBoth(t *testing.T) {
	a, b := makeComparePolicies()
	result := Compare(a, b)
	if len(result.InBoth) != 2 {
		t.Errorf("expected 2 in both, got %d", len(result.InBoth))
	}
}

func TestCompare_ConflictDetected(t *testing.T) {
	a, b := makeComparePolicies()
	result := Compare(a, b)
	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	c := result.Conflicts[0]
	if c.Name != "alpha" || c.Field != "limit" {
		t.Errorf("unexpected conflict: %+v", c)
	}
}

func TestCompare_NoConflicts(t *testing.T) {
	a := []Policy{
		{Name: "x", Endpoint: "/x", Method: "GET", Limit: 10, Window: "1m"},
	}
	b := []Policy{
		{Name: "x", Endpoint: "/x", Method: "GET", Limit: 10, Window: "1m"},
	}
	result := Compare(a, b)
	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", result.Conflicts)
	}
}

func TestFormatCompare_NoDifferences(t *testing.T) {
	result := CompareResult{}
	out := FormatCompare(result)
	if !strings.Contains(out, "no differences") {
		t.Errorf("expected no differences message, got: %s", out)
	}
}

func TestFormatCompare_ContainsConflict(t *testing.T) {
	a, b := makeComparePolicies()
	result := Compare(a, b)
	out := FormatCompare(result)
	if !strings.Contains(out, "conflict") {
		t.Errorf("expected conflict in output, got: %s", out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output, got: %s", out)
	}
}
