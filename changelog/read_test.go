package changelog

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestResolveChangelogFile(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		fileName string
		want     string
	}{
		{
			name:     "simple filename",
			dir:      "/repo",
			fileName: "CHANGELOG.md",
			want:     "/repo/CHANGELOG.md",
		},
		{
			name:     "absolute path with forward slash",
			dir:      "/repo",
			fileName: "/other/CHANGELOG.md",
			want:     "/other/CHANGELOG.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveChangelogFile(tt.dir, tt.fileName); got != tt.want {
				t.Errorf("ResolveChangelogFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "test.md")
	err := os.WriteFile(filePath, []byte("line1\nline2\nline3"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want := []string{"line1", "line2", "line3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadFile() = %v, want %v", got, want)
	}
}

func TestReadFile_nonExistent(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.md")
	if err == nil {
		t.Error("ReadFile() expected error for non-existent file")
	}
}

func TestParseChangelog(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "CHANGELOG.md")
	content := `# Changelog

## [1.0.0] - 2024-01-01
### Added
- feat: something new

## [0.9.0] - 2023-12-01
### Fixed
- fix: old bug
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ParseChangelog(filePath, "1.0.0", false)
	if err != nil {
		t.Fatalf("ParseChangelog() error = %v", err)
	}
	want := []string{"### Added", "- feat: something new"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseChangelog() = %v, want %v", got, want)
	}
}

func TestParseChangelog_withHeader(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "CHANGELOG.md")
	content := `# Changelog

## [1.0.0] - 2024-01-01
### Added
- feat: something new

## [0.9.0] - 2023-12-01
### Fixed
- fix: old bug
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ParseChangelog(filePath, "1.0.0", true)
	if err != nil {
		t.Fatalf("ParseChangelog() error = %v", err)
	}
	if len(got) == 0 {
		t.Fatal("ParseChangelog() returned empty result")
	}
	if got[0] != "## [1.0.0] - 2024-01-01" {
		t.Errorf("ParseChangelog() first line = %v, want version header", got[0])
	}
}

func TestParseChangelog_nonExistent(t *testing.T) {
	_, err := ParseChangelog("/nonexistent/CHANGELOG.md", "1.0.0", false)
	if err == nil {
		t.Error("ParseChangelog() expected error for non-existent file")
	}
}

func Test_readChanges_versionNotFound(t *testing.T) {
	lines := []string{
		"# Changelog",
		"",
		"## [1.0.0] - 2024-01-01",
		"- feat: foo",
	}
	_, err := readChanges(lines, "9.9.9", false)
	if err == nil {
		t.Fatal("readChanges() expected error for missing version")
	}
	want := "could not find version 9.9.9 in changelog"
	if err.Error() != want {
		t.Errorf("readChanges() error = %q, want %q", err.Error(), want)
	}
}

func Test_readChanges_noVersionSections(t *testing.T) {
	lines := []string{
		"# Changelog",
		"",
		"All notable changes are documented here.",
	}
	_, err := readChanges(lines, "", false)
	if err == nil {
		t.Fatal("readChanges() expected error when no version sections present")
	}
	want := "changelog contains no version sections"
	if err.Error() != want {
		t.Errorf("readChanges() error = %q, want %q", err.Error(), want)
	}
}

func Test_readChanges_lastVersionToEndOfFile(t *testing.T) {
	// trailing empty string mirrors ReadFile splitting a file that ends in a
	// newline; readChanges relies on it to bound the final section
	lines := []string{
		"# Changelog",
		"",
		"## [1.0.0] - 2024-01-01",
		"### Added",
		"- feat: foo",
		"",
		"## [0.9.0] - 2023-12-01",
		"### Fixed",
		"- fix: bar",
		"",
	}
	got, err := readChanges(lines, "0.9.0", true)
	if err != nil {
		t.Fatalf("readChanges() error = %v", err)
	}
	want := []string{"## [0.9.0] - 2023-12-01", "### Fixed", "- fix: bar"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readChanges() = %v, want %v", got, want)
	}
}

func Test_readChanges_firstVersion(t *testing.T) {
	lines := []string{
		"# Changelog",
		"",
		"## [1.0.0] - 2024-01-01",
		"### Added",
		"- feat: foo",
		"",
		"## [0.9.0] - 2023-12-01",
		"- fix: bar",
	}
	got, err := readChanges(lines, "", false)
	if err != nil {
		t.Fatalf("readChanges() error = %v", err)
	}
	want := []string{"### Added", "- feat: foo"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readChanges() = %v, want %v", got, want)
	}
}
