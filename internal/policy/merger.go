package policy

import (
	"fmt"
	"strings"
)

// MergeStrategy defines how conflicting policies are resolved.
type MergeStrategy string

const (
	MergeStrategyKeepFirst MergeStrategy = "keep-first"
	MergeStrategyKeepLast  MergeStrategy = "keep-last"
	MergeStrategyHighest   MergeStrategy = "highest-limit"
	MergeStrategyLowest    MergeStrategy = "lowest-limit"
)

// MergeOptions configures how Merge behaves.
type MergeOptions struct {
	Strategy MergeStrategy
}

// Merge combines multiple policy slices into one, resolving conflicts
// according to the provided strategy. Policies are considered conflicting
// when they share the same endpoint and method.
func Merge(strategy MergeStrategy, sets ...[]Policy) ([]Policy, error) {
	if err := validateStrategy(strategy); err != nil {
		return nil, err
	}

	type key struct{ endpoint, method string }
	seen := make(map[key]int) // maps key -> index in result
	result := []Policy{}

	for _, policies := range sets {
		for _, p := range policies {
			k := key{
				endpoint: strings.ToLower(p.Endpoint),
				method:   strings.ToUpper(p.Method),
			}
			idx, exists := seen[k]
			if !exists {
				seen[k] = len(result)
				result = append(result, p)
				continue
			}
			resolved := resolve(result[idx], p, strategy)
			result[idx] = resolved
		}
	}
	return result, nil
}

func resolve(existing, incoming Policy, strategy MergeStrategy) Policy {
	switch strategy {
	case MergeStrategyKeepFirst:
		return existing
	case MergeStrategyKeepLast:
		return incoming
	case MergeStrategyHighest:
		if incoming.Limit > existing.Limit {
			return incoming
		}
		return existing
	case MergeStrategyLowest:
		if incoming.Limit < existing.Limit {
			return incoming
		}
		return existing
	default:
		return existing
	}
}

func validateStrategy(s MergeStrategy) error {
	switch s {
	case MergeStrategyKeepFirst, MergeStrategyKeepLast, MergeStrategyHighwest:
		return nil
	}
	return fmt.Errorf("unknown merge strategy %q", s)
}
