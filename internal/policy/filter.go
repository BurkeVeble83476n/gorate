package policy

import "strings"

// FilterOptions holds criteria for filtering policies.
type FilterOptions struct {
	Method   string
	Endpoint string
	MinLimit int
}

// Filter returns a subset of policies matching the given options.
// An empty string field is treated as a wildcard (match all).
func Filter(policies []Policy, opts FilterOptions) []Policy {
	var result []Policy
	for _, p := range policies {
		if opts.Method != "" && !strings.EqualFold(p.Method, opts.Method) && p.Method != "*" {
			continue
		}
		if opts.Endpoint != "" && !endpointContains(p.Endpoint, opts.Endpoint) {
			continue
		}
		if opts.MinLimit > 0 && p.Limit < opts.MinLimit {
			continue
		}
		result = append(result, p)
	}
	return result
}

// endpointContains reports whether the policy endpoint matches or contains the query string.
func endpointContains(policyEndpoint, query string) bool {
	return strings.Contains(
		normalizePath(policyEndpoint),
		normalizePath(query),
	)
}

// normalizePath ensures consistent path comparison by trimming trailing slashes
// and lowercasing the value.
func normalizePath(p string) string {
	return strings.ToLower(strings.TrimRight(p, "/"))
}
