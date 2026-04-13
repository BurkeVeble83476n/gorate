package policy

import (
	"strings"
	"testing"
	"time"
)

func makeProfilerPolicies() []Policy {
	return []Policy{
		{
			Name:     "safe-policy",
			Endpoint: "/api/health",
			Method:   "GET",
			Limit:    50,
			Window:   time.Minute,
		},
		{
			Name:     "risky-policy",
			Endpoint: "/api/data",
			Method:   "*",
			Limit:    10000,
			Window:   time.Second,
		},
	}
}

func TestProfile_ReturnsSameCount(t *testing.T) {
	policies := makeProfilerPolicies()
	profiles := Profile(policies)
	if len(profiles) != len(policies) {
		t.Errorf("expected %d profiles, got %d", len(policies), len(profiles))
	}
}

func TestProfile_RiskyPolicyFirst(t *testing.T) {
	policies := makeProfilerPolicies()
	profiles := Profile(policies)
	if profiles[0].Name != "risky-policy" {
		t.Errorf("expected risky-policy first, got %s", profiles[0].Name)
	}
}

func TestProfile_IssueCountPopulated(t *testing.T) {
	policies := makeProfilerPolicies()
	profiles := Profile(policies)
	found := false
	for _, p := range profiles {
		if p.Name == "risky-policy" && p.IssueCount > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected risky-policy to have lint issues")
	}
}

func TestProfile_FieldsPopulated(t *testing.T) {
	policies := makeProfilerPolicies()
	profiles := Profile(policies)
	for _, p := range profiles {
		if p.Name == "" || p.Endpoint == "" || p.Method == "" {
			t.Errorf("profile has empty required fields: %+v", p)
		}
	}
}

func TestFormatProfile_ContainsName(t *testing.T) {
	pp := PolicyProfile{
		Name:      "test-policy",
		Endpoint:  "/test",
		Method:    "GET",
		Limit:     100,
		Window:    time.Minute,
		RiskScore: 10,
		Notes:     []string{"high limit"},
	}
	out := FormatProfile(pp)
	if !strings.Contains(out, "test-policy") {
		t.Error("expected output to contain policy name")
	}
	if !strings.Contains(out, "high limit") {
		t.Error("expected output to contain issue note")
	}
}

func TestFormatProfile_NoIssues(t *testing.T) {
	pp := PolicyProfile{
		Name:      "clean",
		Endpoint:  "/ok",
		Method:    "GET",
		Limit:     10,
		Window:    time.Minute,
		RiskScore: 2,
	}
	out := FormatProfile(pp)
	if strings.Contains(out, "Issues") {
		t.Error("expected no issues section for clean policy")
	}
}
