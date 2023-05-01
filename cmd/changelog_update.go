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
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
)

var updateArgs struct {
	orderBy  string
	repoPath string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Write updated changelog based on changes since last release",
	Long: `Updates the existing changelog file with a new release section,
using the commits since the last release.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		updateChangelog(
			changelogArgs.changelogFile,
			vcs.TagOrderBy(updateArgs.orderBy),
			updateArgs.repoPath,
		)
	},
}

func init() {
	changelogCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	updateCmd.Flags().StringVarP(&updateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
}

func updateChangelog(changelogFile string, orderBy vcs.TagOrderBy, repoPath string) {
	_, _, updated := changelog.GetUpdatedChangelog(changelogFile, orderBy, repoPath)

	err := changelog.UpdateChangelog(changelogFile, updated)
	if err != nil {
		panic(fmt.Errorf("failed to update changelog: %w", err))
	}
}
