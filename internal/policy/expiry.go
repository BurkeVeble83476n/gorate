package policy

import (
	"fmt"
	"time"
)

const expiryAnnotationKey = "expiry"

// ExpiryResult holds the expiry evaluation result for a single policy.
type ExpiryResult struct {
	Name      string
	ExpiresAt time.Time
	Expired   bool
}

// SetExpiry sets an expiry timestamp on the named policy via its annotations.
func SetExpiry(policies []Policy, name string, at time.Time) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = make(map[string]string)
			}
			policies[i].Annotations[expiryAnnotationKey] = at.UTC().Format(time.RFC3339)
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// RemoveExpiry clears the expiry annotation from the named policy.
func RemoveExpiry(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			delete(policies[i].Annotations, expiryAnnotationKey)
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetExpiry returns the expiry time for a named policy, and whether one is set.
func GetExpiry(policies []Policy, name string) (time.Time, bool) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return time.Time{}, false
			}
			raw, ok := p.Annotations[expiryAnnotationKey]
			if !ok {
				return time.Time{}, false
			}
			t, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return time.Time{}, false
			}
			return t, true
		}
	}
	return time.Time{}, false
}

// EvaluateExpiry returns expiry results for all policies that have an expiry set.
func EvaluateExpiry(policies []Policy) []ExpiryResult {
	now := time.Now()
	var results []ExpiryResult
	for _, p := range policies {
		if p.Annotations == nil {
			continue
		}
		raw, ok := p.Annotations[expiryAnnotationKey]
		if !ok {
			continue
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			continue
		}
		results = append(results, ExpiryResult{
			Name:      p.Name,
			ExpiresAt: t,
			Expired:   now.After(t),
		})
	}
	return results
}

// FilterExpired returns only policies that have not yet expired.
func FilterExpired(policies []Policy) []Policy {
	now := time.Now()
	var out []Policy
	for _, p := range policies {
		if p.Annotations != nil {
			if raw, ok := p.Annotations[expiryAnnotationKey]; ok {
				if t, err := time.Parse(time.RFC3339, raw); err == nil && now.After(t) {
					continue
				}
			}
		}
		out = append(out, p)
	}
	return out
}
