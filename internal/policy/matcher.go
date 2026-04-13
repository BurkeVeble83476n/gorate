package policy

import (
	"net/http"
	"strings"
)

// Match returns true if the request matches the given policy.
func Match(r *http.Request, p Policy) bool {
	return matchEndpoint(r.URL.Path, p.Endpoint) && matchMethod(r.Method, p.Method)
}

// matchEndpoint checks if the request path matches the policy endpoint.
// Supports wildcard suffix matching with "*".
func matchEndpoint(path, endpoint string) bool {
	if endpoint == "*" {
		return true
	}
	if strings.HasSuffix(endpoint, "*") {
		prefix := strings.TrimSuffix(endpoint, "*")
		return strings.HasPrefix(path, prefix)
	}
	return path == endpoint
}

// matchMethod checks if the request method matches the policy method.
// A wildcard "*" matches any method.
func matchMethod(method, policyMethod string) bool {
	if policyMethod == "*" {
		return true
	}
	return strings.EqualFold(method, policyMethod)
}

// FindMatch returns the first policy that matches the given request.
// Returns nil if no policy matches.
func FindMatch(r *http.Request, policies []Policy) *Policy {
	for i := range policies {
		if Match(r, policies[i]) {
			return &policies[i]
		}
	}
	return nil
}
