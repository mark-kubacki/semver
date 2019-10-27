// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package semver

// twoFieldKeyGonly is part of multikeyRadixSort.
// Please see the *_generic.go file for a detailed description.
//
//go:noescape
func twoFieldKey(v *[14]int32, fieldAdjustment uint64, keyIndex uint8) uint64

// isSorted is called by radixSort and multikeyRadixSort, and won't contain any nil.
func (p VersionPtrs) isSorted() bool {
	if len(p) < 2 {
		return true
	}

	previous := p[0]
	for _, ptr := range p {
		if compare(&previous.version, &ptr.version) > 0 {
			return false
		}
		previous = ptr
	}
	return true
}
