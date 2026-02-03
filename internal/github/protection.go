package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ggfevans/gh-mint/internal/config"
)

func buildProtectionPayload(bp config.BranchProtection) map[string]interface{} {
	payload := map[string]interface{}{
		"enforce_admins":                false,
		"required_status_checks":        nil,
		"restrictions":                  nil,
		"required_pull_request_reviews": nil,
	}

	if bp.RequiredReviews > 0 {
		payload["required_pull_request_reviews"] = map[string]interface{}{
			"dismiss_stale_reviews":           bp.DismissStaleReviews,
			"required_approving_review_count": bp.RequiredReviews,
		}
	}

	if bp.RequireStatusChecks {
		payload["required_status_checks"] = map[string]interface{}{
			"strict":   true,
			"contexts": []string{},
		}
	}

	return payload
}

func (c *Client) SetBranchProtection(nwo string, bp config.BranchProtection) error {
	if bp.Branch == "" {
		return nil
	}

	payload := buildProtectionPayload(bp)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling protection: %w", err)
	}

	endpoint := fmt.Sprintf("repos/%s/branches/%s/protection", nwo, bp.Branch)
	cmd := exec.Command(c.ghPath, "api", endpoint, "-X", "PUT", "--input", "-")
	cmd.Stdin = bytes.NewReader(body)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("setting branch protection: %w\n%s", err, string(out))
	}
	return nil
}
