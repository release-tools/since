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
	"strings"
	"testing"

	"github.com/release-tools/since/vcs"
)

func Test_listCommits(t *testing.T) {
	t.Run("lists unreleased commits since the latest tag", func(t *testing.T) {
		repoDir, _ := createChangelogTestRepo(t)
		commitCfg := vcs.CommitConfig{UniqueOnly: true}

		got, err := listCommits(commitCfg, repoDir, "", vcs.TagOrderSemver)
		if err != nil {
			t.Fatalf("listCommits() error = %v", err)
		}
		if !strings.Contains(got, "add a shiny new feature") {
			t.Errorf("listCommits() output missing the unreleased feature:\n%s", got)
		}
	})

	t.Run("lists commits since an explicit tag", func(t *testing.T) {
		repoDir, _ := createChangelogTestRepo(t)
		commitCfg := vcs.CommitConfig{UniqueOnly: true}

		got, err := listCommits(commitCfg, repoDir, "0.1.0", vcs.TagOrderSemver)
		if err != nil {
			t.Fatalf("listCommits() error = %v", err)
		}
		if !strings.Contains(got, "add a shiny new feature") {
			t.Errorf("listCommits() output missing the unreleased feature:\n%s", got)
		}
	})

	t.Run("returns error for a non-repository path", func(t *testing.T) {
		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		if _, err := listCommits(commitCfg, t.TempDir(), "", vcs.TagOrderSemver); err == nil {
			t.Error("listCommits() expected error for a non-repository path")
		}
	})
}
