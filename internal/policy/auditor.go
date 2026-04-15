package policy

import (
	"fmt"
	"strings"
	"time"
)

// AuditEntry represents a single recorded change to a policy.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	PolicyName string   `json:"policy_name"`
	Detail    string    `json:"detail"`
}

// AuditLog holds a list of audit entries.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(action, policyName, detail string) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp:  time.Now().UTC(),
		Action:     action,
		PolicyName: policyName,
		Detail:     detail,
	})
}

// Filter returns entries matching the given action or policy name (empty string matches all).
func (a *AuditLog) Filter(action, policyName string) []AuditEntry {
	var result []AuditEntry
	for _, e := range a.Entries {
		actionMatch := action == "" || strings.EqualFold(e.Action, action)
		nameMatch := policyName == "" || strings.EqualFold(e.PolicyName, policyName)
		if actionMatch && nameMatch {
			result = append(result, e)
		}
	}
	return result
}

// FormatAuditLog returns a human-readable string of all audit entries.
func FormatAuditLog(entries []AuditEntry) string {
	if len(entries) == 0 {
		return "No audit entries found."
	}
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("[%s] %-10s %-20s %s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Action,
			e.PolicyName,
			e.Detail,
		))
	}
	return sb.String()
}
