package policy

import (
	"fmt"
	"sync"
	"time"
)

// PolicyStats holds counters for a single named policy.
type PolicyStats struct {
	Name     string
	Allowed  int64
	Rejected int64
	LastSeen time.Time
}

// RequestStats tracks rate-limit hit/allow counts per policy name.
type RequestStats struct {
	mu     sync.RWMutex
	counts map[string]*PolicyStats
}

// NewRequestStats creates an initialised RequestStats.
func NewRequestStats() *RequestStats {
	return &RequestStats{
		counts: make(map[string]*PolicyStats),
	}
}

// Record increments the allowed or rejected counter for the named policy.
func (s *RequestStats) Record(name string, allowed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ps, ok := s.counts[name]
	if !ok {
		ps = &PolicyStats{Name: name}
		s.counts[name] = ps
	}
	ps.LastSeen = time.Now()
	if allowed {
		ps.Allowed++
	} else {
		ps.Rejected++
	}
}

// Snapshot returns a point-in-time copy of all collected stats.
func (s *RequestStats) Snapshot() []PolicyStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]PolicyStats, 0, len(s.counts))
	for _, ps := range s.counts {
		out = append(out, *ps)
	}
	return out
}

// PrintStats writes a formatted stats table to stdout.
func (s *RequestStats) PrintStats() {
	snap := s.Snapshot()
	fmt.Printf("%-24s %10s %10s  %s\n", "Policy", "Allowed", "Rejected", "Last Seen")
	fmt.Println("------------------------------------------------------------------")
	for _, ps := range snap {
		fmt.Printf("%-24s %10d %10d  %s\n",
			ps.Name, ps.Allowed, ps.Rejected,
			ps.LastSeen.Format(time.RFC3339))
	}
}
