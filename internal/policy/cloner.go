package policy

import "fmt"

// CloneOptions configures the behavior of a policy clone operation.
type CloneOptions struct {
	NewName  string
	Override bool
}

// Clone creates a deep copy of the policy with the given name, assigning it
// a new name. Returns an error if the source policy is not found or if a
// policy with the new name already exists (unless Override is set).
func Clone(policies []Policy, sourceName string, opts CloneOptions) ([]Policy, error) {
	if opts.NewName == "" {
		return nil, fmt.Errorf("new name must not be empty")
	}

	var source *Policy
	for i := range policies {
		if policies[i].Name == sourceName {
			source = &policies[i]
			break
		}
	}
	if source == nil {
		return nil, fmt.Errorf("policy %q not found", sourceName)
	}

	for _, p := range policies {
		if p.Name == opts.NewName {
			if !opts.Override {
				return nil, fmt.Errorf("policy %q already exists; use --override to replace it", opts.NewName)
			}
			// Remove the existing policy so we can append the clone.
			policies = removePolicyByName(policies, opts.NewName)
			break
		}
	}

	cloned := deepCopyPolicy(*source)
	cloned.Name = opts.NewName

	return append(policies, cloned), nil
}

func deepCopyPolicy(p Policy) Policy {
	copy := p
	if p.Tags != nil {
		copy.Tags = make([]string, len(p.Tags))
		for i, t := range p.Tags {
			copy.Tags[i] = t
		}
	}
	return copy
}

func removePolicyByName(policies []Policy, name string) []Policy {
	result := make([]Policy, 0, len(policies))
	for _, p := range policies {
		if p.Name != name {
			result = append(result, p)
		}
	}
	return result
}
