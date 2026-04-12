package policy

import (
	"net/http"
	"strings"
)

// Match checks whether a given HTTP request matches the policy's
// endpoint and method constraints.
func Match(p Policy, r *http.Request) bool {
	if !matchEndpoint(p.Endpoint, r.URL.Path) {
		return false
	}
	if !matchMethod(p.Method, r.Method) {
		return false
	}
	return true
}

// matchEndpoint returns true if the request path matches the policy endpoint.
// A trailing wildcard "*" allows prefix matching.
func matchEndpoint(pattern, path string) bool {
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return pattern == path
}

// matchMethod returns true if the request method matches the policy method.
// An empty or wildcard policy method matches any HTTP method.
func matchMethod(policyMethod, requestMethod string) bool {
	if policyMethod == "" || policyMethod == "*" {
		return true
	}
	return strings.EqualFold(policyMethod, requestMethod)
}

// FindMatch returns the first policy in the slice that matches the request,
// or nil if no policy matches.
func FindMatch(policies []Policy, r *http.Request) *Policy {
	for i := range policies {
		if Match(policies[i], r) {
			return &policies[i]
		}
	}
	return nil
}
