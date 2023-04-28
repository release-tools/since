/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var listArgs struct {
	version       string
	changelogFile string
	includeHeader bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List changes for a given version",
	Long: `Lists changes for a given version in a changelog file.
If no version is specified, the most recent version is used.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		printChanges(listArgs.changelogFile, listArgs.version, listArgs.includeHeader, func(s string) { fmt.Println(s) })
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listArgs.version, "version", "v", "", "Version to parse changelog for")
	listCmd.Flags().StringVarP(&listArgs.changelogFile, "changelog", "c", "CHANGELOG.md", "Path to changelog file")
	listCmd.Flags().BoolVar(&listArgs.includeHeader, "header", false, "whether to include the version header in the output")
}

func printChanges(changelogFile string, version string, includeHeader bool, printer func(s string)) {
	changes, err := parseChangelog(changelogFile, version, includeHeader)
	if err != nil {
		panic(err)
	}
	for _, entry := range changes {
		printer(entry)
	}
}

// parseChangelog loads a changelog file at the given path and returns a slice of strings containing changelog entries
// from the specified version. If no version is specified, the most recent is used.
func parseChangelog(path string, version string, includeHeader bool) ([]string, error) {
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
	return readChanges(string(changelog), version, includeHeader), nil
}

// readChanges parses a changelog and returns all content starting with the h2 for the specified version,
// before the next h2, or the end of the file. If no version is specified, the first h2 is used.
func readChanges(changelog string, version string, includeHeader bool) []string {
	// split changelog into lines
	lines := strings.Split(changelog, "\n")
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
