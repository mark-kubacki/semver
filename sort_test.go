// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"testing"

	"sort"
)

func Benchmark_SortPtr(b *testing.B) {
	b.StopTimer()
	unsorted := make([]*Version, len(VersionsFromGentoo))
	{
		var erroneous int
		actual := make([]Version, len(VersionsFromGentoo))
		for n, src := range VersionsFromGentoo {
			if err := actual[n].UnmarshalText(src); err != nil {
				actual[n].UnmarshalText(verForBenchmarks)
				erroneous += 1
			}
			unsorted[n] = &actual[n]
		}
		b.ReportMetric(float64(erroneous)/float64(len(unsorted)), "substitutes/op")
	}
	data := make([]*Version, len(unsorted))

	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		lessFn := VersionPtrs(data).GetLessFunc()
		sort.Slice(data, lessFn)
		b.StopTimer()
	}
}
