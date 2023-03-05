package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_printChanges(t *testing.T) {
	type args struct {
		path    string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "read and print entries from specific version",
			args: args{path: "testdata/simple.md", version: "0.1.0"},
			want: "## [0.1.0] - 2023-03-04\n### Added\n- feat: another change.\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			printer := func(s string) { output.WriteString(s + "\n") }
			printChanges(tt.args.path, tt.args.version, printer)
			if !reflect.DeepEqual(output.String(), tt.want) {
				t.Errorf("printChanges() got = %v, want %v", output.String(), tt.want)
			}
		})
	}
}

func Test_parseChangelog(t *testing.T) {
	type args struct {
		path    string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "read latest version",
			args: args{path: "testdata/simple.md", version: ""},
			want: []string{
				"## [0.2.0] - 2023-03-05",
				"### Added",
				"- feat: some change.",
			},
			wantErr: false,
		},
		{
			name: "read earlier version",
			args: args{path: "testdata/simple.md", version: "0.1.0"},
			want: []string{
				"## [0.1.0] - 2023-03-04",
				"### Added",
				"- feat: another change.",
			},
			wantErr: false,
		},
		{
			name:    "return error for nonexistent changelog file",
			args:    args{path: "/tmp/nonexistent-changelog"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseChangelog(tt.args.path, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseChangelog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseChangelog() got = %v, want %v", got, tt.want)
			}
		})
	}
}
