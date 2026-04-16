package policy

import (
	"strings"
	"testing"
)

func makeInheritancePolicies() []Policy {
	return []Policy{
		{Name: "base", Endpoint: "/api", Method: "GET", Limit: 100, Window: "60s"},
		{Name: "child", Endpoint: "", Method: "", Limit: 0, Window: "",
			Annotations: map[string]string{"inherits": "base"}},
		{Name: "override", Endpoint: "/custom", Method: "POST", Limit: 50, Window: "30s",
			Annotations: map[string]string{"inherits": "base"}},
		{Name: "standalone", Endpoint: "/other", Method: "GET", Limit: 10, Window: "10s"},
	}
}

func TestResolve_InheritsAllFields(t *testing.T) {
	parent := Policy{Name: "parent", Endpoint: "/api", Method: "GET", Limit: 100, Window: "60s"}
	child := Policy{Name: "child"}

	r, err := Resolve(parent, child)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Policy.Endpoint != "/api" {
		t.Errorf("expected endpoint /api, got %q", r.Policy.Endpoint)
	}
	if r.Policy.Limit != 100 {
		t.Errorf("expected limit 100, got %d", r.Policy.Limit)
	}
	if len(r.Inherited) != 4 {
		t.Errorf("expected 4 inherited fields, got %d", len(r.Inherited))
	}
}

func TestResolve_ChildOverridesParent(t *testing.T) {
	parent := Policy{Name: "parent", Endpoint: "/api", Method: "GET", Limit: 100, Window: "60s"}
	child := Policy{Name: "child", Endpoint: "/custom", Method: "POST", Limit: 50, Window: "30s"}

	r, err := Resolve(parent, child)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Policy.Endpoint != "/custom" {
		t.Errorf("expected /custom, got %q", r.Policy.Endpoint)
	}
	if len(r.Inherited) != 0 {
		t.Errorf("expected no inherited fields, got %v", r.Inherited)
	}
}

func TestResolve_MissingParentName(t *testing.T) {
	parent := Policy{Name: ""}
	child := Policy{Name: "child"}
	_, err := Resolve(parent, child)
	if err == nil {
		t.Fatal("expected error for empty parent name")
	}
}

func TestFormatInheritance_WithInherited(t *testing.T) {
	r := InheritanceResult{
		Policy:    Policy{Name: "child"},
		Inherited: []string{"limit", "window"},
	}
	out := FormatInheritance(r, "parent")
	if !strings.Contains(out, "child") || !strings.Contains(out, "parent") {
		t.Errorf("unexpected format output: %s", out)
	}
	if !strings.Contains(out, "limit") {
		t.Errorf("expected limit in output: %s", out)
	}
}

func TestFormatInheritance_NothingInherited(t *testing.T) {
	r := InheritanceResult{Policy: Policy{Name: "child"}, Inherited: nil}
	out := FormatInheritance(r, "parent")
	if !strings.Contains(out, "nothing") {
		t.Errorf("expected 'nothing' in output: %s", out)
	}
}

func TestApplyInheritance_ResolvesChild(t *testing.T) {
	policies := makeInheritancePolicies()
	resolved, notes, err := ApplyInheritance(policies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != len(policies) {
		t.Errorf("expected %d policies, got %d", len(policies), len(resolved))
	}
	if len(notes) == 0 {
		t.Error("expected at least one inheritance note")
	}
	for _, p := range resolved {
		if p.Name == "child" && p.Endpoint != "/api" {
			t.Errorf("child should inherit endpoint /api, got %q", p.Endpoint)
		}
	}
}

func TestApplyInheritance_UnknownParent(t *testing.T) {
	policies := []Policy{
		{Name: "orphan", Annotations: map[string]string{"inherits": "ghost"}},
	}
	_, _, err := ApplyInheritance(policies)
	if err == nil {
		t.Fatal("expected error for unknown parent")
	}
}
