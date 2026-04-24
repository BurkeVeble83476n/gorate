package policy

import (
	"testing"
)

func makeConditionerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/data", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/submit", Method: "POST", Limit: 20, Window: 30},
	}
}

func TestSetCondition_Success(t *testing.T) {
	policies := makeConditionerPolicies()
	updated, err := SetCondition(policies, "alpha", "env", "eq", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cond, _ := GetCondition(updated, "alpha")
	if cond == nil {
		t.Fatal("expected condition, got nil")
	}
	if cond.Field != "env" || cond.Operator != "eq" || cond.Value != "production" {
		t.Errorf("unexpected condition: %+v", cond)
	}
}

func TestSetCondition_InvalidOperator(t *testing.T) {
	policies := makeConditionerPolicies()
	_, err := SetCondition(policies, "alpha", "env", "gt", "production")
	if err == nil {
		t.Fatal("expected error for invalid operator")
	}
}

func TestSetCondition_PolicyNotFound(t *testing.T) {
	policies := makeConditionerPolicies()
	_, err := SetCondition(policies, "ghost", "env", "eq", "staging")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestSetCondition_EmptyName(t *testing.T) {
	policies := makeConditionerPolicies()
	_, err := SetCondition(policies, "", "env", "eq", "staging")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestGetCondition_NoConditionSet(t *testing.T) {
	policies := makeConditionerPolicies()
	cond, err := GetCondition(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cond != nil {
		t.Errorf("expected nil condition, got %+v", cond)
	}
}

func TestGetCondition_PolicyNotFound(t *testing.T) {
	policies := makeConditionerPolicies()
	_, err := GetCondition(policies, "unknown")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestEvaluateCondition_NoCondition(t *testing.T) {
	p := Policy{Name: "alpha", Endpoint: "/api", Method: "GET", Limit: 10, Window: 60}
	if !EvaluateCondition(p, map[string]string{"env": "production"}) {
		t.Error("expected true when no condition is set")
	}
}

func TestEvaluateCondition_EqMatch(t *testing.T) {
	p := Policy{
		Name: "alpha", Endpoint: "/api", Method: "GET", Limit: 10, Window: 60,
		Annotations: map[string]string{
			"condition.field":    "env",
			"condition.operator": "eq",
			"condition.value":    "production",
		},
	}
	if !EvaluateCondition(p, map[string]string{"env": "production"}) {
		t.Error("expected true for eq match")
	}
	if EvaluateCondition(p, map[string]string{"env": "staging"}) {
		t.Error("expected false for eq mismatch")
	}
}

func TestEvaluateCondition_ContainsMatch(t *testing.T) {
	p := Policy{
		Name: "beta", Endpoint: "/api", Method: "POST", Limit: 5, Window: 30,
		Annotations: map[string]string{
			"condition.field":    "path",
			"condition.operator": "contains",
			"condition.value":    "admin",
		},
	}
	if !EvaluateCondition(p, map[string]string{"path": "/admin/users"}) {
		t.Error("expected true for contains match")
	}
	if EvaluateCondition(p, map[string]string{"path": "/public/home"}) {
		t.Error("expected false for contains mismatch")
	}
}
