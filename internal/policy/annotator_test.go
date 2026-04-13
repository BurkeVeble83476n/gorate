package policy

import (
	"testing"
)

func makeAnnotatorPolicies() []Policy {
	return []Policy{
		{Name: "api-get", Endpoint: "/api", Method: "GET", Limit: 100, Window: 60},
		{Name: "api-post", Endpoint: "/api", Method: "POST", Limit: 50, Window: 60},
	}
}

func TestAnnotate_AddsNewKey(t *testing.T) {
	policies := makeAnnotatorPolicies()
	result, err := Annotate(policies, "api-get", "owner", "team-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Annotations["owner"] != "team-a" {
		t.Errorf("expected annotation 'owner'='team-a', got %v", result[0].Annotations)
	}
}

func TestAnnotate_UpdatesExistingKey(t *testing.T) {
	policies := makeAnnotatorPolicies()
	policies[0].Annotations = map[string]string{"owner": "team-a"}
	result, err := Annotate(policies, "api-get", "owner", "team-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Annotations["owner"] != "team-b" {
		t.Errorf("expected updated annotation 'owner'='team-b'")
	}
}

func TestAnnotate_PolicyNotFound(t *testing.T) {
	policies := makeAnnotatorPolicies()
	_, err := Annotate(policies, "missing", "key", "val")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestAnnotate_EmptyKey(t *testing.T) {
	policies := makeAnnotatorPolicies()
	_, err := Annotate(policies, "api-get", "", "value")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestAnnotate_EmptyValue(t *testing.T) {
	policies := makeAnnotatorPolicies()
	_, err := Annotate(policies, "api-get", "owner", "")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestRemoveAnnotation_Success(t *testing.T) {
	policies := makeAnnotatorPolicies()
	policies[0].Annotations = map[string]string{"owner": "team-a", "env": "prod"}
	result, err := RemoveAnnotation(policies, "api-get", "owner")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result[0].Annotations["owner"]; ok {
		t.Error("expected 'owner' annotation to be removed")
	}
	if result[0].Annotations["env"] != "prod" {
		t.Error("expected 'env' annotation to remain")
	}
}

func TestRemoveAnnotation_KeyNotFound(t *testing.T) {
	policies := makeAnnotatorPolicies()
	policies[0].Annotations = map[string]string{"env": "prod"}
	_, err := RemoveAnnotation(policies, "api-get", "missing-key")
	if err == nil {
		t.Fatal("expected error for missing annotation key")
	}
}

func TestGetAnnotations_ReturnsCopy(t *testing.T) {
	policies := makeAnnotatorPolicies()
	policies[0].Annotations = map[string]string{"owner": "team-a"}
	anns, err := GetAnnotations(policies, "api-get")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	anns["owner"] = "mutated"
	if policies[0].Annotations["owner"] != "team-a" {
		t.Error("expected original annotations to be unmodified")
	}
}

func TestGetAnnotations_EmptyWhenNil(t *testing.T) {
	policies := makeAnnotatorPolicies()
	anns, err := GetAnnotations(policies, "api-get")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(anns) != 0 {
		t.Errorf("expected empty annotations, got %v", anns)
	}
}
