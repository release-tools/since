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

var generateArgs struct {
	orderBy  string
	repoPath string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Print generated changelog based on changes since last release",
	Long: `Generates a new changelog based on an existing changelog file,
adding a new release section using the commits since the last release,
then prints it to stdout.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generateChangelog(
			changelogArgs.changelogFile,
			vcs.TagOrderBy(generateArgs.orderBy),
			generateArgs.repoPath,
		)
	},
}

func init() {
	changelogCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generateArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	generateCmd.Flags().StringVarP(&generateArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
}

func generateChangelog(changelogFile string, orderBy vcs.TagOrderBy, repoPath string) {
	_, _, updated := changelog.GetUpdatedChangelog(changelogFile, orderBy, repoPath)
	fmt.Println(updated)
}
