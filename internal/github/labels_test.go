package github

import (
	"strings"
	"testing"

	"github.com/gvns/gh-repo-defaults/internal/config"
)

func TestCreateLabelArgs(t *testing.T) {
	c := NewClient()
	label := config.Label{Name: "bug", Color: "d73a4a", Description: "Something isn't working"}
	args := c.createLabelArgs("owner/repo", label)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "label") || !strings.Contains(joined, "create") {
		t.Errorf("missing label create: %v", args)
	}
	if !strings.Contains(joined, "bug") {
		t.Errorf("missing label name: %v", args)
	}
	if !strings.Contains(joined, "d73a4a") {
		t.Errorf("missing color: %v", args)
	}
}

func TestDeleteLabelArgs(t *testing.T) {
	c := NewClient()
	args := c.deleteLabelArgs("owner/repo", "bug")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "label") || !strings.Contains(joined, "delete") {
		t.Errorf("missing label delete: %v", args)
	}
	if !strings.Contains(joined, "bug") {
		t.Errorf("missing label name: %v", args)
	}
	if !strings.Contains(joined, "--yes") {
		t.Errorf("missing --yes flag: %v", args)
	}
}

func TestListLabelsArgs(t *testing.T) {
	c := NewClient()
	args := c.listLabelsArgs("owner/repo")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "label") || !strings.Contains(joined, "list") {
		t.Errorf("missing label list: %v", args)
	}
}
