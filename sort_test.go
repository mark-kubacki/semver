// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sort"
)

func makeVersionPtrs(b *testing.B) []*Version {
	unsorted := make([]*Version, len(VersionsFromGentoo))

	var erroneous int
	actual := make([]Version, len(VersionsFromGentoo))
	for n, src := range VersionsFromGentoo {
		if err := actual[n].UnmarshalText(src); err != nil {
			actual[n].UnmarshalText(verForBenchmarks)
			erroneous += 1
		}
		unsorted[n] = &actual[n]
	}
	if b != nil {
		b.ReportMetric(float64(erroneous)/float64(len(unsorted)), "substitutes/op")
	}

	return unsorted
}

func containsAll(reference, dubious []*Version) bool {
examine_reference:
	for _, ptr := range reference {
		if ptr == nil {
			continue
		}
		for _, other := range dubious {
			if ptr == other {
				continue examine_reference
			}
		}
		return false
	}
	return true
}

func TestSortPtr(t *testing.T) {
	Convey("VersionPtrs sorting", t, func() {
		unsorted := makeVersionPtrs(nil)
		data := make([]*Version, len(unsorted)+1)
		// Insert a 'nil' as the item before the last.
		data[len(data)-2], data[len(data)-1] = nil, data[len(data)-2]
		copy(data, unsorted)

		x := VersionPtrs(data)
		lessFn := x.Less
		x.Sort()

		Convey("establishes an ascending order", func() {
			isSorted := sort.SliceIsSorted(data, lessFn)
			if isSorted {
				So(isSorted, ShouldBeTrue)
				return
			}

			precedingVersion := x[0] // Conveniently our test set does not start with a nil.
			for i := range x {
				if x[i] == nil {
					precedingVersion = nil
					continue
				}
				if precedingVersion == nil { // Case [nil, proper].
					t.Error("nil not contiguous at the end, got:", i, *x[i])
					break
				}
				if Compare(*precedingVersion, *x[i]) >= 1 {
					t.Error("Wrong order between:", i, *precedingVersion, *x[i])
					break
				}
				precedingVersion = x[i]
			}
			So(isSorted, ShouldBeTrue)
		})

		Convey("does not lose elements", func() {
			So(containsAll(unsorted, data), ShouldBeTrue)
		})
	})
}

func Benchmark_SortPtr(b *testing.B) {
	b.StopTimer()
	unsorted := makeVersionPtrs(b)
	data := make([]*Version, len(unsorted))
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		x := VersionPtrs(data)
		x.Sort()
		b.StopTimer()
		if !sort.SliceIsSorted(x, x.Less) {
			b.Skip("Resulting slice is not in order.")
			break
		}
	}
}
