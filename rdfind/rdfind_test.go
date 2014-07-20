package rdfind

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/phst/dupremove/dup"
)

func TestParse(t *testing.T) {
	input := `# Automatically generated
# duptype id depth size device inode priority name
DUPTYPE_FIRST_OCCURRENCE 1 0 4 16777218 12719320 1 /file
DUPTYPE_WITHIN_SAME_TREE -1 0 4 16777218 12719321 1 /file with spaces
DUPTYPE_WITHIN_SAME_TREE -1 0 4 16777218 12719355 1 /another file with spaces
DUPTYPE_OUTSIDE_TREE -1 0 4 16777218 12719431 2 /dir with spaces/file
# end of file
`
	want := []dup.Group{
		{"/file", "/file with spaces", "/another file with spaces", "/dir with spaces/file"},
	}
	got, err := parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatalf("parse returned error %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parse returned %#v, want %#v", got, want)
	}
}
