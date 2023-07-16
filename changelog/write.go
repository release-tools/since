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

package changelog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// WriteChangelog updates the changelog file with the given content.
func WriteChangelog(changelogFile string, updatedChangelog string) error {
	tempChangelog := writeTempChangelog(updatedChangelog)
	err := os.Rename(tempChangelog, changelogFile)
	if err != nil {
		panic(fmt.Errorf("failed to rename temp file: %s: %w", tempChangelog, err))
	}
	logrus.Debugf("updated changelog: %s", changelogFile)
	return err
}

// writeTempChangelog writes the updated changelog to a temp file and returns the path to the temp file.
func writeTempChangelog(content string) string {
	temp, err := os.CreateTemp(os.TempDir(), "changelog*.md")
	if err != nil {
		panic(fmt.Errorf("failed to create temp file: %w", err))
	}
	_, err = temp.WriteString(content + "\n")
	if err != nil {
		panic(fmt.Errorf("failed to write to temp file: %w", err))
	}
	_ = temp.Close()
	return temp.Name()
}
