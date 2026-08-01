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

func Test_initChangelog(t *testing.T) {
	t.Run("writes an initialised changelog to the output file", func(t *testing.T) {
		repoDir, _ := createChangelogTestRepo(t)

		outPath := filepath.Join(t.TempDir(), "CHANGELOG.md")
		changelogArgs.outputFile = outPath
		defer func() { changelogArgs.outputFile = "" }()

		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		// pass a changelog file that does not yet exist, as init would in practice
		err := initChangelog(commitCfg, filepath.Join(repoDir, "CHANGELOG.md"), vcs.TagOrderSemver, repoDir)
		if err != nil {
			t.Fatalf("initChangelog() error = %v", err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read initialised changelog: %v", err)
		}
		got := string(content)
		if !strings.Contains(got, "# Changelog") {
			t.Errorf("initialised changelog missing boilerplate heading:\n%s", got)
		}
		if !strings.Contains(got, "## [0.1.0]") {
			t.Errorf("initialised changelog missing the tagged release section:\n%s", got)
		}
		if !strings.Contains(got, "chore: initial commit") {
			t.Errorf("initialised changelog missing the release commit:\n%s", got)
		}
	})

	t.Run("returns error for a non-repository path", func(t *testing.T) {
		tmpDir := t.TempDir()
		commitCfg := vcs.CommitConfig{UniqueOnly: true}
		// use an absolute path within a temp dir: InitChangelog writes the
		// template before hitting the repo error, so a relative path would
		// otherwise pollute the working directory.
		err := initChangelog(commitCfg, filepath.Join(tmpDir, "CHANGELOG.md"), vcs.TagOrderSemver, tmpDir)
		if err == nil {
			t.Error("initChangelog() expected error for a non-repository path")
		}
	})
}
