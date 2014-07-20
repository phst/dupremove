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
