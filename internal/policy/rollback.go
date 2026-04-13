package policy

import "fmt"

// Snapshot holds a named point-in-time copy of a policy list.
type Snapshot struct {
	Label    string
	Policies []Policy
}

// History maintains an ordered stack of policy snapshots for rollback support.
type History struct {
	snapshots []Snapshot
	maxDepth  int
}

// NewHistory creates a History with the given maximum snapshot depth.
func NewHistory(maxDepth int) *History {
	if maxDepth <= 0 {
		maxDepth = 10
	}
	return &History{maxDepth: maxDepth}
}

// Push saves a labelled snapshot of the provided policies.
func (h *History) Push(label string, policies []Policy) {
	copied := make([]Policy, len(policies))
	copy(copied, policies)
	h.snapshots = append(h.snapshots, Snapshot{Label: label, Policies: copied})
	if len(h.snapshots) > h.maxDepth {
		h.snapshots = h.snapshots[len(h.snapshots)-h.maxDepth:]
	}
}

// Pop removes and returns the most recent snapshot.
// Returns an error if the history is empty.
func (h *History) Pop() (Snapshot, error) {
	if len(h.snapshots) == 0 {
		return Snapshot{}, fmt.Errorf("rollback: no snapshots available")
	}
	last := h.snapshots[len(h.snapshots)-1]
	h.snapshots = h.snapshots[:len(h.snapshots)-1]
	return last, nil
}

// Peek returns the most recent snapshot without removing it.
func (h *History) Peek() (Snapshot, error) {
	if len(h.snapshots) == 0 {
		return Snapshot{}, fmt.Errorf("rollback: no snapshots available")
	}
	return h.snapshots[len(h.snapshots)-1], nil
}

// Len returns the number of stored snapshots.
func (h *History) Len() int {
	return len(h.snapshots)
}

// Labels returns the label of every stored snapshot in order (oldest first).
func (h *History) Labels() []string {
	labels := make([]string, len(h.snapshots))
	for i, s := range h.snapshots {
		labels[i] = s.Label
	}
	return labels
}
