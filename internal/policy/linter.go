package policy

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning or error for a policy.
type LintIssue struct {
	PolicyName string
	Field      string
	Message    string
	Severity   string // "warning" or "error"
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s -> %s: %s", i.Severity, i.PolicyName, i.Field, i.Message)
}

// Lint inspects a slice of policies and returns a list of issues found.
func Lint(policies []Policy) []LintIssue {
	var issues []LintIssue

	seen := make(map[string]int)

	for _, p := range policies {
		// Warn on very high limits
		if p.Limit > 10000 {
			issues = append(issues, LintIssue{
				PolicyName: p.Name,
				Field:      "limit",
				Message:    fmt.Sprintf("limit %d is unusually high; consider lowering for local dev", p.Limit),
				Severity:   "warning",
			})
		}

		// Warn on very short windows
		if p.Window > 0 && p.Window < 100 {
			issues = append(issues, LintIssue{
				PolicyName: p.Name,
				Field:      "window",
				Message:    fmt.Sprintf("window %dms is very short and may cause excessive rejections", p.Window),
				Severity:   "warning",
			})
		}

		// Warn on wildcard method with low limit
		if p.Method == "*" && p.Limit < 5 {
			issues = append(issues, LintIssue{
				PolicyName: p.Name,
				Field:      "method",
				Message:    "wildcard method with limit < 5 may block legitimate traffic",
				Severity:   "warning",
			})
		}

		// Error on duplicate name
		normName := strings.ToLower(strings.TrimSpace(p.Name))
		seen[normName]++
		if seen[normName] == 2 {
			issues = append(issues, LintIssue{
				PolicyName: p.Name,
				Field:      "name",
				Message:    fmt.Sprintf("duplicate policy name %q detected", p.Name),
				Severity:   "error",
			})
		}
	}

	return issues
}

// HasErrors returns true if any of the issues have severity "error".
func HasErrors(issues []LintIssue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}
