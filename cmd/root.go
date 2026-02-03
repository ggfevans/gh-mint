package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-repo-defaults",
	Short: "Create GitHub repos with consistent defaults",
	Long:  "A gh CLI extension that creates GitHub repos and applies labels, settings, boilerplate, and branch protection from named profiles.",
}

func Execute() error {
	return rootCmd.Execute()
}
