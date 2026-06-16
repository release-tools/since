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

var updateArgs struct {
	orderBy  string
	repoPath string
	unique   bool
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Write updated changelog based on changes since last release",
	Long: `Updates the existing changelog file with a new release section,
using the commits since the last release.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		changelogFile := changelog.ResolveChangelogFile(updateArgs.repoPath, changelogArgs.changelogFile)
		commitCfg := vcs.CommitConfig{
			ExcludeTagCommits: changelogArgs.excludeTagCommits,
			UniqueOnly:        updateArgs.unique,
		}
		return updateChangelog(
			commitCfg,
			changelogFile,
			vcs.TagOrderBy(updateArgs.orderBy),
			updateArgs.repoPath,
		)
	},
}

func init() {
	changelogCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	updateCmd.Flags().StringVarP(&updateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	updateCmd.Flags().BoolVar(&updateArgs.unique, "unique", true, "De-duplicate commit messages")
}

func updateChangelog(
	commitCfg vcs.CommitConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
) error {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		return err
	}

	latestTag, err := vcs.GetLatestTag(repoPath, orderBy)
	if err != nil {
		return err
	}

	_, updated, err := changelog.GetUpdatedChangelog(config, commitCfg, changelogFile, orderBy, repoPath, "", latestTag)
	if err != nil {
		return err
	}

	if err := changelog.WriteChangelog(changelogFile, updated); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}
	return nil
}
