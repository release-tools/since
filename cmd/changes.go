/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/outofcoffee/since/convcommits"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
)

var changesArgs struct {
	orderBy  string
	repoPath string
	tag      string
}

// changesCmd represents the changes command
var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "List the changes since the last release",
	Long: `Reads the commit history for the current git repository, starting
from the most recent tag. Lists the commits categorised by their type.`,
	Run: func(cmd *cobra.Command, args []string) {
		listCommits(changesArgs.repoPath, changesArgs.tag, vcs.TagOrderBy(semverArgs.orderBy), func(s string) { fmt.Println(s) })
	},
}

func init() {
	rootCmd.AddCommand(changesCmd)

	changesCmd.Flags().StringVarP(&changesArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	changesCmd.Flags().StringVarP(&changesArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	changesCmd.Flags().StringVarP(&changesArgs.tag, "tag", "t", "", "Include commits after this tag")
}

func listCommits(repoPath string, tag string, orderBy vcs.TagOrderBy, printer func(s string)) {
	commits, err := vcs.FetchCommitMessages(repoPath, tag, orderBy)
	if err != nil {
		panic(err)
	}

	printCommits(commits, printer)
}

func printCommits(commits []string, printer func(s string)) {
	categorised := convcommits.CategoriseByType(commits)

	for prefix, category := range categorised {
		printer("\n### " + prefix + "\n")
		for _, commit := range category {
			printer("- " + commit)
		}
	}
}
