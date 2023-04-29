/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/

package cmd

import (
	"github.com/outofcoffee/since/convcommits"
	"github.com/outofcoffee/since/semver"
	"github.com/outofcoffee/since/vcs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var semverArgs struct {
	orderBy  string
	repoPath string
	tag      string
}

// semverCmd represents the semver command
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Get the next semantic version based on changes since last tag",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Returns the next semantic version
based on the changes.

Changes influence the version according to
conventional commits: https://www.conventionalcommits.org/en/v1.0.0/`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		nextVersion := determineNextVersion(semverArgs.repoPath, semverArgs.tag, vcs.TagOrderBy(semverArgs.orderBy))
		if nextVersion == "" {
			os.Exit(1)
		}
		println(nextVersion)
	},
}

func init() {
	rootCmd.AddCommand(semverCmd)

	semverCmd.Flags().StringVarP(&semverArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	semverCmd.Flags().StringVarP(&semverArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	semverCmd.Flags().StringVarP(&semverArgs.tag, "tag", "t", "", "Include commits after this tag")
}

func determineNextVersion(repoPath string, tag string, orderBy vcs.TagOrderBy) string {
	currentVersion, vPrefix := getCurrentVersion(repoPath, orderBy)

	commits, err := vcs.FetchCommitMessages(repoPath, tag, orderBy)
	if err != nil {
		panic(err)
	}
	return getNextVersion(currentVersion, vPrefix, commits)
}

func getCurrentVersion(repoPath string, orderBy vcs.TagOrderBy) (version string, vPrefix bool) {
	version, err := vcs.GetLatestTag(repoPath, orderBy)
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
		vPrefix = true
	}
	logrus.Tracef("current version: %s", version)
	return version, vPrefix
}

func getNextVersion(currentVersion string, vPrefix bool, commits []string) string {
	types := convcommits.DetermineTypes(commits)
	logrus.Debugf("commit types: %v", types)

	changeType := semver.DetermineChangeType(types)
	if changeType == semver.ComponentNone {
		logrus.Warnf("no changes detected")
		return ""
	}

	nextVersion := semver.BumpVersion(currentVersion, changeType)
	if vPrefix {
		nextVersion = "v" + nextVersion
	}
	return nextVersion
}
