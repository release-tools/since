/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
)

var updateArgs struct {
	orderBy  string
	repoPath string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Print updated changelog based on changes since last release",
	Long: `Generates a new changelog based on an existing changelog file,
using the commits since the last release.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, _, updated := changelog.GetUpdatedChangelog(
			changelogArgs.changelogFile,
			vcs.TagOrderBy(updateArgs.orderBy),
			updateArgs.repoPath,
		)
		fmt.Println(updated)
	},
}

func init() {
	changelogCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	updateCmd.Flags().StringVarP(&updateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
}
