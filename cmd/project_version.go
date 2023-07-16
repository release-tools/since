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
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/semver"
	"github.com/release-tools/since/vcs"
	"github.com/spf13/cobra"
	"os"
)

var versionArgs struct {
	current bool
	unique  bool
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the next semantic version based on changes since last tag",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Returns the next semantic version
based on the changes.

Changes influence the version according to
conventional commits: https://www.conventionalcommits.org/en/v1.0.0/`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		version := printVersion(
			projectArgs.repoPath,
			projectArgs.tag,
			vcs.TagOrderBy(projectArgs.orderBy),
			versionArgs.current,
			versionArgs.unique,
		)
		if version == "" {
			os.Exit(1)
		}
		fmt.Println(version)
	},
}

func init() {
	projectCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&versionArgs.current, "current", "c", false, "Just print the current version")
	versionCmd.Flags().BoolVar(&versionArgs.unique, "unique", true, "De-duplicate commit messages")
}

func printVersion(
	repoPath string,
	tag string,
	orderBy vcs.TagOrderBy,
	current bool,
	unique bool,
) string {
	currentVersion, vPrefix := semver.GetCurrentVersion(repoPath, orderBy)
	if current {
		if vPrefix {
			currentVersion = "v" + currentVersion
		}
		return currentVersion
	}

	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		panic(err)
	}

	commits, err := vcs.FetchCommitMessages(config, repoPath, "", tag, orderBy, unique)
	if err != nil {
		panic(err)
	}
	return semver.GetNextVersion(currentVersion, vPrefix, commits)
}
