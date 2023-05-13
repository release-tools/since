package stringutil

import "strings"

// ContainsIgnoreCase returns true if the orig slice contains any of the search strings,
// compared in a case-insensitive manner.
func ContainsIgnoreCase(orig []string, search ...string) bool {
	for _, o := range orig {
		for _, s := range search {
			if strings.ToLower(o) == strings.ToLower(s) {
				return true
			}
		}
	}
	return false
}

// Unique returns a slice of unique strings.
func Unique(s []string) []string {
	var unique []string
	for _, message := range s {
		if !ContainsIgnoreCase(unique, message) {
			unique = append(unique, message)
		}
	}
	return unique
}
