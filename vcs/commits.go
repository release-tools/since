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
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/release-tools/since/cfg"
	"github.com/release-tools/since/stringutil"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// FetchCommitMessages returns a slice of commit messages after the given tag.
func FetchCommitMessages(
	config cfg.SinceConfig,
	repoPath string,
	tag string,
	orderBy TagOrderBy,
	unique bool,
) ([]string, error) {
	if tag == "" {
		latestTag, err := GetLatestTag(repoPath, orderBy)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("most recent tag: %s", latestTag)
		tag = latestTag
	}
	commits, err := fetchCommitsAfter(config, repoPath, tag, unique)
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
func fetchCommitsAfter(
	config cfg.SinceConfig,
	repoPath string,
	tag string,
	unique bool,
) ([]string, error) {
	var excludes []*regexp.Regexp
	for _, i := range config.Ignore {
		excludes = append(excludes, regexp.MustCompile(i))
	}

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
	if unique {
		commitMessages = stringutil.Unique(commitMessages)
	}
	return commitMessages, nil
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
