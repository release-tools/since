package changelog

import "testing"

func TestNoChangesError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  NoChangesError
		want string
	}{
		{
			name: "excluded commits with start tag",
			err:  NoChangesError{StartTag: "0.1.0", ExcludedCount: 3},
			want: "no eligible commits since 0.1.0 — 3 commit(s) were excluded by ignore patterns in since.yaml",
		},
		{
			name: "excluded commits without start tag",
			err:  NoChangesError{ExcludedCount: 2},
			want: "no eligible commits found — 2 commit(s) were excluded by ignore patterns in since.yaml",
		},
		{
			name: "start tag with no exclusions",
			err:  NoChangesError{StartTag: "1.2.3"},
			want: "no commits since 1.2.3",
		},
		{
			name: "no start tag and no exclusions",
			err:  NoChangesError{},
			want: "no commits found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
