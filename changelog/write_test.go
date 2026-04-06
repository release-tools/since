package changelog

import (
	"os"
	"path"
	"testing"
)

func TestWriteChangelog(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "CHANGELOG.md")

	content := "# Changelog\n\n## [1.0.0]\n- feat: foo\n"
	err := WriteChangelog(filePath, content)
	if err != nil {
		t.Fatalf("WriteChangelog() error = %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	// WriteChangelog appends a newline
	want := content + "\n"
	if string(got) != want {
		t.Errorf("WriteChangelog() file content = %q, want %q", string(got), want)
	}
}

func TestWriteChangelog_overwrite(t *testing.T) {
	dir := t.TempDir()
	filePath := path.Join(dir, "CHANGELOG.md")

	// write initial content
	err := os.WriteFile(filePath, []byte("old content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	newContent := "# New Changelog"
	err = WriteChangelog(filePath, newContent)
	if err != nil {
		t.Fatalf("WriteChangelog() error = %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	want := newContent + "\n"
	if string(got) != want {
		t.Errorf("WriteChangelog() file content = %q, want %q", string(got), want)
	}
}
