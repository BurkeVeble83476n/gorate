package policy

import (
	"testing"
)

func makeFilterPolicies() []Policy {
	return []Policy{
		{Name: "get-users", Endpoint: "/api/users", Method: "GET", Limit: 100, Window: 60},
		{Name: "post-users", Endpoint: "/api/users", Method: "POST", Limit: 20, Window: 60},
		{Name: "get-items", Endpoint: "/api/items", Method: "GET", Limit: 50, Window: 30},
		{Name: "wildcard", Endpoint: "/health", Method: "*", Limit: 200, Window: 60},
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	policies := makeFilterPolicies()
	result := Filter(policies, FilterOptions{})
	if len(result) != len(policies) {
		t.Errorf("expected %d policies, got %d", len(policies), len(result))
	}
}

func TestFilter_ByMethod(t *testing.T) {
	result := Filter(makeFilterPolicies(), FilterOptions{Method: "GET"})
	for _, p := range result {
		if p.Method != "GET" && p.Method != "*" {
			t.Errorf("unexpected method %q in filtered results", p.Method)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 results (2 GET + 1 wildcard), got %d", len(result))
	}
}

func TestFilter_ByEndpoint(t *testing.T) {
	result := Filter(makeFilterPolicies(), FilterOptions{Endpoint: "/api/users"})
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestFilter_ByMinLimit(t *testing.T) {
	result := Filter(makeFilterPolicies(), FilterOptions{MinLimit: 100})
	if len(result) != 2 {
		t.Errorf("expected 2 results with limit >= 100, got %d", len(result))
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	result := Filter(makeFilterPolicies(), FilterOptions{
		Method:   "GET",
		Endpoint: "/api/users",
	})
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "get-users" {
		t.Errorf("expected policy 'get-users', got %q", result[0].Name)
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	result := Filter(makeFilterPolicies(), FilterOptions{Endpoint: "/nonexistent"})
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestNormalizePath_TrailingSlash(t *testing.T) {
	if normalizePath("/api/users/") != "/api/users" {
		t.Error("expected trailing slash to be trimmed")
	}
}
