package semver

import "testing"

func TestGetNextVersion(t *testing.T) {
	type args struct {
		currentVersion string
		vPrefix        bool
		commits        []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "major",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"feat!: major change",
					"feat: new feature",
					"fix: all bugs fixed",
				},
			},
			want: "2.0.0",
		},
		{
			name: "minor",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"feat: new feature",
					"fix: all bugs fixed",
				},
			},
			want: "1.3.0",
		},
		{
			name: "patch",
			args: args{
				currentVersion: "1.2.3",
				vPrefix:        false,
				commits: []string{
					"fix: all bugs fixed",
				},
			},
			want: "1.2.4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNextVersion(tt.args.currentVersion, tt.args.vPrefix, tt.args.commits); got != tt.want {
				t.Errorf("GetNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineChangeType(t *testing.T) {
	type args struct {
		types []string
	}
	tests := []struct {
		name string
		args args
		want Component
	}{
		{
			name: "major",
			args: args{
				types: []string{
					"BREAKING CHANGE",
				},
			},
			want: ComponentMajor,
		},
		{
			name: "minor",
			args: args{
				types: []string{
					"feat",
				},
			},
			want: ComponentMinor,
		},
		{
			name: "patch",
			args: args{
				types: []string{
					"fix",
				},
			},
			want: ComponentPatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetermineChangeType(tt.args.types); got != tt.want {
				t.Errorf("DetermineChangeType() = %v, want %v", got, tt.want)
			}
		})
	}
}
