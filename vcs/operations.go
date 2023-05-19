/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vcs

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/release-tools/since/cfg"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/sirupsen/logrus"
	"strings"
)

type ReleaseMetadata struct {
	NewVersion string
	OldVersion string
	RepoPath   string
	Sha        string
	VPrefix    bool
}

type TagOrderBy string

const (
	TagOrderAlphabetical TagOrderBy = "alphabetical"
	TagOrderCommitDate   TagOrderBy = "commit-date"
	TagOrderSemver       TagOrderBy = "semver"
)

var latestTag string

// GetLatestTag returns the latest tag in the repository.
func GetLatestTag(repoPath string, orderBy TagOrderBy) (string, error) {
	if latestTag == "" {
		tag, err := getLatestTag(repoPath, orderBy)
		if err != nil {
			return "", err
		}
		latestTag = tag
	}
	return latestTag, nil
}

// getLatestTag returns the latest tag in the repository, determined by the given order.
func getLatestTag(repoPath string, orderBy TagOrderBy) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	tags, err := r.Tags()
	if err != nil {
		return "", err
	}

	var latestTag *plumbing.Reference
	var latestCommit *object.Commit
	err = tags.ForEach(func(t *plumbing.Reference) error {
		latest := false
		if latestTag == nil {
			latest = true

		} else {
			switch orderBy {
			case TagOrderAlphabetical:
				latest = t.Name().Short() > latestTag.Name().Short()
				break

			case TagOrderCommitDate:
				commit, err := r.CommitObject(t.Hash())
				if err != nil {
					logrus.Tracef("failed to get commit object for tag %s: %v", t.Name().Short(), err)
					return nil
				}
				if latestCommit == nil || commit.Committer.When.After(latestCommit.Committer.When) {
					latestCommit = commit
					latest = true
				}
				break

			case TagOrderSemver:
				latest = semver.Compare(t.Name().Short(), latestTag.Name().Short()) > 0

			default:
				panic("unknown tag order by: " + orderBy)
			}
		}

		if latest {
			latestTag = t
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	tagName := latestTag.Name().Short()
	logrus.Tracef("latest tag ordered by %s: %s", orderBy, tagName)
	return tagName, nil
}

// CommitChangelog commits the changelog file.
func CommitChangelog(repoPath string, changelogFile string, version string) (hash string, err error) {
	// make relative to repo root
	repoPathToChangelog := strings.TrimPrefix(changelogFile, repoPath)
	if strings.HasPrefix(repoPathToChangelog, "/") || strings.HasPrefix(repoPathToChangelog, "\\") {
		repoPathToChangelog = repoPathToChangelog[1:]
	}

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	_, err = w.Add(repoPathToChangelog)
	if err != nil {
		return "", err
	}
	commit, err := w.Commit("build: release "+version+".", &git.CommitOptions{})
	if err != nil {
		return "", err
	}
	sha := commit.String()

	logrus.Debugf("committed changelog %s with %s", changelogFile, sha)
	return sha, nil
}

// TagRelease tags the repository with the given version.
func TagRelease(repoPath string, hash string, version string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	_, err = r.CreateTag(version, plumbing.NewHash(hash), nil)
	if err != nil {
		return err
	}
	logrus.Debugf("tagged %s with %s", hash, version)
	return nil
}

// GetHeadSha returns the SHA of the HEAD commit.
func GetHeadSha(repoPath string) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	head, err := r.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

// CheckBranch checks if the current branch is the required branch.
func CheckBranch(repoPath string, config cfg.SinceConfig) error {
	if config.RequireBranch == "" {
		return nil
	}
	branch, err := getCurrentBranch(repoPath)
	if err != nil {
		return err
	}
	if branch != config.RequireBranch {
		return fmt.Errorf("not on branch %s", config.RequireBranch)
	}
	return nil
}

// getCurrentBranch returns the current branch name.
func getCurrentBranch(repoPath string) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	head, err := r.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}
