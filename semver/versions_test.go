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

package semver

import (
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/release-tools/since/vcs"
)

func TestGetNextVersion(t *testing.T) {
	type args struct {
		currentVersion string
		vPrefix        bool
		commits        []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "major",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"feat!: major change",
					"feat: new feature",
					"fix: all bugs fixed",
				},
			},
			want: "2.0.0",
		},
		{
			name: "minor",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"feat: new feature",
					"fix: all bugs fixed",
				},
			},
			want: "1.3.0",
		},
		{
			name: "patch",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"fix: all bugs fixed",
				},
			},
			want: "1.2.4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNextVersion(tt.args.currentVersion, tt.args.vPrefix, tt.args.commits); got != tt.want {
				t.Errorf("GetNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetCurrentVersion verifies that a "v" prefixed tag is reported with the
// prefix stripped and vPrefix set. Note: vcs caches the latest tag in a
// package-level variable that this package cannot reset, so this test performs
// a single repository lookup to avoid cross-test contamination.
func TestGetCurrentVersion(t *testing.T) {
	repoDir := t.TempDir()

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoDir+"/README.md", []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Add("README.md"); err != nil {
		t.Fatal(err)
	}
	sig := &object.Signature{Name: "user", Email: "user@example.com", When: time.Now()}
	c, err := w.Commit("feat: first", &git.CommitOptions{Author: sig, Committer: sig})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateTag("v2.3.4", c, nil); err != nil {
		t.Fatal(err)
	}

	version, vPrefix, err := GetCurrentVersion(repoDir, vcs.TagOrderSemver)
	if err != nil {
		t.Fatalf("GetCurrentVersion() error = %v", err)
	}
	if version != "2.3.4" {
		t.Errorf("GetCurrentVersion() version = %q, want %q", version, "2.3.4")
	}
	if !vPrefix {
		t.Errorf("GetCurrentVersion() vPrefix = %v, want true", vPrefix)
	}
}

func TestDetermineChangeType(t *testing.T) {
	type args struct {
		types []string
	}
	tests := []struct {
		name string
		args args
		want Component
	}{
		{
			name: "major",
			args: args{
				types: []string{
					"BREAKING CHANGE",
				},
			},
			want: ComponentMajor,
		},
		{
			name: "minor",
			args: args{
				types: []string{
					"feat",
				},
			},
			want: ComponentMinor,
		},
		{
			name: "patch",
			args: args{
				types: []string{
					"fix",
				},
			},
			want: ComponentPatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetermineChangeType(tt.args.types); got != tt.want {
				t.Errorf("DetermineChangeType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNextVersion_withVPrefix(t *testing.T) {
	got := GetNextVersion("1.2.3", true, []string{"feat: new feature"})
	want := "v1.3.0"
	if got != want {
		t.Errorf("GetNextVersion() with vPrefix = %v, want %v", got, want)
	}
}

func TestGetNextVersion_noChanges(t *testing.T) {
	got := GetNextVersion("1.2.3", false, []string{"unknown: something"})
	want := ""
	if got != want {
		t.Errorf("GetNextVersion() with no recognised changes = %v, want %v", got, want)
	}
}

func TestDetermineChangeType_allPatchTypes(t *testing.T) {
	patchTypes := []string{"build", "chore", "ci", "docs", "fix", "refactor", "security", "style", "test"}
	for _, pt := range patchTypes {
		t.Run(pt, func(t *testing.T) {
			got := DetermineChangeType([]string{pt})
			if got != ComponentPatch {
				t.Errorf("DetermineChangeType(%v) = %v, want %v", pt, got, ComponentPatch)
			}
		})
	}
}

func TestDetermineChangeType_none(t *testing.T) {
	got := DetermineChangeType([]string{"unknown"})
	if got != ComponentNone {
		t.Errorf("DetermineChangeType(unknown) = %v, want %v", got, ComponentNone)
	}
}

func TestDetermineChangeType_majorTakesPrecedence(t *testing.T) {
	got := DetermineChangeType([]string{"feat", "BREAKING CHANGE", "fix"})
	if got != ComponentMajor {
		t.Errorf("DetermineChangeType() = %v, want %v", got, ComponentMajor)
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		component Component
		want      string
	}{
		{
			name:      "bump major resets minor and patch",
			version:   "1.5.9",
			component: ComponentMajor,
			want:      "2.0.0",
		},
		{
			name:      "bump minor resets patch",
			version:   "1.5.9",
			component: ComponentMinor,
			want:      "1.6.0",
		},
		{
			name:      "bump patch only",
			version:   "1.5.9",
			component: ComponentPatch,
			want:      "1.5.10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bumpVersion(tt.version, tt.component); got != tt.want {
				t.Errorf("bumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
