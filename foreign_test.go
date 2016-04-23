package semver

import (
	"testing"

	blang "github.com/blang/semver"
	hashicorp "github.com/hashicorp/go-version"
)

var benchHashicorpV, benchHashicorpErr = hashicorp.NewVersion("1.2.3-beta")

func BenchmarkHashicorpNewVersion(b *testing.B) {
	var v, e = hashicorp.NewVersion("1.2.3-beta")
	for n := 0; n < b.N; n++ {
		v, e = hashicorp.NewVersion("1.2.3-beta")
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

var benchBlangV, benchBlangErr = blang.Make("1.2.3-beta")

func BenchmarkBlangMake(b *testing.B) {
	var v, e = blang.Make("1.2.3-beta")
	for n := 0; n < b.N; n++ {
		v, e = blang.Make("1.2.3-beta")
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
