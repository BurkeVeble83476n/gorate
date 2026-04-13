package policy

import (
	"strings"
	"testing"
)

func makeNormalizerPolicies() []Policy {
	return []Policy{
		{Name: "  api-get  ", Endpoint: "  /api/items  ", Method: "get", Limit: 10, Window: 60},
		{Name: "no-slash", Endpoint: "api/no-slash", Method: "POST", Limit: 5, Window: 30},
		{Name: "already-clean", Endpoint: "/api/clean", Method: "GET", Limit: 20, Window: 60},
		{Name: "mixed-method", Endpoint: "/api/mixed", Method: "Delete", Limit: 3, Window: 10},
	}
}

func TestNormalize_ReturnsSameCount(t *testing.T) {
	policies := makeNormalizerPolicies()
	result := Normalize(policies)
	if len(result.Normalized) != len(policies) {
		t.Errorf("expected %d policies, got %d", len(policies), len(result.Normalized))
	}
}

func TestNormalize_TrimsNameWhitespace(t *testing.T) {
	policies := makeNormalizerPolicies()
	result := Normalize(policies)
	for _, p := range result.Normalized {
		if strings.TrimSpace(p.Name) != p.Name {
			t.Errorf("expected trimmed name, got '%s'", p.Name)
		}
	}
}

func TestNormalize_UppercasesMethod(t *testing.T) {
	policies := makeNormalizerPolicies()
	result := Normalize(policies)
	for _, p := range result.Normalized {
		if p.Method != strings.ToUpper(p.Method) {
			t.Errorf("expected uppercase method, got '%s'", p.Method)
		}
	}
}

func TestNormalize_AddsLeadingSlash(t *testing.T) {
	policies := []Policy{
		{Name: "no-slash", Endpoint: "api/resource", Method: "GET", Limit: 5, Window: 30},
	}
	result := Normalize(policies)
	if !strings.HasPrefix(result.Normalized[0].Endpoint, "/") {
		t.Errorf("expected leading slash in endpoint, got '%s'", result.Normalized[0].Endpoint)
	}
}

func TestNormalize_AlreadyCleanPolicyNoChange(t *testing.T) {
	policies := []Policy{
		{Name: "clean", Endpoint: "/api/clean", Method: "GET", Limit: 10, Window: 60},
	}
	result := Normalize(policies)
	if len(result.Changes) != 0 {
		t.Errorf("expected no changes for clean policy, got %d", len(result.Changes))
	}
}

func TestNormalize_ChangesRecorded(t *testing.T) {
	policies := []Policy{
		{Name: "dirty", Endpoint: "api/dirty", Method: "post", Limit: 5, Window: 30},
	}
	result := Normalize(policies)
	if len(result.Changes) == 0 {
		t.Error("expected changes to be recorded for dirty policy")
	}
}

func TestNormalize_TrimsEndpointWhitespace(t *testing.T) {
	policies := []Policy{
		{Name: "ws-endpoint", Endpoint: "  /api/ws  ", Method: "GET", Limit: 5, Window: 30},
	}
	result := Normalize(policies)
	if strings.TrimSpace(result.Normalized[0].Endpoint) != result.Normalized[0].Endpoint {
		t.Errorf("expected trimmed endpoint, got '%s'", result.Normalized[0].Endpoint)
	}
}
