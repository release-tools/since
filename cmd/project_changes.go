/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
)

// changesCmd represents the changes command
var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "List the changes since the last release",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Lists the commits categorised by their type.`,
	Run: func(cmd *cobra.Command, args []string) {
		changes, err := listCommits(projectArgs.repoPath, projectArgs.tag, vcs.TagOrderBy(projectArgs.orderBy))
		if err != nil {
			panic(err)
		}
		fmt.Println(changes)
	},
}

func init() {
	projectCmd.AddCommand(changesCmd)
}

func listCommits(repoPath string, tag string, orderBy vcs.TagOrderBy) (string, error) {
	commits, err := vcs.FetchCommitMessages(repoPath, tag, orderBy)
	if err != nil {
		return "", err
	}
	return changelog.RenderCommits(commits, true), nil
}
