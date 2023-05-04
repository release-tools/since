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
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/outofcoffee/since/cfg"
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

// FetchCommitMessages returns a slice of commit messages after the given tag.
func FetchCommitMessages(repoPath string, tag string, orderBy TagOrderBy) ([]string, error) {
	if tag == "" {
		latestTag, err := GetLatestTag(repoPath, orderBy)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("latest tag: %s", latestTag)
		tag = latestTag
	}
	commits, err := fetchCommitsAfter(repoPath, tag)
	if err != nil {
		return nil, err
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.Tracef("commits: %v", commits)
	} else {
		logrus.Debugf("fetched %d commits\n", len(commits))
	}
	return commits, nil
}

// fetchCommitsAfter returns a slice of commit messages after the given tag.
func fetchCommitsAfter(repoPath string, tag string) ([]string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	afterTag, err := r.Tag(tag)
	if err != nil {
		return nil, err
	}
	afterTagCommit, err := r.CommitObject(afterTag.Hash())
	if err != nil {
		return nil, err
	}
	commits, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}
	var commitMessages []string
	err = commits.ForEach(func(c *object.Commit) error {
		if c.Hash == afterTagCommit.Hash {
			return storer.ErrStop
		}
		message := getShortMessage(c.Message)
		commitMessages = append(commitMessages, message)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commitMessages, nil
}

// getShortMessage returns the first line of a commit message.
func getShortMessage(message string) string {
	var short string
	if strings.Contains(message, "\n") {
		short = strings.Split(message, "\n")[0]
	} else {
		short = message
	}
	return strings.TrimSpace(short)
}

// CommitChangelog commits the changelog file.
func CommitChangelog(repoPath string, changelogFile string, version string) (hash string, err error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	_, err = w.Add(changelogFile)
	if err != nil {
		return "", err
	}
	commit, err := w.Commit("build: release "+version+".", &git.CommitOptions{})
	if err != nil {
		return "", err
	}
	return commit.String(), nil
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
