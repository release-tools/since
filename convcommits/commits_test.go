package convcommits

import (
	"reflect"
	"testing"
)

func TestCategoriseByType(t *testing.T) {
	type args struct {
		commits []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "no changes",
			args: args{
				commits: []string{},
			},
			want: map[string][]string{},
		},
		{
			name: "categorised",
			args: args{
				commits: []string{"feat: foo", "fix: bar"},
			},
			want: map[string][]string{
				"feat": {"feat: foo"},
				"fix":  {"fix: bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CategoriseByType(tt.args.commits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoriseByType() = %v, want %v", got, tt.want)
			}
		})
	}
}
