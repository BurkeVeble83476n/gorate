package policy

import (
	"testing"
)

func makeGrouperPolicies() []Policy {
	return []Policy{
		{Name: "a", Endpoint: "/api/users", Method: "GET", Limit: 100, Window: "1m"},
		{Name: "b", Endpoint: "/api/orders", Method: "POST", Limit: 50, Window: "1m"},
		{Name: "c", Endpoint: "/api/users", Method: "POST", Limit: 30, Window: "5m"},
		{Name: "d", Endpoint: "/api/health", Method: "GET", Limit: 200, Window: "1m"},
	}
}

func TestGroup_ByMethod(t *testing.T) {
	policies := makeGrouperPolicies()
	groups, err := Group(policies, GroupByMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// sorted: GET, POST
	if groups[0].Key != "GET" {
		t.Errorf("expected first key GET, got %s", groups[0].Key)
	}
	if len(groups[0].Policies) != 2 {
		t.Errorf("expected 2 GET policies, got %d", len(groups[0].Policies))
	}
}

func TestGroup_ByEndpoint(t *testing.T) {
	policies := makeGrouperPolicies()
	groups, err := Group(policies, GroupByEndpoint)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestGroup_ByWindow(t *testing.T) {
	policies := makeGrouperPolicies()
	groups, err := Group(policies, GroupByWindow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroup_InvalidField(t *testing.T) {
	policies := makeGrouperPolicies()
	_, err := Group(policies, GroupBy("invalid"))
	if err == nil {
		t.Error("expected error for invalid group-by field")
	}
}

func TestGroup_EmptyPolicies(t *testing.T) {
	groups, err := Group([]Policy{}, GroupByMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroup_SortedKeys(t *testing.T) {
	policies := makeGrouperPolicies()
	groups, err := Group(policies, GroupByEndpoint)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(groups); i++ {
		if groups[i].Key < groups[i-1].Key {
			t.Errorf("groups not sorted: %s before %s", groups[i-1].Key, groups[i].Key)
		}
	}
}
