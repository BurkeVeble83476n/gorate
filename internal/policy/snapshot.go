package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PolicySnapshot represents a saved state of policies at a point in time.
type PolicySnapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Label     string    `json:"label"`
	Policies  []Policy  `json:"policies"`
}

// SaveSnapshot writes a snapshot of the given policies to a file.
func SaveSnapshot(policies []Policy, label, path string) error {
	snap := PolicySnapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Policies:  policies,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a file.
func LoadSnapshot(path string) (*PolicySnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}
	var snap PolicySnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}
	return &snap, nil
}

// DescribeSnapshot returns a human-readable summary of a snapshot.
func DescribeSnapshot(snap *PolicySnapshot) string {
	return fmt.Sprintf("Snapshot %q — %d policies saved at %s",
		snap.Label,
		len(snap.Policies),
		snap.Timestamp.Format(time.RFC3339),
	)
}
