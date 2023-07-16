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

package cmd

import (
	"fmt"
	"github.com/release-tools/since/changelog"
	"github.com/spf13/cobra"
)

var extractArgs struct {
	version       string
	includeHeader bool
}

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract changes for a given version",
	Long: `Extracts changes for a given version in a changelog file.
If no version is specified, the most recent version is used.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		changelogFile := changelog.ResolveChangelogFile(getWorkingDir(), changelogArgs.changelogFile)
		changes, err := printChanges(changelogFile, extractArgs.version, extractArgs.includeHeader)
		if err != nil {
			panic(err)
		}
		fmt.Println(changes)
	},
}

func init() {
	changelogCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&extractArgs.version, "version", "v", "", "Version to parse changelog for")
	extractCmd.Flags().BoolVar(&extractArgs.includeHeader, "header", false, "whether to include the version header in the output")
}

func printChanges(changelogFile string, version string, includeHeader bool) (string, error) {
	changes, err := changelog.ParseChangelog(changelogFile, version, includeHeader)
	if err != nil {
		return "", err
	}
	output := ""
	for _, entry := range changes {
		output += entry + "\n"
	}
	return output, nil
}
