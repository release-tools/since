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

var changesArgs struct {
	unique bool
}

// changesCmd represents the changes command
var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "List the changes since the last release",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Lists the commits categorised by their type.`,
	Run: func(cmd *cobra.Command, args []string) {
		commitCfg := vcs.CommitConfig{
			ExcludeTagCommits: projectArgs.excludeTagCommits,
			UniqueOnly:        changesArgs.unique,
		}
		changes, err := listCommits(
			commitCfg,
			projectArgs.repoPath,
			projectArgs.tag,
			vcs.TagOrderBy(projectArgs.orderBy),
		)
		if err != nil {
			panic(err)
		}
		fmt.Println(changes)
	},
}

func init() {
	projectCmd.AddCommand(changesCmd)

	changesCmd.Flags().BoolVar(&changesArgs.unique, "unique", true, "De-duplicate commit messages")
}

func listCommits(
	commitCfg vcs.CommitConfig,
	repoPath string,
	tag string,
	orderBy vcs.TagOrderBy,
) (string, error) {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		panic(err)
	}

	var afterTag string
	if tag == "" {
		latestTag, err := vcs.GetLatestTag(repoPath, orderBy)
		if err != nil {
			panic(err)
		}
		afterTag = latestTag
	} else {
		afterTag = tag
	}

	commits, err := vcs.FetchCommitsByTag(config, commitCfg, repoPath, "", afterTag)
	if err != nil {
		return "", err
	}
	return changelog.RenderCommits(commits, true, false, vcs.UnreleasedVersionName), nil
}
