package policy

import (
	"fmt"
	"time"
)

// WatchdogRule defines a threshold condition for alerting.
type WatchdogRule struct {
	PolicyName    string
	MaxRejections int
	Window        time.Duration
}

// WatchdogAlert represents a triggered alert.
type WatchdogAlert struct {
	PolicyName  string
	Rejections  int
	Window      time.Duration
	TriggeredAt time.Time
	Message     string
}

// EvaluateWatchdog checks RequestStats against watchdog rules and returns any triggered alerts.
func EvaluateWatchdog(stats map[string]StatsSnapshot, rules []WatchdogRule) []WatchdogAlert {
	var alerts []WatchdogAlert
	for _, rule := range rules {
		snap, ok := stats[rule.PolicyName]
		if !ok {
			continue
		}
		if snap.Rejected >= rule.MaxRejections {
			alerts = append(alerts, WatchdogAlert{
				PolicyName:  rule.PolicyName,
				Rejections:  snap.Rejected,
				Window:      rule.Window,
				TriggeredAt: time.Now(),
				Message:     fmt.Sprintf("policy %q exceeded rejection threshold: %d >= %d", rule.PolicyName, snap.Rejected, rule.MaxRejections),
			})
		}
	}
	return alerts
}

// FormatAlerts returns a human-readable string of all triggered alerts.
func FormatAlerts(alerts []WatchdogAlert) string {
	if len(alerts) == 0 {
		return "no watchdog alerts triggered\n"
	}
	out := ""
	for _, a := range alerts {
		out += fmt.Sprintf("[ALERT] %s (at %s)\n", a.Message, a.TriggeredAt.Format(time.RFC3339))
	}
	return out
}
