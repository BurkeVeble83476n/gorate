package policy

import (
	"sort"
)

// SortField represents a field by which policies can be sorted.
type SortField string

const (
	SortByName     SortField = "name"
	SortByEndpoint SortField = "endpoint"
	SortByMethod   SortField = "method"
	SortByLimit    SortField = "limit"
	SortByWindow   SortField = "window"
)

// SortOrder represents ascending or descending sort direction.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortOptions configures how policies are sorted.
type SortOptions struct {
	Field SortField
	Order SortOrder
}

// Sort returns a sorted copy of the given policy slice based on the provided options.
// If Field is empty or unrecognized, the original order is preserved.
func Sort(policies []Policy, opts SortOptions) []Policy {
	result := make([]Policy, len(policies))
	copy(result, policies)

	var less func(i, j int) bool

	switch opts.Field {
	case SortByName:
		less = func(i, j int) bool { return result[i].Name < result[j].Name }
	case SortByEndpoint:
		less = func(i, j int) bool { return result[i].Endpoint < result[j].Endpoint }
	case SortByMethod:
		less = func(i, j int) bool { return result[i].Method < result[j].Method }
	case SortByLimit:
		less = func(i, j int) bool { return result[i].Limit < result[j].Limit }
	case SortByWindow:
		less = func(i, j int) bool { return result[i].Window < result[j].Window }
	default:
		return result
	}

	if opts.Order == SortDesc {
		orig := less
		less = func(i, j int) bool { return orig(j, i) }
	}

	sort.SliceStable(result, less)
	return result
}
