package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/sirupsen/logrus"
	"strings"
)

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
	logrus.Tracef("commits: %v", commits)
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

func getShortMessage(message string) string {
	var short string
	if strings.Contains(message, "\n") {
		short = strings.Split(message, "\n")[0]
	} else {
		short = message
	}
	return strings.TrimSpace(short)
}
