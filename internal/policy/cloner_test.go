package policy

import (
	"testing"
)

func makeClonerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/v2", Method: "POST", Limit: 50, Window: 30, Tags: []string{"prod"}},
	}
}

func TestClone_Success(t *testing.T) {
	policies := makeClonerPolicies()
	result, err := Clone(policies, "alpha", CloneOptions{NewName: "alpha-copy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 policies, got %d", len(result))
	}
	cloned := result[2]
	if cloned.Name != "alpha-copy" {
		t.Errorf("expected name alpha-copy, got %s", cloned.Name)
	}
	if cloned.Endpoint != "/api/v1" || cloned.Limit != 100 {
		t.Errorf("cloned policy fields do not match source")
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	policies := makeClonerPolicies()
	_, err := Clone(policies, "nonexistent", CloneOptions{NewName: "copy"})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestClone_EmptyNewName(t *testing.T) {
	policies := makeClonerPolicies()
	_, err := Clone(policies, "alpha", CloneOptions{NewName: ""})
	if err == nil {
		t.Fatal("expected error for empty new name")
	}
}

func TestClone_ConflictNoOverride(t *testing.T) {
	policies := makeClonerPolicies()
	_, err := Clone(policies, "alpha", CloneOptions{NewName: "beta"})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestClone_ConflictWithOverride(t *testing.T) {
	policies := makeClonerPolicies()
	result, err := Clone(policies, "alpha", CloneOptions{NewName: "beta", Override: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 policies after override, got %d", len(result))
	}
}

func TestClone_TagsAreDeepCopied(t *testing.T) {
	policies := makeClonerPolicies()
	result, err := Clone(policies, "beta", CloneOptions{NewName: "beta-copy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cloned := result[len(result)-1]
	cloned.Tags[0] = "mutated"
	original := result[1]
	if original.Tags[0] == "mutated" {
		t.Error("tags are not deep copied; mutation affected original")
	}
}
