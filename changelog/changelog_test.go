package changelog

import (
	"reflect"
	"testing"
)

func TestRenderCommits(t *testing.T) {
	type args struct {
		groupIntoSections bool
		commits           []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no commits",
			args: args{
				commits:           []string{},
				groupIntoSections: false,
			},
			want: "",
		},
		{
			name: "print commits",
			args: args{
				commits:           []string{"feat: foo", "fix: bar"},
				groupIntoSections: false,
			},
			want: `### feat
- feat: foo

### fix
- fix: bar`,
		},
		{
			name: "group commits",
			args: args{
				commits:           []string{"feat: foo", "fix: bar", "ci: baz", "chore: qux", "build: quux"},
				groupIntoSections: true,
			},
			want: `### Added
- feat: foo

### Changed
- build: quux
- chore: qux
- ci: baz

### Fixed
- fix: bar`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderCommits(tt.args.commits, tt.args.groupIntoSections); got != tt.want {
				t.Errorf("RenderCommits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitIntoSections(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    ChangelogSections
		wantErr bool
	}{
		{
			name: "no sections",
			args: args{
				lines: []string{},
			},
			want:    ChangelogSections{},
			wantErr: true,
		},
		{
			name: "split into sections",
			args: args{
				lines: []string{"# Change Log", "", "## [0.1.0]", "### feat", "- feat: foo", "", "### fix", "- fix: bar"},
			},
			want: ChangelogSections{
				Boilerplate: "# Change Log\n\n",
				Body:        "## [0.1.0]\n### feat\n- feat: foo\n\n### fix\n- fix: bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitIntoSections(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitIntoSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Boilerplate, tt.want.Boilerplate) {
				t.Errorf("SplitIntoSections() Boilerplate got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Body, tt.want.Body) {
				t.Errorf("SplitIntoSections() Body got = %v, want %v", got, tt.want)
			}
		})
	}
}
