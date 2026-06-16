/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package changelog

import "fmt"

// NoChangesError is returned when no commits are found between the start tag
// and HEAD. ExcludedCount reports how many commits were dropped by the
// configured ignore patterns; when non-zero the user is most likely hitting
// an over-broad ignore regex.
type NoChangesError struct {
	StartTag      string
	ExcludedCount int
}

func (e *NoChangesError) Error() string {
	switch {
	case e.ExcludedCount > 0 && e.StartTag != "":
		return fmt.Sprintf(
			"no eligible commits since %s — %d commit(s) were excluded by ignore patterns in since.yaml",
			e.StartTag, e.ExcludedCount,
		)
	case e.ExcludedCount > 0:
		return fmt.Sprintf(
			"no eligible commits found — %d commit(s) were excluded by ignore patterns in since.yaml",
			e.ExcludedCount,
		)
	case e.StartTag != "":
		return fmt.Sprintf("no commits since %s", e.StartTag)
	default:
		return "no commits found"
	}
}
