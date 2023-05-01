/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/

package changelog

import (
	"fmt"
	"os"
)

// UpdateChangelog updates the changelog file with the given content.
func UpdateChangelog(changelogFile string, updatedChangelog string) error {
	tempChangelog := writeTempChangelog(updatedChangelog)
	err := os.Rename(tempChangelog, changelogFile)
	if err != nil {
		panic(fmt.Errorf("failed to rename temp file: %s: %w", tempChangelog, err))
	}
	return err
}

// writeTempChangelog writes the updated changelog to a temp file and returns the path to the temp file.
func writeTempChangelog(updatedChangelog string) string {
	temp, err := os.CreateTemp(os.TempDir(), "changelog*.md")
	if err != nil {
		panic(fmt.Errorf("failed to create temp file: %w", err))
	}
	_, err = temp.WriteString(updatedChangelog + "\n")
	if err != nil {
		panic(fmt.Errorf("failed to write to temp file: %w", err))
	}
	_ = temp.Close()
	return temp.Name()
}
