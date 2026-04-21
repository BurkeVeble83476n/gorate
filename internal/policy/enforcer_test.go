package policy

import (
	"strings"
	"testing"
)

func makeEnforcerPolicies() []Policy {
	return []Policy{
		{Name: "safe", Endpoint: "/api/safe", Method: "GET", Limit: 100, Window: 60},
		{Name: "high-limit", Endpoint: "/api/bulk", Method: "POST", Limit: 20000, Window: 60},
		{Name: "wildcard-heavy", Endpoint: "/api/data", Method: "*", Limit: 6000, Window: 30},
		{Name: "short-window", Endpoint: "/api/fast", Method: "GET", Limit: 50, Window: 0},
	}
}

func TestEnforce_StrictBlocksViolations(t *testing.T) {
	policies := makeEnforcerPolicies()
	results := Enforce(policies, EnforceModeStrict)
	if len(results) != len(policies) {
		t.Fatalf("expected %d results, got %d", len(policies), len(results))
	}
	byName := map[string]EnforcementResult{}
	for _, r := range results {
		byName[r.PolicyName] = r
	}
	if !byName["safe"].Enforced {
		t.Error("expected 'safe' to be enforced")
	}
	if byName["high-limit"].Enforced {
		t.Error("expected 'high-limit' to be blocked in strict mode")
	}
	if byName["short-window"].Enforced {
		t.Error("expected 'short-window' to be blocked due to window=0")
	}
}

func TestEnforce_WarnAllowsButRecordsViolations(t *testing.T) {
	policies := makeEnforcerPolicies()
	results := Enforce(policies, EnforceModeWarn)
	for _, r := range results {
		if !r.Enforced {
			t.Errorf("expected policy %q to be enforced in warn mode", r.PolicyName)
		}
	}
	byName := map[string]EnforcementResult{}
	for _, r := range results {
		byName[r.PolicyName] = r
	}
	if len(byName["high-limit"].Violations) == 0 {
		t.Error("expected violations recorded for 'high-limit' in warn mode")
	}
}

func TestEnforce_DisableSkipsAll(t *testing.T) {
	policies := makeEnforcerPolicies()
	results := Enforce(policies, EnforceModeDisable)
	for _, r := range results {
		if !r.Enforced {
			t.Errorf("expected policy %q to be enforced in disable mode", r.PolicyName)
		}
		if len(r.Violations) != 0 {
			t.Errorf("expected no violations collected in disable mode for %q", r.PolicyName)
		}
	}
}

func TestParseEnforcementMode_Valid(t *testing.T) {
	for _, tc := range []string{"strict", "warn", "disable"} {
		m, err := ParseEnforcementMode(tc)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc, err)
		}
		if string(m) != tc {
			t.Errorf("expected mode %q, got %q", tc, m)
		}
	}
}

func TestParseEnforcementMode_Invalid(t *testing.T) {
	_, err := ParseEnforcementMode("unknown")
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}

func TestFormatEnforcement_ContainsPolicyName(t *testing.T) {
	policies := []Policy{
		{Name: "my-policy", Endpoint: "/api/x", Method: "GET", Limit: 10, Window: 60},
	}
	results := Enforce(policies, EnforceModeStrict)
	out := FormatEnforcement(results)
	if !strings.Contains(out, "my-policy") {
		t.Error("expected output to contain policy name")
	}
}

func TestFormatEnforcement_EmptyReturnsMessage(t *testing.T) {
	out := FormatEnforcement(nil)
	if !strings.Contains(out, "No policies") {
		t.Error("expected empty message for nil results")
	}
}
