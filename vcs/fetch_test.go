package vcs

import (
	"github.com/release-tools/since/cfg"
	"testing"
)

func TestFetchCommitMessages(t *testing.T) {
	repoDir := createTestRepo(t)

	commits, err := FetchCommitMessages(cfg.SinceConfig{}, CommitConfig{}, repoDir, "", "0.0.1")
	if err != nil {
		t.Fatalf("FetchCommitMessages() error = %v", err)
	}
	if len(commits) == 0 {
		t.Fatal("FetchCommitMessages() returned no commits")
	}
	found := false
	for _, c := range commits {
		if c == "second update" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("FetchCommitMessages() did not contain 'second update', got %v", commits)
	}
}

func TestFetchCommitMessages_allCommits(t *testing.T) {
	repoDir := createTestRepo(t)

	commits, err := FetchCommitMessages(cfg.SinceConfig{}, CommitConfig{}, repoDir, "", "")
	if err != nil {
		t.Fatalf("FetchCommitMessages() error = %v", err)
	}
	if len(commits) < 2 {
		t.Errorf("FetchCommitMessages() returned %d commits, want at least 2", len(commits))
	}
}

func TestFetchCommitsByTag(t *testing.T) {
	repoDir := createTestRepo(t)

	tagCommits, err := FetchCommitsByTag(cfg.SinceConfig{}, CommitConfig{}, repoDir, "", "0.0.1")
	if err != nil {
		t.Fatalf("FetchCommitsByTag() error = %v", err)
	}
	if tagCommits == nil || len(*tagCommits) == 0 {
		t.Fatal("FetchCommitsByTag() returned no tag commits")
	}
}

func TestFetchCommitMessages_withExcludes(t *testing.T) {
	repoDir := createTestRepo(t)

	config := cfg.SinceConfig{
		Ignore: []string{"^second"},
	}
	commits, err := FetchCommitMessages(config, CommitConfig{}, repoDir, "", "0.0.1")
	if err != nil {
		t.Fatalf("FetchCommitMessages() error = %v", err)
	}
	for _, c := range commits {
		if c == "second update" {
			t.Error("FetchCommitMessages() should have excluded 'second update'")
		}
	}
}

func TestFetchCommitMessages_excludeTagCommits(t *testing.T) {
	repoDir := createTestRepo(t)

	commitCfg := CommitConfig{ExcludeTagCommits: true}
	commits, err := FetchCommitMessages(cfg.SinceConfig{}, commitCfg, repoDir, "", "")
	if err != nil {
		t.Fatalf("FetchCommitMessages() error = %v", err)
	}
	// with tag commits excluded, should have fewer commits
	allCommits, _ := FetchCommitMessages(cfg.SinceConfig{}, CommitConfig{}, repoDir, "", "")
	if len(commits) > len(allCommits) {
		t.Errorf("excluding tag commits should not increase count: excluded=%d, all=%d", len(commits), len(allCommits))
	}
}

func TestFetchCommitMessages_uniqueOnly(t *testing.T) {
	repoDir := createTestRepo(t)

	commitCfg := CommitConfig{UniqueOnly: true}
	commits, err := FetchCommitMessages(cfg.SinceConfig{}, commitCfg, repoDir, "", "")
	if err != nil {
		t.Fatalf("FetchCommitMessages() error = %v", err)
	}
	// verify no duplicates
	seen := make(map[string]bool)
	for _, c := range commits {
		if seen[c] {
			t.Errorf("FetchCommitMessages() with UniqueOnly has duplicate: %v", c)
		}
		seen[c] = true
	}
}

func TestFetchCommitMessages_invalidRepo(t *testing.T) {
	_, err := FetchCommitMessages(cfg.SinceConfig{}, CommitConfig{}, t.TempDir(), "", "")
	if err == nil {
		t.Error("FetchCommitMessages() expected error for invalid repo")
	}
}
