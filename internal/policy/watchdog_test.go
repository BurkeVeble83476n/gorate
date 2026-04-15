package policy

import (
	"strings"
	"testing"
	"time"
)

func makeWatchdogStats(allowed, rejected int) map[string]StatsSnapshot {
	return map[string]StatsSnapshot{
		"api-limit": {Allowed: allowed, Rejected: rejected},
	}
}

func TestEvaluateWatchdog_NoAlert(t *testing.T) {
	stats := makeWatchdogStats(100, 2)
	rules := []WatchdogRule{{PolicyName: "api-limit", MaxRejections: 5, Window: time.Minute}}
	alerts := EvaluateWatchdog(stats, rules)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}

func TestEvaluateWatchdog_TriggersAlert(t *testing.T) {
	stats := makeWatchdogStats(50, 10)
	rules := []WatchdogRule{{PolicyName: "api-limit", MaxRejections: 5, Window: time.Minute}}
	alerts := EvaluateWatchdog(stats, rules)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].PolicyName != "api-limit" {
		t.Errorf("expected policy name api-limit, got %s", alerts[0].PolicyName)
	}
	if alerts[0].Rejections != 10 {
		t.Errorf("expected rejections 10, got %d", alerts[0].Rejections)
	}
}

func TestEvaluateWatchdog_UnknownPolicy(t *testing.T) {
	stats := makeWatchdogStats(10, 10)
	rules := []WatchdogRule{{PolicyName: "unknown", MaxRejections: 1, Window: time.Minute}}
	alerts := EvaluateWatchdog(stats, rules)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts for unknown policy, got %d", len(alerts))
	}
}

func TestEvaluateWatchdog_ExactThreshold(t *testing.T) {
	stats := makeWatchdogStats(20, 5)
	rules := []WatchdogRule{{PolicyName: "api-limit", MaxRejections: 5, Window: time.Minute}}
	alerts := EvaluateWatchdog(stats, rules)
	if len(alerts) != 1 {
		t.Fatalf("expected alert at exact threshold, got %d alerts", len(alerts))
	}
}

func TestFormatAlerts_NoAlerts(t *testing.T) {
	out := FormatAlerts(nil)
	if !strings.Contains(out, "no watchdog alerts") {
		t.Errorf("expected no-alert message, got: %s", out)
	}
}

func TestFormatAlerts_WithAlerts(t *testing.T) {
	alerts := []WatchdogAlert{
		{PolicyName: "api-limit", Rejections: 10, TriggeredAt: time.Now(), Message: "policy \"api-limit\" exceeded rejection threshold: 10 >= 5"},
	}
	out := FormatAlerts(alerts)
	if !strings.Contains(out, "[ALERT]") {
		t.Errorf("expected ALERT tag in output, got: %s", out)
	}
	if !strings.Contains(out, "api-limit") {
		t.Errorf("expected policy name in output, got: %s", out)
	}
}
