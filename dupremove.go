// Written in 2014 by Philipp Stephani <p.stephani2@gmail.com>.
//
// To the extent possible under law, the author has dedicated all copyright and
// related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along with
// this software.  If not, see http://creativecommons.org/publicdomain/zero/1.0/.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/phst/dupremove/dup"
	"github.com/phst/dupremove/filter"
	"github.com/phst/dupremove/rdfind"
)

var dryRun = flag.Bool("n", false, "dry-run mode: don't remove any files")

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	keep := []string{}
	dirs := []string{}
	mode := ""
	for _, arg := range flag.Args() {
		if arg == "keep" || arg == "remove" {
			mode = arg
		} else {
			if mode == "" {
				log.Fatalf("command line arguments need to start with 'keep' or 'remove'")
			} else {
				dirs = append(dirs, arg)
				if mode == "keep" {
					keep = append(keep, arg)
				}
			}
		}
	}
	if len(dirs) == 0 {
		log.Fatalf("no directories specified")
	}

	groups, err := rdfind.Run(dirs)
	if err != nil {
		log.Fatalf("error running rdfind: %s", err)
	}
	log.Printf("found %d file groups", len(groups))

	removed := 0
	for _, group := range groups {
		files := filter.RemovableFiles(group, keep)
		for _, file := range files {
			if err := remove(file); err != nil {
				log.Printf("could not remove file %s: %s", file, err)
			} else {
				log.Printf("removed file %s", file)
				removed++
			}
		}
	}
	log.Printf("removed %d files", removed)
}

func remove(f dup.FileName) error {
	if *dryRun {
		return nil
	} else {
		return os.Remove(string(f))
	}
}
