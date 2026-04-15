package policy

import (
	"fmt"

	"github.com/jdoe/gorate/internal/policy"
)

// Pin marks a policy as pinned, preventing it from being modified by
// automated operations such as merges, deduplication, or normalization.
func Pin(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				policies[i].Annotations = map[string]string{}
			}
			policies[i].Annotations["pinned"] = "true"
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// Unpin removes the pinned annotation from a policy.
func Unpin(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			if p.Annotations != nil {
				delete(policies[i].Annotations, "pinned")
			}
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// IsPinned returns true if the policy with the given name has the pinned annotation.
func IsPinned(policies []Policy, name string) bool {
	for _, p := range policies {
		if p.Name == name {
			return p.Annotations != nil && p.Annotations["pinned"] == "true"
		}
	}
	return false
}

// ListPinned returns all policies that are currently pinned.
func ListPinned(policies []Policy) []Policy {
	var pinned []Policy
	for _, p := range policies {
		if p.Annotations != nil && p.Annotations["pinned"] == "true" {
			pinned = append(pinned, p)
		}
	}
	return pinned
}
