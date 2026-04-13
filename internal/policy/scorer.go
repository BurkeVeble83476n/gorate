package policy

import (
	"fmt"
	"strings"
)

// PolicyScore holds the computed score and breakdown for a policy.
type PolicyScore struct {
	Name       string
	Score      int
	Breakdown  []string
}

// Score computes a quality/risk score for each policy.
// Higher scores indicate more permissive or potentially risky configurations.
func Score(policies []Policy) []PolicyScore {
	results := make([]PolicyScore, 0, len(policies))
	for _, p := range policies {
		results = append(results, scoreOne(p))
	}
	return results
}

func scoreOne(p Policy) PolicyScore {
	ps := PolicyScore{Name: p.Name}

	// High limit increases risk score
	if p.Limit > 1000 {
		ps.Score += 40
		ps.Breakdown = append(ps.Breakdown, fmt.Sprintf("high limit (%d > 1000): +40", p.Limit))
	} else if p.Limit > 500 {
		ps.Score += 20
		ps.Breakdown = append(ps.Breakdown, fmt.Sprintf("elevated limit (%d > 500): +20", p.Limit))
	}

	// Wildcard method is more permissive
	if strings.ToUpper(p.Method) == "*" {
		ps.Score += 20
		ps.Breakdown = append(ps.Breakdown, "wildcard method (*): +20")
	}

	// Very short window with high limit is risky
	if p.Window <= 5 && p.Limit > 100 {
		ps.Score += 25
		ps.Breakdown = append(ps.Breakdown, fmt.Sprintf("short window (%ds) with high limit (%d): +25", p.Window, p.Limit))
	}

	// Wildcard or root endpoint
	if p.Endpoint == "/*" || p.Endpoint == "/" {
		ps.Score += 15
		ps.Breakdown = append(ps.Breakdown, fmt.Sprintf("broad endpoint (%s): +15", p.Endpoint))
	}

	return ps
}

// FormatScores returns a human-readable summary of policy scores.
func FormatScores(scores []PolicyScore) string {
	var sb strings.Builder
	for _, s := range scores {
		sb.WriteString(fmt.Sprintf("[%s] score=%d\n", s.Name, s.Score))
		for _, b := range s.Breakdown {
			sb.WriteString(fmt.Sprintf("  - %s\n", b))
		}
	}
	return sb.String()
}
