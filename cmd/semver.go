/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/

package cmd

import (
	"github.com/outofcoffee/changelog-parser/semver"
	"github.com/outofcoffee/changelog-parser/vcs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"strconv"
	"strings"
)

var semverArgs struct {
	repoPath string
	tag      string
}

// semverCmd represents the semver command
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Get the next semantic version based on changelog",
	Long: `Parse changelog file and return the next semantic version
based on the changes. Changes influence the version according to
conventional commits: https://www.conventionalcommits.org/en/v1.0.0/`,
	Run: func(cmd *cobra.Command, args []string) {
		determineNextVersion(semverArgs.repoPath, semverArgs.tag)
	},
}

func init() {
	rootCmd.AddCommand(semverCmd)

	semverCmd.Flags().StringVarP(&semverArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	semverCmd.Flags().StringVarP(&semverArgs.tag, "tag", "t", "", "Search for changes after this tag")
}

func determineNextVersion(repoPath string, tag string) {
	latestVersion := getLatestVersion(repoPath)

	changes, err := fetchCommitMessages(repoPath, tag)
	if err != nil {
		panic(err)
	}
	types := determineTypes(changes)
	logrus.Debugf("types: %v", types)

	if semver.HasMajor(types) {
		bumpVersion(latestVersion, semver.SemverMajor)
	} else if semver.HasMinor(types) {
		bumpVersion(latestVersion, semver.SemverMinor)
	} else if semver.HasPatch(types) {
		bumpVersion(latestVersion, semver.SemverPatch)
	} else {
		logrus.Warnf("unable to determine next version from changes")
	}
}

func bumpVersion(version string, component semver.SemverComponent) string {
	logrus.Debugf("bumping %v version", component)
	components := strings.Split(version, ".")
	switch component {
	case semver.SemverMajor:
		components[0] = bump(components[0])
		components[1] = "0"
		components[2] = "0"
		break

	case semver.SemverMinor:
		components[1] = bump(components[1])
		components[2] = "0"
		break

	case semver.SemverPatch:
		components[2] = bump(components[2])
		break
	}
	nextVersion := strings.Join(components, ".")
	logrus.Infof("next version: %s", nextVersion)
	return nextVersion
}

func bump(v string) string {
	num, _ := strconv.Atoi(v)
	return strconv.Itoa(num + 1)
}

func getLatestVersion(repoPath string) string {
	latestVersion, err := vcs.GetLatestTag(repoPath)
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(latestVersion, "v") {
		latestVersion = strings.TrimPrefix(latestVersion, "v")
	}
	logrus.Tracef("latest version: %s", latestVersion)
	return latestVersion
}

func fetchCommitMessages(repoPath string, tag string) ([]string, error) {
	if tag == "" {
		latestTag, err := vcs.GetLatestTag(repoPath)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("latest tag: %s", latestTag)
		tag = latestTag
	}
	commits, err := vcs.FetchCommitsAfter(repoPath, tag)
	if err != nil {
		return nil, err
	}
	logrus.Tracef("commits: %v", commits)
	return commits, nil
}

func determineTypes(changes []string) []string {
	types := make(map[string]bool)
	for _, change := range changes {
		parts := strings.Split(change, ":")
		if len(parts) < 2 {
			continue
		}
		prefix := strings.TrimSpace(parts[0])
		if strings.Contains(prefix, "(") {
			prefix = strings.Split(prefix, "(")[0]
		}
		if !types[prefix] {
			types[prefix] = true
		}
	}
	return maps.Keys(types)
}
