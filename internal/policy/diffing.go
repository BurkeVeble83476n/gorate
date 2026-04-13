package policy

import "fmt"

// DiffResult represents the difference between two policy slices.
type DiffResult struct {
	Added   []Policy
	Removed []Policy
	Changed []PolicyChange
}

// PolicyChange describes a single policy that changed between two sets.
type PolicyChange struct {
	Name   string
	Before Policy
	After  Policy
}

// Diff compares two slices of policies and returns what was added, removed, or changed.
func Diff(base, updated []Policy) DiffResult {
	result := DiffResult{}

	baseMap := make(map[string]Policy, len(base))
	for _, p := range base {
		baseMap[p.Name] = p
	}

	updatedMap := make(map[string]Policy, len(updated))
	for _, p := range updated {
		updatedMap[p.Name] = p
	}

	for _, p := range updated {
		if orig, found := baseMap[p.Name]; !found {
			result.Added = append(result.Added, p)
		} else if policyChanged(orig, p) {
			result.Changed = append(result.Changed, PolicyChange{
				Name:   p.Name,
				Before: orig,
				After:  p,
			})
		}
	}

	for _, p := range base {
		if _, found := updatedMap[p.Name]; !found {
			result.Removed = append(result.Removed, p)
		}
	}

	return result
}

func policyChanged(a, b Policy) bool {
	return a.Endpoint != b.Endpoint ||
		a.Method != b.Method ||
		a.Limit != b.Limit ||
		a.Window != b.Window
}

// FormatDiff returns a human-readable summary of a DiffResult.
func FormatDiff(d DiffResult) string {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return "No differences found."
	}

	out := ""
	for _, p := range d.Added {
		out += fmt.Sprintf("+ [added]   %s (%s %s, limit=%d, window=%s)\n",
			p.Name, p.Method, p.Endpoint, p.Limit, p.Window)
	}
	for _, p := range d.Removed {
		out += fmt.Sprintf("- [removed] %s (%s %s, limit=%d, window=%s)\n",
			p.Name, p.Method, p.Endpoint, p.Limit, p.Window)
	}
	for _, c := range d.Changed {
		out += fmt.Sprintf("~ [changed] %s: (%s %s, limit=%d, window=%s) -> (%s %s, limit=%d, window=%s)\n",
			c.Name,
			c.Before.Method, c.Before.Endpoint, c.Before.Limit, c.Before.Window,
			c.After.Method, c.After.Endpoint, c.After.Limit, c.After.Window)
	}
	return out
}
