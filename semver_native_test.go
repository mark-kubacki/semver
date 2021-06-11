// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386
// +build !purego
// +build !go1.17

package semver

import (
	"testing"
)

func BenchmarkLess_asm(b *testing.B) {
	t := Version{}
	o := Version{}
	o.version[benchCompareIdx] = benchCompareIdx
	r := t.Less(o)

	for n := 0; n < b.N; n++ {
		r = less(&t.version, &o.version)
	}
	benchResult = r
}
