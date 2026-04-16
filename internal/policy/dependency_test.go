package policy

import (
	"testing"
)

func makeDependencyPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 5, Window: 30},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 20, Window: 120},
	}
}

func TestAddDependency_Success(t *testing.T) {
	policies := makeDependencyPolicies()
	updated, err := AddDependency(policies, "alpha", "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, _ := GetDependencies(updated, "alpha")
	if len(deps) != 1 || deps[0] != "beta" {
		t.Errorf("expected [beta], got %v", deps)
	}
}

func TestAddDependency_Idempotent(t *testing.T) {
	policies := makeDependencyPolicies()
	updated, _ := AddDependency(policies, "alpha", "beta")
	updated, err := AddDependency(updated, "alpha", "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, _ := GetDependencies(updated, "alpha")
	if len(deps) != 1 {
		t.Errorf("expected 1 dep, got %d", len(deps))
	}
}

func TestAddDependency_SelfReference(t *testing.T) {
	policies := makeDependencyPolicies()
	_, err := AddDependency(policies, "alpha", "alpha")
	if err == nil {
		t.Error("expected error for self-reference")
	}
}

func TestAddDependency_FromNotFound(t *testing.T) {
	policies := makeDependencyPolicies()
	_, err := AddDependency(policies, "unknown", "beta")
	if err == nil {
		t.Error("expected error for missing from policy")
	}
}

func TestAddDependency_ToNotFound(t *testing.T) {
	policies := makeDependencyPolicies()
	_, err := AddDependency(policies, "alpha", "unknown")
	if err == nil {
		t.Error("expected error for missing to policy")
	}
}

func TestGetDependencies_Empty(t *testing.T) {
	policies := makeDependencyPolicies()
	deps, err := GetDependencies(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected no deps, got %v", deps)
	}
}

func TestGetDependencies_PolicyNotFound(t *testing.T) {
	policies := makeDependencyPolicies()
	_, err := GetDependencies(policies, "missing")
	if err == nil {
		t.Error("expected error for missing policy")
	}
}

func TestBuildGraph_EdgesPopulated(t *testing.T) {
	policies := makeDependencyPolicies()
	policies, _ = AddDependency(policies, "alpha", "beta")
	policies, _ = AddDependency(policies, "alpha", "gamma")
	graph := BuildGraph(policies)
	if len(graph.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(graph.Edges))
	}
}

func TestBuildGraph_EmptyPolicies(t *testing.T) {
	graph := BuildGraph([]Policy{})
	if len(graph.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(graph.Edges))
	}
}
