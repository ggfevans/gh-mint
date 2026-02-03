package github

import (
	"testing"

	"github.com/ggfevans/gh-mint/internal/config"
)

func TestBranchProtectionPayload(t *testing.T) {
	bp := config.BranchProtection{
		Branch:              "main",
		RequiredReviews:     1,
		DismissStaleReviews: true,
		RequireStatusChecks: false,
	}
	payload := buildProtectionPayload(bp)
	reviews, ok := payload["required_pull_request_reviews"].(map[string]interface{})
	if !ok {
		t.Fatal("missing required_pull_request_reviews")
	}
	if reviews["dismiss_stale_reviews"] != true {
		t.Error("dismiss_stale_reviews should be true")
	}
	if reviews["required_approving_review_count"] != 1 {
		t.Errorf("required_approving_review_count = %v, want 1", reviews["required_approving_review_count"])
	}
}

func TestBranchProtectionPayload_NoReviews(t *testing.T) {
	bp := config.BranchProtection{
		Branch:          "main",
		RequiredReviews: 0,
	}
	payload := buildProtectionPayload(bp)
	if payload["required_pull_request_reviews"] != nil {
		t.Error("expected nil required_pull_request_reviews when RequiredReviews is 0")
	}
}
