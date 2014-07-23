// Written in 2014 by Philipp Stephani <p.stephani2@gmail.com>.
//
// To the extent possible under law, the author has dedicated all copyright and
// related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along with
// this software.  If not, see http://creativecommons.org/publicdomain/zero/1.0/.

package rdfind

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"

	"github.com/phst/dupremove/dup"
)

// Run executes the rdfind program with the given list of diretories, parses
// its output, and returns a list of groups of duplicate files.
func Run(dirs []string) ([]dup.Group, error) {
	outf, err := ioutil.TempFile("", "rdfind")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file for output: %s", err)
	}
	defer outf.Close()

	glog.Infof("running rdfind for %d directories %s", len(dirs), dirs)
	cmd := exec.Command("rdfind", "-outputname", outf.Name())
	cmd.Args = append(cmd.Args, dirs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	glog.V(1).Infof("running rdfind with arguments %v", cmd.Args)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if _, err := outf.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error seeking in output file: %s", err)
	}
	glog.Infof("parsing rdfind output from %s", outf.Name())
	res, err := parse(outf)
	if err != nil {
		return nil, fmt.Errorf("error parsing rdfind output from %s: %s", outf.Name(), err)
	}
	glog.V(1).Infof("removing temporary file %s", outf.Name())
	if err := os.Remove(outf.Name()); err != nil {
		glog.Warningf("error removing temporary file %s: %s", outf.Name(), err)
	}
	return res, nil
}

func parse(r io.Reader) ([]dup.Group, error) {
	const (
		first       = "DUPTYPE_FIRST_OCCURRENCE"
		withinTree  = "DUPTYPE_WITHIN_SAME_TREE"
		outsideTree = "DUPTYPE_OUTSIDE_TREE"
	)
	duptypes := map[string]bool{first: true, withinTree: true, outsideTree: true}
	res := []dup.Group{}
	var group dup.Group
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		s := scanner.Text()
		glog.V(3).Infof("parsing line %d: %s", i+1, s)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		fields := strings.SplitN(s, " ", 8)
		if len(fields) != 8 {
			return nil, fmt.Errorf("line %d: expected eight fields, got %d", i+1, len(fields))
		}
		duptype := fields[0]
		if !duptypes[duptype] {
			return nil, fmt.Errorf("line %d: unknown duptype %s", i+1, duptype)
		}
		name := dup.FileName(fields[7])
		if duptype == first {
			if group != nil {
				glog.V(3).Infof("finishing group with %d files", len(group))
				res = append(res, group)
			}
			glog.V(3).Info("starting new group")
			group = dup.Group{name}
		} else {
			glog.V(4).Infof("appending file %s to current group", name)
			group = append(group, name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if group != nil {
		glog.V(3).Infof("finishing last group with %d files", len(group))
		res = append(res, group)
	}
	return res, nil
}
