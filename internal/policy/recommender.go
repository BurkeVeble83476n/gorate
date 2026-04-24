package policy

import "fmt"

// Recommendation holds a suggestion for improving a policy.
type Recommendation struct {
	PolicyName string
	Field      string
	Current    string
	Suggested  string
	Reason     string
}

// Recommend analyzes a slice of policies and returns improvement suggestions.
func Recommend(policies []Policy) []Recommendation {
	var recs []Recommendation
	for _, p := range policies {
		recs = append(recs, recommendOne(p)...)
	}
	return recs
}

func recommendOne(p Policy) []Recommendation {
	var recs []Recommendation

	if p.Limit > 1000 {
		recs = append(recs, Recommendation{
			PolicyName: p.Name,
			Field:      "limit",
			Current:    fmt.Sprintf("%d", p.Limit),
			Suggested:  "<= 1000",
			Reason:     "very high limits may mask upstream issues during local development",
		})
	}

	if p.Window < 5 {
		recs = append(recs, Recommendation{
			PolicyName: p.Name,
			Field:      "window",
			Current:    fmt.Sprintf("%ds", p.Window),
			Suggested:  ">= 5s",
			Reason:     "very short windows can cause unpredictable rate-limit behaviour",
		})
	}

	if p.Method == "*" && p.Limit > 500 {
		recs = append(recs, Recommendation{
			PolicyName: p.Name,
			Field:      "method",
			Current:    "*",
			Suggested:  "specific HTTP method",
			Reason:     "wildcard method with high limit applies broadly; consider narrowing scope",
		})
	}

	if p.Endpoint == "*" {
		recs = append(recs, Recommendation{
			PolicyName: p.Name,
			Field:      "endpoint",
			Current:    "*",
			Suggested:  "specific path prefix",
			Reason:     "wildcard endpoint matches all routes; a specific path is safer",
		})
	}

	return recs
}

// FormatRecommendations returns a human-readable summary of recommendations.
func FormatRecommendations(recs []Recommendation) string {
	if len(recs) == 0 {
		return "No recommendations — all policies look good.\n"
	}
	out := fmt.Sprintf("%d recommendation(s):\n", len(recs))
	for _, r := range recs {
		out += fmt.Sprintf("  [%s] %s: current=%s suggested=%s — %s\n",
			r.PolicyName, r.Field, r.Current, r.Suggested, r.Reason)
	}
	return out
}
