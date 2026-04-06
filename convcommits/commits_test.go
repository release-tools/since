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

package convcommits

import (
	"reflect"
	"sort"
	"testing"
)

func TestCategoriseByType(t *testing.T) {
	type args struct {
		commits []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "no changes",
			args: args{
				commits: []string{},
			},
			want: map[string][]string{},
		},
		{
			name: "categorised",
			args: args{
				commits: []string{"feat: foo", "fix: bar"},
			},
			want: map[string][]string{
				"feat": {"feat: foo"},
				"fix":  {"fix: bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CategoriseByType(tt.args.commits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoriseByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoriseByType_scopedCommits(t *testing.T) {
	commits := []string{"feat(api): add endpoint", "fix(ui): correct layout"}
	got := CategoriseByType(commits)

	if len(got) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(got))
	}
	if _, ok := got["feat"]; !ok {
		t.Errorf("expected 'feat' category, got keys: %v", got)
	}
	if _, ok := got["fix"]; !ok {
		t.Errorf("expected 'fix' category, got keys: %v", got)
	}
}

func TestCategoriseByType_breakingChange(t *testing.T) {
	commits := []string{"feat!: breaking feature"}
	got := CategoriseByType(commits)

	if _, ok := got["BREAKING CHANGE"]; !ok {
		t.Errorf("expected 'BREAKING CHANGE' category, got keys: %v", got)
	}
}

func TestCategoriseByType_noPrefix(t *testing.T) {
	commits := []string{"some random commit message"}
	got := CategoriseByType(commits)

	if _, ok := got[""]; !ok {
		t.Errorf("expected empty prefix category, got keys: %v", got)
	}
}

func TestDetermineTypes(t *testing.T) {
	commits := []string{"feat: foo", "fix: bar", "feat: baz"}
	got := DetermineTypes(commits)
	sort.Strings(got)

	want := []string{"feat", "fix"}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("DetermineTypes() length = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("DetermineTypes()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
