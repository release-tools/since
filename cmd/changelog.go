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
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

var changelogArgs struct {
	changelogFile     string
	excludeTagCommits bool
	outputFile        string
}

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Commands related to changelog files",
	Long:  `Parse and update changelog files.`,
}

func init() {
	rootCmd.AddCommand(changelogCmd)

	changelogCmd.PersistentFlags().StringVarP(&changelogArgs.changelogFile, "changelog", "c", "CHANGELOG.md", "Path to changelog file")
	changelogCmd.PersistentFlags().StringVar(&changelogArgs.outputFile, "output-file", "", "Path to output file (otherwise stdout)")
	changelogCmd.PersistentFlags().BoolVar(&changelogArgs.excludeTagCommits, "exclude-tag-commits", false, "Exclude tag commits in the changelog")
}

func getWorkingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %v", err))
	}
	return workingDir
}

// writeOutput writes the output to the output file, or stdout if not set.
func writeOutput(output string) {
	outputFile := changelogArgs.outputFile
	if outputFile == "" {
		logrus.Warn("no output file specified, writing to stdout")
		fmt.Println(output)
	} else {
		logrus.Tracef("writing output to file: %s", outputFile)
		file, err := os.Create(outputFile)
		if err != nil {
			panic(fmt.Errorf("failed to create output file: %s: %v", outputFile, err))
		}
		defer file.Close()

		_, err = file.WriteString(output)
		if err != nil {
			panic(fmt.Errorf("failed to write output to file: %v", err))
		}
	}
}
