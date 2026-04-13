package policy

import (
	"strings"
	"testing"
)

func TestValidateEndpoint_Valid(t *testing.T) {
	cases := []string{"/api/v1", "/health", "*", "/"}
	for _, c := range cases {
		if err := ValidateEndpoint(c); err != nil {
			t.Errorf("expected no error for %q, got: %v", c, err)
		}
	}
}

func TestValidateEndpoint_Empty(t *testing.T) {
	err := ValidateEndpoint("")
	if err == nil {
		t.Fatal("expected error for empty endpoint")
	}
	if !strings.Contains(err.Error(), "must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateMethod_Valid(t *testing.T) {
	cases := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "*", "get", "post"}
	for _, c := range cases {
		if err := ValidateMethod(c); err != nil {
			t.Errorf("expected no error for %q, got: %v", c, err)
		}
	}
}

func TestValidateMethod_Invalid(t *testing.T) {
	cases := []string{"FETCH", "SEND", "", "CONNECT"}
	for _, c := range cases {
		if err := ValidateMethod(c); err == nil {
			t.Errorf("expected error for method %q", c)
		}
	}
}

func TestValidateLimit_Valid(t *testing.T) {
	errs := ValidateLimit(100, 60)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidateLimit_ZeroLimit(t *testing.T) {
	errs := ValidateLimit(0, 60)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !strings.Contains(errs[0], "limit must be > 0") {
		t.Errorf("unexpected error: %s", errs[0])
	}
}

func TestValidateLimit_ZeroWindow(t *testing.T) {
	errs := ValidateLimit(10, 0)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !strings.Contains(errs[0], "window_seconds must be > 0") {
		t.Errorf("unexpected error: %s", errs[0])
	}
}

func TestValidateLimit_BothInvalid(t *testing.T) {
	errs := ValidateLimit(-1, -5)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestValidationError_Message(t *testing.T) {
	e := &ValidationError{
		Name:   "my-policy",
		Errors: []string{"limit must be > 0", "endpoint is empty"},
	}
	msg := e.Error()
	if !strings.Contains(msg, "my-policy") {
		t.Errorf("expected policy name in error, got: %s", msg)
	}
	if !strings.Contains(msg, "2 validation error") {
		t.Errorf("expected error count in message, got: %s", msg)
	}
}
