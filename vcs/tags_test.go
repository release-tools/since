package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func Test_isTagSigningEnabled(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  bool
	}{
		{name: "unset", key: "", want: false},
		{name: "tag.gpgSign=true", key: "tag.gpgSign", value: "true", want: true},
		{name: "tag.gpgSign=false", key: "tag.gpgSign", value: "false", want: false},
		{name: "tag.forceSignAnnotated=true", key: "tag.forceSignAnnotated", value: "true", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := createTestRepo(t)
			if tt.key != "" {
				repo, err := git.PlainOpen(repoDir)
				if err != nil {
					t.Fatal(err)
				}
				cfg, err := repo.Config()
				if err != nil {
					t.Fatal(err)
				}
				cfg.Raw.Section("tag").SetOption(strings.TrimPrefix(tt.key, "tag."), tt.value)
				if err := repo.SetConfig(cfg); err != nil {
					t.Fatal(err)
				}
			}

			got, err := isTagSigningEnabled(repoDir)
			if err != nil {
				t.Fatalf("isTagSigningEnabled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("isTagSigningEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// setGitConfig sets a raw git config option in the given section for the repo.
func setGitConfig(t *testing.T, repoDir, section, option, value string) {
	t.Helper()
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatal(err)
	}
	conf, err := repo.Config()
	if err != nil {
		t.Fatal(err)
	}
	conf.Raw.Section(section).SetOption(option, value)
	if err := repo.SetConfig(conf); err != nil {
		t.Fatal(err)
	}
}

func TestTagRelease_signedTagInvokesGitCli(t *testing.T) {
	repoDir := createTestRepo(t)

	// enable signing, but point gpg at a program that always fails so the
	// signing attempt errors immediately rather than prompting for a passphrase.
	setGitConfig(t, repoDir, "tag", "gpgSign", "true")
	setGitConfig(t, repoDir, "gpg", "program", "/nonexistent-gpg-program")

	sha, err := GetHeadSha(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	err = TagRelease(repoDir, sha, "v2.0.0")
	if err == nil {
		t.Fatal("TagRelease() expected error when signing is enabled but gpg is unavailable")
	}
	if !strings.Contains(err.Error(), "git tag -s failed") {
		t.Errorf("TagRelease() error = %v, want it to mention 'git tag -s failed'", err)
	}
}

func TestTagRelease_invalidRepo(t *testing.T) {
	// a bare temp dir is not a git repository, so opening it should fail after
	// the (unsigned) signing check reports no signing config.
	err := TagRelease(t.TempDir(), "0000000000000000000000000000000000000000", "v1.0.0")
	if err == nil {
		t.Error("TagRelease() expected error for a non-repository path")
	}
}

func Test_createSignedTag_failsWithoutSigningKey(t *testing.T) {
	repoDir := createTestRepo(t)
	setGitConfig(t, repoDir, "gpg", "program", "/nonexistent-gpg-program")

	sha, err := GetHeadSha(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	err = createSignedTag(repoDir, sha, "v3.0.0")
	if err == nil {
		t.Fatal("createSignedTag() expected error when gpg program is unavailable")
	}
	if !strings.Contains(err.Error(), "git tag -s failed") {
		t.Errorf("createSignedTag() error = %v, want it to mention 'git tag -s failed'", err)
	}
}

func Test_getEndTag(t *testing.T) {
	repoDir := createTestRepo(t)

	type args struct {
		repoPath string
		endType  endTagType
		orderBy  TagOrderBy
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get latest tag by semver",
			args: args{
				repoPath: repoDir,
				endType:  endTagLatest,
				orderBy:  TagOrderSemver,
			},
			want:    "0.1.0",
			wantErr: false,
		},
		{
			name: "get earliest tag by semver",
			args: args{
				repoPath: repoDir,
				endType:  endTagEarliest,
				orderBy:  TagOrderSemver,
			},
			want:    "0.0.1",
			wantErr: false,
		},
		{
			name: "get latest tag by alphabetical sort",
			args: args{
				repoPath: repoDir,
				endType:  endTagLatest,
				orderBy:  TagOrderAlphabetical,
			},
			want:    "0.1.0",
			wantErr: false,
		},
		{
			name: "get earliest tag by alphabetical sort",
			args: args{
				repoPath: repoDir,
				endType:  endTagEarliest,
				orderBy:  TagOrderAlphabetical,
			},
			want:    "0.0.1",
			wantErr: false,
		},
		{
			name: "get latest tag by date",
			args: args{
				repoPath: repoDir,
				endType:  endTagLatest,
				orderBy:  TagOrderCommitDate,
			},
			want:    "0.1.0",
			wantErr: false,
		},
		{
			name: "get earliest tag by date",
			args: args{
				repoPath: repoDir,
				endType:  endTagEarliest,
				orderBy:  TagOrderCommitDate,
			},
			want:    "0.0.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEndTag(tt.args.repoPath, tt.args.endType, tt.args.orderBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEndTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEndTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// createTestRepo creates a test repo with two tags:
// 0.0.1 and 0.1.0
// The first tag is created 10 seconds before the second tag.
func createTestRepo(t *testing.T) string {
	repoDir := t.TempDir()
	t.Logf("created repo dir: %s", repoDir)

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
	c1sig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  time.UnixMilli(time.Now().UnixMilli() - 10000),
	}
	c1, err := w.Commit("first update", &git.CommitOptions{
		Author:    c1sig,
		Committer: c1sig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// lightweight tag
	_, err = repo.CreateTag("0.0.1", c1, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = readme.WriteString("second update")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}
	c2Sig := &object.Signature{
		Name:  "user",
		Email: "user@example.com",
		When:  time.Now(),
	}
	c2, err := w.Commit("second update", &git.CommitOptions{
		Author:    c2Sig,
		Committer: c2Sig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// annotated tag
	_, err = repo.CreateTag("0.1.0", c2, &git.CreateTagOptions{
		Tagger:  c2Sig,
		Message: "annotated tag 0.1.0",
	})
	if err != nil {
		t.Fatal(err)
	}

	return repoDir
}
