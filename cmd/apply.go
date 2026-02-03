package cmd

import (
	"fmt"

	"github.com/ggfevans/gh-mint/internal/config"
	ghclient "github.com/ggfevans/gh-mint/internal/github"
	"github.com/spf13/cobra"
)

var applyProfile string

var applyCmd = &cobra.Command{
	Use:   "apply [owner/repo]",
	Short: "Apply a profile to an existing repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nwo := args[0]

		if err := config.ValidateNWO(nwo); err != nil {
			return err
		}

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

		var failures int
		progress := func(name string, err error) {
			if err == nil {
				fmt.Printf("  ✓ %s\n", name)
			} else {
				fmt.Printf("  ✗ %s: %s\n", name, err.Error())
				failures++
			}
		}

		fmt.Printf("Applying profile %q to %s...\n", profileName, nwo)

		// Apply settings
		settings, err := ghclient.SettingsFromRepoSettings(profile.Settings)
		if err != nil {
			progress("Applied repo settings", err)
		} else {
			err = client.UpdateSettings(nwo, settings)
			progress("Applied repo settings", err)
		}

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

		if failures > 0 {
			fmt.Printf("\nCompleted with %d error(s).\n", failures)
			return fmt.Errorf("%d step(s) failed", failures)
		}
		fmt.Println("\nDone!")
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyProfile, "profile", "p", "", "Profile to apply (default: from config)")
	rootCmd.AddCommand(applyCmd)
}
