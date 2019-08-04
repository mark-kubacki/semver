// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package semver

// less returns true if t is lexically smaller than o.
// As side effect, the adjacent 'build' gets compared as well.
//
//go:noescape
func less(t, o *[14]int32) bool

// Less is a convenience function for sorting.
func (t Version) Less(o Version) bool {
	return less(&t.version, &o.version)
}
