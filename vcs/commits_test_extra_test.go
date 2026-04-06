package vcs

import (
	"reflect"
	"regexp"
	"testing"
)

func TestFlattenCommits(t *testing.T) {
	tags := []TagCommits{
		{
			TagMeta: TagMeta{Name: "v1"},
			Commits: []string{"feat: a", "fix: b"},
		},
		{
			TagMeta: TagMeta{Name: "v2"},
			Commits: []string{"chore: c"},
		},
	}
	got := FlattenCommits(&tags)
	want := []string{"feat: a", "fix: b", "chore: c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FlattenCommits() = %v, want %v", got, want)
	}
}

func TestFlattenCommits_empty(t *testing.T) {
	tags := []TagCommits{}
	got := FlattenCommits(&tags)
	if len(got) != 0 {
		t.Errorf("FlattenCommits() = %v, want empty", got)
	}
}

func Test_shouldInclude(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		excludes []*regexp.Regexp
		want     bool
	}{
		{
			name:     "no excludes",
			message:  "feat: something",
			excludes: nil,
			want:     true,
		},
		{
			name:     "matching exclude",
			message:  "build: release v1.0.0",
			excludes: []*regexp.Regexp{regexp.MustCompile(`^build: release`)},
			want:     false,
		},
		{
			name:     "non-matching exclude",
			message:  "feat: new feature",
			excludes: []*regexp.Regexp{regexp.MustCompile(`^build: release`)},
			want:     true,
		},
		{
			name:    "multiple excludes with one match",
			message: "chore: bump deps",
			excludes: []*regexp.Regexp{
				regexp.MustCompile(`^build: release`),
				regexp.MustCompile(`^chore: bump`),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldInclude(tt.message, tt.excludes); got != tt.want {
				t.Errorf("shouldInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getShortMessage_singleLine(t *testing.T) {
	got := getShortMessage("simple message")
	want := "simple message"
	if got != want {
		t.Errorf("getShortMessage() = %v, want %v", got, want)
	}
}

func Test_getShortMessage_withLeadingWhitespace(t *testing.T) {
	got := getShortMessage("  message with spaces  ")
	want := "message with spaces"
	if got != want {
		t.Errorf("getShortMessage() = %v, want %v", got, want)
	}
}
