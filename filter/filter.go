// Copyright 2014 Philipp Stephani <phst@google.com>
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.  You may obtain a copy
// of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

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
	if len(candidates) == 0 {
		glog.V(2).Info("no candidates for removal")
		return candidates
	}
	if kept > 0 {
		glog.V(2).Infof("all removal candidates can be removed because %d files will be held back", kept)
		return candidates
	}
	glog.V(2).Infof("file %s will be kept because it would be the only remaining file in current group", candidates[0])
	return candidates[1:]
}
