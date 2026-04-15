package policy

import (
	"testing"
)

func makeAliaserPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/v2", Method: "POST", Limit: 50, Window: 30},
	}
}

func TestAddAlias_Success(t *testing.T) {
	policies := makeAliaserPolicies()
	updated, err := AddAlias(policies, "alpha", "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aliases := GetAliases(updated[0])
	if len(aliases) != 1 || aliases[0] != "a1" {
		t.Errorf("expected alias 'a1', got %v", aliases)
	}
}

func TestAddAlias_Idempotent(t *testing.T) {
	policies := makeAliaserPolicies()
	policies, _ = AddAlias(policies, "alpha", "a1")
	policies, err := AddAlias(policies, "alpha", "a1")
	if err != nil {
		t.Fatalf("unexpected error on duplicate add: %v", err)
	}
	if len(GetAliases(policies[0])) != 1 {
		t.Errorf("expected exactly 1 alias after idempotent add")
	}
}

func TestAddAlias_Collision(t *testing.T) {
	policies := makeAliaserPolicies()
	policies, _ = AddAlias(policies, "alpha", "shared")
	_, err := AddAlias(policies, "beta", "shared")
	if err == nil {
		t.Error("expected error for alias collision, got nil")
	}
}

func TestAddAlias_PolicyNotFound(t *testing.T) {
	policies := makeAliaserPolicies()
	_, err := AddAlias(policies, "ghost", "g1")
	if err == nil {
		t.Error("expected error for missing policy, got nil")
	}
}

func TestAddAlias_EmptyName(t *testing.T) {
	policies := makeAliaserPolicies()
	_, err := AddAlias(policies, "", "a1")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestAddAlias_EmptyAlias(t *testing.T) {
	policies := makeAliaserPolicies()
	_, err := AddAlias(policies, "alpha", "")
	if err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestRemoveAlias_Success(t *testing.T) {
	policies := makeAliaserPolicies()
	policies, _ = AddAlias(policies, "alpha", "a1")
	policies, err := RemoveAlias(policies, "alpha", "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(GetAliases(policies[0])) != 0 {
		t.Error("expected no aliases after removal")
	}
}

func TestRemoveAlias_PolicyNotFound(t *testing.T) {
	policies := makeAliaserPolicies()
	_, err := RemoveAlias(policies, "ghost", "g1")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestFindByAlias_Found(t *testing.T) {
	policies := makeAliaserPolicies()
	policies, _ = AddAlias(policies, "beta", "b-alias")
	p, err := FindByAlias(policies, "b-alias")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "beta" {
		t.Errorf("expected policy 'beta', got %q", p.Name)
	}
}

func TestFindByAlias_NotFound(t *testing.T) {
	policies := makeAliaserPolicies()
	_, err := FindByAlias(policies, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown alias")
	}
}
