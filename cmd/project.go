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
	"github.com/outofcoffee/since/vcs"
	"github.com/spf13/cobra"
)

var projectArgs struct {
	orderBy  string
	repoPath string
	tag      string
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Commands related to the project",
	Long: `List the changes since the last release in the project
repository, or determine the next semantic version based on
those changes.`,
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.PersistentFlags().StringVarP(&projectArgs.orderBy, "order-by", "o", string(vcs.TagOrderSemver), "How to determine the latest tag (alphabetical|commit-date|semver))")
	projectCmd.PersistentFlags().StringVarP(&projectArgs.repoPath, "git-repo", "g", ".", "Path to git repository")
	projectCmd.PersistentFlags().StringVarP(&projectArgs.tag, "tag", "t", "", "Include commits after this tag")
}
