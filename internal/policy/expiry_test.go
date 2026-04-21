package policy_test

import (
	"testing"
	"time"

	"github.com/yourusername/gorate/internal/policy"
)

func makeExpiryPolicies() []policy.Policy {
	return []policy.Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 20, Window: 30},
		{Name: "gamma", Endpoint: "/api/gamma", Method: "*", Limit: 5, Window: 120},
	}
}

func TestSetExpiry_Success(t *testing.T) {
	policies := makeExpiryPolicies()
	at := time.Now().Add(24 * time.Hour)
	updated, err := policy.SetExpiry(policies, "alpha", at)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exp, _ := policy.GetExpiry(updated, "alpha")
	if !exp.Equal(at) {
		t.Errorf("expected expiry %v, got %v", at, exp)
	}
}

func TestSetExpiry_PolicyNotFound(t *testing.T) {
	policies := makeExpiryPolicies()
	_, err := policy.SetExpiry(policies, "nonexistent", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestRemoveExpiry_ClearsExpiry(t *testing.T) {
	policies := makeExpiryPolicies()
	policies, _ = policy.SetExpiry(policies, "beta", time.Now().Add(time.Hour))
	updated, err := policy.RemoveExpiry(policies, "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, found := policy.GetExpiry(updated, "beta")
	if found {
		t.Error("expected expiry to be removed")
	}
}

func TestEvaluateExpiry_ExpiredPolicy(t *testing.T) {
	policies := makeExpiryPolicies()
	past := time.Now().Add(-time.Hour)
	policies, _ = policy.SetExpiry(policies, "gamma", past)
	results := policy.EvaluateExpiry(policies)
	if len(results) == 0 {
		t.Fatal("expected at least one expired result")
	}
	if results[0].Name != "gamma" {
		t.Errorf("expected gamma, got %s", results[0].Name)
	}
	if !results[0].Expired {
		t.Error("expected policy to be marked expired")
	}
}

func TestEvaluateExpiry_ActivePolicy(t *testing.T) {
	policies := makeExpiryPolicies()
	future := time.Now().Add(48 * time.Hour)
	policies, _ = policy.SetExpiry(policies, "alpha", future)
	results := policy.EvaluateExpiry(policies)
	for _, r := range results {
		if r.Name == "alpha" && r.Expired {
			t.Error("alpha should not be expired")
		}
	}
}

func TestFilterExpired_RemovesExpired(t *testing.T) {
	policies := makeExpiryPolicies()
	policies, _ = policy.SetExpiry(policies, "beta", time.Now().Add(-time.Minute))
	filtered := policy.FilterExpired(policies)
	for _, p := range filtered {
		if p.Name == "beta" {
			t.Error("expected beta to be removed as expired")
		}
	}
	if len(filtered) != 2 {
		t.Errorf("expected 2 policies, got %d", len(filtered))
	}
}
