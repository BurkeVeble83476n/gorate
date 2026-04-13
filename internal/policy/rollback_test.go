package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRollbackPolicies(names ...string) []Policy {
	policies := make([]Policy, len(names))
	for i, n := range names {
		policies[i] = Policy{Name: n, Endpoint: "/api", Method: "GET", Limit: 10, Window: 60}
	}
	return policies
}

func TestNewHistory_DefaultDepth(t *testing.T) {
	h := NewHistory(0)
	assert.Equal(t, 10, h.maxDepth)
}

func TestPush_StoresSnapshot(t *testing.T) {
	h := NewHistory(5)
	policies := makeRollbackPolicies("alpha", "beta")
	h.Push("initial", policies)
	assert.Equal(t, 1, h.Len())
}

func TestPush_EnforcesMaxDepth(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Push("snap", makeRollbackPolicies("p"))
	}
	assert.Equal(t, 3, h.Len())
}

func TestPop_ReturnsLastSnapshot(t *testing.T) {
	h := NewHistory(5)
	p1 := makeRollbackPolicies("first")
	p2 := makeRollbackPolicies("second")
	h.Push("snap1", p1)
	h.Push("snap2", p2)

	snap, err := h.Pop()
	require.NoError(t, err)
	assert.Equal(t, "snap2", snap.Label)
	assert.Equal(t, "second", snap.Policies[0].Name)
	assert.Equal(t, 1, h.Len())
}

func TestPop_EmptyHistoryReturnsError(t *testing.T) {
	h := NewHistory(5)
	_, err := h.Pop()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no snapshots available")
}

func TestPeek_DoesNotRemove(t *testing.T) {
	h := NewHistory(5)
	h.Push("only", makeRollbackPolicies("a"))

	snap, err := h.Peek()
	require.NoError(t, err)
	assert.Equal(t, "only", snap.Label)
	assert.Equal(t, 1, h.Len())
}

func TestPeek_EmptyHistoryReturnsError(t *testing.T) {
	h := NewHistory(5)
	_, err := h.Peek()
	require.Error(t, err)
}

func TestLabels_OrderedOldestFirst(t *testing.T) {
	h := NewHistory(5)
	h.Push("first", makeRollbackPolicies("a"))
	h.Push("second", makeRollbackPolicies("b"))
	h.Push("third", makeRollbackPolicies("c"))

	labels := h.Labels()
	assert.Equal(t, []string{"first", "second", "third"}, labels)
}

func TestPush_IsolatesCopy(t *testing.T) {
	h := NewHistory(5)
	policies := makeRollbackPolicies("original")
	h.Push("snap", policies)

	// Mutate original slice after push
	policies[0].Name = "mutated"

	snap, err := h.Peek()
	require.NoError(t, err)
	assert.Equal(t, "original", snap.Policies[0].Name)
}
