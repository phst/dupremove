package filter

import (
	"strings"

	"github.com/phst/dupremove/dup"
)

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
	} else {
		return removeCandidates[1:]
	}
}
