package policy

import (
	"testing"
)

func makePatcherPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: "1m"},
		{Name: "beta", Endpoint: "/api/v2", Method: "POST", Limit: 50, Window: "30s"},
	}
}

func TestPatch_UpdateLimit(t *testing.T) {
	policies := makePatcherPolicies()
	newLimit := 200
	result, err := Patch(policies, "alpha", PatchOptions{Limit: &newLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Limit != 200 {
		t.Errorf("expected limit 200, got %d", result[0].Limit)
	}
	if result[0].Window != "1m" {
		t.Errorf("window should be unchanged, got %s", result[0].Window)
	}
}

func TestPatch_UpdateMethod(t *testing.T) {
	policies := makePatcherPolicies()
	newMethod := "POST"
	result, err := Patch(policies, "alpha", PatchOptions{Method: &newMethod})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Method != "POST" {
		t.Errorf("expected method POST, got %s", result[0].Method)
	}
}

func TestPatch_UpdateWindow(t *testing.T) {
	policies := makePatcherPolicies()
	newWindow := "5m"
	result, err := Patch(policies, "beta", PatchOptions{Window: &newWindow})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Window != "5m" {
		t.Errorf("expected window 5m, got %s", result[1].Window)
	}
}

func TestPatch_PolicyNotFound(t *testing.T) {
	policies := makePatcherPolicies()
	newLimit := 10
	_, err := Patch(policies, "nonexistent", PatchOptions{Limit: &newLimit})
	if err == nil {
		t.Fatal("expected error for missing policy, got nil")
	}
}

func TestPatch_EmptyName(t *testing.T) {
	policies := makePatcherPolicies()
	newLimit := 10
	_, err := Patch(policies, "", PatchOptions{Limit: &newLimit})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestPatch_InvalidResultFails(t *testing.T) {
	policies := makePatcherPolicies()
	badLimit := -1
	_, err := Patch(policies, "alpha", PatchOptions{Limit: &badLimit})
	if err == nil {
		t.Fatal("expected validation error for negative limit, got nil")
	}
}

func TestPatch_DoesNotMutateOriginal(t *testing.T) {
	policies := makePatcherPolicies()
	newLimit := 999
	_, err := Patch(policies, "alpha", PatchOptions{Limit: &newLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policies[0].Limit != 100 {
		t.Errorf("original slice was mutated, expected 100 got %d", policies[0].Limit)
	}
}
