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
