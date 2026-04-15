package policy

import (
	"testing"
)

func makePinnerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 5, Window: 30},
		{Name: "gamma", Endpoint: "/api/gamma", Method: "*", Limit: 100, Window: 120,
			Annotations: map[string]string{"owner": "team-a"}},
	}
}

func TestPin_Success(t *testing.T) {
	policies := makePinnerPolicies()
	result, err := Pin(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsPinned(result, "alpha") {
		t.Error("expected alpha to be pinned")
	}
}

func TestPin_PreservesExistingAnnotations(t *testing.T) {
	policies := makePinnerPolicies()
	result, err := Pin(policies, "gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range result {
		if p.Name == "gamma" {
			if p.Annotations["owner"] != "team-a" {
				t.Error("expected existing annotation to be preserved")
			}
			if p.Annotations["pinned"] != "true" {
				t.Error("expected pinned annotation to be set")
			}
		}
	}
}

func TestPin_NotFound(t *testing.T) {
	policies := makePinnerPolicies()
	_, err := Pin(policies, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestPin_EmptyName(t *testing.T) {
	policies := makePinnerPolicies()
	_, err := Pin(policies, "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUnpin_Success(t *testing.T) {
	policies := makePinnerPolicies()
	policies, _ = Pin(policies, "beta")
	result, err := Unpin(policies, "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsPinned(result, "beta") {
		t.Error("expected beta to be unpinned")
	}
}

func TestUnpin_NotFound(t *testing.T) {
	policies := makePinnerPolicies()
	_, err := Unpin(policies, "missing")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestListPinned_ReturnsOnlyPinned(t *testing.T) {
	policies := makePinnerPolicies()
	policies, _ = Pin(policies, "alpha")
	policies, _ = Pin(policies, "gamma")
	pinned := ListPinned(policies)
	if len(pinned) != 2 {
		t.Errorf("expected 2 pinned policies, got %d", len(pinned))
	}
}

func TestListPinned_EmptyWhenNonePinned(t *testing.T) {
	policies := makePinnerPolicies()
	pinned := ListPinned(policies)
	if len(pinned) != 0 {
		t.Errorf("expected 0 pinned policies, got %d", len(pinned))
	}
}
