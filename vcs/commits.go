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
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/stringutil"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

const UnreleasedVersionName = "Unreleased"

// FetchCommitMessages returns a slice of commit messages between the given tags.
// If beforeTag is empty, then HEAD is used.
// If afterTag is empty, the oldest commit is used.
func FetchCommitMessages(
	config cfg.SinceConfig,
	repoPath string,
	beforeTag string,
	afterTag string,
	unique bool,
) ([]string, error) {
	commits, err := FetchCommitsByTag(config, repoPath, beforeTag, afterTag, unique)
	if err != nil {
		return nil, err
	}
	return FlattenCommits(commits), nil
}

// FetchCommitsByTag returns a map of commit messages between the given tags.
// The key is the tag metadata, and the value is a slice of commit messages.
// If beforeTag is empty, then HEAD is used.
// If afterTag is empty, the oldest commit is used.
func FetchCommitsByTag(
	config cfg.SinceConfig,
	repoPath string,
	beforeTag string,
	afterTag string,
	unique bool,
) (*[]TagCommits, error) {
	commits, err := fetchCommitsBetween(config, repoPath, beforeTag, afterTag, unique)
	if err != nil {
		return nil, err
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.Tracef("commits by tag: %v", commits)
	} else {
		logrus.Debugf("fetched %d tags\n", len(*commits))
	}
	return commits, nil
}

func FlattenCommits(tags *[]TagCommits) []string {
	var messages []string
	for _, tag := range *tags {
		messages = append(messages, tag.Commits...)
	}
	return messages
}

// fetchCommitsBetween returns a slice of commit messages between the given tags.
// If beforeTag is empty, then HEAD is used.
// If afterTag is empty, the oldest commit is used.
func fetchCommitsBetween(
	config cfg.SinceConfig,
	repoPath string,
	beforeTag string,
	afterTag string,
	unique bool,
) (*[]TagCommits, error) {
	var excludes []*regexp.Regexp
	for _, i := range config.Ignore {
		excludes = append(excludes, regexp.MustCompile(i))
	}

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	var beforeCommit *object.Commit
	if beforeTag != "" {
		beforeTagMeta, err := r.Tag(beforeTag)
		if err != nil {
			return nil, err
		}
		beforeCommit, err = r.CommitObject(beforeTagMeta.Hash())
		if err != nil {
			return nil, err
		}
	}

	var afterCommit *object.Commit
	if afterTag != "" {
		afterTagMeta, err := r.Tag(afterTag)
		if err != nil {
			return nil, err
		}
		afterCommit, err = r.CommitObject(afterTagMeta.Hash())
		if err != nil {
			return nil, err
		}
	}

	allTags, err := listAllTags(r)

	commits, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	tagCommits, err := processCommits(beforeCommit, afterCommit, commits, allTags, unique, excludes)
	if err != nil {
		return nil, err
	}

	return tagCommits, nil
}

func processCommits(
	beforeCommit *object.Commit,
	afterCommit *object.Commit,
	commits object.CommitIter,
	allTags map[string]*TagMeta,
	unique bool,
	excludes []*regexp.Regexp,
) (*[]TagCommits, error) {
	var tagCommits []TagCommits

	currentTag := TagMeta{
		Name: UnreleasedVersionName,
		Date: time.Now(),
	}

	var commitMessages []string

	appendCurrentTag := func() {
		if len(commitMessages) > 0 {
			if unique {
				commitMessages = stringutil.Unique(commitMessages)
			}

			tag := TagCommits{
				TagMeta: currentTag,
				Commits: commitMessages,
			}
			tagCommits = append(tagCommits, tag)
			commitMessages = nil
		}
	}

	// skip commits until reaching beforeTag
	skip := beforeCommit != nil

	err := commits.ForEach(func(c *object.Commit) error {
		if beforeCommit != nil && c.Hash == beforeCommit.Hash {
			skip = false
		}
		if skip {
			return nil
		}

		tagCommit := allTags[c.Hash.String()]
		if tagCommit != nil {
			appendCurrentTag()
			currentTag = *tagCommit
		}

		// stop after appending tag commits for previous tag
		if afterCommit != nil && c.Hash == afterCommit.Hash {
			return storer.ErrStop
		}

		longMessage := c.Message
		if !shouldInclude(longMessage, excludes) {
			return nil
		}
		message := getShortMessage(longMessage)
		commitMessages = append(commitMessages, message)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// final tag
	appendCurrentTag()

	return &tagCommits, nil
}

// listAllTags returns a map of tag hashes to tag metadata.
func listAllTags(r *git.Repository) (map[string]*TagMeta, error) {
	tags := make(map[string]*TagMeta)
	tagRefs, err := r.Tags()
	if err != nil {
		return nil, err
	}
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		commit, err := r.CommitObject(t.Hash())
		if err != nil {
			return err
		}
		tags[t.Hash().String()] = &TagMeta{
			Name: t.Name().Short(),
			Date: commit.Committer.When,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// shouldInclude returns true if the commit message does not match any of the excludes.
func shouldInclude(message string, excludes []*regexp.Regexp) bool {
	for _, exclude := range excludes {
		if exclude.MatchString(message) {
			return false
		}
	}
	return true
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
