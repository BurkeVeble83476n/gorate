package policy

import "fmt"

// DuplicateError describes a duplicate policy conflict.
type DuplicateError struct {
	Name     string
	Endpoint string
	Method   string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate policy %q: endpoint %s method %s already defined",
		e.Name, e.Endpoint, e.Method)
}

// FindDuplicates returns a list of errors for any policies that share the
// same endpoint+method combination. The first occurrence is considered the
// canonical entry; all subsequent ones are reported as duplicates.
func FindDuplicates(policies []Policy) []error {
	type key struct {
		Endpoint string
		Method   string
	}

	seen := make(map[key]string) // key -> first policy name
	var errs []error

	for _, p := range policies {
		method := p.Method
		if method == "" {
			method = "*"
		}
		k := key{Endpoint: p.Endpoint, Method: method}
		if first, exists := seen[k]; exists {
			errs = append(errs, &DuplicateError{
				Name:     p.Name,
				Endpoint: p.Endpoint,
				Method:   method,
			})
			_ = first
		} else {
			seen[k] = p.Name
		}
	}

	return errs
}

// DeduplicatePolicies returns a new slice with duplicate endpoint+method
// combinations removed, keeping only the first occurrence.
func DeduplicatePolicies(policies []Policy) []Policy {
	type key struct {
		Endpoint string
		Method   string
	}

	seen := make(map[key]bool)
	result := make([]Policy, 0, len(policies))

	for _, p := range policies {
		method := p.Method
		if method == "" {
			method = "*"
		}
		k := key{Endpoint: p.Endpoint, Method: method}
		if !seen[k] {
			seen[k] = true
			result = append(result, p)
		}
	}

	return result
}
