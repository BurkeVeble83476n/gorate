package policy

import (
	"testing"
)

func makeReviewerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 5, Window: 30},
	}
}

func TestSetReview_Success(t *testing.T) {
	policies := makeReviewerPolicies()
	updated, err := SetReview(policies, "alpha", ReviewApproved, "alice", "looks good")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, _ := GetReview(updated, "alpha")
	if entry.Status != ReviewApproved {
		t.Errorf("expected approved, got %s", entry.Status)
	}
	if entry.Reviewer != "alice" {
		t.Errorf("expected reviewer alice, got %s", entry.Reviewer)
	}
	if entry.Comment != "looks good" {
		t.Errorf("unexpected comment: %s", entry.Comment)
	}
}

func TestSetReview_InvalidStatus(t *testing.T) {
	policies := makeReviewerPolicies()
	_, err := SetReview(policies, "alpha", ReviewStatus("unknown"), "alice", "")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestSetReview_PolicyNotFound(t *testing.T) {
	policies := makeReviewerPolicies()
	_, err := SetReview(policies, "nonexistent", ReviewApproved, "", "")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestSetReview_EmptyName(t *testing.T) {
	policies := makeReviewerPolicies()
	_, err := SetReview(policies, "", ReviewApproved, "", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestGetReview_DefaultsPending(t *testing.T) {
	policies := makeReviewerPolicies()
	entry, err := GetReview(policies, "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Status != ReviewPending {
		t.Errorf("expected pending default, got %s", entry.Status)
	}
}

func TestGetReview_NotFound(t *testing.T) {
	policies := makeReviewerPolicies()
	_, err := GetReview(policies, "ghost")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestListByReviewStatus_ReturnsMatching(t *testing.T) {
	policies := makeReviewerPolicies()
	policies, _ = SetReview(policies, "alpha", ReviewApproved, "alice", "ok")
	policies, _ = SetReview(policies, "beta", ReviewRejected, "bob", "needs work")

	approved := ListByReviewStatus(policies, ReviewApproved)
	if len(approved) != 1 || approved[0].PolicyName != "alpha" {
		t.Errorf("expected 1 approved policy (alpha), got %+v", approved)
	}
}

func TestListByReviewStatus_PendingByDefault(t *testing.T) {
	policies := makeReviewerPolicies()
	pending := ListByReviewStatus(policies, ReviewPending)
	if len(pending) != 2 {
		t.Errorf("expected 2 pending policies, got %d", len(pending))
	}
}

func TestFormatReview_ContainsFields(t *testing.T) {
	e := ReviewEntry{PolicyName: "alpha", Status: ReviewApproved, Reviewer: "alice", Comment: "ok"}
	out := FormatReview(e)
	for _, want := range []string{"alpha", "approved", "alice", "ok"} {
		if !containsStr(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
