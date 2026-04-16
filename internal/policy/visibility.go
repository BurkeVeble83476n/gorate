package policy

import "fmt"

const (
	VisibilityPublic   = "public"
	VisibilityInternal = "internal"
	VisibilityPrivate  = "private"
)

var validVisibilities = map[string]bool{
	VisibilityPublic:   true,
	VisibilityInternal: true,
	VisibilityPrivate:  true,
}

// SetVisibility sets the visibility annotation on a named policy.
func SetVisibility(policies []Policy, name, visibility string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	if !validVisibilities[visibility] {
		return nil, fmt.Errorf("invalid visibility %q: must be public, internal, or private", visibility)
	}
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = map[string]string{}
			}
			policies[i].Annotations["visibility"] = visibility
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetVisibility returns the visibility of a named policy.
func GetVisibility(policies []Policy, name string) (string, error) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations != nil {
				if v, ok := p.Annotations["visibility"]; ok {
					return v, nil
				}
			}
			return "", nil
		}
	}
	return "", fmt.Errorf("policy %q not found", name)
}

// FilterByVisibility returns policies matching the given visibility.
func FilterByVisibility(policies []Policy, visibility string) []Policy {
	var out []Policy
	for _, p := range policies {
		v := ""
		if p.Annotations != nil {
			v = p.Annotations["visibility"]
		}
		if v == visibility {
			out = append(out, p)
		}
	}
	return out
}
