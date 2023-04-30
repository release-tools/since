package changelog

import (
	"fmt"
	"github.com/outofcoffee/since/convcommits"
	"golang.org/x/exp/maps"
	"io"
	"os"
	"sort"
	"strings"
)

type ChangelogSections struct {
	Boilerplate string
	Body        string
}

// ParseChangelog loads a changelog file at the given path and returns a slice of strings containing changelog entries
// from the specified version. If no version is specified, the most recent is used.
func ParseChangelog(path string, version string, includeHeader bool) ([]string, error) {
	lines, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return readChanges(lines, version, includeHeader), nil
}

// ReadFile loads a changelog file at the given path and returns a slice of strings containing each line.
func ReadFile(path string) ([]string, error) {
	// load changelog file
	changelogfile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// convert file to string
	changelog, err := io.ReadAll(changelogfile)
	if err != nil {
		return nil, err
	}

	// split changelog into lines
	lines := strings.Split(string(changelog), "\n")
	return lines, nil
}

// readChanges parses a changelog and returns all content starting with the h2 for the specified version,
// before the next h2, or the end of the file. If no version is specified, the first h2 is used.
func readChanges(lines []string, version string, includeHeader bool) []string {
	// find the first h2
	firstH2 := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if len(version) == 0 || strings.Contains(line, "["+version+"]") {
				if includeHeader {
					firstH2 = i
				} else {
					firstH2 = i + 1
				}
				break
			}
		}
	}
	if firstH2 == 0 {
		panic(fmt.Sprintf("could not find version %s in changelog", version))
	}
	// find the next h2, or the end of the file
	nextH2 := len(lines) - firstH2 - 1
	for i, line := range lines[firstH2+1:] {
		if strings.HasPrefix(line, "## ") {
			nextH2 = i
			break
		}
	}
	return lines[firstH2 : firstH2+nextH2]
}

// RenderCommits takes a slice of commits and returns a markdown-formatted string,
// including the category header.
func RenderCommits(commits []string) string {
	categorised := convcommits.CategoriseByType(commits)

	categories := maps.Keys(categorised)
	sort.Strings(categories)

	var output string
	for _, category := range categories {
		output += "### " + category + "\n"
		items := categorised[category]
		for _, commit := range items {
			output += "- " + commit + "\n"
		}
		output += "\n"
	}
	return output
}

func SplitIntoSections(lines []string) ChangelogSections {
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
		panic("could not find h2 in changelog")
	}

	var body string
	for _, line := range lines[firstH2:] {
		body += line + "\n"
	}

	return ChangelogSections{
		Boilerplate: boilerplate,
		Body:        body,
	}
}
