package policy

import (
	"fmt"
	"sort"
)

// GroupBy defines the field to group policies by.
type GroupBy string

const (
	GroupByMethod   GroupBy = "method"
	GroupByEndpoint GroupBy = "endpoint"
	GroupByWindow   GroupBy = "window"
)

// PolicyGroup holds a group key and the policies belonging to it.
type PolicyGroup struct {
	Key      string
	Policies []Policy
}

// Group partitions policies by the given field and returns sorted groups.
func Group(policies []Policy, by GroupBy) ([]PolicyGroup, error) {
	if err := validateGroupBy(by); err != nil {
		return nil, err
	}

	buckets := make(map[string][]Policy)
	for _, p := range policies {
		key := groupKey(p, by)
		buckets[key] = append(buckets[key], p)
	}

	groups := make([]PolicyGroup, 0, len(buckets))
	for k, v := range buckets {
		groups = append(groups, PolicyGroup{Key: k, Policies: v})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})

	return groups, nil
}

func groupKey(p Policy, by GroupBy) string {
	switch by {
	case GroupByMethod:
		return p.Method
	case GroupByEndpoint:
		return p.Endpoint
	case GroupByWindow:
		return p.Window
	default:
		return ""
	}
}

func validateGroupBy(by GroupBy) error {
	switch by {
	case GroupByMethod, GroupByEndpoint, GroupByWindow:
		return nil
	}
	return fmt.Errorf("unsupported group-by field: %q (valid: method, endpoint, window)", by)
}
