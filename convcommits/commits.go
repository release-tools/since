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

package convcommits

import (
	"golang.org/x/exp/maps"
	"strings"
)

func CategoriseByType(commits []string) map[string][]string {
	categorised := make(map[string][]string)
	for _, commit := range commits {
		parts := strings.Split(commit, ":")

		var prefix string
		if len(parts) >= 2 {
			prefix = strings.TrimSpace(parts[0])
			if strings.HasSuffix(prefix, "!") {
				prefix = "BREAKING CHANGE"
			}
			if strings.Contains(prefix, "(") {
				prefix = strings.Split(prefix, "(")[0]
			}
		} else {
			prefix = ""
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
