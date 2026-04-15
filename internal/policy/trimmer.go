package policy

import (
	"fmt"
	"strings"
)

// TrimResult holds the result of a trim operation on a single policy.
type TrimResult struct {
	Name    string
	Removed bool
	Reason  string
}

// TrimOptions controls which policies are removed during trimming.
type TrimOptions struct {
	RemoveWildcard bool
	MaxLimit       int // 0 means no limit check
	RequireTag     string
}

// Trim removes policies from the list based on the provided options.
// It returns the filtered list and a slice of TrimResults describing what was removed.
func Trim(policies []Policy, opts TrimOptions) ([]Policy, []TrimResult) {
	var kept []Policy
	var results []TrimResult

	for _, p := range policies {
		removed, reason := shouldTrim(p, opts)
		results = append(results, TrimResult{
			Name:    p.Name,
			Removed: removed,
			Reason:  reason,
		})
		if !removed {
			kept = append(kept, p)
		}
	}

	return kept, results
}

func shouldTrim(p Policy, opts TrimOptions) (bool, string) {
	if opts.RemoveWildcard && strings.EqualFold(p.Method, "*") {
		return true, "wildcard method"
	}
	if opts.MaxLimit > 0 && p.Limit > opts.MaxLimit {
		return true, fmt.Sprintf("limit %d exceeds max %d", p.Limit, opts.MaxLimit)
	}
	if opts.RequireTag != "" && !HasTag(p, opts.RequireTag) {
		return true, fmt.Sprintf("missing required tag %q", opts.RequireTag)
	}
	return false, ""
}

// FormatTrimResults returns a human-readable summary of trim results.
func FormatTrimResults(results []TrimResult) string {
	var sb strings.Builder
	removed := 0
	for _, r := range results {
		if r.Removed {
			removed++
			sb.WriteString(fmt.Sprintf("  - removed %q: %s\n", r.Name, r.Reason))
		}
	}
	if removed == 0 {
		return "no policies trimmed\n"
	}
	return fmt.Sprintf("trimmed %d polic(ies):\n", removed) + sb.String()
}
