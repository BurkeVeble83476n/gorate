package policy

import (
	"strings"
	"testing"
)

func makeRecommenderPolicies() []Policy {
	return []Policy{
		{Name: "safe", Endpoint: "/api/users", Method: "GET", Limit: 100, Window: 60},
		{Name: "high-limit", Endpoint: "/api/data", Method: "POST", Limit: 2000, Window: 30},
		{Name: "short-window", Endpoint: "/api/ping", Method: "GET", Limit: 50, Window: 2},
		{Name: "wildcard-method", Endpoint: "/api/bulk", Method: "*", Limit: 600, Window: 60},
		{Name: "wildcard-endpoint", Endpoint: "*", Method: "GET", Limit: 100, Window: 30},
	}
}

func TestRecommend_NoIssues(t *testing.T) {
	policies := []Policy{
		{Name: "safe", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
	}
	recs := Recommend(policies)
	if len(recs) != 0 {
		t.Errorf("expected 0 recommendations, got %d", len(recs))
	}
}

func TestRecommend_HighLimit(t *testing.T) {
	policies := []Policy{
		{Name: "high-limit", Endpoint: "/api/data", Method: "POST", Limit: 2000, Window: 30},
	}
	recs := Recommend(policies)
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Field != "limit" {
		t.Errorf("expected field=limit, got %s", recs[0].Field)
	}
}

func TestRecommend_ShortWindow(t *testing.T) {
	policies := []Policy{
		{Name: "short-window", Endpoint: "/api/ping", Method: "GET", Limit: 50, Window: 2},
	}
	recs := Recommend(policies)
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Field != "window" {
		t.Errorf("expected field=window, got %s", recs[0].Field)
	}
}

func TestRecommend_WildcardMethodHighLimit(t *testing.T) {
	policies := []Policy{
		{Name: "wildcard-method", Endpoint: "/api/bulk", Method: "*", Limit: 600, Window: 60},
	}
	recs := Recommend(policies)
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Field != "method" {
		t.Errorf("expected field=method, got %s", recs[0].Field)
	}
}

func TestRecommend_WildcardEndpoint(t *testing.T) {
	policies := []Policy{
		{Name: "wildcard-endpoint", Endpoint: "*", Method: "GET", Limit: 100, Window: 30},
	}
	recs := Recommend(policies)
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}
	if recs[0].Field != "endpoint" {
		t.Errorf("expected field=endpoint, got %s", recs[0].Field)
	}
}

func TestRecommend_MultiplePolicies(t *testing.T) {
	policies := makeRecommenderPolicies()
	recs := Recommend(policies)
	// high-limit(1) + short-window(1) + wildcard-method(1) + wildcard-endpoint(1) = 4
	if len(recs) != 4 {
		t.Errorf("expected 4 recommendations, got %d", len(recs))
	}
}

func TestFormatRecommendations_Empty(t *testing.T) {
	out := FormatRecommendations(nil)
	if !strings.Contains(out, "No recommendations") {
		t.Errorf("expected no-recommendation message, got: %s", out)
	}
}

func TestFormatRecommendations_ContainsPolicyName(t *testing.T) {
	recs := []Recommendation{
		{PolicyName: "my-policy", Field: "limit", Current: "2000", Suggested: "<= 1000", Reason: "too high"},
	}
	out := FormatRecommendations(recs)
	if !strings.Contains(out, "my-policy") {
		t.Errorf("expected policy name in output, got: %s", out)
	}
	if !strings.Contains(out, "1 recommendation") {
		t.Errorf("expected count in output, got: %s", out)
	}
}
