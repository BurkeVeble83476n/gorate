package policy

import (
	"fmt"
	"strings"
)

// ReviewStatus represents the outcome of a policy review.
type ReviewStatus string

const (
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "rejected"
	ReviewPending  ReviewStatus = "pending"
)

// ReviewEntry holds review metadata for a single policy.
type ReviewEntry struct {
	PolicyName string
	Status     ReviewStatus
	Reviewer   string
	Comment    string
}

// SetReview sets the review status, reviewer, and comment on a named policy.
func SetReview(policies []Policy, name string, status ReviewStatus, reviewer, comment string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	valid := map[ReviewStatus]bool{ReviewApproved: true, ReviewRejected: true, ReviewPending: true}
	if !valid[status] {
		return nil, fmt.Errorf("invalid review status %q: must be approved, rejected, or pending", status)
	}
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = map[string]string{}
			}
			policies[i].Annotations["review/status"] = string(status)
			policies[i].Annotations["review/reviewer"] = reviewer
			policies[i].Annotations["review/comment"] = comment
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetReview returns the ReviewEntry for the named policy.
func GetReview(policies []Policy, name string) (ReviewEntry, error) {
	for _, p := range policies {
		if p.Name == name {
			status := ReviewStatus(p.Annotations["review/status"])
			if status == "" {
				status = ReviewPending
			}
			return ReviewEntry{
				PolicyName: p.Name,
				Status:     status,
				Reviewer:   p.Annotations["review/reviewer"],
				Comment:    p.Annotations["review/comment"],
			}, nil
		}
	}
	return ReviewEntry{}, fmt.Errorf("policy %q not found", name)
}

// ListByReviewStatus returns all policies matching the given review status.
func ListByReviewStatus(policies []Policy, status ReviewStatus) []ReviewEntry {
	var results []ReviewEntry
	for _, p := range policies {
		s := ReviewStatus(p.Annotations["review/status"])
		if s == "" {
			s = ReviewPending
		}
		if s == status {
			results = append(results, ReviewEntry{
				PolicyName: p.Name,
				Status:     s,
				Reviewer:   p.Annotations["review/reviewer"],
				Comment:    p.Annotations["review/comment"],
			})
		}
	}
	return results
}

// FormatReview returns a human-readable summary of a ReviewEntry.
func FormatReview(e ReviewEntry) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Policy:   %s\n", e.PolicyName)
	fmt.Fprintf(&sb, "Status:   %s\n", e.Status)
	if e.Reviewer != "" {
		fmt.Fprintf(&sb, "Reviewer: %s\n", e.Reviewer)
	}
	if e.Comment != "" {
		fmt.Fprintf(&sb, "Comment:  %s\n", e.Comment)
	}
	return sb.String()
}
