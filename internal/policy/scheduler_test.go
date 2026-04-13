package policy

import (
	"testing"
	"time"
)

func makeSchedule(name string, fromOffset, toOffset time.Duration, now time.Time) Schedule {
	return Schedule{
		PolicyName: name,
		ActiveFrom: now.Add(fromOffset),
		ActiveTo:   now.Add(toOffset),
	}
}

func TestEvaluateSchedules_Active(t *testing.T) {
	now := time.Now()
	s := makeSchedule("api-limit", -time.Hour, time.Hour, now)
	results := EvaluateSchedules([]Schedule{s}, now)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Active {
		t.Errorf("expected policy to be active, got reason: %s", results[0].Reason)
	}
}

func TestEvaluateSchedules_NotYetActive(t *testing.T) {
	now := time.Now()
	s := makeSchedule("future-policy", time.Hour, 2*time.Hour, now)
	results := EvaluateSchedules([]Schedule{s}, now)
	if results[0].Active {
		t.Error("expected policy to be inactive (not yet started)")
	}
}

func TestEvaluateSchedules_Expired(t *testing.T) {
	now := time.Now()
	s := makeSchedule("old-policy", -2*time.Hour, -time.Hour, now)
	results := EvaluateSchedules([]Schedule{s}, now)
	if results[0].Active {
		t.Error("expected policy to be inactive (expired)")
	}
}

func TestFilterActiveBySchedule_ReturnsOnlyActive(t *testing.T) {
	now := time.Now()
	policies := []Policy{
		{Name: "active-policy", Endpoint: "/a", Method: "GET", Limit: 10, Window: "1m"},
		{Name: "inactive-policy", Endpoint: "/b", Method: "GET", Limit: 5, Window: "1m"},
	}
	schedules := []Schedule{
		makeSchedule("active-policy", -time.Hour, time.Hour, now),
		makeSchedule("inactive-policy", time.Hour, 2*time.Hour, now),
	}
	active := FilterActiveBySchedule(policies, schedules, now)
	if len(active) != 1 {
		t.Fatalf("expected 1 active policy, got %d", len(active))
	}
	if active[0].Name != "active-policy" {
		t.Errorf("unexpected active policy: %s", active[0].Name)
	}
}

func TestFilterActiveBySchedule_NoneActive(t *testing.T) {
	now := time.Now()
	policies := []Policy{
		{Name: "p1", Endpoint: "/x", Method: "GET", Limit: 1, Window: "1m"},
	}
	schedules := []Schedule{
		makeSchedule("p1", time.Hour, 2*time.Hour, now),
	}
	active := FilterActiveBySchedule(policies, schedules, now)
	if len(active) != 0 {
		t.Errorf("expected no active policies, got %d", len(active))
	}
}

func TestValidateSchedule_Valid(t *testing.T) {
	now := time.Now()
	s := makeSchedule("valid", -time.Hour, time.Hour, now)
	if err := ValidateSchedule(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSchedule_MissingName(t *testing.T) {
	now := time.Now()
	s := makeSchedule("", -time.Hour, time.Hour, now)
	if err := ValidateSchedule(s); err == nil {
		t.Error("expected error for missing policy name")
	}
}

func TestValidateSchedule_InvalidWindow(t *testing.T) {
	now := time.Now()
	s := makeSchedule("bad-window", time.Hour, -time.Hour, now)
	if err := ValidateSchedule(s); err == nil {
		t.Error("expected error for active_to before active_from")
	}
}
