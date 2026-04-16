package policy

import (
	"fmt"
	"strings"
)

// Dependency represents a directional dependency between two policies.
type Dependency struct {
	From string
	To   string
}

// DependencyGraph holds the dependency relationships between policies.
type DependencyGraph struct {
	Edges []Dependency
}

// AddDependency adds a dependency from one policy to another.
func AddDependency(policies []Policy, from, to string) ([]Policy, error) {
	if from == "" || to == "" {
		return nil, fmt.Errorf("both 'from' and 'to' policy names are required")
	}
	if from == to {
		return nil, fmt.Errorf("a policy cannot depend on itself")
	}
	fromFound, toFound := false, false
	for _, p := range policies {
		if p.Name == from {
			fromFound = true
		}
		if p.Name == to {
			toFound = true
		}
	}
	if !fromFound {
		return nil, fmt.Errorf("policy %q not found", from)
	}
	if !toFound {
		return nil, fmt.Errorf("policy %q not found", to)
	}
	updated := make([]Policy, len(policies))
	copy(updated, policies)
	for i, p := range updated {
		if p.Name == from {
			if p.Annotations == nil {
				updated[i].Annotations = map[string]string{}
			}
			existing := updated[i].Annotations["depends_on"]
			deps := splitDeps(existing)
			for _, d := range deps {
				if d == to {
					return updated, nil
				}
			}
			deps = append(deps, to)
			updated[i].Annotations["depends_on"] = strings.Join(deps, ",")
			break
		}
	}
	return updated, nil
}

// GetDependencies returns the list of policy names that the given policy depends on.
func GetDependencies(policies []Policy, name string) ([]string, error) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return []string{}, nil
			}
			return splitDeps(p.Annotations["depends_on"]), nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// BuildGraph constructs a DependencyGraph from all policies.
func BuildGraph(policies []Policy) DependencyGraph {
	var edges []Dependency
	for _, p := range policies {
		if p.Annotations == nil {
			continue
		}
		for _, dep := range splitDeps(p.Annotations["depends_on"]) {
			if dep != "" {
				edges = append(edges, Dependency{From: p.Name, To: dep})
			}
		}
	}
	return DependencyGraph{Edges: edges}
}

func splitDeps(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
