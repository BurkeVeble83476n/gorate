package policy

import (
	"testing"
)

func makeDupPolicies() []Policy {
	return []Policy{
		{Name: "a", Endpoint: "/api/users", Method: "GET", Limit: 10, Window: 60},
		{Name: "b", Endpoint: "/api/orders", Method: "POST", Limit: 5, Window: 60},
		{Name: "c", Endpoint: "/api/users", Method: "GET", Limit: 20, Window: 30},
		{Name: "d", Endpoint: "/api/orders", Method: "POST", Limit: 8, Window: 60},
		{Name: "e", Endpoint: "/api/health", Method: "GET", Limit: 100, Window: 60},
	}
}

func TestFindDuplicates_NoDuplicates(t *testing.T) {
	policies := []Policy{
		{Name: "a", Endpoint: "/api/users", Method: "GET", Limit: 10, Window: 60},
		{Name: "b", Endpoint: "/api/orders", Method: "POST", Limit: 5, Window: 60},
	}
	errs := FindDuplicates(policies)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d", len(errs))
	}
}

func TestFindDuplicates_WithDuplicates(t *testing.T) {
	policies := makeDupPolicies()
	errs := FindDuplicates(policies)
	if len(errs) != 2 {
		t.Errorf("expected 2 duplicate errors, got %d", len(errs))
	}
}

func TestFindDuplicates_ErrorMessage(t *testing.T) {
	policies := []Policy{
		{Name: "first", Endpoint: "/api/x", Method: "GET", Limit: 10, Window: 60},
		{Name: "second", Endpoint: "/api/x", Method: "GET", Limit: 20, Window: 60},
	}
	errs := FindDuplicates(policies)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	de, ok := errs[0].(*DuplicateError)
	if !ok {
		t.Fatal("expected *DuplicateError")
	}
	if de.Name != "second" {
		t.Errorf("expected duplicate name 'second', got %q", de.Name)
	}
}

func TestFindDuplicates_WildcardMethod(t *testing.T) {
	policies := []Policy{
		{Name: "a", Endpoint: "/api/x", Method: "", Limit: 10, Window: 60},
		{Name: "b", Endpoint: "/api/x", Method: "*", Limit: 20, Window: 60},
	}
	errs := FindDuplicates(policies)
	if len(errs) != 1 {
		t.Errorf("expected 1 error for wildcard collision, got %d", len(errs))
	}
}

func TestDeduplicatePolicies_KeepsFirst(t *testing.T) {
	policies := makeDupPolicies()
	result := DeduplicatePolicies(policies)
	if len(result) != 3 {
		t.Errorf("expected 3 unique policies, got %d", len(result))
	}
	if result[0].Name != "a" {
		t.Errorf("expected first policy 'a', got %q", result[0].Name)
	}
}

func TestDeduplicatePolicies_Empty(t *testing.T) {
	result := DeduplicatePolicies([]Policy{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
