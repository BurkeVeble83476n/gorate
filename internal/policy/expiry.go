package policy

import (
	"fmt"
	"time"

	"github.com/jdotcurs/gorate/internal/policy"
)

// ExpiryResult holds the result of an expiry evaluation for a single policy.
type ExpiryResult struct {
	PolicyName string
	ExpiresAt  time.Time
	Expired    bool
	TTL        time.Duration
}

// SetExpiry attaches an expiry timestamp to a policy via its annotations.
// The expiry is stored as a RFC3339 string under the key "expiry".
func SetExpiry(policies []Policy, name string, expiresAt time.Time) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = make(map[string]string)
			}
			policies[i].Annotations["expiry"] = expiresAt.UTC().Format(time.RFC3339)
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// RemoveExpiry removes the expiry annotation from a policy.
func RemoveExpiry(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			delete(policies[i].Annotations, "expiry")
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetExpiry returns the expiry time for a named policy.
// Returns an error if the policy is not found or has no expiry set.
func GetExpiry(policies []Policy, name string) (time.Time, error) {
	for _, p := range policies {
		if p.Name == name {
			raw, ok := p.Annotations["expiry"]
			if !ok || raw == "" {
				return time.Time{}, fmt.Errorf("policy %q has no expiry set", name)
			}
			t, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return time.Time{}, fmt.Errorf("policy %q has invalid expiry format: %w", name, err)
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("policy %q not found", name)
}

// EvaluateExpiry checks all policies for expiry annotations and returns
// an ExpiryResult for each policy that has one, relative to now.
func EvaluateExpiry(policies []Policy, now time.Time) []ExpiryResult {
	var results []ExpiryResult
	for _, p := range policies {
		raw, ok := p.Annotations["expiry"]
		if !ok || raw == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			continue
		}
		expired := now.After(t)
		ttl := time.Duration(0)
		if !expired {
			ttl = t.Sub(now).Truncate(time.Second)
		}
		results = append(results, ExpiryResult{
			PolicyName: p.Name,
			ExpiresAt:  t,
			Expired:    expired,
			TTL:        ttl,
		})
	}
	return results
}

// FilterExpired returns only the policies whose expiry annotation has passed
// relative to now. Policies without an expiry annotation are not included.
func FilterExpired(policies []Policy, now time.Time) []Policy {
	var expired []Policy
	for _, p := range policies {
		raw, ok := p.Annotations["expiry"]
		if !ok || raw == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			continue
		}
		if now.After(t) {
			expired = append(expired, p)
		}
	}
	return expired
}
