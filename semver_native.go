// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package semver

//go:noescape
func compare(t, o *[14]int32) int

// Compare computes the difference between two Versions and returns its signum.
//
//   1  if a > b
//   0  if a == b
//   -1 if a < b
//
// The 'build' is not compared.
func Compare(a, b Version) int {
	return compare(&a.version, &b.version)
}

// less returns true if t is lexically smaller than o.
// As side effect, the adjacent 'build' gets compared as well.
//
//go:noescape
func less(t, o *[14]int32) bool

// Less is a convenience function for sorting.
func (t Version) Less(o Version) bool {
	return less(&t.version, &o.version)
}