// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"sort"
	"sync"
)

const (
	// Less than this many elements and a "residual" sort gets called,
	// which usually is sort.Sort or sort.Slice.
	// To figure out this particular value I've run benchmarks,
	// but got a range of close results; you could go as low
	// as 64 or 32 on some architectures.
	thresholdForResidualSort = 128
)

// Radix sort—and variants will be used below—needs some scratch space,
// which this pool will provide.
//
// Don't rely on the initial size for new arrays. Expand the capacity if need be.
var versionPointerBuffer = sync.Pool{
	New: func() interface{} {
		b := make([]*Version, 40*1024)
		return &b
	},
}

// VersionPtrs represents an array with elements derived from but smaller than Versions.
// Go through this to sort large collections of Versions to minimize bytes written to memory.
type VersionPtrs []*Version

// VersionPtrs.Less calls specialized functions.
// Find it in files *_native.go and *_generic.go.
// As of Go 1.13 inlining didn't work across two levels.

// Len implements the sort.Interface.
func (p VersionPtrs) Len() int {
	return len(p)
}

// Swap implements the sort.Interface.
func (p VersionPtrs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Sort reorders the pointers so that the Versions appear in ascending order.
//
// For that it will use optimized algorithms, usually less time-complex than
// the generic ones found in package 'Sort'.
// Specifically, variants of radix sort expected to run in O(n);
// worst case in O(n*log(n)) —which is unlikely— deferring to 'sort.*'
// on degenerated collections.
//
// Allocates a copy of VersionPtrs.
func (p VersionPtrs) Sort() {
	if len(p) < thresholdForResidualSort {
		sort.Sort(p)
		return
	}

	buf := versionPointerBuffer.Get().(*[]*Version)
	tmp := *buf
	p.multikeyRadixSort(tmp, 0)
	for i := range tmp {
		tmp[i] = nil
	}
	versionPointerBuffer.Put(buf)
}

// magnitudeAwareKey is part of twoFieldKey below, and returns small figures verbatim,
// else signals magnitudes in a way susceptible to sorting.
// Parameter 'x' must be non-negative.
func magnitudeAwareKey(x int32) uint8 {
	if x <= 11 {
		return uint8(x)
	}
	// For all larger numbers, store the number of bytes +11.
	if x <= 0xffff {
		if x <= 0xff {
			return 12
		}
		return 13
	}
	if x <= 0xffffff {
		return 14
	}
	return 15
}

// twoFieldKey is part of multikeyRadixSort and derives a key from two fields in 'v'.
// The order established by the keys is ascending but not total:
// fields with great values map to a low-resolution key.
// Fields must be non-negative.
func twoFieldKey(v *[14]int32, keyIndex uint8) uint8 {
	n1 := magnitudeAwareKey(v[keyIndex]) << 4
	if n1 >= (12 << 4) {
		return n1
	}
	return (n1 | magnitudeAwareKey(v[keyIndex+1]))
}

// multikeyRadixSort exploits the typical distribution of Version values
// to use  two keys at once  in a radix-sort run.
func (p VersionPtrs) multikeyRadixSort(tmp []*Version, keyIndex uint8) {
	keyIndex &= 0x07 // Signals the compiler we expect a limited set of values for this.

	// Collate the histogram.
	var offset [256]int
	for _, v := range p {
		if v == nil {
			continue
		}
		k := twoFieldKey(&v.version, uint8(keyIndex))
		offset[k]++
	}
	watermark := offset[0] - offset[0] // 'watermark' will finally be the total tally.
	for i, count := range offset {
		offset[i] = watermark
		watermark += count
	}

	// Setup an unordered copy.
	// The allocated space will subsequently be recycled as scratch space.
	if len(tmp) >= len(p) {
		tmp = tmp[:len(p)]
		copy(tmp, p)
	} else {
		tmp = append(tmp[:0], p...)
	}
	for i := watermark; i < len(p); i++ {
		p[i] = nil // Fill the tail end with the 'nil' we'll be skipping.
	}

	// Order from 'tmp' into 'p'.
	for _, v := range tmp {
		if v == nil {
			continue
		}
		k := twoFieldKey(&v.version, uint8(keyIndex))
		p[offset[k]] = v
		offset[k]++
	}

	p.multikeyRadixSortDescent(tmp, keyIndex, offset)
}

// multikeyRadixSortDescent is multikeyRadixSort's outsourced descent- and recurse steps.
// Split for easier profiling.
func (p VersionPtrs) multikeyRadixSortDescent(tmp []*Version, keyIndex uint8, offset [256]int) {
	// Any tailing nil are beyond offsets, henceforth no longer considered.
	watermark := offset[0] - offset[0]
	for k, ceiling := range offset {
		subsliceLen := ceiling - watermark // aka "stride"
		if subsliceLen < 2 {
			watermark = ceiling
			continue
		}

		subslice := p[watermark:ceiling]
		watermark = ceiling
		if subsliceLen < thresholdForResidualSort || keyIndex >= 3 {
			sort.Sort(subslice)
			continue
		}

		switch k := uint8(k); {
		case (k & 0x0f) >= 12: // This key is in order, the next is not: descent.
			maxBits := ((k & 0x0f) - 11) * 8 // 12 → 1 → 8
			subslice.radixSort(tmp, keyIndex+1, maxBits)
		case k >= (12 << 4): // Unsorted trailer with values that keyFn did not resolve.
			maxBits := ((k >> 4) - 11) * 8
			subslice.radixSort(tmp, keyIndex, maxBits)
		case keyIndex >= 1: // Guards the below call to multikeyRadixSort.
			sort.Sort(subslice)
		default:
			subslice.multikeyRadixSort(tmp, keyIndex+2)
		}
	}
}

// radixSort sorts on the one field indicated by keyIndex.
// maxBits really denominates the octets (bytes) to consider, and any excess MSB are assumed to be zero.
//
// Tailing nil are expected to have been stripped.
func (p VersionPtrs) radixSort(tmp []*Version, keyIndex, maxBits uint8) {
	if keyIndex > 3 {
		panic("keyIndex out of bounds")
	}
	from, to := p, tmp[:len(p)] // Have the compiler check this once.
	var offset [256]uint

	for fromBits := maxBits - maxBits; fromBits < maxBits; fromBits += 8 {
		// Building the histogram again.
		// Although this can be done for all bytes in one run,
		// which would need a [1024], I found it's slower in Golang.
		for i := range offset {
			offset[i] = 0
		}
		for _, v := range from {
			if v == nil {
				continue
			}
			k := uint8(v.version[keyIndex] >> fromBits)
			offset[k]++
		}
		watermark := offset[0] - offset[0]
		for i, count := range offset {
			offset[i] = watermark
			watermark += count
		}

		// Now comes the ordering, which is stable of course.
		for _, v := range from {
			if v == nil {
				continue
			}
			k := uint8(v.version[keyIndex] >> fromBits)
			to[offset[k]] = v
			offset[k]++
		}
		to, from = from, to // Prepare the next run.
	}
	if maxBits%16 != 0 {
		copy(to, from)
	}

	p.radixSortDescent(tmp, keyIndex)
}

// radixSortDescent is radixSort's outsourced descent- and recurse steps.
// Split for easier profiling.
func (p VersionPtrs) radixSortDescent(tmp []*Version, keyIndex uint8) {
	// The descent. multikeyRadixSort has only one run, hence
	// is able to read strides from its histogram ("offset[]").
	// As classical radix sort cannot (even if optimized to one run for the histogram),
	// the collection needs to be visited once more.
	startIdx := 0
	lastValue := p[0].version[keyIndex]
	for i, v := range p {
		value := v.version[keyIndex]
		if lastValue == value { // Accumulate spans of the same value.
			continue
		}
		if i-startIdx < 2 {
			startIdx, lastValue = i, value
			continue
		}

		subslice := p[startIdx:i]
		if i-startIdx < thresholdForResidualSort || keyIndex >= 2 {
			sort.Sort(subslice)
		} else {
			subslice.multikeyRadixSort(tmp, keyIndex+1)
		}
		startIdx, lastValue = i, value
	}
	// Capture trailer of same values (such as 250.100, 250.0).
	if residualLength := len(p) - startIdx; residualLength > 1 {
		subslice := p[startIdx:]
		if residualLength < thresholdForResidualSort || keyIndex >= 2 {
			sort.Sort(subslice)
		} else {
			subslice.multikeyRadixSort(tmp, keyIndex+1)
		}
	}
}
