package policy

import (
	"testing"
	"time"
)

func TestValidate_ValidPolicy(t *testing.T) {
	p := &Policy{
		Name:     "test-policy",
		Endpoint: "/api/v1/resource",
		Method:   "GET",
		Limit:    100,
		Window:   time.Minute,
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_DefaultsMethod(t *testing.T) {
	p := &Policy{
		Name:     "test-policy",
		Endpoint: "/api/v1/resource",
		Limit:    10,
		Window:   time.Second * 30,
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if p.Method != "*" {
		t.Errorf("expected method to default to '*', got: %s", p.Method)
	}
}

func TestValidate_MissingName(t *testing.T) {
	p := &Policy{Endpoint: "/api", Limit: 10, Window: time.Second}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestValidate_MissingEndpoint(t *testing.T) {
	p := &Policy{Name: "p", Limit: 10, Window: time.Second}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestValidate_InvalidLimit(t *testing.T) {
	p := &Policy{Name: "p", Endpoint: "/api", Limit: 0, Window: time.Second}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestValidate_InvalidWindow(t *testing.T) {
	p := &Policy{Name: "p", Endpoint: "/api", Limit: 5, Window: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestRatePerSecond(t *testing.T) {
	p := &Policy{Limit: 60, Window: time.Minute}
	got := p.RatePerSecond()
	if got != 1.0 {
		t.Errorf("expected 1.0 req/s, got: %f", got)
	}
}
