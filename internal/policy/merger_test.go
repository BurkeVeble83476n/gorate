package policy

import (
	"testing"
)

func makeMergePolicies(name, endpoint, method string, limit int) Policy {
	return Policy{
		Name:     name,
		Endpoint: endpoint,
		Method:   method,
		Limit:    limit,
		Window:   "1m",
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	a := []Policy{makeMergePolicies("p1", "/a", "GET", 10)}
	b := []Policy{makeMergePolicies("p2", "/b", "POST", 20)}

	result, err := Merge(MergeStrategyKeepFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(result))
	}
}

func TestMerge_KeepFirst(t *testing.T) {
	a := []Policy{makeMergePolicies("first", "/api", "GET", 10)}
	b := []Policy{makeMergePolicies("second", "/api", "GET", 50)}

	result, err := Merge(MergeStrategyKeepFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(result))
	}
	if result[0].Name != "first" {
		t.Errorf("expected 'first', got %q", result[0].Name)
	}
}

func TestMerge_KeepLast(t *testing.T) {
	a := []Policy{makeMergePolicies("first", "/api", "GET", 10)}
	b := []Policy{makeMergePolicies("second", "/api", "GET", 50)}

	result, err := Merge(MergeStrategyKeepLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Name != "second" {
		t.Errorf("expected 'second', got %q", result[0].Name)
	}
}

func TestMerge_HighestLimit(t *testing.T) {
	a := []Policy{makeMergePolicies("low", "/api", "GET", 5)}
	b := []Policy{makeMergePolicies("high", "/api", "GET", 100)}

	result, err := Merge(MergeStrategyHighest, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Limit != 100 {
		t.Errorf("expected limit 100, got %d", result[0].Limit)
	}
}

func TestMerge_LowestLimit(t *testing.T) {
	a := []Policy{makeMergePolicies("low", "/api", "GET", 5)}
	b := []Policy{makeMergePolicies("high", "/api", "GET", 100)}

	result, err := Merge(MergeStrategyLowest, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Limit != 5 {
		t.Errorf("expected limit 5, got %d", result[0].Limit)
	}
}

func TestMerge_InvalidStrategy(t *testing.T) {
	a := []Policy{makeMergePolicies("p1", "/api", "GET", 10)}
	_, err := Merge(MergeStrategy("bogus"), a)
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}

func TestMerge_CaseInsensitiveKey(t *testing.T) {
	a := []Policy{makeMergePolicies("p1", "/API", "get", 10)}
	b := []Policy{makeMergePolicies("p2", "/api", "GET", 99)}

	result, err := Merge(MergeStrategyKeepFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 deduplicated policy, got %d", len(result))
	}
}
