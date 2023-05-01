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
	"github.com/outofcoffee/since/semver"
	"testing"
)

func Test_getNextVersion(t *testing.T) {
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
			name: "no changes",
			args: args{
				currentVersion: "1.0.0",
				vPrefix:        false,
				commits:        []string{},
			},
			want: "",
		},
		{
			name: "patch",
			args: args{
				currentVersion: "1.0.1",
				vPrefix:        false,
				commits:        []string{"fix: foo"},
			},
			want: "1.0.2",
		},
		{
			name: "minor",
			args: args{
				currentVersion: "1.0.1",
				vPrefix:        false,
				commits:        []string{"feat: foo"},
			},
			want: "1.1.0",
		},
		{
			name: "major",
			args: args{
				currentVersion: "1.0.1",
				vPrefix:        false,
				commits:        []string{"feat!: foo"},
			},
			want: "2.0.0",
		},
		{
			name: "breaking change",
			args: args{
				currentVersion: "1.0.1",
				vPrefix:        false,
				commits:        []string{"BREAKING CHANGE: foo"},
			},
			want: "2.0.0",
		},
		{
			name: "v prefix",
			args: args{
				currentVersion: "1.0.1",
				vPrefix:        true,
				commits:        []string{"feat: foo"},
			},
			want: "v1.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.GetNextVersion(tt.args.currentVersion, tt.args.vPrefix, tt.args.commits); got != tt.want {
				t.Errorf("getNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
