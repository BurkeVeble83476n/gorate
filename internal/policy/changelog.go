package policy

import (
	"fmt"
	"strings"
	"time"
)

// ChangeEntry records a single change made to a policy.
type ChangeEntry struct {
	Timestamp time.Time
	PolicyName string
	Field      string
	OldValue   string
	NewValue   string
	Author     string
}

// Changelog holds an ordered list of change entries.
type Changelog struct {
	Entries []ChangeEntry
}

// NewChangelog creates an empty Changelog.
func NewChangelog() *Changelog {
	return &Changelog{}
}

// Record appends a new change entry to the changelog.
func (c *Changelog) Record(policyName, field, oldValue, newValue, author string) {
	c.Entries = append(c.Entries, ChangeEntry{
		Timestamp:  time.Now().UTC(),
		PolicyName: policyName,
		Field:      field,
		OldValue:   oldValue,
		NewValue:   newValue,
		Author:     author,
	})
}

// FilterByPolicy returns only entries matching the given policy name.
func (c *Changelog) FilterByPolicy(name string) []ChangeEntry {
	var out []ChangeEntry
	for _, e := range c.Entries {
		if e.PolicyName == name {
			out = append(out, e)
		}
	}
	return out
}

// FilterByField returns only entries matching the given field name.
func (c *Changelog) FilterByField(field string) []ChangeEntry {
	var out []ChangeEntry
	for _, e := range c.Entries {
		if strings.EqualFold(e.Field, field) {
			out = append(out, e)
		}
	}
	return out
}

// FormatChangelog returns a human-readable string of all entries.
func FormatChangelog(entries []ChangeEntry) string {
	if len(entries) == 0 {
		return "no changelog entries"
	}
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("[%s] %s.%s: %q -> %q (by %s)\n",
			e.Timestamp.Format(time.RFC3339),
			e.PolicyName,
			e.Field,
			e.OldValue,
			e.NewValue,
			e.Author,
		))
	}
	return strings.TrimRight(sb.String(), "\n")
}
