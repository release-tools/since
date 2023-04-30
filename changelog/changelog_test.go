package changelog

import "testing"

func TestRenderCommits(t *testing.T) {
	type args struct {
		commits []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no commits",
			args: args{
				commits: []string{},
			},
			want: "",
		},
		{
			name: "print commits",
			args: args{
				commits: []string{"feat: foo", "fix: bar"},
			},
			want: `### feat
- feat: foo

### fix
- fix: bar

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderCommits(tt.args.commits); got != tt.want {
				t.Errorf("RenderCommits() got = %v, want %v", got, tt.want)
			}
		})
	}
}
