package policy

import (
	"testing"
)

func makeTemplate() Template {
	return Template{
		Name:     "{{SERVICE}}-limit",
		Endpoint: "/api/{{SERVICE}}/{{RESOURCE}}",
		Method:   "GET",
		Limit:    100,
		Window:   60,
	}
}

func TestApplyTemplate_Success(t *testing.T) {
	tmpl := makeTemplate()
	vars := map[string]string{"SERVICE": "users", "RESOURCE": "profile"}

	p, err := ApplyTemplate(tmpl, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "users-limit" {
		t.Errorf("expected name %q, got %q", "users-limit", p.Name)
	}
	if p.Endpoint != "/api/users/profile" {
		t.Errorf("expected endpoint %q, got %q", "/api/users/profile", p.Endpoint)
	}
	if p.Limit != 100 {
		t.Errorf("expected limit 100, got %d", p.Limit)
	}
}

func TestApplyTemplate_DefaultVarsUsed(t *testing.T) {
	tmpl := makeTemplate()
	tmpl.Vars = map[string]string{"SERVICE": "orders", "RESOURCE": "list"}

	p, err := ApplyTemplate(tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "orders-limit" {
		t.Errorf("expected name %q, got %q", "orders-limit", p.Name)
	}
}

func TestApplyTemplate_VarsOverrideDefaults(t *testing.T) {
	tmpl := makeTemplate()
	tmpl.Vars = map[string]string{"SERVICE": "default", "RESOURCE": "items"}
	vars := map[string]string{"SERVICE": "payments"}

	p, err := ApplyTemplate(tmpl, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "payments-limit" {
		t.Errorf("expected name %q, got %q", "payments-limit", p.Name)
	}
}

func TestApplyTemplate_MissingName(t *testing.T) {
	tmpl := makeTemplate()
	tmpl.Name = ""

	_, err := ApplyTemplate(tmpl, nil)
	if err == nil {
		t.Fatal("expected error for missing template name")
	}
}

func TestApplyTemplate_InvalidPolicy(t *testing.T) {
	tmpl := makeTemplate()
	tmpl.Limit = -1

	_, err := ApplyTemplate(tmpl, map[string]string{"SERVICE": "s", "RESOURCE": "r"})
	if err == nil {
		t.Fatal("expected error for invalid limit")
	}
}

func TestApplyTemplate_AppliesDefaults(t *testing.T) {
	tmpl := makeTemplate()
	tmpl.Method = ""

	p, err := ApplyTemplate(tmpl, map[string]string{"SERVICE": "svc", "RESOURCE": "res"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Method == "" {
		t.Error("expected default method to be applied")
	}
}
