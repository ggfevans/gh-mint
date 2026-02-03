# gh-repo-defaults Design Document

**Date:** 2026-02-03
**Status:** Draft

## Problem

Creating new GitHub repos with consistent configuration is tedious. Each new repo requires manually:
- Deleting default labels and creating custom ones
- Disabling wiki, projects, and other unused features
- Setting merge strategies and branch protection
- Adding boilerplate files (LICENSE, .gitignore, CI workflows)

This is especially painful when creating multiple repo types: personal projects, open source libraries, and GitHub Actions.

## Solution

A `gh` CLI extension called `gh-repo-defaults` that creates a new GitHub repo and applies a full set of defaults from a named profile in one command, with an interactive TUI for guided creation.

## User Flow

### TUI Mode

```
$ gh repo-defaults
```

1. Mode selection: Create new repo / Apply to existing / Manage profiles
2. Create screen: name, description, visibility, profile selection
3. Preview: labels, settings, boilerplate to be applied
4. Confirm and execute
5. Progress screen with per-step status

### CLI Mode (Scripting)

```
$ gh repo-defaults create my-tool --profile oss --public
$ gh repo-defaults apply my-existing-repo --profile action
$ gh repo-defaults profiles list
$ gh repo-defaults profiles show oss
```

## TUI Screens

### Screen 1 - Mode Selection

```
gh repo-defaults
-----------------
> Create new repo
  Apply to existing repo
  Manage profiles
  Edit config
```

### Screen 2 - Create New Repo

```
Create Repository
-----------------
Name:        [my-new-tool          ]
Description: [A CLI tool for...     ]
Visibility:  * Public  o Private
Profile:     [oss v]

  Labels:           12 labels
  Settings:         wiki off, projects off, squash merge
  Boilerplate:      LICENSE, .gitignore, CI workflow
  Branch protection: main (require PR reviews)

  [Create]  [Back]
```

### Screen 3 - Progress

```
Creating repository...
  ok Created gvns/my-new-tool
  ok Applied repo settings
  ok Removed default labels
  ok Created 12 custom labels
  ok Pushed boilerplate files
  ok Set branch protection

Done! https://github.com/gvns/my-new-tool
```

## Configuration

### Location

```
~/.config/gh-repo-defaults/config.yaml       # Main config + profiles
~/.config/gh-repo-defaults/templates/         # User template overrides
```

Config files are created with `0600` permissions. The tool warns if existing config has overly permissive permissions (anything beyond owner read/write).

### Config Schema

```yaml
# Default profile to use when none specified
default_profile: oss

# GitHub org to create repos in (optional, defaults to personal account)
default_owner: ""

profiles:
  oss:
    description: "Open source project defaults"
    settings:
      has_wiki: false
      has_projects: false
      has_discussions: false
      delete_branch_on_merge: true
      allow_squash_merge: true
      allow_merge_commit: false
      allow_rebase_merge: false
      squash_merge_commit_title: "PR_TITLE"
      squash_merge_commit_message: "PR_BODY"
    labels:
      # Set clear_existing: true to remove GitHub's default labels first
      clear_existing: true
      items:
        - name: bug
          color: "d73a4a"
          description: "Something isn't working"
        - name: enhancement
          color: "a2eeef"
          description: "New feature or request"
        - name: documentation
          color: "0075ca"
          description: "Improvements or additions to docs"
        - name: good first issue
          color: "7057ff"
          description: "Good for newcomers"
        - name: help wanted
          color: "008672"
          description: "Extra attention is needed"
        - name: wontfix
          color: "ffffff"
          description: "This will not be worked on"
    boilerplate:
      license: MIT
      gitignore: Go
      files:
        - src: contributing.md
          dest: CONTRIBUTING.md
        - src: ci.yml
          dest: .github/workflows/ci.yml
    branch_protection:
      branch: main
      required_reviews: 0
      dismiss_stale_reviews: true
      require_status_checks: false

  personal:
    description: "Personal project - minimal config"
    settings:
      has_wiki: false
      has_projects: false
      delete_branch_on_merge: true
      allow_squash_merge: true
      allow_merge_commit: true
      allow_rebase_merge: false
    labels:
      clear_existing: true
      items:
        - name: bug
          color: "d73a4a"
        - name: enhancement
          color: "a2eeef"
        - name: chore
          color: "fef2c0"
    boilerplate:
      license: MIT
      gitignore: Go

  action:
    description: "GitHub Action defaults"
    settings:
      has_wiki: false
      has_projects: false
      delete_branch_on_merge: true
      allow_squash_merge: true
    labels:
      clear_existing: true
      items:
        - name: bug
          color: "d73a4a"
        - name: enhancement
          color: "a2eeef"
        - name: breaking change
          color: "e11d48"
          description: "Introduces a breaking change"
    boilerplate:
      license: MIT
      files:
        - src: action.yml
          dest: action.yml
        - src: action-ci.yml
          dest: .github/workflows/ci.yml
        - src: action-release.yml
          dest: .github/workflows/release.yml
```

