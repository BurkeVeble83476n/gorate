package policy

import (
	"testing"
)

func makeFreezePolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 5, Window: 30},
		{Name: "gamma", Endpoint: "/api/gamma", Method: "*", Limit: 20, Window: 120},
	}
}

func TestFreeze_Success(t *testing.T) {
	policies := makeFreezePolicies()
	updated, err := Freeze(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsFrozen(updated, "alpha") {
		t.Error("expected alpha to be frozen")
	}
}

func TestFreeze_NotFound(t *testing.T) {
	policies := makeFreezePolicies()
	_, err := Freeze(policies, "nonexistent")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestFreeze_EmptyName(t *testing.T) {
	policies := makeFreezePolicies()
	_, err := Freeze(policies, "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestUnfreeze_Success(t *testing.T) {
	policies := makeFreezePolicies()
	updated, _ := Freeze(policies, "beta")
	if !IsFrozen(updated, "beta") {
		t.Fatal("expected beta to be frozen before unfreeze")
	}
	updated, err := Unfreeze(updated, "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsFrozen(updated, "beta") {
		t.Error("expected beta to be unfrozen")
	}
}

func TestUnfreeze_NotFound(t *testing.T) {
	policies := makeFreezePolicies()
	_, err := Unfreeze(policies, "missing")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestIsFrozen_FalseWhenNotFrozen(t *testing.T) {
	policies := makeFreezePolicies()
	if IsFrozen(policies, "gamma") {
		t.Error("expected gamma to not be frozen")
	}
}

func TestIsFrozen_FalseForUnknownPolicy(t *testing.T) {
	policies := makeFreezePolicies()
	if IsFrozen(policies, "unknown") {
		t.Error("expected false for unknown policy")
	}
}

func TestListFrozen_ReturnsOnlyFrozen(t *testing.T) {
	policies := makeFreezePolicies()
	updated, _ := Freeze(policies, "alpha")
	updated, _ = Freeze(updated, "gamma")
	frozen := ListFrozen(updated)
	if len(frozen) != 2 {
		t.Fatalf("expected 2 frozen policies, got %d", len(frozen))
	}
}

func TestListFrozen_EmptyWhenNoneFrozen(t *testing.T) {
	policies := makeFreezePolicies()
	frozen := ListFrozen(policies)
	if len(frozen) != 0 {
		t.Errorf("expected 0 frozen policies, got %d", len(frozen))
	}
}
