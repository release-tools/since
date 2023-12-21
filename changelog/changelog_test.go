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
	"github.com/release-tools/since/vcs"
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
