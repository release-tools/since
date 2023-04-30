/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/

package cmd

import (
	"github.com/outofcoffee/since/semver"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
	"os"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the next semantic version based on changes since last tag",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Returns the next semantic version
based on the changes.

Changes influence the version according to
conventional commits: https://www.conventionalcommits.org/en/v1.0.0/`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		nextVersion := determineNextVersion(projectArgs.repoPath, projectArgs.tag, vcs.TagOrderBy(projectArgs.orderBy))
		if nextVersion == "" {
			os.Exit(1)
		}
		println(nextVersion)
	},
}

func init() {
	projectCmd.AddCommand(versionCmd)
}

func determineNextVersion(repoPath string, tag string, orderBy vcs.TagOrderBy) string {
	currentVersion, vPrefix := semver.GetCurrentVersion(repoPath, orderBy)

	commits, err := vcs.FetchCommitMessages(repoPath, tag, orderBy)
	if err != nil {
		panic(err)
	}
	return semver.GetNextVersion(currentVersion, vPrefix, commits)
}
