package policy

import (
	"testing"
)

func makeLintPolicies() []Policy {
	return []Policy{
		{Name: "normal", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 1000},
		{Name: "strict", Endpoint: "/api/v2", Method: "POST", Limit: 50, Window: 500},
	}
}

func TestLint_NoIssues(t *testing.T) {
	policies := makeLintPolicies()
	issues := Lint(policies)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_HighLimit(t *testing.T) {
	policies := []Policy{
		{Name: "heavy", Endpoint: "/flood", Method: "GET", Limit: 99999, Window: 1000},
	}
	issues := Lint(policies)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Field != "limit" || issues[0].Severity != "warning" {
		t.Errorf("unexpected issue: %v", issues[0])
	}
}

func TestLint_ShortWindow(t *testing.T) {
	policies := []Policy{
		{Name: "fast", Endpoint: "/ping", Method: "GET", Limit: 10, Window: 50},
	}
	issues := Lint(policies)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Field != "window" {
		t.Errorf("expected window issue, got field=%s", issues[0].Field)
	}
}

func TestLint_WildcardMethodLowLimit(t *testing.T) {
	policies := []Policy{
		{Name: "wildlow", Endpoint: "/api", Method: "*", Limit: 2, Window: 1000},
	}
	issues := Lint(policies)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Field != "method" {
		t.Errorf("expected method issue, got field=%s", issues[0].Field)
	}
}

func TestLint_DuplicateName(t *testing.T) {
	policies := []Policy{
		{Name: "dup", Endpoint: "/a", Method: "GET", Limit: 10, Window: 1000},
		{Name: "dup", Endpoint: "/b", Method: "POST", Limit: 20, Window: 1000},
	}
	issues := Lint(policies)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestHasErrors_WithError(t *testing.T) {
	issues := []LintIssue{
		{PolicyName: "p", Field: "name", Message: "dup", Severity: "error"},
	}
	if !HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_OnlyWarnings(t *testing.T) {
	issues := []LintIssue{
		{PolicyName: "p", Field: "limit", Message: "high", Severity: "warning"},
	}
	if HasErrors(issues) {
		t.Error("expected HasErrors to return false")
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{PolicyName: "mypolicy", Field: "limit", Message: "too high", Severity: "warning"}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
