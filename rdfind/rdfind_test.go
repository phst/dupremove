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
