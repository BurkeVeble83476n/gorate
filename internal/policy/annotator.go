package policy

import (
	"fmt"
	"strings"
	"time"
)

// Annotation holds a key-value note attached to a policy.
type Annotation struct {
	Key       string    `json:"key" yaml:"key"`
	Value     string    `json:"value" yaml:"value"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

// Annotate adds or updates an annotation on the named policy.
// Returns an error if the policy is not found or key/value are empty.
func Annotate(policies []Policy, name, key, value string) ([]Policy, error) {
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("annotation key must not be empty")
	}
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("annotation value must not be empty")
	}

	for i, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				p.Annotations = make(map[string]string)
			}
			p.Annotations[key] = value
			policies[i] = p
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// RemoveAnnotation removes an annotation key from the named policy.
func RemoveAnnotation(policies []Policy, name, key string) ([]Policy, error) {
	for i, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return nil, fmt.Errorf("policy %q has no annotations", name)
			}
			if _, ok := p.Annotations[key]; !ok {
				return nil, fmt.Errorf("annotation key %q not found on policy %q", key, name)
			}
			delete(p.Annotations, key)
			policies[i] = p
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetAnnotations returns the annotations map for the named policy.
func GetAnnotations(policies []Policy, name string) (map[string]string, error) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return map[string]string{}, nil
			}
			copy := make(map[string]string, len(p.Annotations))
			for k, v := range p.Annotations {
				copy[k] = v
			}
			return copy, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}
