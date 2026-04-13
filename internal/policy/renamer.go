package policy

import "fmt"

// RenameOptions configures how a policy rename is applied.
type RenameOptions struct {
	OldName string
	NewName string
}

// Rename finds a policy by OldName and renames it to NewName.
// Returns an error if OldName is not found or NewName is already taken.
func Rename(policies []Policy, opts RenameOptions) ([]Policy, error) {
	if opts.OldName == "" {
		return nil, fmt.Errorf("old name must not be empty")
	}
	if opts.NewName == "" {
		return nil, fmt.Errorf("new name must not be empty")
	}
	if opts.OldName == opts.NewName {
		return nil, fmt.Errorf("old name and new name are the same: %q", opts.OldName)
	}

	foundOld := false
	for _, p := range policies {
		if p.Name == opts.NewName {
			return nil, fmt.Errorf("a policy with name %q already exists", opts.NewName)
		}
		if p.Name == opts.OldName {
			foundOld = true
		}
	}
	if !foundOld {
		return nil, fmt.Errorf("policy %q not found", opts.OldName)
	}

	updated := make([]Policy, len(policies))
	for i, p := range policies {
		if p.Name == opts.OldName {
			p.Name = opts.NewName
		}
		updated[i] = p
	}
	return updated, nil
}
