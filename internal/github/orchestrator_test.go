package github

import (
	"testing"

	"github.com/ggfevans/gh-mint/internal/config"
)

func TestSettingsFromRepoSettings(t *testing.T) {
	f := false
	tr := true
	s := config.RepoSettings{
		HasWiki:             &f,
		HasProjects:         &f,
		DeleteBranchOnMerge: &tr,
		AllowSquashMerge:    &tr,
	}
	m, err := SettingsFromRepoSettings(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["has_wiki"] != false {
		t.Error("has_wiki should be false")
	}
	if m["delete_branch_on_merge"] != true {
		t.Error("delete_branch_on_merge should be true")
	}
	// Verify unset fields are omitted
	if _, ok := m["allow_merge_commit"]; ok {
		t.Error("allow_merge_commit should be omitted when nil")
	}
	if _, ok := m["allow_rebase_merge"]; ok {
		t.Error("allow_rebase_merge should be omitted when nil")
	}
}

func TestSettingsFromRepoSettings_OmitsNilFields(t *testing.T) {
	s := config.RepoSettings{} // all nil
	m, err := SettingsFromRepoSettings(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map for zero-value settings, got %d keys: %v", len(m), m)
	}
}

func TestSplitRepoURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://github.com/owner/repo", "owner/repo"},
		{"https://github.com/org/my-tool", "org/my-tool"},
		{"https://github.com/owner/repo/", "owner/repo"},
		{"something-else", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := splitRepoURL(tt.input)
		if got != tt.want {
			t.Errorf("splitRepoURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
