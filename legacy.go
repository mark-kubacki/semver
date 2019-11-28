// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The methods contained herein are considered “legacy” and
// are marked for removal if any errors come up and nobody steps up
// to submit patches.

package semver

// NextVersions returns a list of possible next versions after t. For each of
// the three version points, pre-releases are given as options starting with
// the minimum release type (-4 <= 0), and those release types are numbered
// if numberedPre is true. Release types:
//
//   alpha: -4
//   beta:  -3
//   pre:   -2
//   rc:    -1
//   common: 0
//
// Thus, if you don't want any pre-release options, set minReleaseType to 0.
//
// Deprecated: This is a legacy method for the Caddyserver's build infrastructure.
// Do not rely on it, they are free to~ and can change it anytime.
func (t Version) NextVersions(minReleaseType int, numberedPre bool) []*Version {
	var next []*Version

	if minReleaseType < alpha || minReleaseType > common {
		return next
	}

	// if this is a pre-release, suggest next pre-releases or
	// common of same version
	for releaseType := t.version[idxReleaseType]; releaseType < common; releaseType++ {
		if releaseType == t.version[idxReleaseType] {
			if !numberedPre {
				continue
			}
			ver := t
			ver.version[idxRelease]++
			next = append(next, &ver)
		} else {
			ver := t
			ver.version[idxReleaseType] = releaseType
			if numberedPre {
				ver.version[idxRelease] = 1
			} else {
				ver.version[idxRelease] = 0
			}
			next = append(next, &ver)
		}
	}
	if t.version[idxReleaseType] < common {
		ver := t
		ver.version[idxReleaseType] = common
		ver.version[idxRelease] = 0
		next = append(next, &ver)
	}

	// if the current version is at least common release type,
	// suggest patch or revision if not one of those already
	if t.version[idxReleaseType] == common ||
		t.version[idxReleaseType] == patch {
		ver := t
		ver.version[idxReleaseType] = revision
		ver.version[idxRelease] = 1
		next = append(next, &ver)
	}
	if t.version[idxReleaseType] == common ||
		t.version[idxReleaseType] == revision {
		ver := t
		ver.version[idxReleaseType] = patch
		ver.version[idxRelease] = 1
		next = append(next, &ver)
	}

	for i := idxReleaseType - 2; 0 <= i; i-- {
		// for each version point, iterate the release types within desired bounds
		for releaseType := int32(minReleaseType); releaseType <= common; releaseType++ {
			ver := t
			ver.version[i]++
			for j := i + 1; j < len(ver.version); j++ {
				ver.version[j] = 0 // when incrementing, reset next points to 0
			}
			if i == 2 && releaseType < common {
				continue // patches seldom have pre-releases
			}
			ver.version[idxReleaseType] = releaseType
			if releaseType < common {
				if numberedPre {
					ver.version[idxRelease] = 1
				} else {
					ver.version[idxRelease] = 0
				}
			}
			next = append(next, &ver)
		}
	}

	return next
}
