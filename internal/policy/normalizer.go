package policy

import (
	"strings"

	"github.com/isoment/gorate/internal/policy"
)

// NormalizeResult holds the result of normalizing a set of policies.
type NormalizeResult struct {
	Normalized []Policy
	Changes    []string
}

// Normalize applies consistent formatting and casing rules to all policies in the slice.
// It trims whitespace, lowercases method names, and cleans endpoint paths.
func Normalize(policies []Policy) NormalizeResult {
	result := NormalizeResult{
		Normalized: make([]Policy, 0, len(policies)),
		Changes:    []string{},
	}

	for _, p := range policies {
		original := p
		changed := false

		// Trim and normalize name
		trimmedName := strings.TrimSpace(p.Name)
		if trimmedName != p.Name {
			p.Name = trimmedName
			changed = true
		}

		// Uppercase method
		upperMethod := strings.ToUpper(strings.TrimSpace(p.Method))
		if upperMethod != p.Method {
			p.Method = upperMethod
			changed = true
		}

		// Normalize endpoint: trim spaces, ensure leading slash
		normEndpoint := strings.TrimSpace(p.Endpoint)
		if normEndpoint != "" && !strings.HasPrefix(normEndpoint, "/") {
			normEndpoint = "/" + normEndpoint
			changed = true
		}
		if normEndpoint != p.Endpoint {
			p.Endpoint = normEndpoint
			changed = true
		}

		if changed {
			result.Changes = append(result.Changes,
				formatNormalizeChange(original, p),
			)
		}

		result.Normalized = append(result.Normalized, p)
	}

	return result
}

func formatNormalizeChange(before, after Policy) string {
	return "normalized policy '" + after.Name + "': " +
		"method='" + before.Method + "'->" + "'" + after.Method + "' " +
		"endpoint='" + before.Endpoint + "'->" + "'" + after.Endpoint + "'"
}
