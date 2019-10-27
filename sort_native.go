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
