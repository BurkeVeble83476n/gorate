package policy

import "fmt"

// PatchOptions holds the fields that can be updated on a policy.
type PatchOptions struct {
	Limit  *int
	Window *string
	Method *string
}

// Patch applies a partial update to the named policy in the slice.
// Only non-nil fields in opts are applied. Returns an error if the
// policy is not found or if the resulting policy fails validation.
func Patch(policies []Policy, name string, opts PatchOptions) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}

	idx := -1
	for i, p := range policies {
		if p.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return nil, fmt.Errorf("policy %q not found", name)
	}

	updated := policies[idx]

	if opts.Limit != nil {
		updated.Limit = *opts.Limit
	}
	if opts.Window != nil {
		updated.Window = *opts.Window
	}
	if opts.Method != nil {
		updated.Method = *opts.Method
	}

	if err := Validate(updated); err != nil {
		return nil, fmt.Errorf("patched policy is invalid: %w", err)
	}

	result := make([]Policy, len(policies))
	copy(result, policies)
	result[idx] = updated

	return result, nil
}
