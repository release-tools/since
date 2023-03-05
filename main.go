package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var printer = fmt.Println

func main() {
	version := flag.String("version", "", "version to parse changelog for")
	changelogFile := flag.String("changelog", "CHANGELOG.md", "path to changelog file")
	flag.Parse()
	changes, err := parseChangelog(*changelogFile, *version)
	if err != nil {
		panic(err)
	}
	for _, entry := range changes {
		printer(entry)
	}
}

// parseChangelog loads a changelog file at the given path and returns a slice of strings containing changelog entries
// from the specified version. If no version is specified, the most recent is used.
func parseChangelog(path string, version string) ([]string, error) {
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
	// parse changelog file
	return readChanges(string(changelog), version), nil
}

// processChangLog parses a changelog and returns all content starting with the h2 for the specified version,
// before the next h2, or the end of the file. If no version is specified, the first h2 is used.
func readChanges(changelog string, version string) []string {
	// split changelog into lines
	lines := strings.Split(changelog, "\n")
	// find the first h2
	firstH2 := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if len(version) == 0 || strings.Contains(line, "["+version+"]") {
				firstH2 = i
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
