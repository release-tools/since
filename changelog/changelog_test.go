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

package changelog

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/vcs"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestRenderCommits(t *testing.T) {
	type args struct {
		groupIntoSections bool
		releaseUnreleased bool
		commits           *[]vcs.TagCommits
	}

	var fewCommits []vcs.TagCommits
	fewCommits = append(fewCommits, vcs.TagCommits{
		TagMeta: vcs.TagMeta{
			Name: vcs.UnreleasedVersionName,
			Date: time.Date(2023, 8, 28, 0, 0, 0, 0, time.UTC),
		},
		Commits: []string{"feat: foo", "fix: bar"},
	})

	var manyCommits []vcs.TagCommits
	manyCommits = append(manyCommits, vcs.TagCommits{
		TagMeta: vcs.TagMeta{
			Name: vcs.UnreleasedVersionName,
			Date: time.Date(2023, 8, 28, 0, 0, 0, 0, time.UTC),
		},
		Commits: []string{"feat: foo", "fix: bar", "chore: qux"},
	})
	manyCommits = append(manyCommits, vcs.TagCommits{
		TagMeta: vcs.TagMeta{
			Name: "0.1.0",
			Date: time.Date(2023, 8, 27, 0, 0, 0, 0, time.UTC),
		},
		Commits: []string{"ci: baz", "build: quux", "feat: corge"},
	})

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no commits",
			args: args{
				commits:           nil,
				groupIntoSections: false,
			},
			want: "",
		},
		{
			name: "print commits",
			args: args{
				commits:           &fewCommits,
				groupIntoSections: false,
			},
			want: `## [Unreleased]
### feat
- feat: foo

### fix
- fix: bar`,
		},
		{
			name: "group commits",
			args: args{
				commits:           &manyCommits,
				groupIntoSections: true,
			},
			want: `## [Unreleased]
### Added
- feat: foo

### Changed
- chore: qux

### Fixed
- fix: bar

## [0.1.0] - 2023-08-27
### Added
- feat: corge

### Changed
- build: quux
- ci: baz`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderCommits(tt.args.commits, tt.args.groupIntoSections, tt.args.releaseUnreleased, vcs.UnreleasedVersionName); got != tt.want {
				t.Errorf("RenderCommits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitIntoSections(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want Sections
	}{
		{
			name: "no sections",
			args: args{
				lines: []string{},
			},
			want: Sections{},
		},
		{
			name: "split into sections",
			args: args{
				lines: []string{"# Change Log", "", "## [0.1.0]", "### feat", "- feat: foo", "", "### fix", "- fix: bar"},
			},
			want: Sections{
				Boilerplate: "# Change Log\n\n",
				Body:        "## [0.1.0]\n### feat\n- feat: foo\n\n### fix\n- fix: bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitIntoSections(tt.args.lines)
			if !reflect.DeepEqual(got.Boilerplate, tt.want.Boilerplate) {
				t.Errorf("SplitIntoSections() Boilerplate got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Body, tt.want.Body) {
				t.Errorf("SplitIntoSections() Body got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUpdatedChangelog(t *testing.T) {
	repoWithTagsAndNoUnreleasedChanges := createTestRepo(t)

	repoWithTagsAndUnreleasedChanges := createTestRepo(t)
	commitChange(
		t,
		repoWithTagsAndUnreleasedChanges,
		"CHANGELOG.md",
		changelogTemplate,
		"docs: adds changelog",
		time.Now(),
	)
	unreleasedCommitSha := commitChange(
		t,
		repoWithTagsAndUnreleasedChanges,
		"README.md",
		"unreleased change\r\n",
		"feat: unreleased change",
		time.Now(),
	)

	today := time.Now().Format("2006-01-02")

	type args struct {
		config        cfg.SinceConfig
		changelogFile string
		orderBy       vcs.TagOrderBy
		repoPath      string
		beforeTag     string
		afterTag      string
		unique        bool
	}
	tests := []struct {
		name                 string
		args                 args
		wantMetadata         vcs.ReleaseMetadata
		wantUpdatedChangelog string
		wantErr              bool
		errMessage           string
	}{
		{
			name: "no changes",
			args: args{
				orderBy:  vcs.TagOrderSemver,
				repoPath: repoWithTagsAndNoUnreleasedChanges,
				afterTag: "0.1.0",
			},
			wantMetadata:         vcs.ReleaseMetadata{},
			wantUpdatedChangelog: "",
			wantErr:              true,
			errMessage:           "no changes since start tag",
		},
		{
			name: "unreleased changes",
			args: args{
				changelogFile: path.Join(repoWithTagsAndUnreleasedChanges, "CHANGELOG.md"),
				orderBy:       vcs.TagOrderSemver,
				repoPath:      repoWithTagsAndUnreleasedChanges,
				afterTag:      "0.1.0",
			},
			wantMetadata: vcs.ReleaseMetadata{
				NewVersion: "0.2.0",
				OldVersion: "0.1.0",
				RepoPath:   repoWithTagsAndUnreleasedChanges,
				Sha:        unreleasedCommitSha.String(),
				VPrefix:    false,
			},
			wantUpdatedChangelog: fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - %v
### Added
- feat: unreleased change

### Changed
- docs: adds changelog

`, today),
			wantErr: false,
		},
		{
			name: "multiple versions",
			args: args{
				changelogFile: path.Join(repoWithTagsAndUnreleasedChanges, "CHANGELOG.md"),
				orderBy:       vcs.TagOrderSemver,
				repoPath:      repoWithTagsAndUnreleasedChanges,
				afterTag:      "0.0.1",
			},
			wantMetadata: vcs.ReleaseMetadata{
				NewVersion: "0.2.0",
				OldVersion: "0.1.0",
				RepoPath:   repoWithTagsAndUnreleasedChanges,
				Sha:        unreleasedCommitSha.String(),
				VPrefix:    false,
			},
			wantUpdatedChangelog: fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - %[1]v
### Added
- feat: unreleased change

### Changed
- docs: adds changelog

## [0.1.0] - %[1]v
### Added
- feat: second update

`, today),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetadata, gotUpdatedChangelog, err := GetUpdatedChangelog(tt.args.config, tt.args.changelogFile, tt.args.orderBy, tt.args.repoPath, tt.args.beforeTag, tt.args.afterTag, tt.args.unique)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpdatedChangelog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && tt.errMessage != err.Error() {
				t.Errorf("GetUpdatedChangelog() error message = '%v', want '%v'", err.Error(), tt.errMessage)
				return
			}
			if !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("GetUpdatedChangelog() gotMetadata = %v, want %v", gotMetadata, tt.wantMetadata)
			}
			if gotUpdatedChangelog != tt.wantUpdatedChangelog {
				t.Errorf("GetUpdatedChangelog() gotUpdatedChangelog = %v, want %v", gotUpdatedChangelog, tt.wantUpdatedChangelog)
			}
		})
	}
}

// createTestRepo creates a test repo with two tags:
// 0.0.1 and 0.1.0
// The first tag is created 10 seconds before the second tag.
func createTestRepo(t *testing.T) string {
	repoDir := t.TempDir()
	t.Logf("created repo dir: %s", repoDir)

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}

	c1 := commitChange(t, repoDir, "README.md", "first update\r\n", "feat: first update", time.UnixMilli(time.Now().UnixMilli()-10000))
	_, err = repo.CreateTag("0.0.1", c1, nil)
	if err != nil {
		t.Fatal(err)
	}

	c2 := commitChange(t, repoDir, "README.md", "second update\r\n", "feat: second update", time.Now())
	_, err = repo.CreateTag("0.1.0", c2, nil)
	if err != nil {
		t.Fatal(err)
	}

	return repoDir
}

func commitChange(
	t *testing.T,
	repoDir string,
	filename string,
	appendText string,
	msg string,
	when time.Time,
) plumbing.Hash {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	repoPathToFile := path.Join(repoDir, filename)
	file, err := os.OpenFile(repoPathToFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.WriteString(appendText)
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Add(filename)
	if err != nil {
		t.Fatal(err)
	}
	commitSig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  when,
	}
	commit, err := w.Commit(msg, &git.CommitOptions{
		Author:    commitSig,
		Committer: commitSig,
	})
	if err != nil {
		t.Fatal(err)
	}
	return commit
}
