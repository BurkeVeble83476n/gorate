package policy

import (
	"fmt"
	"strings"
)

// AliasKey is the annotation key used to store aliases.
const AliasKey = "alias"

// AddAlias adds an alias to the policy with the given name.
// Returns an error if the policy is not found, the alias is empty,
// or the alias already exists on another policy.
func AddAlias(policies []Policy, name, alias string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return nil, fmt.Errorf("alias must not be empty")
	}

	// Check for alias collision across all policies.
	for _, p := range policies {
		for _, a := range GetAliases(p) {
			if a == alias && p.Name != name {
				return nil, fmt.Errorf("alias %q already used by policy %q", alias, p.Name)
			}
		}
	}

	found := false
	for i, p := range policies {
		if p.Name != name {
			continue
		}
		found = true
		existing := GetAliases(p)
		for _, a := range existing {
			if a == alias {
				return policies, nil // idempotent
			}
		}
		updated := append(existing, alias)
		policies[i] = Annotate(policies, name, AliasKey, strings.Join(updated, ","))[i]
		return policies, nil
	}
	if !found {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	return policies, nil
}

// RemoveAlias removes an alias from the policy with the given name.
func RemoveAlias(policies []Policy, name, alias string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name != name {
			continue
		}
		existing := GetAliases(p)
		filtered := make([]string, 0, len(existing))
		for _, a := range existing {
			if a != alias {
				filtered = append(filtered, a)
			}
		}
		if len(filtered) == 0 {
			policies[i] = RemoveAnnotation(policies, name, AliasKey)[i]
		} else {
			policies[i] = Annotate(policies, name, AliasKey, strings.Join(filtered, ","))[i]
		}
		return policies, nil
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetAliases returns the list of aliases for the given policy.
func GetAliases(p Policy) []string {
	if p.Annotations == nil {
		return nil
	}
	val, ok := p.Annotations[AliasKey]
	if !ok || val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if t := strings.TrimSpace(part); t != "" {
			result = append(result, t)
		}
	}
	return result
}

// FindByAlias returns the policy that has the given alias, or an error if none found.
func FindByAlias(policies []Policy, alias string) (*Policy, error) {
	for _, p := range policies {
		for _, a := range GetAliases(p) {
			if a == alias {
				copy := p
				return &copy, nil
			}
		}
	}
	return nil, fmt.Errorf("no policy found with alias %q", alias)
}
