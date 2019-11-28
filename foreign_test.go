// +build ignore

package semver

import (
	"sort"
	"testing"

	blang "github.com/blang/semver"
	hashicorp "github.com/hashicorp/go-version"
)

var benchLibraryOneV, benchLibraryOneErr = hashicorp.NewVersion(strForBenchmarks)

func BenchmarkLibraryOne_NewVersion(b *testing.B) {
	var v, e = hashicorp.NewVersion(strForBenchmarks)
	lim := len(VersionsFromGentoo)

	for n := 0; n < b.N; n++ {
		v, e = hashicorp.NewVersion(string(VersionsFromGentoo[n%lim]))
	}
	benchLibraryOneV, benchLibraryOneErr = v, e
}

var benchLibraryOneR, benchLibraryOneRErr = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")

func BenchmarkLibraryOne_NewConstraint(b *testing.B) {
	var r, e = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")
	for n := 0; n < b.N; n++ {
		r, e = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")
	}
	benchLibraryOneR, benchLibraryOneRErr = r, e
}

var benchLibraryOneResult = 5

func BenchmarkLibraryOne_Compare(b *testing.B) {
	var v, _ = hashicorp.NewVersion(strForBenchmarks)
	r := benchLibraryOneV.Compare(v)
	for n := 0; n < b.N; n++ {
		r = benchLibraryOneV.Compare(v)
	}
	benchLibraryOneResult = r
}

func BenchmarkLibraryOne_SortPtr(b *testing.B) {
	b.StopTimer()
	var erroneous int
	unsorted := make([]*hashicorp.Version, len(VersionsFromGentoo))
	for n, src := range VersionsFromGentoo {
		if v, err := hashicorp.NewVersion(string(src)); err == nil {
			unsorted[n] = v
		} else {
			unsorted[n], _ = hashicorp.NewVersion(strForBenchmarks)
			erroneous += 1
		}
	}
	b.ReportMetric(float64(erroneous)/float64(len(unsorted)), "substitutes/op")
	data := make([]*hashicorp.Version, len(unsorted))

	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		sort.Sort(hashicorp.Collection(data))
		b.StopTimer()
	}
}

// Blang published their library after mine, and doing so even did imitate parts
// of my first release.
// Yet, as of writing this, their error rate is a staggering 30%.

var benchLibraryTwoV, benchLibraryTwoErr = blang.Make(strForBenchmarks)

func BenchmarkLibraryTwo_Make(b *testing.B) {
	var v, e = blang.Make(strForBenchmarks)
	lim := len(VersionsFromGentoo)

	for n := 0; n < b.N; n++ {
		v, e = blang.Make(string(VersionsFromGentoo[n%lim]))
	}
	benchLibraryTwoV, benchLibraryTwoErr = v, e
}

var benchLibraryTwoR, benchLibraryTwoRErr = blang.ParseRange(">=1.2.3 <=1.3.0")

func BenchmarkLibraryTwo_ParseRange(b *testing.B) {
	var r, e = blang.ParseRange(">=1.2.3 <=1.3.0")
	for n := 0; n < b.N; n++ {
		r, e = blang.ParseRange(">=1.2.3 <=1.3.0")
	}
	benchLibraryTwoR, benchLibraryTwoRErr = r, e
}

var benchLibraryTwoResult = 5

func BenchmarkLibraryTwo_Compare(b *testing.B) {
	var v, _ = blang.Make(strForBenchmarks)
	r := benchLibraryTwoV.Compare(v)
	for n := 0; n < b.N; n++ {
		r = benchLibraryTwoV.Compare(v)
	}
	benchLibraryTwoResult = r
}

func BenchmarkLibraryTwo_Sort(b *testing.B) {
	b.StopTimer()
	var erroneous int
	unsorted := make([]blang.Version, len(VersionsFromGentoo))
	for n, src := range VersionsFromGentoo {
		if v, err := blang.ParseTolerant(string(src)); err == nil {
			unsorted[n] = v
		} else {
			unsorted[n], _ = blang.ParseTolerant(strForBenchmarks)
			erroneous += 1
		}
	}
	b.ReportMetric(float64(erroneous)/float64(len(unsorted)), "substitutes/op")
	data := make([]blang.Version, len(unsorted))

	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		sort.Sort(blang.Versions(data))
		b.StopTimer()
	}
}
