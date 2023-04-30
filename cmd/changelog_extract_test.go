package cmd

import (
	"reflect"
	"testing"
)

func Test_printChanges(t *testing.T) {
	type args struct {
		path          string
		version       string
		includeHeader bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "read and print entries from specific version",
			args: args{path: "testdata/simple.md", version: "0.1.0", includeHeader: true},
			want: "## [0.1.0] - 2023-03-04\n### Added\n- feat: another change.\n",
		},
		{
			name: "read latest version",
			args: args{path: "testdata/simple.md", includeHeader: true},
			want: `## [0.2.0] - 2023-03-05
### Added
- feat: some change.
`,
			wantErr: false,
		},
		{
			name: "skip version header",
			args: args{path: "testdata/simple.md", includeHeader: false},
			want: `### Added
- feat: some change.
`,
			wantErr: false,
		},
		{
			name: "read earlier version",
			args: args{path: "testdata/simple.md", version: "0.1.0", includeHeader: true},
			want: `## [0.1.0] - 2023-03-04
### Added
- feat: another change.
`,
			wantErr: false,
		},
		{
			name:    "return error for nonexistent changelog file",
			args:    args{path: "/tmp/nonexistent-changelog"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := printChanges(tt.args.path, tt.args.version, tt.args.includeHeader)
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
