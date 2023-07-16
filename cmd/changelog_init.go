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
	"github.com/release-tools/since/changelog"
	"github.com/release-tools/since/vcs"
	"github.com/spf13/cobra"
)

var initArgs struct {
	orderBy  string
	repoPath string
	unique   bool
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new changelog file",
	Long:  `Initialises a new changelog file with a placeholder entry.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		workingDir := getWorkingDir()
		changelogFile := changelog.ResolveChangelogFile(
			workingDir,
			changelogArgs.changelogFile,
		)
		initChangelog(
			changelogFile,
			vcs.TagOrderBy(generateArgs.orderBy),
			generateArgs.repoPath,
			generateArgs.unique,
		)
	},
}

func init() {
	changelogCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	initCmd.Flags().StringVarP(&initArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	initCmd.Flags().BoolVar(&initArgs.unique, "unique", true, "De-duplicate commit messages")
}

func initChangelog(
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
	unique bool,
) {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		panic(err)
	}
	newChangelog, err := changelog.InitChangelog(config, changelogFile, orderBy, repoPath, unique)
	if err != nil {
		panic(err)
	}
	fmt.Println(newChangelog)
}
