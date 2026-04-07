package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/release-tools/since/cfg"
	"os"
	"path"
	"testing"
	"time"
)

func TestGetHeadSha(t *testing.T) {
	repoDir := createTestRepo(t)

	sha, err := GetHeadSha(repoDir)
	if err != nil {
		t.Fatalf("GetHeadSha() error = %v", err)
	}
	if len(sha) != 40 {
		t.Errorf("GetHeadSha() sha length = %d, want 40", len(sha))
	}
}

func TestGetHeadSha_invalidRepo(t *testing.T) {
	_, err := GetHeadSha(t.TempDir())
	if err == nil {
		t.Error("GetHeadSha() expected error for invalid repo")
	}
}

func TestCheckBranch_noBranchRequired(t *testing.T) {
	repoDir := createTestRepo(t)
	config := cfg.SinceConfig{}

	err := CheckBranch(repoDir, config)
	if err != nil {
		t.Errorf("CheckBranch() with no required branch error = %v", err)
	}
}

func TestCheckBranch_wrongBranch(t *testing.T) {
	repoDir := createTestRepo(t)
	config := cfg.SinceConfig{RequireBranch: "release"}

	err := CheckBranch(repoDir, config)
	if err == nil {
		t.Error("CheckBranch() expected error for wrong branch")
	}
}

func TestCommitChangelog(t *testing.T) {
	repoDir := createTestRepo(t)

	// configure git user so CommitChangelog can create a commit without
	// relying on the global git config (which may not exist in CI)
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := repo.Config()
	if err != nil {
		t.Fatal(err)
	}
	cfg.User.Name = "user"
	cfg.User.Email = "user@example.com"
	err = repo.SetConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	changelogPath := path.Join(repoDir, "CHANGELOG.md")
	err = os.WriteFile(changelogPath, []byte("# Changelog\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sha, err := CommitChangelog(repoDir, changelogPath, "1.0.0")
	if err != nil {
		t.Fatalf("CommitChangelog() error = %v", err)
	}
	if len(sha) != 40 {
		t.Errorf("CommitChangelog() sha length = %d, want 40", len(sha))
	}
}

func TestTagRelease(t *testing.T) {
	repoDir := createTestRepo(t)

	sha, err := GetHeadSha(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	err = TagRelease(repoDir, sha, "v1.0.0")
	if err != nil {
		t.Fatalf("TagRelease() error = %v", err)
	}

	// reset cached tags
	earliestTag = ""
	latestTag = ""

	// verify tag exists
	got, err := GetLatestTag(repoDir, TagOrderSemver)
	if err != nil {
		t.Fatal(err)
	}
	if got != "v1.0.0" {
		t.Errorf("TagRelease() latest tag = %v, want v1.0.0", got)
	}

	// reset cached tags for other tests
	earliestTag = ""
	latestTag = ""
}

func TestGetEarliestTag(t *testing.T) {
	repoDir := createTestRepo(t)

	// reset cached tags
	earliestTag = ""

	got, err := GetEarliestTag(repoDir, TagOrderSemver)
	if err != nil {
		t.Fatalf("GetEarliestTag() error = %v", err)
	}
	if got != "0.0.1" {
		t.Errorf("GetEarliestTag() = %v, want 0.0.1", got)
	}

	// reset cached tags for other tests
	earliestTag = ""
}

func TestGetLatestTag(t *testing.T) {
	repoDir := createTestRepo(t)

	// reset cached tags
	latestTag = ""

	got, err := GetLatestTag(repoDir, TagOrderSemver)
	if err != nil {
		t.Fatalf("GetLatestTag() error = %v", err)
	}
	if got != "0.1.0" {
		t.Errorf("GetLatestTag() = %v, want 0.1.0", got)
	}

	// reset cached tags for other tests
	latestTag = ""
}

// createTestRepoForOps creates a minimal test repo with two tags.
func createTestRepoForOps(t *testing.T) string {
	repoDir := t.TempDir()

	repoPathToReadme := path.Join(repoDir, "README.md")
	readme, err := os.Create(repoPathToReadme)
	if err != nil {
		t.Fatal(err)
	}

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	_, err = readme.WriteString("first update")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}
	sig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  time.Now(),
	}
	c, err := w.Commit("initial commit", &git.CommitOptions{
		Author:    sig,
		Committer: sig,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.CreateTag("0.0.1", c, nil)
	if err != nil {
		t.Fatal(err)
	}

	return repoDir
}
