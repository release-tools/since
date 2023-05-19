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

package hooks

import (
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/vcs"
	"testing"
)

func TestExecuteHooks(t *testing.T) {
	type args struct {
		config   cfg.SinceConfig
		hookType HookType
		metadata vcs.ReleaseMetadata
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no hooks",
			args:    args{config: cfg.SinceConfig{}, hookType: Before, metadata: vcs.ReleaseMetadata{}},
			wantErr: false,
		},
		{
			name: "successful command hook",
			args: args{
				config: cfg.SinceConfig{
					Before: []cfg.Hook{
						{
							Command: "echo",
							Args:    []string{"hello world"},
						},
					},
				},
				metadata: vcs.ReleaseMetadata{},
				hookType: Before,
			},
			wantErr: false,
		},
		{
			name: "successful script hook",
			args: args{
				config: cfg.SinceConfig{
					Before: []cfg.Hook{
						{
							Script: `echo "hello world"`,
						},
					},
				},
				metadata: vcs.ReleaseMetadata{},
				hookType: Before,
			},
			wantErr: false,
		},
		{
			name: "failing hook",
			args: args{
				config: cfg.SinceConfig{
					Before: []cfg.Hook{
						{
							Command: "false",
							Args:    []string{},
						},
					},
				},
				metadata: vcs.ReleaseMetadata{},
				hookType: Before,
			},
			wantErr: true,
		},
		{
			name: "env substitution",
			args: args{
				config: cfg.SinceConfig{
					Before: []cfg.Hook{
						{
							Command: "bash",
							Args:    []string{"-c", `[ "$SINCE_SHA" == "1234567890" ]`},
						},
					},
				},
				metadata: vcs.ReleaseMetadata{
					Sha: "1234567890",
				},
				hookType: Before,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecuteHooks(tt.args.config, tt.args.hookType, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteHooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
