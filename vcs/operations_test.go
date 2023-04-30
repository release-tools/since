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
