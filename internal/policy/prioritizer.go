package policy

import (
	"fmt"
	"sort"
)

// Priority levels for policies.
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

var priorityRank = map[string]int{
	PriorityHigh:   3,
	PriorityMedium: 2,
	PriorityLow:    1,
}

// SetPriority assigns a priority level to a named policy via its annotations.
func SetPriority(policies []Policy, name, level string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	if _, ok := priorityRank[level]; !ok {
		return nil, fmt.Errorf("invalid priority level %q: must be high, medium, or low", level)
	}
	found := false
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = map[string]string{}
			}
			policies[i].Annotations["priority"] = level
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	return policies, nil
}

// GetPriority returns the priority level of a named policy.
func GetPriority(policies []Policy, name string) (string, error) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations != nil {
				if lvl, ok := p.Annotations["priority"]; ok {
					return lvl, nil
				}
			}
			return PriorityMedium, nil
		}
	}
	return "", fmt.Errorf("policy %q not found", name)
}

// SortByPriority returns a copy of policies sorted by priority descending (high first).
func SortByPriority(policies []Policy) []Policy {
	result := make([]Policy, len(policies))
	copy(result, policies)
	sort.SliceStable(result, func(i, j int) bool {
		pi := priorityOf(result[i])
		pj := priorityOf(result[j])
		return pi > pj
	})
	return result
}

func priorityOf(p Policy) int {
	if p.Annotations != nil {
		if lvl, ok := p.Annotations["priority"]; ok {
			if rank, ok := priorityRank[lvl]; ok {
				return rank
			}
		}
	}
	return priorityRank[PriorityMedium]
}
