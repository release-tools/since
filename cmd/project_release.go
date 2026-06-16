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
	"github.com/release-tools/since/hooks"
	"github.com/release-tools/since/vcs"
	"github.com/spf13/cobra"
)

var releaseArgs struct {
	changelogFile string
	unique        bool
}

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Update the changelog, commit the changes and tag the release",
	Long: `Generates a new changelog based on an existing changelog file,
using the commits since the last release.

The changelog is then committed and a new tag is created
with the new version.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		changelogFile := changelog.ResolveChangelogFile(projectArgs.repoPath, changelogArgs.changelogFile)
		commitCfg := vcs.CommitConfig{
			ExcludeTagCommits: projectArgs.excludeTagCommits,
			UniqueOnly:        releaseArgs.unique,
		}
		return release(
			commitCfg,
			changelogFile,
			vcs.TagOrderBy(projectArgs.orderBy),
			projectArgs.repoPath,
		)
	},
}

func init() {
	projectCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringVarP(&releaseArgs.changelogFile, "changelog", "c", "CHANGELOG.md", "Path to changelog file")
	releaseCmd.Flags().BoolVar(&releaseArgs.unique, "unique", true, "De-duplicate commit messages")
}

func release(
	commitCfg vcs.CommitConfig,
	changelogFile string,
	orderBy vcs.TagOrderBy,
	repoPath string,
) error {
	config, err := cfg.LoadConfig(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := vcs.CheckBranch(repoPath, config); err != nil {
		return err
	}

	latestTag, err := vcs.GetLatestTag(repoPath, orderBy)
	if err != nil {
		return err
	}

	metadata, updatedChangelog, err := changelog.GetUpdatedChangelog(config, commitCfg, changelogFile, orderBy, repoPath, "", latestTag)
	if err != nil {
		return err
	}

	version := metadata.NewVersion
	if metadata.VPrefix {
		version = "v" + version
	}

	if err := hooks.ExecuteHooks(config, hooks.Before, metadata); err != nil {
		return fmt.Errorf("failed to execute hooks before release: %w", err)
	}

	if err := changelog.WriteChangelog(changelogFile, updatedChangelog); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	hash, err := vcs.CommitChangelog(repoPath, changelogFile, version)
	if err != nil {
		return fmt.Errorf("failed to commit changelog: %w", err)
	}

	if err := vcs.TagRelease(repoPath, hash, version); err != nil {
		return fmt.Errorf("failed to tag release commit: %s: %w", hash, err)
	}

	if err := hooks.ExecuteHooks(config, hooks.After, metadata); err != nil {
		return fmt.Errorf("failed to execute hooks after release: %w", err)
	}

	fmt.Printf("released version %s\n", version)
	return nil
}
