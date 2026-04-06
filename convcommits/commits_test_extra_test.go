package convcommits

import (
	"sort"
	"testing"
)

func TestCategoriseByType_scopedCommits(t *testing.T) {
	commits := []string{"feat(api): add endpoint", "fix(ui): correct layout"}
	got := CategoriseByType(commits)

	if len(got) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(got))
	}
	if _, ok := got["feat"]; !ok {
		t.Errorf("expected 'feat' category, got keys: %v", got)
	}
	if _, ok := got["fix"]; !ok {
		t.Errorf("expected 'fix' category, got keys: %v", got)
	}
}

func TestCategoriseByType_breakingChange(t *testing.T) {
	commits := []string{"feat!: breaking feature"}
	got := CategoriseByType(commits)

	if _, ok := got["BREAKING CHANGE"]; !ok {
		t.Errorf("expected 'BREAKING CHANGE' category, got keys: %v", got)
	}
}

func TestCategoriseByType_noPrefix(t *testing.T) {
	commits := []string{"some random commit message"}
	got := CategoriseByType(commits)

	if _, ok := got[""]; !ok {
		t.Errorf("expected empty prefix category, got keys: %v", got)
	}
}

func TestDetermineTypes(t *testing.T) {
	commits := []string{"feat: foo", "fix: bar", "feat: baz"}
	got := DetermineTypes(commits)
	sort.Strings(got)

	want := []string{"feat", "fix"}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("DetermineTypes() length = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("DetermineTypes()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
