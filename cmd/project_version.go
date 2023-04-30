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

var versionArgs struct {
	orderBy  string
	repoPath string
	tag      string
}

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
		nextVersion := determineNextVersion(versionArgs.repoPath, versionArgs.tag, vcs.TagOrderBy(versionArgs.orderBy))
		if nextVersion == "" {
			os.Exit(1)
		}
		println(nextVersion)
	},
}

func init() {
	projectCmd.AddCommand(versionCmd)

	versionCmd.Flags().StringVarP(&versionArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	versionCmd.Flags().StringVarP(&versionArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	versionCmd.Flags().StringVarP(&versionArgs.tag, "tag", "t", "", "Include commits after this tag")
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
