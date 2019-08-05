// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64

package semver

// Compare computes the difference between two Versions and returns its signum.
//
//   1  if a > b
//   0  if a == b
//   -1 if a < b
//
// The 'build' is not compared.
func Compare(a, b Version) int {
	for i := 0; i < len(a.version); i++ {
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
		if t.version[i] < o.version[i] {
			return true
		}
		return false
	}
	return t.build < o.build
}
