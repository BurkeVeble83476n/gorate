package policy

import (
	"testing"
)

func makeReorderPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 20, Window: 60},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 30, Window: 60},
		{Name: "delta", Endpoint: "/d", Method: "PUT", Limit: 40, Window: 60},
	}
}

func TestReorder_MoveToMiddle(t *testing.T) {
	policies := makeReorderPolicies()
	result, err := Reorder(policies, "delta", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Name != "delta" {
		t.Errorf("expected delta at index 1, got %s", result[1].Name)
	}
	if len(result) != len(policies) {
		t.Errorf("expected same length %d, got %d", len(policies), len(result))
	}
}

func TestReorder_SameIndex(t *testing.T) {
	policies := makeReorderPolicies()
	result, err := Reorder(policies, "beta", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[1].Name != "beta" {
		t.Errorf("expected beta at index 1, got %s", result[1].Name)
	}
}

func TestReorder_PolicyNotFound(t *testing.T) {
	policies := makeReorderPolicies()
	_, err := Reorder(policies, "missing", 0)
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestReorder_OutOfRange(t *testing.T) {
	policies := makeReorderPolicies()
	_, err := Reorder(policies, "alpha", 10)
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestReorder_EmptyName(t *testing.T) {
	policies := makeReorderPolicies()
	_, err := Reorder(policies, "", 0)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestReorderToFront(t *testing.T) {
	policies := makeReorderPolicies()
	result, err := ReorderToFront(policies, "gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Name != "gamma" {
		t.Errorf("expected gamma at index 0, got %s", result[0].Name)
	}
}

func TestReorderToBack(t *testing.T) {
	policies := makeReorderPolicies()
	result, err := ReorderToBack(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	last := result[len(result)-1]
	if last.Name != "alpha" {
		t.Errorf("expected alpha at last index, got %s", last.Name)
	}
}

func TestReorder_OriginalUnmodified(t *testing.T) {
	policies := makeReorderPolicies()
	originalFirst := policies[0].Name
	_, _ = Reorder(policies, "delta", 0)
	if policies[0].Name != originalFirst {
		t.Errorf("original slice was modified")
	}
}
