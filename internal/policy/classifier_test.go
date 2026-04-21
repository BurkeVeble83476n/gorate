package policy

import (
	"strings"
	"testing"
	"time"
)

func makeClassifierPolicies() []Policy {
	return []Policy{
		{Name: "strict-api", Endpoint: "/api/data", Method: "GET", Limit: 5, Window: 30 * time.Second},
		{Name: "moderate-api", Endpoint: "/api/submit", Method: "POST", Limit: 10, Window: 5 * time.Second},
		{Name: "permissive-api", Endpoint: "/api/ping", Method: "GET", Limit: 100, Window: 5 * time.Second},
		{Name: "wildcard-strict", Endpoint: "/admin", Method: "*", Limit: 2, Window: 10 * time.Second},
		{Name: "wildcard-perm", Endpoint: "/public", Method: "*", Limit: 50, Window: 5 * time.Second},
	}
}

func TestClassify_StrictPolicy(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify(policies[:1])
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Class != "strict" {
		t.Errorf("expected strict, got %s", results[0].Class)
	}
}

func TestClassify_ModeratePolicy(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify([]Policy{policies[1]})
	if results[0].Class != "moderate" {
		t.Errorf("expected moderate, got %s", results[0].Class)
	}
}

func TestClassify_PermissivePolicy(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify([]Policy{policies[2]})
	if results[0].Class != "permissive" {
		t.Errorf("expected permissive, got %s", results[0].Class)
	}
}

func TestClassify_WildcardStrict(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify([]Policy{policies[3]})
	if results[0].Class != "wildcard-strict" {
		t.Errorf("expected wildcard-strict, got %s", results[0].Class)
	}
}

func TestClassify_WildcardPermissive(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify([]Policy{policies[4]})
	if results[0].Class != "wildcard-permissive" {
		t.Errorf("expected wildcard-permissive, got %s", results[0].Class)
	}
}

func TestClassify_ReasonContainsRate(t *testing.T) {
	policies := makeClassifierPolicies()
	results := Classify(policies)
	for _, r := range results {
		if !strings.Contains(r.Reason, "req/s") {
			t.Errorf("expected reason to contain req/s for policy %s, got: %s", r.PolicyName, r.Reason)
		}
	}
}

func TestClassify_EmptyPolicies(t *testing.T) {
	results := Classify([]Policy{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFormatClassify_ContainsHeaders(t *testing.T) {
	results := Classify(makeClassifierPolicies())
	out := FormatClassify(results)
	for _, header := range []string{"POLICY", "CLASS", "REASON"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected output to contain header %q", header)
		}
	}
}

func TestFormatClassify_Empty(t *testing.T) {
	out := FormatClassify([]ClassifyResult{})
	if !strings.Contains(out, "no policies") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
