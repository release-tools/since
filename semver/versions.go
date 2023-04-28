package semver

import (
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

// BumpVersion bumps the version based on the component.
func BumpVersion(version string, component Component) string {
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
