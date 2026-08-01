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
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// baseTimeMillis is a fixed epoch used to make test commit timestamps
// deterministic, avoiding any dependency on the wall clock.
const baseTimeMillis int64 = 1_700_000_000_000

func Test_writeOutput(t *testing.T) {
	t.Run("writes to stdout when no output file is set", func(t *testing.T) {
		changelogArgs.outputFile = ""
		if err := writeOutput("some content"); err != nil {
			t.Errorf("writeOutput() to stdout error = %v", err)
		}
	})

	t.Run("writes to the configured output file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outPath := filepath.Join(tmpDir, "out.md")

		changelogArgs.outputFile = outPath
		defer func() { changelogArgs.outputFile = "" }()

		if err := writeOutput("some content"); err != nil {
			t.Fatalf("writeOutput() error = %v", err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}
		if string(content) != "some content" {
			t.Errorf("output file content = %q, want %q", string(content), "some content")
		}
	})

	t.Run("returns error when output file cannot be created", func(t *testing.T) {
		tmpDir := t.TempDir()
		// a path under a missing directory cannot be created
		changelogArgs.outputFile = filepath.Join(tmpDir, "does-not-exist", "out.md")
		defer func() { changelogArgs.outputFile = "" }()

		if err := writeOutput("some content"); err == nil {
			t.Error("writeOutput() expected error for uncreatable output file, got nil")
		}
	})
}

func Test_getWorkingDir(t *testing.T) {
	dir, err := getWorkingDir()
	if err != nil {
		t.Fatalf("getWorkingDir() error = %v", err)
	}
	if dir == "" {
		t.Error("getWorkingDir() returned an empty path")
	}
}

// createChangelogTestRepo builds a repository tagged 0.1.0 with a single
// unreleased conventional commit after the tag, and writes a changelog file
// containing the standard boilerplate plus the 0.1.0 section. It returns the
// repo directory and the path to the changelog file.
//
// The latest tag is always 0.1.0 so that vcs's package-level tag cache stays
// valid across the changelog handler tests regardless of execution order.
func createChangelogTestRepo(t *testing.T) (repoDir, changelogFile string) {
	t.Helper()
	repoDir = t.TempDir()

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	commit := func(message string, offsetMillis int64) plumbing.Hash {
		path := filepath.Join(repoDir, "README.md")
		if err := os.WriteFile(path, []byte(message), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := w.Add("README.md"); err != nil {
			t.Fatal(err)
		}
		sig := &object.Signature{
			Name:  "user",
			Email: "user@example.com",
			When:  time.UnixMilli(baseTimeMillis + offsetMillis),
		}
		h, err := w.Commit(message, &git.CommitOptions{Author: sig, Committer: sig})
		if err != nil {
			t.Fatal(err)
		}
		return h
	}

	first := commit("chore: initial commit", 0)
	if _, err := repo.CreateTag("0.1.0", first, nil); err != nil {
		t.Fatal(err)
	}
	// an unreleased feature commit so the next version resolves to a minor bump
	commit("feat: add a shiny new feature", 10000)

	changelogFile = filepath.Join(repoDir, "CHANGELOG.md")
	body := `# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.1.0] - 2023-03-04
### Added
- chore: initial commit.
`
	if err := os.WriteFile(changelogFile, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return repoDir, changelogFile
}
