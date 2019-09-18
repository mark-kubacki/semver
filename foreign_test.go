// +build ignore

package semver

import (
	"testing"

	blang "github.com/blang/semver"
	hashicorp "github.com/hashicorp/go-version"
)

var benchHashicorpV, benchHashicorpErr = hashicorp.NewVersion(strForBenchmarks)

func BenchmarkHashicorpNewVersion(b *testing.B) {
	var v, e = hashicorp.NewVersion(strForBenchmarks)
	lim := len(VersionsFromGentoo)

	for n := 0; n < b.N; n++ {
		v, e = hashicorp.NewVersion(string(VersionsFromGentoo[n%lim]))
	}
	benchHashicorpV, benchHashicorpErr = v, e
}

var benchHashicorpR, benchHashicorpRErr = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")

func BenchmarkHashicorpNewConstraint(b *testing.B) {
	var r, e = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")
	for n := 0; n < b.N; n++ {
		r, e = hashicorp.NewConstraint(">=1.2.3, <=1.3.0")
	}
	benchHashicorpR, benchHashicorpRErr = r, e
}

var benchHashicorpResult = 5

func BenchmarkHashicorpCompare(b *testing.B) {
	var v, _ = hashicorp.NewVersion(strForBenchmarks)
	r := benchHashicorpV.Compare(v)
	for n := 0; n < b.N; n++ {
		r = benchHashicorpV.Compare(v)
	}
	benchHashicorpResult = r
}

var benchBlangV, benchBlangErr = blang.Make(strForBenchmarks)

func BenchmarkBlangMake(b *testing.B) {
	var v, e = blang.Make(strForBenchmarks)
	lim := len(VersionsFromGentoo)

	for n := 0; n < b.N; n++ {
		v, e = blang.Make(string(VersionsFromGentoo[n%lim]))
	}
	benchBlangV, benchBlangErr = v, e
}

var benchBlangR, benchBlangRErr = blang.ParseRange(">=1.2.3 <=1.3.0")

func BenchmarkBlangParseRange(b *testing.B) {
	var r, e = blang.ParseRange(">=1.2.3 <=1.3.0")
	for n := 0; n < b.N; n++ {
		r, e = blang.ParseRange(">=1.2.3 <=1.3.0")
	}
	benchBlangR, benchBlangRErr = r, e
}

var benchBlangResult = 5

func BenchmarkBlangCompare(b *testing.B) {
	var v, _ = blang.Make(strForBenchmarks)
	r := benchBlangV.Compare(v)
	for n := 0; n < b.N; n++ {
		r = benchBlangV.Compare(v)
	}
	benchBlangResult = r
}
