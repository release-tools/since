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
	"os"
	"strings"
	"testing"

	"github.com/release-tools/since/vcs"
)

func Test_updateChangelog(t *testing.T) {
	t.Run("writes a new release section into the changelog file", func(t *testing.T) {
		repoDir, changelogFile := createChangelogTestRepo(t)
		commitCfg := vcs.CommitConfig{UniqueOnly: true}

		err := updateChangelog(commitCfg, changelogFile, vcs.TagOrderSemver, repoDir)
		if err != nil {
			t.Fatalf("updateChangelog() error = %v", err)
		}

		content, err := os.ReadFile(changelogFile)
		if err != nil {
			t.Fatalf("failed to read changelog file: %v", err)
		}
		got := string(content)

		// the unreleased feat commit should produce a new 0.2.0 minor section
		if !strings.Contains(got, "## [0.2.0]") {
			t.Errorf("changelog does not contain new 0.2.0 section:\n%s", got)
		}
		if !strings.Contains(got, "add a shiny new feature") {
			t.Errorf("changelog does not contain the unreleased feature:\n%s", got)
		}
		// the pre-existing 0.1.0 section should be preserved
		if !strings.Contains(got, "## [0.1.0]") {
			t.Errorf("changelog lost the existing 0.1.0 section:\n%s", got)
		}
	})

	t.Run("returns error for a non-repository path", func(t *testing.T) {
		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		err := updateChangelog(commitCfg, "CHANGELOG.md", vcs.TagOrderSemver, t.TempDir())
		if err == nil {
			t.Error("updateChangelog() expected error for a non-repository path")
		}
	})
}
