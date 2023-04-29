package convcommits

import (
	"golang.org/x/exp/maps"
	"strings"
)

func CategoriseByType(commits []string) map[string][]string {
	categorised := make(map[string][]string)
	for _, commit := range commits {
		parts := strings.Split(commit, ":")
		if len(parts) < 2 {
			continue
		}
		prefix := strings.TrimSpace(parts[0])
		if strings.HasSuffix(prefix, "!") {
			prefix = "BREAKING CHANGE"
		}
		if strings.Contains(prefix, "(") {
			prefix = strings.Split(prefix, "(")[0]
		}

		category := categorised[prefix]
		category = append(category, commit)
		categorised[prefix] = category
	}
	return categorised
}

func DetermineTypes(commits []string) []string {
	return maps.Keys(CategoriseByType(commits))
}
