package policy

import (
	"testing"
)

func makeDeprecatorPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/alpha", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/beta", Method: "POST", Limit: 50, Window: 30},
		{Name: "gamma", Endpoint: "/gamma", Method: "*", Limit: 200, Window: 120},
	}
}

func TestDeprecate_Success(t *testing.T) {
	policies := makeDeprecatorPolicies()
	updated, err := Deprecate(policies, "alpha", "use /v2/alpha", "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsDeprecated(updated[0]) {
		t.Error("expected alpha to be deprecated")
	}
}

func TestDeprecate_SetsReason(t *testing.T) {
	policies := makeDeprecatorPolicies()
	updated, err := Deprecate(policies, "beta", "legacy endpoint", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, ok := GetDeprecationInfo(updated[1])
	if !ok {
		t.Fatal("expected deprecation info")
	}
	if info.Reason != "legacy endpoint" {
		t.Errorf("expected reason 'legacy endpoint', got %q", info.Reason)
	}
}

func TestDeprecate_SetsReplacement(t *testing.T) {
	policies := makeDeprecatorPolicies()
	updated, err := Deprecate(policies, "alpha", "", "gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, _ := GetDeprecationInfo(updated[0])
	if info.Replacement != "gamma" {
		t.Errorf("expected replacement 'gamma', got %q", info.Replacement)
	}
}

func TestDeprecate_PolicyNotFound(t *testing.T) {
	policies := makeDeprecatorPolicies()
	_, err := Deprecate(policies, "nonexistent", "", "")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestDeprecate_EmptyName(t *testing.T) {
	policies := makeDeprecatorPolicies()
	_, err := Deprecate(policies, "", "", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestUndeprecate_Success(t *testing.T) {
	policies := makeDeprecatorPolicies()
	policies, _ = Deprecate(policies, "alpha", "old", "beta")
	policies, err := Undeprecate(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsDeprecated(policies[0]) {
		t.Error("expected alpha to no longer be deprecated")
	}
}

func TestUndeprecate_NotFound(t *testing.T) {
	policies := makeDeprecatorPolicies()
	_, err := Undeprecate(policies, "missing")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestListDeprecated_ReturnsOnlyDeprecated(t *testing.T) {
	policies := makeDeprecatorPolicies()
	policies, _ = Deprecate(policies, "alpha", "", "")
	policies, _ = Deprecate(policies, "gamma", "", "")
	deprecated := ListDeprecated(policies)
	if len(deprecated) != 2 {
		t.Errorf("expected 2 deprecated, got %d", len(deprecated))
	}
}

func TestListDeprecated_NoneDeprecated(t *testing.T) {
	policies := makeDeprecatorPolicies()
	deprecated := ListDeprecated(policies)
	if len(deprecated) != 0 {
		t.Errorf("expected 0 deprecated, got %d", len(deprecated))
	}
}
