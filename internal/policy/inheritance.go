package policy

import (
	"fmt"
	"strings"
)

// Policy is assumed to be defined elsewhere in the package.
// Inheritance allows a policy to extend another, inheriting its fields.

// InheritanceResult holds the resolved policy and a note about what was inherited.
type InheritanceResult struct {
	Policy Policy
	Inherited []string
}

// Resolve applies inheritance: child fields override parent fields.
// If a child field is zero-valued, the parent's value is used.
func Resolve(parent, child Policy) (InheritanceResult, error) {
	if parent.Name == "" {
		return InheritanceResult{}, fmt.Errorf("parent policy must have a name")
	}
	if child.Name == "" {
		return InheritanceResult{}, fmt.Errorf("child policy must have a name")
	}

	result := child
	var inherited []string

	if result.Endpoint == "" && parent.Endpoint != "" {
		result.Endpoint = parent.Endpoint
		inherited = append(inherited, "endpoint")
	}
	if result.Method == "" && parent.Method != "" {
		result.Method = parent.Method
		inherited = append(inherited, "method")
	}
	if result.Limit == 0 && parent.Limit != 0 {
		result.Limit = parent.Limit
		inherited = append(inherited, "limit")
	}
	if result.Window == "" && parent.Window != "" {
		result.Window = parent.Window
		inherited = append(inherited, "window")
	}

	return InheritanceResult{Policy: result, Inherited: inherited}, nil
}

// FormatInheritance returns a human-readable description of what was inherited.
func FormatInheritance(r InheritanceResult, parentName string) string {
	if len(r.Inherited) == 0 {
		return fmt.Sprintf("policy %q inherits nothing from %q", r.Policy.Name, parentName)
	}
	return fmt.Sprintf("policy %q inherits [%s] from %q",
		r.Policy.Name, strings.Join(r.Inherited, ", "), parentName)
}

// ApplyInheritance resolves inheritance for all policies that declare a parent
// via the annotation key "inherits".
func ApplyInheritance(policies []Policy) ([]Policy, []string, error) {
	index := make(map[string]Policy, len(policies))
	for _, p := range policies {
		index[p.Name] = p
	}

	var notes []string
	resolved := make([]Policy, 0, len(policies))

	for _, p := range policies {
		parentName := ""
		if p.Annotations != nil {
			parentName = p.Annotations["inherits"]
		}
		if parentName == "" {
			resolved = append(resolved, p)
			continue
		}
		parent, ok := index[parentName]
		if !ok {
			return nil, nil, fmt.Errorf("policy %q references unknown parent %q", p.Name, parentName)
		}
		r, err := Resolve(parent, p)
		if err != nil {
			return nil, nil, err
		}
		notes = append(notes, FormatInheritance(r, parentName))
		resolved = append(resolved, r.Policy)
	}

	return resolved, notes, nil
}
