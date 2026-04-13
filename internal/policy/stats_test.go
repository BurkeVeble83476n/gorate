package policy

import (
	"testing"
)

func TestRecord_AllowedIncrement(t *testing.T) {
	s := NewRequestStats()
	s.Record("api-limit", true)
	s.Record("api-limit", true)

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].Allowed != 2 {
		t.Errorf("expected Allowed=2, got %d", snap[0].Allowed)
	}
	if snap[0].Rejected != 0 {
		t.Errorf("expected Rejected=0, got %d", snap[0].Rejected)
	}
}

func TestRecord_RejectedIncrement(t *testing.T) {
	s := NewRequestStats()
	s.Record("api-limit", false)

	snap := s.Snapshot()
	if snap[0].Rejected != 1 {
		t.Errorf("expected Rejected=1, got %d", snap[0].Rejected)
	}
}

func TestRecord_MultiplePolicies(t *testing.T) {
	s := NewRequestStats()
	s.Record("policy-a", true)
	s.Record("policy-b", false)
	s.Record("policy-a", false)

	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(snap))
	}

	byName := make(map[string]PolicyStats)
	for _, ps := range snap {
		byName[ps.Name] = ps
	}

	if byName["policy-a"].Allowed != 1 {
		t.Errorf("policy-a Allowed: want 1, got %d", byName["policy-a"].Allowed)
	}
	if byName["policy-a"].Rejected != 1 {
		t.Errorf("policy-a Rejected: want 1, got %d", byName["policy-a"].Rejected)
	}
	if byName["policy-b"].Rejected != 1 {
		t.Errorf("policy-b Rejected: want 1, got %d", byName["policy-b"].Rejected)
	}
}

func TestSnapshot_IsACopy(t *testing.T) {
	s := NewRequestStats()
	s.Record("x", true)

	snap1 := s.Snapshot()
	s.Record("x", true)
	snap2 := s.Snapshot()

	if snap1[0].Allowed == snap2[0].Allowed {
		t.Error("expected snapshots to differ after additional record")
	}
}

func TestPrintStats_DoesNotPanic(t *testing.T) {
	s := NewRequestStats()
	s.Record("demo", true)
	s.Record("demo", false)
	// Should not panic
	s.PrintStats()
}
