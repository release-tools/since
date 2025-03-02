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
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/changelog"
	"github.com/release-tools/since/vcs"
	"github.com/spf13/cobra"
)

var generateArgs struct {
	orderBy  string
	repoPath string
	unique   bool
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Print generated changelog based on changes since last release",
	Long: `Generates a new changelog based on an existing changelog file,
adding a new release section using the commits since the last release,
then prints it to stdout, or output-file, if specified.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		changelogFile := changelog.ResolveChangelogFile(
			generateArgs.repoPath,
			changelogArgs.changelogFile,
		)
		commitCfg := vcs.CommitConfig{
			ExcludeTagCommits: changelogArgs.excludeTagCommits,
			UniqueOnly:        generateArgs.unique,
		}
		generateChangelog(
			commitCfg,
			changelogFile,
			vcs.TagOrderBy(generateArgs.orderBy),
			generateArgs.repoPath,
		)
	},
}

func init() {
	changelogCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	generateCmd.Flags().StringVarP(&generateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	generateCmd.Flags().BoolVar(&generateArgs.unique, "unique", true, "De-duplicate commit messages")
}

func generateChangelog(
	commitCfg vcs.CommitConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
) {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		panic(err)
	}

	latestTag, err := vcs.GetLatestTag(repoPath, orderBy)
	if err != nil {
		panic(err)
	}

	_, updated, err := changelog.GetUpdatedChangelog(config, commitCfg, changelogFile, orderBy, repoPath, "", latestTag)
	if err != nil {
		panic(err)
	}
	writeOutput(updated)
}
