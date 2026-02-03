package config

import (
	"strings"
	"testing"
)

func TestValidateRepoName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "my-repo", false},
		{"valid dots", "my.repo", false},
		{"valid underscores", "my_repo", false},
		{"empty", "", true},
		{"spaces", "my repo", true},
		{"special chars", "my@repo", true},
		{"path traversal", "../evil", true},
		{"leading dot", ".hidden", true},
		{"leading hyphen", "-repo", true},
		{"too long", strings.Repeat("a", 101), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRepoName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRepoName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLabelColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid lowercase", "d73a4a", false},
		{"valid uppercase", "D73A4A", false},
		{"valid mixed", "aaBB11", false},
		{"too short", "d73a4", true},
		{"too long", "d73a4aa", true},
		{"invalid hex", "zzzzzz", true},
		{"with hash", "#d73a4a", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLabelColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLabelColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProfileName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "oss", false},
		{"valid with dash", "my-profile", false},
		{"valid with underscore", "my_profile", false},
		{"empty", "", true},
		{"spaces", "my profile", true},
		{"special", "my@profile", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProfileName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProfileName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLabelName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "bug", false},
		{"valid with spaces", "good first issue", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 51), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLabelName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLabelName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateNWO(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "owner/repo", false},
		{"valid with dots", "my.org/my.repo", false},
		{"valid with dashes", "my-org/my-repo", false},
		{"empty", "", true},
		{"no slash", "ownerrepo", true},
		{"double slash", "owner//repo", true},
		{"trailing slash", "owner/repo/", true},
		{"path traversal", "owner/../../etc", true},
		{"just slash", "/", true},
		{"leading slash", "/repo", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNWO(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNWO(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid short", "A CLI tool", false},
		{"empty", "", false},
		{"max length", strings.Repeat("a", 350), false},
		{"too long", strings.Repeat("a", 351), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid main", "main", false},
		{"valid with slash", "feature/my-branch", false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 256), true},
		{"spaces", "my branch", true},
		{"special chars", "branch;evil", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBranchName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBranchName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProfile(t *testing.T) {
	t.Run("valid profile", func(t *testing.T) {
		p := Profile{
			Labels: LabelConfig{
				Items: []Label{
					{Name: "bug", Color: "d73a4a"},
				},
			},
		}
		if err := ValidateProfile("oss", p); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid label color", func(t *testing.T) {
		p := Profile{
			Labels: LabelConfig{
				Items: []Label{
					{Name: "bug", Color: "invalid"},
				},
			},
		}
		err := ValidateProfile("oss", p)
		if err == nil {
			t.Error("expected error for invalid label color")
		}
	})

	t.Run("invalid branch protection reviews", func(t *testing.T) {
		p := Profile{
			BranchProtection: BranchProtection{
				Branch:          "main",
				RequiredReviews: 7,
			},
		}
		err := ValidateProfile("oss", p)
		if err == nil {
			t.Error("expected error for required_reviews > 6")
		}
	})

	t.Run("empty boilerplate src", func(t *testing.T) {
		p := Profile{
			Boilerplate: BoilerplateConfig{
				Files: []BoilerplateFile{
					{Src: "", Dest: "README.md"},
				},
			},
		}
		err := ValidateProfile("oss", p)
		if err == nil {
			t.Error("expected error for empty boilerplate src")
		}
	})
}
