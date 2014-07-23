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

	"github.com/golang/glog"

	"github.com/phst/dupremove/dup"
)

// RemovableFiles returns all files from the given duplicate file group that
// can be removed under the condition that at least one file should remain and
// the directories listed in keep should be left untouched.
func RemovableFiles(group dup.Group, keep []string) []dup.FileName {
	glog.V(2).Infof("searching duplicate group with %d files for candidates for removal", len(group))
	kept := 0
	candidates := []dup.FileName{}
	for _, file := range group {
		glog.V(3).Infof("testing whether %s can be removed", file)
		remove := true
		for _, dir := range keep {
			if strings.HasPrefix(string(file), dir) {
				glog.V(2).Infof("file %s cannot be removed because it is in precious directory %s", file, dir)
				remove = false
				break
			}
		}
		if remove {
			glog.V(3).Infof("file %s is a candidate for removal", file)
			candidates = append(candidates, file)
		} else {
			glog.V(3).Infof("file %s will be kept because it is precious", file)
			kept++
		}
	}
	if len(candidates) == 0 || kept > 0 {
		glog.V(2).Infof("all removal candidates can be removed")
		return candidates
	}
	glog.V(2).Infof("file %s will be kept because it would be the only remaining file in current group", candidates[0])
	return candidates[1:]
}
