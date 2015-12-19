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

package main

import (
	"flag"
	"os"

	"github.com/golang/glog"

	"github.com/phst/dupremove/dup"
	"github.com/phst/dupremove/filter"
	"github.com/phst/dupremove/rdfind"
)

var dryRun = flag.Bool("n", false, "dry-run mode: don't remove any files")

func main() {
	flag.Parse()
	keep := []string{}
	dirs := []string{}
	mode := ""
	for _, arg := range flag.Args() {
		if arg == "keep" || arg == "remove" {
			mode = arg
		} else {
			if mode == "" {
				glog.Fatal("command line arguments need to start with 'keep' or 'remove'")
			} else {
				dirs = append(dirs, arg)
				if mode == "keep" {
					keep = append(keep, arg)
				}
			}
		}
	}
	if len(dirs) == 0 {
		glog.Fatal("no directories specified")
	}

	groups, err := rdfind.Run(dirs)
	if err != nil {
		glog.Fatalf("error running rdfind: %s", err)
	}
	glog.Infof("found %d file groups", len(groups))

	removed := 0
	errors := 0
	for _, group := range groups {
		files := filter.RemovableFiles(group, keep)
		for _, file := range files {
			if err := remove(file); err != nil {
				glog.Errorf("could not remove file %s: %s", file, err)
				errors++
			} else {
				glog.V(1).Infof("removed file %s", file)
				removed++
			}
		}
	}
	glog.Infof("removed %d files", removed)
	if errors > 0 {
		glog.Errorf("could not remove %d files due to errors", errors)
		os.Exit(1)
	}
}

func remove(f dup.FileName) error {
	if *dryRun {
		glog.V(2).Infof("skipping file %s in dry-run mode", f)
		return nil
	}
	glog.V(2).Infof("removing file %s", f)
	return os.Remove(string(f))
}
