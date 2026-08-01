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
	"path/filepath"
	"strings"
	"testing"

	"github.com/release-tools/since/vcs"
)

func Test_generateChangelog(t *testing.T) {
	t.Run("writes the generated changelog to the output file", func(t *testing.T) {
		repoDir, changelogFile := createChangelogTestRepo(t)

		outPath := filepath.Join(t.TempDir(), "generated.md")
		changelogArgs.outputFile = outPath
		defer func() { changelogArgs.outputFile = "" }()

		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		err := generateChangelog(commitCfg, changelogFile, vcs.TagOrderSemver, repoDir)
		if err != nil {
			t.Fatalf("generateChangelog() error = %v", err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read generated output: %v", err)
		}
		got := string(content)
		if !strings.Contains(got, "## [0.2.0]") {
			t.Errorf("generated changelog does not contain new 0.2.0 section:\n%s", got)
		}
		if !strings.Contains(got, "## [0.1.0]") {
			t.Errorf("generated changelog does not preserve existing 0.1.0 section:\n%s", got)
		}
	})

	t.Run("leaves the source changelog file untouched", func(t *testing.T) {
		repoDir, changelogFile := createChangelogTestRepo(t)

		before, err := os.ReadFile(changelogFile)
		if err != nil {
			t.Fatal(err)
		}

		changelogArgs.outputFile = filepath.Join(t.TempDir(), "generated.md")
		defer func() { changelogArgs.outputFile = "" }()

		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		if err := generateChangelog(commitCfg, changelogFile, vcs.TagOrderSemver, repoDir); err != nil {
			t.Fatalf("generateChangelog() error = %v", err)
		}

		after, err := os.ReadFile(changelogFile)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Error("generateChangelog() modified the source changelog file, expected it untouched")
		}
	})

	t.Run("returns error for a non-repository path", func(t *testing.T) {
		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		err := generateChangelog(commitCfg, "CHANGELOG.md", vcs.TagOrderSemver, t.TempDir())
		if err == nil {
			t.Error("generateChangelog() expected error for a non-repository path")
		}
	})
}
