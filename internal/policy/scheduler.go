package policy

import (
	"fmt"
	"time"
)

// Schedule represents a time-based activation window for a policy.
type Schedule struct {
	PolicyName string
	ActiveFrom time.Time
	ActiveTo   time.Time
}

// ScheduleResult holds the outcome of a schedule evaluation.
type ScheduleResult struct {
	PolicyName string
	Active     bool
	Reason     string
}

// EvaluateSchedules checks which policies are currently active based on their schedules.
func EvaluateSchedules(schedules []Schedule, now time.Time) []ScheduleResult {
	results := make([]ScheduleResult, 0, len(schedules))
	for _, s := range schedules {
		var result ScheduleResult
		result.PolicyName = s.PolicyName
		if now.Before(s.ActiveFrom) {
			result.Active = false
			result.Reason = fmt.Sprintf("not yet active, starts at %s", s.ActiveFrom.Format(time.RFC3339))
		} else if now.After(s.ActiveTo) {
			result.Active = false
			result.Reason = fmt.Sprintf("expired at %s", s.ActiveTo.Format(time.RFC3339))
		} else {
			result.Active = true
			result.Reason = fmt.Sprintf("active until %s", s.ActiveTo.Format(time.RFC3339))
		}
		results = append(results, result)
	}
	return results
}

// FilterActiveBySchedule returns only the policies whose schedules are currently active.
func FilterActiveBySchedule(policies []Policy, schedules []Schedule, now time.Time) []Policy {
	activeNames := map[string]bool{}
	for _, r := range EvaluateSchedules(schedules, now) {
		if r.Active {
			activeNames[r.PolicyName] = true
		}
	}
	var active []Policy
	for _, p := range policies {
		if activeNames[p.Name] {
			active = append(active, p)
		}
	}
	return active
}

// ValidateSchedule returns an error if the schedule window is invalid.
func ValidateSchedule(s Schedule) error {
	if s.PolicyName == "" {
		return fmt.Errorf("schedule must reference a policy name")
	}
	if !s.ActiveTo.After(s.ActiveFrom) {
		return fmt.Errorf("schedule active_to must be after active_from for policy %q", s.PolicyName)
	}
	return nil
}
