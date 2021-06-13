// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386
// +build !purego
// +build !go1.16

package semver

// Compare computes the difference between two Versions and returns its signum.
//
//   1  if a > b
//   0  if a == b
//   -1 if a < b
//
// The 'build' is not compared.
//go:noescape
func Compare(a, b *Version) int

// less returns true if t is lexically smaller than o.
// As side effect, the adjacent 'build' gets compared as well.
//
//go:noescape
func less(a, b *Version) bool

// Less is a convenience function for sorting.
func (t *Version) Less(o *Version) bool {
	return less(t, o)
}

// Less implements the sort.Interface.
func (p VersionPtrs) Less(i, j int) bool {
	if p[i] == nil {
		return false
	} else if p[j] == nil {
		return true
	}
	return less(p[i], p[j])
}
