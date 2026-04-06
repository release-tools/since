package stringutil

import (
	"reflect"
	"testing"
)

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		orig   []string
		search []string
		want   bool
	}{
		{
			name:   "exact match",
			orig:   []string{"foo", "bar"},
			search: []string{"foo"},
			want:   true,
		},
		{
			name:   "case insensitive match",
			orig:   []string{"Foo", "Bar"},
			search: []string{"foo"},
			want:   true,
		},
		{
			name:   "no match",
			orig:   []string{"foo", "bar"},
			search: []string{"baz"},
			want:   false,
		},
		{
			name:   "empty orig",
			orig:   []string{},
			search: []string{"foo"},
			want:   false,
		},
		{
			name:   "empty search",
			orig:   []string{"foo"},
			search: []string{},
			want:   false,
		},
		{
			name:   "multiple search terms with one match",
			orig:   []string{"foo", "bar"},
			search: []string{"baz", "bar"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsIgnoreCase(tt.orig, tt.search...); got != tt.want {
				t.Errorf("ContainsIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{
			name: "no duplicates",
			s:    []string{"foo", "bar"},
			want: []string{"foo", "bar"},
		},
		{
			name: "with duplicates",
			s:    []string{"foo", "bar", "foo"},
			want: []string{"foo", "bar"},
		},
		{
			name: "case insensitive duplicates",
			s:    []string{"Foo", "foo"},
			want: []string{"Foo"},
		},
		{
			name: "empty slice",
			s:    []string{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
