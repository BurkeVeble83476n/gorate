package policy

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that modifies a Policy in place.
type TransformFunc func(p *Policy) error

// TransformResult holds the outcome of a transformation for one policy.
type TransformResult struct {
	Name    string
	Changed bool
	Note    string
}

// Transform applies one or more TransformFuncs to each policy in the slice.
// It returns a list of TransformResults describing what changed.
func Transform(policies []Policy, fns ...TransformFunc) ([]Policy, []TransformResult, error) {
	results := make([]TransformResult, 0, len(policies))
	out := make([]Policy, len(policies))

	for i, p := range policies {
		before := fmt.Sprintf("%s|%s|%d|%s", p.Name, p.Method, p.Limit, p.Window)
		copy := p
		for _, fn := range fns {
			if err := fn(&copy); err != nil {
				return nil, nil, fmt.Errorf("transform failed for policy %q: %w", p.Name, err)
			}
		}
		after := fmt.Sprintf("%s|%s|%d|%s", copy.Name, copy.Method, copy.Limit, copy.Window)
		out[i] = copy
		results = append(results, TransformResult{
			Name:    copy.Name,
			Changed: before != after,
			Note:    buildTransformNote(p, copy),
		})
	}
	return out, results, nil
}

func buildTransformNote(before, after Policy) string {
	var parts []string
	if before.Method != after.Method {
		parts = append(parts, fmt.Sprintf("method: %s→%s", before.Method, after.Method))
	}
	if before.Limit != after.Limit {
		parts = append(parts, fmt.Sprintf("limit: %d→%d", before.Limit, after.Limit))
	}
	if before.Window != after.Window {
		parts = append(parts, fmt.Sprintf("window: %s→%s", before.Window, after.Window))
	}
	if len(parts) == 0 {
		return "no change"
	}
	return strings.Join(parts, ", ")
}

// UppercaseMethod returns a TransformFunc that uppercases the Method field.
func UppercaseMethod() TransformFunc {
	return func(p *Policy) error {
		p.Method = strings.ToUpper(p.Method)
		return nil
	}
}

// CapLimit returns a TransformFunc that caps the Limit to the given maximum.
func CapLimit(max int) TransformFunc {
	return func(p *Policy) error {
		if max <= 0 {
			return fmt.Errorf("cap limit must be positive")
		}
		if p.Limit > max {
			p.Limit = max
		}
		return nil
	}
}

// SetDefaultWindow returns a TransformFunc that sets Window if it is empty.
func SetDefaultWindow(window string) TransformFunc {
	return func(p *Policy) error {
		if p.Window == "" {
			p.Window = window
		}
		return nil
	}
}
