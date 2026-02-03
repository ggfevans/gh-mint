package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/gvns/gh-repo-defaults/internal/config"
	"github.com/spf13/cobra"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
}

var profilesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tLABELS\tDEFAULT")
		for name, p := range cfg.Profiles {
			def := ""
			if name == cfg.DefaultProfile {
				def = "*"
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", name, p.Description, len(p.Labels.Items), def)
		}
		return w.Flush()
	},
}

var profilesShowCmd = &cobra.Command{
	Use:   "show [profile]",
	Short: "Show profile details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		name := args[0]
		p, ok := cfg.Profiles[name]
		if !ok {
			return fmt.Errorf("profile %q not found", name)
		}
		fmt.Printf("Profile: %s\n", name)
		fmt.Printf("Description: %s\n\n", p.Description)

		fmt.Println("Settings:")
		printBoolSetting := func(label string, v *bool) {
			if v != nil {
				fmt.Printf("  %s: %v\n", label, *v)
			}
		}
		printBoolSetting("Wiki", p.Settings.HasWiki)
		printBoolSetting("Projects", p.Settings.HasProjects)
		printBoolSetting("Delete branch on merge", p.Settings.DeleteBranchOnMerge)
		printBoolSetting("Allow squash merge", p.Settings.AllowSquashMerge)
		printBoolSetting("Allow merge commit", p.Settings.AllowMergeCommit)
		printBoolSetting("Allow rebase merge", p.Settings.AllowRebaseMerge)
		fmt.Println()

		fmt.Printf("Labels (%d):\n", len(p.Labels.Items))
		for _, l := range p.Labels.Items {
			desc := ""
			if l.Description != "" {
				desc = " - " + l.Description
			}
			fmt.Printf("  #%s %s%s\n", l.Color, l.Name, desc)
		}
		fmt.Println()

		if p.Boilerplate.License != "" {
			fmt.Printf("License: %s\n", p.Boilerplate.License)
		}
		if p.Boilerplate.Gitignore != "" {
			fmt.Printf("Gitignore: %s\n", p.Boilerplate.Gitignore)
		}
		if len(p.Boilerplate.Files) > 0 {
			fmt.Println("Boilerplate files:")
			for _, f := range p.Boilerplate.Files {
				fmt.Printf("  %s -> %s\n", f.Src, f.Dest)
			}
		}

		if p.BranchProtection.Branch != "" {
			fmt.Printf("\nBranch protection: %s\n", p.BranchProtection.Branch)
			fmt.Printf("  Required reviews: %d\n", p.BranchProtection.RequiredReviews)
		}

		return nil
	},
}

func loadConfig() (*config.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".config", "gh-repo-defaults", "config.yaml")
	return config.LoadFromFile(configPath)
}

func init() {
	profilesCmd.AddCommand(profilesListCmd)
	profilesCmd.AddCommand(profilesShowCmd)
	rootCmd.AddCommand(profilesCmd)
}
