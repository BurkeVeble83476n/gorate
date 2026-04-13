package policy_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/gorate/internal/policy"
)

func makePolicies() []policy.Policy {
	return []policy.Policy{
		{Name: "alpha", Endpoint: "/api/users", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/orders", Method: "POST", Limit: 20, Window: 30},
	}
}

func TestToSummary_FieldsMatch(t *testing.T) {
	p := policy.Policy{Name: "test", Endpoint: "/x", Method: "PUT", Limit: 5, Window: 10}
	s := policy.ToSummary(p)
	if s.Name != p.Name || s.Endpoint != p.Endpoint || s.Method != p.Method ||
		s.Limit != p.Limit || s.Window != p.Window {
		t.Errorf("ToSummary fields mismatch: got %+v", s)
	}
}

func TestPrintTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	policy.PrintTable(&buf, makePolicies())
	out := buf.String()
	for _, header := range []string{"NAME", "ENDPOINT", "METHOD", "LIMIT", "WINDOW"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestPrintTable_ContainsPolicyData(t *testing.T) {
	var buf bytes.Buffer
	policies := makePolicies()
	policy.PrintTable(&buf, policies)
	out := buf.String()
	for _, p := range policies {
		if !strings.Contains(out, p.Name) {
			t.Errorf("expected policy name %q in output", p.Name)
		}
		if !strings.Contains(out, p.Endpoint) {
			t.Errorf("expected endpoint %q in output", p.Endpoint)
		}
	}
}

func TestPrintTable_EmptyPolicies(t *testing.T) {
	var buf bytes.Buffer
	policy.PrintTable(&buf, []policy.Policy{})
	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Error("expected header row even for empty policies")
	}
}
