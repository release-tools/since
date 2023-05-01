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

package vcs

import "testing"

func Test_getShortMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "short message",
			args: args{
				message: "long message\nwith multiple lines",
			},
			want: "long message",
		},
		{
			name: "remove trailing newlines",
			args: args{
				message: "commit message\n\n",
			},
			want: "commit message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getShortMessage(tt.args.message); got != tt.want {
				t.Errorf("getShortMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
