package policy

import (
	"strings"
	"testing"
)

func makeScorerPolicies() []Policy {
	return []Policy{
		{Name: "low-risk", Endpoint: "/api/data", Method: "GET", Limit: 100, Window: 60},
		{Name: "high-limit", Endpoint: "/api/bulk", Method: "POST", Limit: 1500, Window: 60},
		{Name: "wildcard-method", Endpoint: "/api/resource", Method: "*", Limit: 200, Window: 30},
		{Name: "short-window-high-limit", Endpoint: "/api/fast", Method: "GET", Limit: 500, Window: 3},
		{Name: "root-endpoint", Endpoint: "/", Method: "GET", Limit: 50, Window: 60},
	}
}

func TestScore_LowRiskPolicy(t *testing.T) {
	policies := makeScorerPolicies()
	scores := Score(policies[:1])
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].Score != 0 {
		t.Errorf("expected score 0 for low-risk policy, got %d", scores[0].Score)
	}
	if len(scores[0].Breakdown) != 0 {
		t.Errorf("expected no breakdown for low-risk policy")
	}
}

func TestScore_HighLimitPolicy(t *testing.T) {
	p := Policy{Name: "high-limit", Endpoint: "/api/bulk", Method: "POST", Limit: 1500, Window: 60}
	scores := Score([]Policy{p})
	if scores[0].Score < 40 {
		t.Errorf("expected score >= 40 for high limit, got %d", scores[0].Score)
	}
}

func TestScore_WildcardMethod(t *testing.T) {
	p := Policy{Name: "wildcard", Endpoint: "/api/x", Method: "*", Limit: 10, Window: 60}
	scores := Score([]Policy{p})
	if scores[0].Score < 20 {
		t.Errorf("expected score >= 20 for wildcard method, got %d", scores[0].Score)
	}
}

func TestScore_ShortWindowHighLimit(t *testing.T) {
	p := Policy{Name: "fast", Endpoint: "/api/fast", Method: "GET", Limit: 500, Window: 3}
	scores := Score([]Policy{p})
	found := false
	for _, b := range scores[0].Breakdown {
		if strings.Contains(b, "short window") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'short window' in breakdown")
	}
}

func TestScore_RootEndpoint(t *testing.T) {
	p := Policy{Name: "root", Endpoint: "/", Method: "GET", Limit: 50, Window: 60}
	scores := Score([]Policy{p})
	if scores[0].Score < 15 {
		t.Errorf("expected score >= 15 for root endpoint, got %d", scores[0].Score)
	}
}

func TestFormatScores_ContainsName(t *testing.T) {
	policies := makeScorerPolicies()
	scores := Score(policies)
	out := FormatScores(scores)
	for _, p := range policies {
		if !strings.Contains(out, p.Name) {
			t.Errorf("expected output to contain policy name %q", p.Name)
		}
	}
}

func TestScore_MultipleRiskFactors(t *testing.T) {
	p := Policy{Name: "max-risk", Endpoint: "/*", Method: "*", Limit: 2000, Window: 2}
	scores := Score([]Policy{p})
	if scores[0].Score < 80 {
		t.Errorf("expected score >= 80 for max-risk policy, got %d", scores[0].Score)
	}
}
