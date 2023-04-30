/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Commands related to the project",
	Long: `List the changes since the last release in the project
repository, or determine the next semantic version based on
those changes.`,
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
