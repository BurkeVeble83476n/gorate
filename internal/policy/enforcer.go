package policy

import "fmt"

// EnforcementMode controls how violations are handled.
type EnforcementMode string

const (
	EnforceModeStrict  EnforcementMode = "strict"
	EnforceModeWarn    EnforcementMode = "warn"
	EnforceModeDisable EnforcementMode = "disable"
)

// EnforcementResult describes the outcome of enforcing a single policy.
type EnforcementResult struct {
	PolicyName string
	Mode       EnforcementMode
	Violations []string
	Enforced   bool
}

// Enforce evaluates each policy against its enforcement mode and returns results.
// Strict mode marks the policy as not enforced if violations exist.
// Warn mode records violations but still marks the policy as enforced.
// Disable mode skips enforcement entirely.
func Enforce(policies []Policy, mode EnforcementMode) []EnforcementResult {
	results := make([]EnforcementResult, 0, len(policies))
	for _, p := range policies {
		result := EnforcementResult{
			PolicyName: p.Name,
			Mode:       mode,
		}
		if mode == EnforceModeDisable {
			result.Enforced = true
			results = append(results, result)
			continue
		}
		result.Violations = collectViolations(p)
		switch mode {
		case EnforceModeStrict:
			result.Enforced = len(result.Violations) == 0
		case EnforceModeWarn:
			result.Enforced = true
		}
		results = append(results, result)
	}
	return results
}

// ParseEnforcementMode parses a string into an EnforcementMode.
func ParseEnforcementMode(s string) (EnforcementMode, error) {
	switch EnforcementMode(s) {
	case EnforceModeStrict, EnforceModeWarn, EnforceModeDisable:
		return EnforcementMode(s), nil
	}
	return "", fmt.Errorf("unknown enforcement mode %q: must be strict, warn, or disable", s)
}

// FormatEnforcement returns a human-readable summary of enforcement results.
func FormatEnforcement(results []EnforcementResult) string {
	if len(results) == 0 {
		return "No policies evaluated.\n"
	}
	out := ""
	for _, r := range results {
		status := "OK"
		if !r.Enforced {
			status = "BLOCKED"
		} else if len(r.Violations) > 0 {
			status = "WARN"
		}
		out += fmt.Sprintf("[%s] %s (mode: %s)\n", status, r.PolicyName, r.Mode)
		for _, v := range r.Violations {
			out += fmt.Sprintf("  - %s\n", v)
		}
	}
	return out
}

func collectViolations(p Policy) []string {
	var violations []string
	if p.Limit > 10000 {
		violations = append(violations, fmt.Sprintf("limit %d exceeds maximum allowed 10000", p.Limit))
	}
	if p.Window < 1 {
		violations = append(violations, "window must be at least 1 second")
	}
	if p.Method == "*" && p.Limit > 5000 {
		violations = append(violations, "wildcard method with limit >5000 is discouraged")
	}
	return violations
}
