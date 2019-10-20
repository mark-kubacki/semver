// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func makeVersionCollection(b *testing.B) ([]Version, []*Version) {
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

	return actual, unsorted
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
		_, unsorted := makeVersionCollection(nil)
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
	_, unsorted := makeVersionCollection(b)
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

func TestTwoFieldKey(t *testing.T) {
	Convey("twoFieldKey correctly derives keys", t, FailureContinues, func() {
		// {input, expected output}
		for _, testcase := range []struct {
			version     string
			keyIndex    uint8
			expectedKey []uint8
		}{
			// Examine output beyond thresholds for resolved/unresolved fields.
			{"9.16777216", 0, []uint8{(9<<4 | 15)}},
			{"9.65536", 0, []uint8{(9<<4 | 14)}},
			{"9.256", 0, []uint8{(9<<4 | 13)}},
			{"9.250", 0, []uint8{(9<<4 | 12)}},
			{"11.11", 0, []uint8{(11<<4 | 11)}},
			{"250.9", 0, []uint8{(12<<4 | 9), (12 << 4)}},
			{"256.9", 0, []uint8{(13<<4 | 9), (13 << 4)}},
			{"65536.9", 0, []uint8{(14<<4 | 9), (14 << 4)}},
			{"16777216.9", 0, []uint8{(15<<4 | 9), (15 << 4)}},
		} {
			given, _ := NewVersion([]byte(testcase.version))
			gotKey := uint8(twoFieldKey(&given.version, testcase.keyIndex))
			// The keyFn could've already collapsed lower fields below unresolved larger fields.
			So(gotKey, ShouldBeIn, testcase.expectedKey)
		}
	})
}

// By running multiple versions through key-derivation functions
// the cpu's branch predictor is utilized "realistically."
// That is, merely using one version might appear to be faster.

var tmpForTwoFieldKey = twoFieldKey(&benchV.version, 0) // To inherit its return type.

func BenchmarkTwoFieldKey(b *testing.B) {
	b.StopTimer()
	versions, _ := makeVersionCollection(b)
	versionsLen := len(versions)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		v := versions[i%versionsLen]
		tmpForTwoFieldKey |= twoFieldKey(&v.version, 0)
	}
}
