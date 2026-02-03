package cmd

import (
	"fmt"

	"github.com/gvns/gh-repo-defaults/internal/config"
	ghclient "github.com/gvns/gh-repo-defaults/internal/github"
	"github.com/spf13/cobra"
)

var applyProfile string

var applyCmd = &cobra.Command{
	Use:   "apply [owner/repo]",
	Short: "Apply a profile to an existing repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nwo := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		profileName := applyProfile
		if profileName == "" {
			profileName = cfg.DefaultProfile
		}
		profile, ok := cfg.Profiles[profileName]
		if !ok {
			return fmt.Errorf("profile %q not found", profileName)
		}
		if err := config.ValidateProfile(profileName, profile); err != nil {
			return err
		}

		client := ghclient.NewClient()
		if err := client.CheckInstalled(); err != nil {
			return err
		}

		progress := func(name string, err error) {
			if err == nil {
				fmt.Printf("  ✓ %s\n", name)
			} else {
				fmt.Printf("  ✗ %s: %s\n", name, err.Error())
			}
		}

		fmt.Printf("Applying profile %q to %s...\n", profileName, nwo)

		// Apply settings
		settings := ghclient.SettingsFromRepoSettings(profile.Settings)
		err = client.UpdateSettings(nwo, settings)
		progress("Applied repo settings", err)

		// Sync labels
		deleted, created, labelErrs := client.SyncLabels(nwo, profile.Labels)
		var labelErr error
		if len(labelErrs) > 0 {
			labelErr = fmt.Errorf("%d label errors", len(labelErrs))
		}
		progress(fmt.Sprintf("Synced labels (-%d/+%d)", deleted, created), labelErr)

		// Branch protection
		if profile.BranchProtection.Branch != "" {
			err = client.SetBranchProtection(nwo, profile.BranchProtection)
			progress("Set branch protection", err)
		}

		fmt.Println("\nDone!")
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyProfile, "profile", "p", "", "Profile to apply (default: from config)")
	rootCmd.AddCommand(applyCmd)
}
