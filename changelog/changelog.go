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

package changelog

import (
	"fmt"
	"github.com/outofcoffee/since/cfg"
	"github.com/outofcoffee/since/convcommits"
	"github.com/outofcoffee/since/semver"
	"github.com/outofcoffee/since/vcs"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
	"time"
)

type ChangelogSections struct {
	Boilerplate string
	Body        string
}

var sectionMap map[string][]string

func init() {
	sectionMap = make(map[string][]string)
	sectionMap["Added"] = []string{"feat"}
	sectionMap["Fixed"] = []string{"fix"}
	sectionMap["Changed"] = []string{"build", "chore", "ci", "docs", "refactor", "style", "test"}
}

// RenderCommits takes a slice of commits and returns a markdown-formatted string,
// including the category header.
func RenderCommits(commits []string, groupIntoSections bool) string {
	categorised := convcommits.CategoriseByType(commits)
	if groupIntoSections {
		categorised = groupBySection(categorised)
	}

	categories := maps.Keys(categorised)
	sort.Strings(categories)

	var output string
	for _, category := range categories {
		output += "### " + category + "\n"
		items := categorised[category]
		sort.Strings(items)

		for _, commit := range items {
			output += "- " + commit + "\n"
		}
		output += "\n"
	}
	output = strings.TrimSpace(output)
	logrus.Debugf("grouped commits into %d sections\n", len(maps.Keys(categorised)))
	return output
}

// SplitIntoSections takes a slice of changelog lines and splits it into
// boilerplate and body sections.
func SplitIntoSections(lines []string) (ChangelogSections, error) {
	var boilerplate string

	// find the first h2
	firstH2 := 0
	skipping := false

	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if strings.Contains(line, "[Unreleased]") {
				skipping = true
			} else {
				firstH2 = i
				break
			}
		} else {
			if !skipping {
				boilerplate += line + "\n"
			}
		}
	}
	if firstH2 == 0 {
		return ChangelogSections{}, fmt.Errorf("could not find h2 in changelog")
	}

	var body string
	for _, line := range lines[firstH2:] {
		body += line + "\n"
	}
	sections := ChangelogSections{
		Boilerplate: boilerplate,
		Body:        strings.TrimSpace(body),
	}
	return sections, nil
}

// groupBySection maps the commit prefixes to sections.
func groupBySection(input map[string][]string) map[string][]string {
	output := make(map[string][]string)
	for prefix, commits := range input {
		prefix = mapTypeToSection(prefix)

		existing := output[prefix]
		commits = append(existing, commits...)
		output[prefix] = commits
	}
	return output
}

// mapTypeToSection maps a commit prefix to a section.
func mapTypeToSection(prefix string) string {
	mapped := false
	for section, types := range sectionMap {
		for _, t := range types {
			if prefix == t {
				prefix = section
				mapped = true
				break
			}
		}
		if mapped {
			break
		}
	}
	if !mapped {
		prefix = "Other"
	}
	return prefix
}

// GetUpdatedChangelog returns the updated changelog, including the new version header.
func GetUpdatedChangelog(
	config cfg.SinceConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
) (metadata vcs.ReleaseMetadata, updatedChangelog string) {
	commits, err := vcs.FetchCommitMessages(config, repoPath, "", orderBy)
	if err != nil {
		panic(fmt.Errorf("failed to fetch commit messages from repo: %s: %v", repoPath, err))
	}
	rendered := RenderCommits(commits, true)

	currentVersion, vPrefix := semver.GetCurrentVersion(repoPath, orderBy)

	// always disable vPrefix for changelog heading
	nextVersion := semver.GetNextVersion(currentVersion, false, commits)
	if nextVersion == "" {
		panic("Could not determine next version")
	}
	versionHeader := "## [" + nextVersion + "] - " + time.Now().UTC().Format("2006-01-02") + "\n"

	lines, err := ReadFile(changelogFile)
	if err != nil {
		panic(fmt.Errorf("failed to read changelog file: %s: %v", changelogFile, err))
	}
	sections, err := SplitIntoSections(lines)
	if err != nil {
		panic(err)
	}

	output := sections.Boilerplate + versionHeader + rendered + "\n\n" + sections.Body

	sha, err := vcs.GetHeadSha(repoPath)
	if err != nil {
		panic(fmt.Errorf("failed to get head sha: %v", err))
	}
	metadata = vcs.ReleaseMetadata{
		OldVersion: currentVersion,
		NewVersion: nextVersion,
		RepoPath:   repoPath,
		Sha:        sha,
		VPrefix:    vPrefix,
	}
	return metadata, output
}
