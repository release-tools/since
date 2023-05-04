package hooks

import (
	"github.com/outofcoffee/since/vcs"
	"testing"
)

func TestExecuteHooks(t *testing.T) {
	type args struct {
		config   SinceConfig
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
			args:    args{config: SinceConfig{}, hookType: Before, metadata: vcs.ReleaseMetadata{}},
			wantErr: false,
		},
		{
			name: "successful hook",
			args: args{
				config: SinceConfig{
					Before: []Hook{
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
			name: "failing hook",
			args: args{
				config: SinceConfig{
					Before: []Hook{
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
				config: SinceConfig{
					Before: []Hook{
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
