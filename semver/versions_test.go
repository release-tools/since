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

import "testing"

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
