package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gvns/gh-repo-defaults/internal/config"
)

func (c *Client) createLabelArgs(nwo string, label config.Label) []string {
	args := []string{"label", "create", label.Name, "--color", label.Color, "--repo", nwo}
	if label.Description != "" {
		args = append(args, "--description", label.Description)
	}
	return args
}

func (c *Client) deleteLabelArgs(nwo string, name string) []string {
	return []string{"label", "delete", name, "--repo", nwo, "--yes"}
}

func (c *Client) listLabelsArgs(nwo string) []string {
	return []string{"label", "list", "--repo", nwo, "--json", "name"}
}

func (c *Client) CreateLabel(nwo string, label config.Label) error {
	args := c.createLabelArgs(nwo, label)
	if _, err := c.run(args...); err != nil {
		return fmt.Errorf("creating label %q: %w", label.Name, err)
	}
	return nil
}

func (c *Client) DeleteLabel(nwo string, name string) error {
	args := c.deleteLabelArgs(nwo, name)
	if _, err := c.run(args...); err != nil {
		return fmt.Errorf("deleting label %q: %w", name, err)
	}
	return nil
}

func (c *Client) ListLabels(nwo string) ([]string, error) {
	args := c.listLabelsArgs(nwo)
	out, err := c.run(args...)
	if err != nil {
		return nil, fmt.Errorf("listing labels: %w", err)
	}
	var labels []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(out), &labels); err != nil {
		return nil, fmt.Errorf("parsing labels: %w", err)
	}
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return names, nil
}

func (c *Client) SyncLabels(nwo string, cfg config.LabelConfig) (deleted int, created int, errs []error) {
	if cfg.ClearExisting {
		existing, err := c.ListLabels(nwo)
		if err != nil {
			errs = append(errs, err)
			return
		}
		for _, name := range existing {
			if err := c.DeleteLabel(nwo, name); err != nil {
				errs = append(errs, err)
			} else {
				deleted++
			}
		}
	}
	for _, label := range cfg.Items {
		if err := c.CreateLabel(nwo, label); err != nil {
			errs = append(errs, err)
		} else {
			created++
		}
	}
	return
}

func LabelSummary(cfg config.LabelConfig) string {
	names := make([]string, len(cfg.Items))
	for i, l := range cfg.Items {
		names[i] = l.Name
	}
	return fmt.Sprintf("%d labels (%s)", len(cfg.Items), strings.Join(names, ", "))
}
