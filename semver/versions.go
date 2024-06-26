/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package semver

import (
	"github.com/release-tools/since/convcommits"
	"github.com/release-tools/since/stringutil"
	"github.com/release-tools/since/vcs"
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
	logrus.Tracef("bumping %v component", component)
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
	logrus.Debugf("bumped %v component - new version %v", component, nextVersion)
	return nextVersion
}

// bump bumps the version by 1.
func bump(v string) string {
	num, _ := strconv.Atoi(v)
	return strconv.Itoa(num + 1)
}

// DetermineChangeType determines the type of change based on the commit messages.
func DetermineChangeType(types []string) Component {
	if stringutil.ContainsIgnoreCase(types, "breaking change") {
		return ComponentMajor
	} else if stringutil.ContainsIgnoreCase(types, "feat") {
		return ComponentMinor
	} else if stringutil.ContainsIgnoreCase(types, "build", "chore", "ci", "docs", "fix", "refactor", "style", "test") {
		return ComponentPatch
	} else {
		logrus.Warnf("unable to determine next version from changes")
		return ComponentNone
	}
}
