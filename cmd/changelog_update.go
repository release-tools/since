/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/semver"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
	"time"
)

var updateArgs struct {
	changelogFile string
	orderBy       string
	repoPath      string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Print updated changelog based on changes since last release",
	Long: `Generates a new changelog files based on an existing file,
using the changes since the last release in the given project repository.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		updated := getUpdatedChangelog(updateArgs.changelogFile, vcs.TagOrderBy(updateArgs.orderBy), updateArgs.repoPath)
		cmd.Println(updated)
	},
}

func init() {
	changelogCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.changelogFile, "changelog", "c", "CHANGELOG.md", "Path to changelog file")
	updateCmd.PersistentFlags().StringVarP(&updateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	updateCmd.PersistentFlags().StringVarP(&updateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
}

func getUpdatedChangelog(changelogFile string, orderBy vcs.TagOrderBy, repoPath string) string {
	changes, err := vcs.FetchCommitMessages(repoPath, "", orderBy)
	if err != nil {
		panic(err)
	}
	rendered := changelog.RenderCommits(changes)

	lines, err := changelog.ReadFile(changelogFile)
	if err != nil {
		panic(err)
	}

	currentVersion, _ := semver.GetCurrentVersion(repoPath, orderBy)

	commits, err := vcs.FetchCommitMessages(repoPath, "", orderBy)
	if err != nil {
		panic(err)
	}
	nextVersion := semver.GetNextVersion(currentVersion, false, commits)

	if nextVersion == "" {
		panic("Could not determine next version")
	}
	versionHeader := "## [" + nextVersion + "] - " + time.Now().UTC().Format("2006-01-02") + "\n"

	sections := changelog.SplitIntoSections(lines)
	output := sections.Boilerplate + versionHeader + rendered + sections.Body

	return output
}
