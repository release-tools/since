package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/sirupsen/logrus"
	"time"
)

type endTagType string

const (
	endTagEarliest endTagType = "earliest"
	endTagLatest   endTagType = "latest"
)

type TagOrderBy string

const (
	TagOrderAlphabetical TagOrderBy = "alphabetical"
	TagOrderCommitDate   TagOrderBy = "commit-date"
	TagOrderSemver       TagOrderBy = "semver"
)

type TagMeta struct {
	Name string
	Date time.Time
}

type TagCommits struct {
	TagMeta
	Commits []string
}

// Cache the earliest and latest tags in the repository.
var earliestTag, latestTag string

// GetEarliestTag returns the earliest tag in the repository, determined by the given order.
func GetEarliestTag(repoPath string, orderBy TagOrderBy) (string, error) {
	if earliestTag == "" {
		tag, err := getEndTag(repoPath, endTagEarliest, orderBy)
		if err != nil {
			return "", err
		}
		earliestTag = tag
	}
	return earliestTag, nil
}

// GetLatestTag returns the latest tag in the repository, determined by the given order.
func GetLatestTag(repoPath string, orderBy TagOrderBy) (string, error) {
	if latestTag == "" {
		tag, err := getEndTag(repoPath, endTagLatest, orderBy)
		if err != nil {
			return "", err
		}
		latestTag = tag
	}
	return latestTag, nil
}

// getEndTag returns an end tag in the repository, of the given
// end type, determined by the given order.
func getEndTag(repoPath string, endType endTagType, orderBy TagOrderBy) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	tags, err := r.Tags()
	if err != nil {
		return "", err
	}

	var candidateTag *plumbing.Reference
	var candidateCommit *object.Commit
	err = tags.ForEach(func(t *plumbing.Reference) error {
		candidate := false
		if candidateTag == nil {
			candidate = true

		} else {
			switch orderBy {
			case TagOrderAlphabetical:
				switch endType {
				case endTagLatest:
					candidate = t.Name().Short() > candidateTag.Name().Short()
				case endTagEarliest:
					candidate = t.Name().Short() < candidateTag.Name().Short()
				}
				break

			case TagOrderCommitDate:
				commit, err := r.CommitObject(t.Hash())
				if err != nil {
					logrus.Tracef("failed to get commit object for tag %s: %v", t.Name().Short(), err)
					return nil
				}
				var commitDateMatch bool
				switch endType {
				case endTagLatest:
					commitDateMatch = candidateCommit == nil || commit.Committer.When.After(candidateCommit.Committer.When)
				case endTagEarliest:
					commitDateMatch = candidateCommit == nil || commit.Committer.When.Before(candidateCommit.Committer.When)
				}
				if commitDateMatch {
					candidateCommit = commit
					candidate = true
				}
				break

			case TagOrderSemver:
				switch endType {
				case endTagLatest:
					candidate = semver.Compare(t.Name().Short(), candidateTag.Name().Short()) > 0
				case endTagEarliest:
					candidate = semver.Compare(t.Name().Short(), candidateTag.Name().Short()) < 0
				}

			default:
				panic("unknown tag order by: " + orderBy)
			}
		}

		if candidate {
			candidateTag = t
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	tagName := candidateTag.Name().Short()
	logrus.Tracef("%s tag ordered by %s: %s", endType, orderBy, tagName)
	return tagName, nil
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