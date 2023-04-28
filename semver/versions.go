package semver

import "strings"

type SemverComponent string

const (
	SemverMajor SemverComponent = "major"
	SemverMinor                 = "minor"
	SemverPatch                 = "patch"
)

func HasMajor(types []string) bool {
	return containsIgnoreCase(types, "breaking change")
}

func HasMinor(types []string) bool {
	return containsIgnoreCase(types, "feat")
}

func HasPatch(types []string) bool {
	return containsIgnoreCase(types, "build", "chore", "ci", "docs", "fix", "refactor", "style", "test")
}

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
