package policy

import "fmt"

const frozenAnnotationKey = "frozen"

// Freeze marks a policy as frozen, preventing it from being modified.
// A frozen policy has the "frozen" annotation set to "true".
func Freeze(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			updated, err := Annotate(policies, name, frozenAnnotationKey, "true")
			if err != nil {
				return nil, err
			}
			_ = i
			return updated, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// Unfreeze removes the frozen annotation from a policy.
func Unfreeze(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for _, p := range policies {
		if p.Name == name {
			return RemoveAnnotation(policies, name, frozenAnnotationKey)
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// IsFrozen returns true if the named policy has the frozen annotation set.
func IsFrozen(policies []Policy, name string) bool {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return false
			}
			return p.Annotations[frozenAnnotationKey] == "true"
		}
	}
	return false
}

// ListFrozen returns all policies that are currently frozen.
func ListFrozen(policies []Policy) []Policy {
	var frozen []Policy
	for _, p := range policies {
		if p.Annotations != nil && p.Annotations[frozenAnnotationKey] == "true" {
			frozen = append(frozen, p)
		}
	}
	return frozen
}
