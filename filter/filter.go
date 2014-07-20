// Written in 2014 by Philipp Stephani <p.stephani2@gmail.com>.
//
// To the extent possible under law, the author has dedicated all copyright and
// related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along with
// this software.  If not, see http://creativecommons.org/publicdomain/zero/1.0/.

package filter

import (
	"strings"

	"github.com/phst/dupremove/dup"
)

// RemovableFiles returns all files from the given duplicate file group that
// can be removed under the condition that at least one file should remain and
// the directories listed in keep should be left untouched.
func RemovableFiles(group dup.Group, keep []string) []dup.FileName {
	keepCandidates := []dup.FileName{}
	removeCandidates := []dup.FileName{}
	for _, file := range group {
		remove := true
		for _, dir := range keep {
			if strings.HasPrefix(string(file), dir) {
				remove = false
				break
			}
		}
		if remove {
			removeCandidates = append(removeCandidates, file)
		} else {
			keepCandidates = append(keepCandidates, file)
		}
	}
	if len(removeCandidates) == 0 || len(keepCandidates) > 0 {
		return removeCandidates
	}
	return removeCandidates[1:]
}
