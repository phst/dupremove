// Written in 2014 by Philipp Stephani <p.stephani2@gmail.com>.
//
// To the extent possible under law, the author has dedicated all copyright and
// related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along with
// this software.  If not, see http://creativecommons.org/publicdomain/zero/1.0/.

package dup

type (
	// Group denotes a set of files that are duplicates of each other.
	Group []FileName
	// FileName represents a file name of a duplicate file.
	FileName string
)
