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
