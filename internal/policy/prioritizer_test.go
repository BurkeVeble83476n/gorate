package policy

import (
	"testing"
)

func makePrioritizerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 20, Window: 60},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 5, Window: 60},
	}
}

func TestSetPriority_Success(t *testing.T) {
	policies := makePrioritizerPolicies()
	updated, err := SetPriority(policies, "alpha", PriorityHigh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range updated {
		if p.Name == "alpha" {
			if p.Annotations["priority"] != PriorityHigh {
				t.Errorf("expected priority %q, got %q", PriorityHigh, p.Annotations["priority"])
			}
			return
		}
	}
	t.Fatal("policy alpha not found after update")
}

func TestSetPriority_InvalidLevel(t *testing.T) {
	policies := makePrioritizerPolicies()
	_, err := SetPriority(policies, "alpha", "critical")
	if err == nil {
		t.Fatal("expected error for invalid priority level")
	}
}

func TestSetPriority_PolicyNotFound(t *testing.T) {
	policies := makePrioritizerPolicies()
	_, err := SetPriority(policies, "unknown", PriorityLow)
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestSetPriority_EmptyName(t *testing.T) {
	policies := makePrioritizerPolicies()
	_, err := SetPriority(policies, "", PriorityLow)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestGetPriority_ReturnsSet(t *testing.T) {
	policies := makePrioritizerPolicies()
	policies, _ = SetPriority(policies, "beta", PriorityLow)
	lvl, err := GetPriority(policies, "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lvl != PriorityLow {
		t.Errorf("expected %q, got %q", PriorityLow, lvl)
	}
}

func TestGetPriority_DefaultsMedium(t *testing.T) {
	policies := makePrioritizerPolicies()
	lvl, err := GetPriority(policies, "gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lvl != PriorityMedium {
		t.Errorf("expected default %q, got %q", PriorityMedium, lvl)
	}
}

func TestSortByPriority_OrderIsCorrect(t *testing.T) {
	policies := makePrioritizerPolicies()
	policies, _ = SetPriority(policies, "gamma", PriorityHigh)
	policies, _ = SetPriority(policies, "alpha", PriorityLow)
	// beta has no explicit priority -> medium
	sorted := SortByPriority(policies)
	if sorted[0].Name != "gamma" {
		t.Errorf("expected gamma first (high), got %q", sorted[0].Name)
	}
	if sorted[1].Name != "beta" {
		t.Errorf("expected beta second (medium), got %q", sorted[1].Name)
	}
	if sorted[2].Name != "alpha" {
		t.Errorf("expected alpha last (low), got %q", sorted[2].Name)
	}
}

func TestSortByPriority_DoesNotMutateOriginal(t *testing.T) {
	policies := makePrioritizerPolicies()
	policies, _ = SetPriority(policies, "alpha", PriorityHigh)
	originalFirst := policies[0].Name
	_ = SortByPriority(policies)
	if policies[0].Name != originalFirst {
		t.Error("SortByPriority mutated the original slice")
	}
}
