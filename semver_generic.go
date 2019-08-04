// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64

package semver

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
