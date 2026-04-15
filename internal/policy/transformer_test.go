package policy

import (
	"testing"
)

func makeTransformerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "get", Limit: 50, Window: "1m"},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 200, Window: ""},
		{Name: "gamma", Endpoint: "/c", Method: "put", Limit: 5, Window: "30s"},
	}
}

func TestTransform_UppercaseMethod(t *testing.T) {
	policies := makeTransformerPolicies()
	out, results, err := Transform(policies, UppercaseMethod())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Method != "GET" {
		t.Errorf("expected GET, got %s", out[0].Method)
	}
	if out[2].Method != "PUT" {
		t.Errorf("expected PUT, got %s", out[2].Method)
	}
	if !results[0].Changed {
		t.Error("expected alpha to be marked changed")
	}
	if results[1].Changed {
		t.Error("expected beta to be unchanged")
	}
}

func TestTransform_CapLimit(t *testing.T) {
	policies := makeTransformerPolicies()
	out, _, err := Transform(policies, CapLimit(100))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Limit != 100 {
		t.Errorf("expected beta limit capped to 100, got %d", out[1].Limit)
	}
	if out[0].Limit != 50 {
		t.Errorf("expected alpha limit unchanged at 50, got %d", out[0].Limit)
	}
}

func TestTransform_CapLimit_InvalidMax(t *testing.T) {
	policies := makeTransformerPolicies()
	_, _, err := Transform(policies, CapLimit(0))
	if err == nil {
		t.Error("expected error for cap limit of 0")
	}
}

func TestTransform_SetDefaultWindow(t *testing.T) {
	policies := makeTransformerPolicies()
	out, _, err := Transform(policies, SetDefaultWindow("60s"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Window != "60s" {
		t.Errorf("expected beta window set to 60s, got %s", out[1].Window)
	}
	if out[0].Window != "1m" {
		t.Errorf("expected alpha window unchanged, got %s", out[0].Window)
	}
}

func TestTransform_MultipleTransforms(t *testing.T) {
	policies := makeTransformerPolicies()
	out, results, err := Transform(policies, UppercaseMethod(), CapLimit(10), SetDefaultWindow("2m"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Method != "GET" {
		t.Errorf("expected GET, got %s", out[0].Method)
	}
	if out[1].Limit != 10 {
		t.Errorf("expected beta capped at 10, got %d", out[1].Limit)
	}
	if out[1].Window != "2m" {
		t.Errorf("expected beta window set to 2m, got %s", out[1].Window)
	}
	for _, r := range results {
		if !r.Changed {
			t.Errorf("expected policy %s to be changed", r.Name)
		}
	}
}

func TestTransform_NoteContainsFields(t *testing.T) {
	policies := []Policy{
		{Name: "p1", Endpoint: "/x", Method: "get", Limit: 500, Window: "1m"},
	}
	_, results, err := Transform(policies, UppercaseMethod(), CapLimit(100))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Note == "no change" {
		t.Error("expected a non-empty change note")
	}
}
