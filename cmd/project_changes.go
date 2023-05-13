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
	"github.com/outofcoffee/since/cfg"
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/vcs"
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
		changes, err := listCommits(
			projectArgs.repoPath,
			projectArgs.tag,
			vcs.TagOrderBy(projectArgs.orderBy),
			changesArgs.unique,
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
	repoPath string,
	tag string,
	orderBy vcs.TagOrderBy,
	unique bool,
) (string, error) {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		panic(err)
	}

	commits, err := vcs.FetchCommitMessages(config, repoPath, tag, orderBy, unique)
	if err != nil {
		return "", err
	}
	return changelog.RenderCommits(commits, true), nil
}
