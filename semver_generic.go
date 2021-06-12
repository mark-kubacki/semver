// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.17 purego !amd64,!386

package semver

// Compare computes the difference between two Versions and returns its signum.
//
//   1  if a > b
//   0  if a == b
//   -1 if a < b
//
// The 'build' is not compared.
func Compare(a, b *Version) int {
	for i := 0; i < len(a.version); i++ {
		if a.version[i] == b.version[i] {
			continue
		}
		x := a.version[i] - b.version[i]
		return int((x >> 31) - (-x >> 31))
	}
	return 0
}

// compare works like the exported Compare,
// only that it allows to skip fields for performance reasons.
func compare(a, b *Version, skipFields uint) int {
	for i := int(skipFields); i < len(a.version); i++ {
		if a.version[i] == b.version[i] {
			continue
		}
		x := a.version[i] - b.version[i]
		return int((x >> 31) - (-x >> 31))
	}
	return 0
}

// Less is a convenience function for sorting.
func (t Version) Less(o Version) bool {
	for i := 0; i < len(t.version); i++ {
		if t.version[i] == o.version[i] {
			continue
		}
		return t.version[i] < o.version[i]
	}
	return t.build < o.build
}

// Less implements the sort.Interface.
func (p VersionPtrs) Less(i, j int) bool {
	if p[i] == nil {
		return false
	} else if p[j] == nil {
		return true
	}
	return p[i].Less(*p[j])
}
