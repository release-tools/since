package semver

import (
	"github.com/outofcoffee/since/convcommits"
	"github.com/outofcoffee/since/vcs"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Component string

const (
	ComponentMajor Component = "major"
	ComponentMinor Component = "minor"
	ComponentPatch Component = "patch"
	ComponentNone  Component = "none"
)

// GetCurrentVersion gets the current version from the repo.
func GetCurrentVersion(repoPath string, orderBy vcs.TagOrderBy) (version string, vPrefix bool) {
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

// GetNextVersion gets the next version based on the current version and the commit messages.
func GetNextVersion(currentVersion string, vPrefix bool, commits []string) string {
	types := convcommits.DetermineTypes(commits)
	logrus.Debugf("commit types: %v", types)

	changeType := DetermineChangeType(types)
	if changeType == ComponentNone {
		logrus.Warnf("no changes detected")
		return ""
	}

	nextVersion := bumpVersion(currentVersion, changeType)
	if vPrefix {
		nextVersion = "v" + nextVersion
	}
	return nextVersion
}

// bumpVersion bumps the version based on the component.
func bumpVersion(version string, component Component) string {
	logrus.Debugf("bumping %v version", component)
	components := strings.Split(version, ".")

	switch component {
	case ComponentMajor:
		components[0] = bump(components[0])
		components[1] = "0"
		components[2] = "0"
		break

	case ComponentMinor:
		components[1] = bump(components[1])
		components[2] = "0"
		break

	case ComponentPatch:
		components[2] = bump(components[2])
		break
	}
	nextVersion := strings.Join(components, ".")
	return nextVersion
}

// bump bumps the version by 1.
func bump(v string) string {
	num, _ := strconv.Atoi(v)
	return strconv.Itoa(num + 1)
}

// DetermineChangeType determines the type of change based on the commit messages.
func DetermineChangeType(types []string) Component {
	if containsIgnoreCase(types, "breaking change") {
		return ComponentMajor
	} else if containsIgnoreCase(types, "feat") {
		return ComponentMinor
	} else if containsIgnoreCase(types, "build", "chore", "ci", "docs", "fix", "refactor", "style", "test") {
		return ComponentPatch
	} else {
		logrus.Warnf("unable to determine next version from changes")
		return ComponentNone
	}
}

// containsIgnoreCase returns true if the orig slice contains any of the search strings,
// compared in a case-insensitive manner.
func containsIgnoreCase(orig []string, search ...string) bool {
	for _, o := range orig {
		for _, s := range search {
			if strings.ToLower(o) == strings.ToLower(s) {
				return true
			}
		}
	}
	return false
}
