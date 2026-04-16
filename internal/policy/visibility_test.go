package policy

import (
	"testing"
)

func makeVisibilityPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 20, Window: 60},
		{Name: "gamma", Endpoint: "/c", Method: "GET", Limit: 5, Window: 30,
			Annotations: map[string]string{"visibility": "internal"}},
	}
}

func TestSetVisibility_Success(t *testing.T) {
	policies := makeVisibilityPolicies()
	out, err := SetVisibility(policies, "alpha", "public")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Annotations["visibility"] != "public" {
		t.Errorf("expected public, got %s", out[0].Annotations["visibility"])
	}
}

func TestSetVisibility_InvalidLevel(t *testing.T) {
	policies := makeVisibilityPolicies()
	_, err := SetVisibility(policies, "alpha", "secret")
	if err == nil {
		t.Fatal("expected error for invalid visibility")
	}
}

func TestSetVisibility_PolicyNotFound(t *testing.T) {
	policies := makeVisibilityPolicies()
	_, err := SetVisibility(policies, "missing", "public")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestSetVisibility_EmptyName(t *testing.T) {
	policies := makeVisibilityPolicies()
	_, err := SetVisibility(policies, "", "public")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestGetVisibility_ReturnsValue(t *testing.T) {
	policies := makeVisibilityPolicies()
	v, err := GetVisibility(policies, "gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "internal" {
		t.Errorf("expected internal, got %s", v)
	}
}

func TestGetVisibility_NotSet(t *testing.T) {
	policies := makeVisibilityPolicies()
	v, err := GetVisibility(policies, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "" {
		t.Errorf("expected empty, got %s", v)
	}
}

func TestFilterByVisibility_ReturnsMatching(t *testing.T) {
	policies := makeVisibilityPolicies()
	out := FilterByVisibility(policies, "internal")
	if len(out) != 1 || out[0].Name != "gamma" {
		t.Errorf("expected gamma, got %+v", out)
	}
}

func TestFilterByVisibility_NoMatch(t *testing.T) {
	policies := makeVisibilityPolicies()
	out := FilterByVisibility(policies, "private")
	if len(out) != 0 {
		t.Errorf("expected empty, got %d", len(out))
	}
}