## Architecture

### Project Structure

```
gh-repo-defaults/
  main.go                 # Entry point, cobra root command
  cmd/
    root.go               # Root command, launches TUI if no subcommand
    create.go             # `create` subcommand
    apply.go              # `apply` subcommand (existing repos)
    profiles.go           # `profiles list/show` subcommand
  internal/
    config/
      config.go           # YAML config loading, validation, defaults
      validation.go       # Input validation (repo names, colors, paths)
    github/
      client.go           # Wraps gh CLI execution
      repo.go             # Create repo, update settings
      labels.go           # Delete/create labels
      protection.go       # Branch protection rules
    scaffold/
      boilerplate.go      # Generate/copy boilerplate files
      templates.go        # Embedded + user template resolution
    tui/
      app.go              # Bubble Tea app model, screen routing
      mode.go             # Mode selection screen
      create.go           # Create repo form screen
      progress.go         # Progress/execution screen
      styles.go           # Lip Gloss styles
  templates/              # Embedded default templates (via go:embed)
    contributing.md
    ci.yml
    action.yml
    action-ci.yml
    action-release.yml
  go.mod
  go.sum
```

### Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | TUI styling |
| `github.com/charmbracelet/huh` | TUI form components |
| `gopkg.in/yaml.v3` | Config parsing |

### GitHub Interaction

All GitHub API calls go through the `gh` CLI rather than using the GitHub API directly:

- `gh repo create` - Create repository
- `gh api repos/{owner}/{repo}` - Update repo settings (PATCH)
- `gh api repos/{owner}/{repo}/labels` - Manage labels
- `gh api repos/{owner}/{repo}/branches/{branch}/protection` - Branch protection
- `gh label list` / `gh label create` / `gh label delete` - Label management

This approach piggybacks on the user's existing `gh` authentication. No token management needed.

All `gh` commands are executed via Go's `exec.Command` with argument arrays (never string interpolation through a shell).

### Boilerplate Strategy

1. Default templates are embedded in the binary using Go's `embed` package
2. Users can override any template by placing a file with the same name in `~/.config/gh-repo-defaults/templates/`
3. User templates take precedence over embedded defaults
4. Template paths are canonicalized and validated to prevent path traversal

### Execution Flow (Create)

```
1. Load config from ~/.config/gh-repo-defaults/config.yaml
2. Validate config (repo name, label colors, template paths)
3. Check gh auth status and token scopes
4. Create repo: gh repo create {name} --public/--private
5. Apply settings: gh api PATCH /repos/{owner}/{repo}
6. Clear existing labels (if configured): gh label delete ...
7. Create labels: gh label create ...
8. Clone repo locally, copy boilerplate, commit, push
9. Apply branch protection: gh api PUT /repos/{owner}/{repo}/branches/{branch}/protection
```

## Security Requirements

### Command Injection Prevention
- All `gh` CLI calls use `exec.Command` with argument arrays
- Never use `exec.Command("sh", "-c", ...)` with string interpolation
- Each argument is passed as a separate string to `exec.Command`

### Input Validation
- Repo names: `^[a-zA-Z0-9._-]+$`
- Label colors: `^[0-9a-fA-F]{6}$`
- Label names: max 50 chars, printable characters only
- Description: max 350 chars
- Profile names: `^[a-zA-Z0-9_-]+$`

### Path Traversal Prevention
- All template file paths are canonicalized with `filepath.Abs` + `filepath.EvalSymlinks`
- Resolved paths must share a common prefix with the template directory
- Reject paths containing `..` before resolution as an extra layer

### Config File Security
- Created with `0600` permissions (owner read/write only)
- Tool warns on startup if config has permissions beyond `0600`
- Config directory created with `0700`

### Token Scope Awareness
- Check `gh auth status` at startup
- Warn if `repo` scope is missing (needed for repo creation/settings)
- Warn if `admin:repo` scope is missing when branch protection is configured
- Fail gracefully with clear messages rather than cryptic API errors

### YAML Parsing
- Use `gopkg.in/yaml.v3` (safe against entity expansion)
- Set reasonable size limit on config file (e.g., 1MB) before parsing

## Edge Cases

- **Repo already exists:** Detect and offer to apply profile to existing repo instead
- **Partial failure:** If label creation fails mid-way, report which succeeded and which failed; don't leave repo in unknown state
- **No gh installed:** Detect early with clear error message
- **No config file:** Use embedded defaults, offer to generate config on first run
- **Rate limiting:** GitHub API rate limits apply; report remaining quota if close to limit

## Future Considerations (Not in v1)

- `gh repo-defaults init` to generate a starter config interactively
- Profile inheritance (e.g., `action` extends `oss`)
- Org-level shared configs (fetch profiles from a central repo)
- Dry-run mode to preview all API calls without executing
