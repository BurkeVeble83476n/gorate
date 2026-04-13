package policy

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// PolicyProfile holds a runtime profile summary for a single policy.
type PolicyProfile struct {
	Name        string
	Endpoint    string
	Method      string
	Limit       int
	Window      time.Duration
	RiskScore   int
	Tags        []string
	IssueCount  int
	Notes       []string
}

// Profile generates a detailed profile for each policy, combining scoring,
// linting, and metadata into a single summary view.
func Profile(policies []Policy) []PolicyProfile {
	scores := scoreOne // reuse internal scorer
	_ = scores

	lintResults := Lint(policies)
	issues := make(map[string][]LintIssue)
	for _, issue := range lintResults {
		issues[issue.PolicyName] = append(issues[issue.PolicyName], issue)
	}

	var profiles []PolicyProfile
	for _, p := range policies {
		notes := []string{}
		for _, iss := range issues[p.Name] {
			notes = append(notes, iss.Message)
		}

		risk := scoreOne(p)

		profiles = append(profiles, PolicyProfile{
			Name:       p.Name,
			Endpoint:   p.Endpoint,
			Method:     p.Method,
			Limit:      p.Limit,
			Window:     p.Window,
			RiskScore:  risk,
			Tags:       GetTags(p),
			IssueCount: len(issues[p.Name]),
			Notes:      notes,
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].RiskScore > profiles[j].RiskScore
	})

	return profiles
}

// FormatProfile returns a human-readable string for a single PolicyProfile.
func FormatProfile(pp PolicyProfile) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Policy : %s\n", pp.Name))
	sb.WriteString(fmt.Sprintf("Endpoint: %s [%s]\n", pp.Endpoint, pp.Method))
	sb.WriteString(fmt.Sprintf("Limit   : %d req / %s\n", pp.Limit, pp.Window))
	sb.WriteString(fmt.Sprintf("Risk    : %d\n", pp.RiskScore))
	if len(pp.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags    : %s\n", strings.Join(pp.Tags, ", ")))
	}
	if len(pp.Notes) > 0 {
		sb.WriteString("Issues  :\n")
		for _, n := range pp.Notes {
			sb.WriteString(fmt.Sprintf("  - %s\n", n))
		}
	}
	return sb.String()
}
