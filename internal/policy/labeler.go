package policy

import "fmt"

// Label represents a simple string label attached to a policy.
type Label = string

// AddLabel adds a label to the given policy by name.
// Returns an error if the policy is not found or the label already exists.
func AddLabel(policies []Policy, name, label string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	if label == "" {
		return nil, fmt.Errorf("label must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			for _, existing := range p.Labels {
				if existing == label {
					return policies, nil // idempotent
				}
			}
			policies[i].Labels = append(policies[i].Labels, label)
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// RemoveLabel removes a label from the given policy by name.
// Returns an error if the policy is not found.
func RemoveLabel(policies []Policy, name, label string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	for i, p := range policies {
		if p.Name == name {
			updated := p.Labels[:0]
			for _, l := range p.Labels {
				if l != label {
					updated = append(updated, l)
				}
			}
			policies[i].Labels = updated
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetLabels returns the labels for the given policy by name.
func GetLabels(policies []Policy, name string) ([]string, error) {
	for _, p := range policies {
		if p.Name == name {
			result := make([]string, len(p.Labels))
			copy(result, p.Labels)
			return result, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// FilterByLabel returns all policies that have the given label.
func FilterByLabel(policies []Policy, label string) []Policy {
	var result []Policy
	for _, p := range policies {
		for _, l := range p.Labels {
			if l == label {
				result = append(result, p)
				break
			}
		}
	}
	return result
}
