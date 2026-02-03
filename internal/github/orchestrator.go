package github

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gvns/gh-repo-defaults/internal/config"
	"github.com/gvns/gh-repo-defaults/internal/scaffold"
)

// StepStatus represents the outcome of a single step.
type StepStatus struct {
	Name    string
	Success bool
	Message string
	Err     error
}

// ProgressFunc is called after each step completes.
type ProgressFunc func(StepStatus)

// CreateOpts holds all options for creating a repo with defaults.
type CreateOpts struct {
	Name        string
	Description string
	Public      bool
	Profile     config.Profile
	Owner       string // optional, for org repos
	OnProgress  ProgressFunc
}

func (o *CreateOpts) nwo() string {
	if o.Owner != "" {
		return o.Owner + "/" + o.Name
	}
	return o.Name
}

func (o *CreateOpts) report(name string, err error) {
	if o.OnProgress == nil {
		return
	}
	s := StepStatus{Name: name, Success: err == nil}
	if err != nil {
		s.Err = err
		s.Message = err.Error()
	}
	o.OnProgress(s)
}

// CreateWithDefaults creates a repo and applies all profile defaults.
func (c *Client) CreateWithDefaults(opts CreateOpts) (string, error) {
	repoArg := opts.Name
	if opts.Owner != "" {
		repoArg = opts.Owner + "/" + opts.Name
	}
	url, err := c.CreateRepo(repoArg, opts.Description, opts.Public)
	opts.report("Created repository", err)
	if err != nil {
		return "", err
	}

	nwo := opts.nwo()
	if url != "" && nwo == opts.Name {
		parts := splitRepoURL(url)
		if parts != "" {
			nwo = parts
		}
	}

	var errs []error

	// Apply settings
	settings, err := SettingsFromRepoSettings(opts.Profile.Settings)
	if err != nil {
		opts.report("Applied repo settings", err)
		errs = append(errs, err)
	} else {
		err = c.UpdateSettings(nwo, settings)
		opts.report("Applied repo settings", err)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Sync labels
	deleted, created, labelErrs := c.SyncLabels(nwo, opts.Profile.Labels)
	var labelErr error
	if len(labelErrs) > 0 {
		labelErr = fmt.Errorf("%d label errors", len(labelErrs))
		errs = append(errs, labelErr)
	}
	opts.report(fmt.Sprintf("Synced labels (-%d/+%d)", deleted, created), labelErr)

	// Scaffold boilerplate
	if len(opts.Profile.Boilerplate.Files) > 0 {
		err = c.scaffoldAndPush(nwo, opts.Profile.Boilerplate, opts.Name)
		opts.report("Pushed boilerplate files", err)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Branch protection
	if opts.Profile.BranchProtection.Branch != "" {
		err = c.SetBranchProtection(nwo, opts.Profile.BranchProtection)
		opts.report("Set branch protection", err)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return url, fmt.Errorf("%d step(s) failed after repo creation", len(errs))
	}
	return url, nil
}

func (c *Client) scaffoldAndPush(nwo string, bp config.BoilerplateConfig, repoName string) error {
	tmpDir, err := os.MkdirTemp("", "gh-repo-defaults-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cloneDir := filepath.Join(tmpDir, repoName)
	if _, err := c.run("repo", "clone", nwo, cloneDir); err != nil {
		return fmt.Errorf("cloning repo: %w", err)
	}

	var userTemplateDir string
	if home, err := os.UserHomeDir(); err == nil {
		userTemplateDir = filepath.Join(home, ".config", "gh-repo-defaults", "templates")
	}

	if _, err := scaffold.PrepareBoilerplate(bp, cloneDir, userTemplateDir); err != nil {
		return fmt.Errorf("preparing boilerplate: %w", err)
	}

	gitCmd := func(args ...string) error {
		cmd := exec.Command("git", args...)
		cmd.Dir = cloneDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git %v: %w\n%s", args, err, string(out))
		}
		return nil
	}
	if err := gitCmd("add", "-A"); err != nil {
		return err
	}
	if err := gitCmd("commit", "-m", "chore: add boilerplate files"); err != nil {
		return err
	}
	if err := gitCmd("push"); err != nil {
		return err
	}
	return nil
}

func splitRepoURL(url string) string {
	nwo, ok := strings.CutPrefix(url, "https://github.com/")
	if !ok || nwo == "" {
		return ""
	}
	return strings.TrimSuffix(nwo, "/")
}

// SettingsFromRepoSettings converts config settings to a map for the API.
func SettingsFromRepoSettings(s config.RepoSettings) (map[string]interface{}, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("marshaling settings: %w", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("converting settings: %w", err)
	}
	return m, nil
}
