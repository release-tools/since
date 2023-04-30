/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/outofcoffee/since/changelog"
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
	"os"
)

var releaseArgs struct {
	changelogFile string
}

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Update the changelog, commit the changes and tag the release",
	Long: `Generates a new changelog based on an existing changelog file,
using the commits since the last release.

The changelog is then committed and tagged with the new version.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		release(releaseArgs.changelogFile, vcs.TagOrderBy(projectArgs.orderBy), projectArgs.repoPath)
	},
}

func init() {
	projectCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringVarP(&releaseArgs.changelogFile, "changelog", "c", "CHANGELOG.md", "Path to changelog file")
}

func release(changelogFile string, orderBy vcs.TagOrderBy, repoPath string) {
	version, vPrefix, updatedChangelog := changelog.GetUpdatedChangelog(changelogFile, orderBy, repoPath)
	if vPrefix {
		version = "v" + version
	}

	tempChangelog := writeTempChangelog(updatedChangelog)
	err := os.Rename(tempChangelog, changelogFile)
	if err != nil {
		panic(fmt.Errorf("failed to rename temp file: %s: %w", tempChangelog, err))
	}

	hash, err := vcs.CommitChangelog(repoPath, changelogFile, version)
	if err != nil {
		panic(fmt.Errorf("failed to commit changelog: %w", err))
	}

	err = vcs.TagRelease(repoPath, hash, version)
	if err != nil {
		panic(fmt.Errorf("failed to tag release commit: %s: %w", hash, err))
	}

	fmt.Printf("released version %s\n", version)
}

func writeTempChangelog(updatedChangelog string) string {
	temp, err := os.CreateTemp(os.TempDir(), "changelog*.md")
	if err != nil {
		panic(fmt.Errorf("failed to create temp file: %w", err))
	}
	_, err = temp.WriteString(updatedChangelog + "\n")
	if err != nil {
		panic(fmt.Errorf("failed to write to temp file: %w", err))
	}
	_ = temp.Close()
	return temp.Name()
}
