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
	Short: "Write updated changelog based on changes since last release",
	Long: `Updates the existing changelog file with a new release section,
using the commits since the last release.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		updateChangelog(
			changelogArgs.changelogFile,
			vcs.TagOrderBy(updateArgs.orderBy),
			updateArgs.repoPath,
		)
	},
}

func init() {
	changelogCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	updateCmd.Flags().StringVarP(&updateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
}

func updateChangelog(changelogFile string, orderBy vcs.TagOrderBy, repoPath string) {
	_, _, updated := changelog.GetUpdatedChangelog(changelogFile, orderBy, repoPath)

	err := changelog.UpdateChangelog(changelogFile, updated)
	if err != nil {
		panic(fmt.Errorf("failed to update changelog: %w", err))
	}
}
