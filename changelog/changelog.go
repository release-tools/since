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
	_ "embed"
	"fmt"
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/convcommits"
	"github.com/release-tools/since/semver"
	"github.com/release-tools/since/vcs"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

type Sections struct {
	Boilerplate string
	Body        string
}

//go:embed templates/changelog.md
var changelogTemplate string

var sectionMap map[string][]string

func init() {
	sectionMap = make(map[string][]string)
	sectionMap["Added"] = []string{"feat"}
	sectionMap["Fixed"] = []string{"fix"}
	sectionMap["Changed"] = []string{"build", "chore", "ci", "docs", "refactor", "style", "test"}
}

// RenderCommits takes a slice of commits and returns a markdown-formatted string,
// including the category header.
func RenderCommits(commits *[]vcs.TagCommits, groupIntoSections bool, unreleasedVersionName string) string {
	if commits == nil {
		logrus.Debug("no commits to render")
		return ""
	}
	var output string
	for _, tagCommits := range *commits {
		var versionName string
		if tagCommits.Name == vcs.UnreleasedVersionName {
			versionName = unreleasedVersionName
		} else {
			version := tagCommits.Name
			if strings.HasPrefix(version, "v") {
				version = strings.TrimPrefix(version, "v")
			}
			versionName = version
		}
		if len(output) > 0 {
			output += "\n\n"
		}
		output += "## [" + versionName + "] - " + tagCommits.Date.Format("2006-01-02") + "\n\n"
		categorised := convcommits.CategoriseByType(tagCommits.Commits)
		if groupIntoSections {
			categorised = groupBySection(categorised)
		}

		categories := maps.Keys(categorised)
		sort.Strings(categories)

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
		logrus.Debugf("grouped %d commits for version %s into %d sections\n", len(tagCommits.Commits), tagCommits.Name, len(maps.Keys(categorised)))
	}
	return output
}

// SplitIntoSections takes a slice of changelog lines and splits it into
// boilerplate and body sections.
func SplitIntoSections(lines []string) (Sections, error) {
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

	var body string
	if firstH2 > 0 {
		for _, line := range lines[firstH2:] {
			body += line + "\n"
		}
	}
	sections := Sections{
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

// GetUpdatedChangelog returns the updated changelog, grouped by version headers.
func GetUpdatedChangelog(
	config cfg.SinceConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
	beforeTag string,
	afterTag string,
	unique bool,
) (metadata vcs.ReleaseMetadata, updatedChangelog string) {
	commits, err := vcs.FetchCommitsByTag(config, repoPath, beforeTag, afterTag, unique)
	if err != nil {
		panic(fmt.Errorf("failed to fetch commit messages from repo: %s: %v", repoPath, err))
	}

	currentVersion, vPrefix := semver.GetCurrentVersion(repoPath, orderBy)

	var nextVersion string
	if beforeTag == "" {
		// determine next version only based on unreleased commits
		unreleasedCommits := (*commits)[0].Commits

		// always disable vPrefix for changelog heading
		nextVersion = semver.GetNextVersion(currentVersion, false, unreleasedCommits)
		if nextVersion == "" {
			panic("Could not determine next version")
		}
	} else {
		nextVersion = vcs.UnreleasedVersionName
	}

	rendered := RenderCommits(commits, true, nextVersion)

	lines, err := ReadFile(changelogFile)
	if err != nil {
		panic(fmt.Errorf("failed to read changelog file: %s: %v", changelogFile, err))
	}
	sections, err := SplitIntoSections(lines)
	if err != nil {
		panic(err)
	}

	output := sections.Boilerplate + rendered + "\n\n" + sections.Body

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

// InitChangelog creates a new changelog file with a placeholder entry.
func InitChangelog(
	config cfg.SinceConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
	unique bool,
) (newChangelog string, err error) {
	err = WriteChangelog(changelogFile, changelogTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to initialise changelog: %s: %v", changelogFile, err)
	}

	latestTag, err := vcs.GetLatestTag(repoPath, orderBy)
	if err != nil {
		return "", fmt.Errorf("failed to get latest tag: %v", err)
	}

	_, updated := GetUpdatedChangelog(config, changelogFile, orderBy, repoPath, latestTag, "", unique)
	err = WriteChangelog(changelogFile, updated)
	if err != nil {
		return "", fmt.Errorf("failed to write to initialised changelog: %s: %v", changelogFile, err)
	}

	return updated, nil
}
