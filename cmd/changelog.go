/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Commands related to changelog files",
	Long:  `Parse and update changelog files.`,
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
