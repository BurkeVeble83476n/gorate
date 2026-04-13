package policy

import (
	"fmt"
	"slices"
)

// Reorder repositions a named policy to a specific index within the slice.
// The index is zero-based. If the index is out of range, an error is returned.
// The original slice is not modified; a new slice is returned.
func Reorder(policies []Policy, name string, targetIndex int) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}

	if targetIndex < 0 || targetIndex >= len(policies) {
		return nil, fmt.Errorf("target index %d is out of range (0..%d)", targetIndex, len(policies)-1)
	}

	srcIndex := -1
	for i, p := range policies {
		if p.Name == name {
			srcIndex = i
			break
		}
	}

	if srcIndex == -1 {
		return nil, fmt.Errorf("policy %q not found", name)
	}

	if srcIndex == targetIndex {
		result := make([]Policy, len(policies))
		copy(result, policies)
		return result, nil
	}

	result := slices.Clone(policies)
	moved := result[srcIndex]
	result = slices.Delete(result, srcIndex, srcIndex+1)
	result = slices.Insert(result, targetIndex, moved)
	return result, nil
}

// ReorderToFront moves the named policy to index 0.
func ReorderToFront(policies []Policy, name string) ([]Policy, error) {
	return Reorder(policies, name, 0)
}

// ReorderToBack moves the named policy to the last position.
func ReorderToBack(policies []Policy, name string) ([]Policy, error) {
	if len(policies) == 0 {
		return nil, fmt.Errorf("policy list is empty")
	}
	return Reorder(policies, name, len(policies)-1)
}
