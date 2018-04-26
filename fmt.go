// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"fmt"
)

// String returns the string representation of t.
func (t *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", t.version[0], t.version[1], t.version[2])
	if t.version[idxReleaseType] != common {
		s += fmt.Sprintf("-%s", releaseDesc[int(t.version[idxReleaseType])])
		if t.version[idxRelease] > 0 {
			s += fmt.Sprintf(".%d", t.version[idxRelease])
		}
	}
	if t.build != 0 {
		s += fmt.Sprintf("+build%d", t.build)
	}
	return s
}
