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
	"github.com/go-git/go-git/v5"
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
			if got := RenderCommits(tt.args.commits, tt.args.groupIntoSections, vcs.UnreleasedVersionName); got != tt.want {
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
	repoDir := createTestRepo(t)

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
				repoPath: repoDir,
				afterTag: "0.1.0",
			},
			wantMetadata:         vcs.ReleaseMetadata{},
			wantUpdatedChangelog: "",
			wantErr:              true,
			errMessage:           "no changes since start tag",
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

	repoPathToReadme := path.Join(repoDir, "README.md")
	readme, err := os.Create(repoPathToReadme)
	if err != nil {
		t.Fatal(err)
	}

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	_, err = readme.WriteString("first update")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}
	c1sig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  time.UnixMilli(time.Now().UnixMilli() - 10000),
	}
	c1, err := w.Commit("first update", &git.CommitOptions{
		Author:    c1sig,
		Committer: c1sig,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CreateTag("0.0.1", c1, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = readme.WriteString("second update")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}
	c2Sig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  time.Now(),
	}
	c2, err := w.Commit("second update", &git.CommitOptions{
		Author:    c2Sig,
		Committer: c2Sig,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CreateTag("0.1.0", c2, nil)
	if err != nil {
		t.Fatal(err)
	}

	return repoDir
}
