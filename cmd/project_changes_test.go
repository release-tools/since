package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func Test_printCommits(t *testing.T) {
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
			want: `
### feat

- feat: foo

### fix

- fix: bar
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			printer := func(s string) { output.WriteString(s + "\n") }
			printCommits(tt.args.commits, printer)
			if !reflect.DeepEqual(output.String(), tt.want) {
				t.Errorf("printChanges() got = %v, want %v", output.String(), tt.want)
			}
		})
	}
}
