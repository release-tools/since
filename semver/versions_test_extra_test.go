package semver

import "testing"

func TestGetNextVersion_withVPrefix(t *testing.T) {
	got := GetNextVersion("1.2.3", true, []string{"feat: new feature"})
	want := "v1.3.0"
	if got != want {
		t.Errorf("GetNextVersion() with vPrefix = %v, want %v", got, want)
	}
}

func TestGetNextVersion_noChanges(t *testing.T) {
	got := GetNextVersion("1.2.3", false, []string{"unknown: something"})
	want := ""
	if got != want {
		t.Errorf("GetNextVersion() with no recognised changes = %v, want %v", got, want)
	}
}

func TestDetermineChangeType_allPatchTypes(t *testing.T) {
	patchTypes := []string{"build", "chore", "ci", "docs", "fix", "refactor", "security", "style", "test"}
	for _, pt := range patchTypes {
		t.Run(pt, func(t *testing.T) {
			got := DetermineChangeType([]string{pt})
			if got != ComponentPatch {
				t.Errorf("DetermineChangeType(%v) = %v, want %v", pt, got, ComponentPatch)
			}
		})
	}
}

func TestDetermineChangeType_none(t *testing.T) {
	got := DetermineChangeType([]string{"unknown"})
	if got != ComponentNone {
		t.Errorf("DetermineChangeType(unknown) = %v, want %v", got, ComponentNone)
	}
}

func TestDetermineChangeType_majorTakesPrecedence(t *testing.T) {
	got := DetermineChangeType([]string{"feat", "BREAKING CHANGE", "fix"})
	if got != ComponentMajor {
		t.Errorf("DetermineChangeType() = %v, want %v", got, ComponentMajor)
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		component Component
		want      string
	}{
		{
			name:      "bump major resets minor and patch",
			version:   "1.5.9",
			component: ComponentMajor,
			want:      "2.0.0",
		},
		{
			name:      "bump minor resets patch",
			version:   "1.5.9",
			component: ComponentMinor,
			want:      "1.6.0",
		},
		{
			name:      "bump patch only",
			version:   "1.5.9",
			component: ComponentPatch,
			want:      "1.5.10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bumpVersion(tt.version, tt.component); got != tt.want {
				t.Errorf("bumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
