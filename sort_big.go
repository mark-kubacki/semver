// Copyright 2021 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build mips mips64 ppc64 s390x

package semver

import (
	"sort"
)

// Sort for bigendian is not optimized and resorts to Go's own sort.
func (p VersionPtrs) Sort() {
	sort.Sort(p)
}
