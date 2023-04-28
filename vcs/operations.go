package vcs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// GetLatestTag returns the latest tag in the repository.
func GetLatestTag(repoPath string) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	tags, err := r.Tags()
	if err != nil {
		return "", err
	}
	var latestTag string
	err = tags.ForEach(func(t *plumbing.Reference) error {
		latestTag = t.Name().Short()
		return nil
	})
	if err != nil {
		return "", err
	}
	return latestTag, nil
}

// FetchCommitsAfter returns a slice of commit messages after the given tag.
func FetchCommitsAfter(repoPath string, tag string) ([]string, error) {
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
		commitMessages = append(commitMessages, c.Message)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commitMessages, nil
}