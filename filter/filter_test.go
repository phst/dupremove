package filter

import (
	"reflect"
	"testing"

	"github.com/phst/dupremove/dup"
)

func TestRemovableFiles(t *testing.T) {
	tests := []struct {
		input dup.Group
		want  []dup.FileName
	}{
		{
			input: dup.Group{"/a/b", "/a/c", "/b/d"},
			want:  []dup.FileName{"/b/d"},
		}, {
			input: dup.Group{},
			want:  []dup.FileName{},
		}, {
			input: dup.Group{"/a/e", "/a/f"},
			want:  []dup.FileName{},
		}, {
			input: dup.Group{"/b/g"},
			want:  []dup.FileName{},
		}, {
			input: dup.Group{"/b/h", "/b/i"},
			want:  []dup.FileName{"/b/i"},
		},
	}
	keep := []string{"/a"}
	for _, test := range tests {
		input := test.input
		want := test.want
		got := RemovableFiles(input, keep)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("RemovableFiles(%#v, %#v) returned %#v, want %#v", input, keep, got, want)
		}
	}
}
