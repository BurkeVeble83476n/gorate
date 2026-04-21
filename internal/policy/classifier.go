package policy

import (
	"fmt"
	"strings"
)

// ClassifyResult holds the classification label and reasoning for a policy.
type ClassifyResult struct {
	PolicyName string
	Class      string
	Reason     string
}

// Classify assigns a semantic class to each policy based on its limit, window, and method.
// Classes: "strict", "moderate", "permissive", "wildcard-strict", "wildcard-permissive"
func Classify(policies []Policy) []ClassifyResult {
	results := make([]ClassifyResult, 0, len(policies))
	for _, p := range policies {
		results = append(results, classifyOne(p))
	}
	return results
}

func classifyOne(p Policy) ClassifyResult {
	isWildcard := p.Method == "*" || strings.EqualFold(p.Method, "any")
	ratePerSecond := float64(p.Limit) / p.Window.Seconds()

	var class, reason string
	switch {
	case isWildcard && ratePerSecond < 1.0:
		class = "wildcard-strict"
		reason = fmt.Sprintf("wildcard method with low rate (%.2f req/s)", ratePerSecond)
	case isWildcard:
		class = "wildcard-permissive"
		reason = fmt.Sprintf("wildcard method with high rate (%.2f req/s)", ratePerSecond)
	case ratePerSecond < 0.5:
		class = "strict"
		reason = fmt.Sprintf("very low rate (%.2f req/s)", ratePerSecond)
	case ratePerSecond < 5.0:
		class = "moderate"
		reason = fmt.Sprintf("moderate rate (%.2f req/s)", ratePerSecond)
	default:
		class = "permissive"
		reason = fmt.Sprintf("high rate (%.2f req/s)", ratePerSecond)
	}

	return ClassifyResult{
		PolicyName: p.Name,
		Class:      class,
		Reason:     reason,
	}
}

// FormatClassify returns a human-readable table of classification results.
func FormatClassify(results []ClassifyResult) string {
	if len(results) == 0 {
		return "no policies to classify\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %-22s %s\n", "POLICY", "CLASS", "REASON"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-24s %-22s %s\n", r.PolicyName, r.Class, r.Reason))
	}
	return sb.String()
}
