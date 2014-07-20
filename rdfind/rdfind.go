package rdfind

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/phst/dupremove/dup"
)

func Run(dirs []string) ([]dup.Group, error) {
	outf, err := ioutil.TempFile("", "rdfind")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file for output: %s", err)
	}
	defer outf.Close()

	log.Printf("running rdfind for %d directories %s", len(dirs), dirs)
	cmd := exec.Command("rdfind", "-outputname", outf.Name())
	cmd.Args = append(cmd.Args, dirs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if _, err := outf.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error seeking in output file: %s", err)
	}
	log.Printf("parsing rdfind output from %s", outf.Name())
	res, err := parse(outf)
	if err != nil {
		return nil, fmt.Errorf("error parsing rdfind output from %s: %s", outf.Name(), err)
	}
	if err := os.Remove(outf.Name()); err != nil {
		log.Printf("error removing temporary file %s: %s", outf.Name(), err)
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
				res = append(res, group)
			}
			group = dup.Group{name}
		} else {
			group = append(group, name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if group != nil {
		res = append(res, group)
	}
	return res, nil
}
