// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewVersion(t *testing.T) {
	Convey("NewVersion works with…", t, FailureContinues, func() {
		Convey("1.23.8", func() {
			refVer, err := NewVersion([]byte("1.23.8"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, common, 0, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("v1.23.8", func() {
			refVer, err := NewVersion([]byte("v1.23.8"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, common, 0, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("1.23.8-alpha", func() {
			refVer, err := NewVersion([]byte("1.23.8-alpha"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, alpha, 0, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("1.23.8-alpha.6.7", func() {
			refVer, err := NewVersion([]byte("1.23.8-alpha.6.7"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, alpha, 6, 7, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("1.23.8-p.3", func() {
			refVer, err := NewVersion([]byte("1.23.8-p.3"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, patch, 3, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("1.23.8-p3", func() {
			refVer, err := NewVersion([]byte("1.23.8-p3"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, patch, 3, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("1.23.8-3", func() {
			refVer, err := NewVersion([]byte("1.23.8-3"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{1, 23, 8, 0, common, 3, 0, 0, 0, common, 0, 0, 0, 0})
		})

		Convey("0-0-0.0.0.4", func() {
			refVer, err := NewVersion([]byte("0-0-0.0.0.4"))
			So(err, ShouldBeNil)
			So(refVer.version, ShouldResemble, [...]int32{0, 0, 0, 0, common, 0, 0, 0, 0, common, 0, 0, 0, 4})
		})
	})
}

func TestVersion(t *testing.T) {
	Convey("Version 1.3.8 should be part of Version…", t, FailureContinues, func() {
		v := []int32{1, 3, 8, 0}

		Convey("1.3.8", func() {
			refVer, err := NewVersion([]byte("1.3.8"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

		Convey("1.3.8+build20140722", func() {
			refVer, err := NewVersion([]byte("1.3.8+build20140722"))
			So(refVer.version[:4], ShouldResemble, v)
			So(refVer.build, ShouldEqual, 20140722)
			So(err, ShouldBeNil)
		})

		Convey("1.3.8+build2014", func() {
			refVer, err := NewVersion([]byte("1.3.8+build2014"))
			So(refVer.version[:4], ShouldResemble, v)
			So(refVer.build, ShouldEqual, 2014)
			So(err, ShouldBeNil)
		})

		Convey("1.3.8-alpha", func() {
			refVer, err := NewVersion([]byte("1.3.8-alpha"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

		Convey("1.3.8-beta", func() {
			refVer, err := NewVersion([]byte("1.3.8-beta"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

		Convey("1.3.8-pre", func() {
			refVer, err := NewVersion([]byte("1.3.8-pre"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

		Convey("1.3.8-r3", func() {
			refVer, err := NewVersion([]byte("1.3.8-r3"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

		Convey("1.3.8-3", func() {
			refVer, err := NewVersion([]byte("1.3.8-3"))
			So(err, ShouldBeNil)
			So(refVer.version[:4], ShouldResemble, v)
		})

	})

	Convey("Working order between Versions", t, func() {

		Convey("equality", func() {
			v1, _ := NewVersion([]byte("1.3.8"))
			v2, _ := NewVersion([]byte("1.3.8"))
			So(v1, ShouldResemble, v2)
		})

		Convey("between different release types", func() {
			Convey("1.0.0 < 2.0.0", func() {
				v1, _ := NewVersion([]byte("1.0.0"))
				v2, _ := NewVersion([]byte("2.0.0"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v2.Less(v1), ShouldBeFalse)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("2.2.1 < 2.4.0-beta", func() {
				v1, _ := NewVersion([]byte("2.2.1"))
				v2, _ := NewVersion([]byte("2.4.0-beta"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v2.Less(v1), ShouldBeFalse)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0 < 1.0.0-p", func() {
				v1, _ := NewVersion([]byte("1.0.0"))
				v2, _ := NewVersion([]byte("1.0.0-p"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v2.Less(v1), ShouldBeFalse)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0-rc < 1.0.0", func() {
				v1, _ := NewVersion([]byte("1.0.0-rc"))
				v2, _ := NewVersion([]byte("1.0.0"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0-pre < 1.0.0-rc", func() {
				v1, _ := NewVersion([]byte("1.0.0-pre"))
				v2, _ := NewVersion([]byte("1.0.0-rc"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0-beta < 1.0.0-pre", func() {
				v1, _ := NewVersion([]byte("1.0.0-beta"))
				v2, _ := NewVersion([]byte("1.0.0-pre"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0-alpha < 1.0.0-beta", func() {
				v1, _ := NewVersion([]byte("1.0.0-alpha"))
				v2, _ := NewVersion([]byte("1.0.0-beta"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})
		})

		Convey("between same release types", func() {
			Convey("1.0.0-p0 < 1.0.0-p1", func() {
				v1, _ := NewVersion([]byte("1.0.0-p0"))
				v2, _ := NewVersion([]byte("1.0.0-p1"))

				So(v1.version, ShouldResemble, [...]int32{1, 0, 0, 0, patch, 0, 0, 0, 0, common, 0, 0, 0, 0})
				So(v2.version, ShouldResemble, [...]int32{1, 0, 0, 0, patch, 1, 0, 0, 0, common, 0, 0, 0, 0})

				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})
		})

		Convey("with release type specifier", func() {
			Convey("1.0.0-rc4-alpha1 < 1.0.0-rc4", func() {
				v1, _ := NewVersion([]byte("1.0.0-rc4-alpha1"))
				v2, _ := NewVersion([]byte("1.0.0-rc4"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})
		})

		Convey("with builds", func() {
			Convey("1.0.0+build1 < 1.0.0+build2", func() {
				v1, _ := NewVersion([]byte("1.0.0+build1"))
				v2, _ := NewVersion([]byte("1.0.0+build2"))
				So(v1.Less(v2), ShouldBeTrue)
				So(v1, ShouldNotResemble, v2)
			})

			Convey("1.0.0_pre20140722+build14 < 1.0.0_pre20140722+build15", func() {
				v1, _ := NewVersion([]byte("1.0.0_pre20140722+build14"))
				v2, _ := NewVersion([]byte("1.0.0_pre20140722+build15"))
				So(v1, ShouldNotResemble, v2)
				So(v1.Less(v2), ShouldBeTrue)
			})
		})

	})

	// see http://devmanual.gentoo.org/ebuild-writing/file-format/
	Convey("Gentoo's example of order works.", t, func() {
		v1, _ := NewVersion([]byte("1.0.0_alpha_pre"))
		v2, _ := NewVersion([]byte("1.0.0_alpha_rc1"))
		v3, _ := NewVersion([]byte("1.0.0_beta_pre"))
		v4, _ := NewVersion([]byte("1.0.0_beta_p1"))
		So(v1.version, ShouldResemble, [...]int32{1, 0, 0, 0, alpha, 0, 0, 0, 0, pre, 0, 0, 0, 0})
		So(v2.version, ShouldResemble, [...]int32{1, 0, 0, 0, alpha, 0, 0, 0, 0, rc, 1, 0, 0, 0})
		So(v3.version, ShouldResemble, [...]int32{1, 0, 0, 0, beta, 0, 0, 0, 0, pre, 0, 0, 0, 0})
		So(v4.version, ShouldResemble, [...]int32{1, 0, 0, 0, beta, 0, 0, 0, 0, patch, 1, 0, 0, 0})

		So(v1, ShouldNotResemble, v2)
		So(v2, ShouldNotResemble, v3)
		So(v3, ShouldNotResemble, v4)
		So(v1.Less(v2), ShouldBeTrue)
		So(v2.Less(v3), ShouldBeTrue)
		So(v3.Less(v4), ShouldBeTrue)
	})

	Convey("Reject invalid Versions.", t, func() {
		Convey("with surplus digits", func() {
			_, err := NewVersion([]byte("1.0.0.0.4"))
			So(err, ShouldNotBeNil)
		})

		Convey("with surplus dots", func() {
			_, err := NewVersion([]byte("1..8"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.8.rc2"))
			So(err, ShouldNotBeNil)
		})

		Convey("with unknown tags", func() {
			_, err := NewVersion([]byte("1.8-gazilla"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.8-+build4"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.8-a"))
			So(err, ShouldNotBeNil)
		})

		Convey("with fringe builds", func() {
			_, err := NewVersion([]byte("10.0.17763.253+build19H3"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("10.0.17763.253+19H3"))
			So(err, ShouldNotBeNil)
			e := err.(InvalidStringValue)
			So(e.IsInvalid(), ShouldBeTrue)
		})

		Convey("with excessive tags", func() {
			_, err := NewVersion([]byte("1.8-alpha-beta-rc"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.8-alpha-beta3rc"))
			So(err, ShouldNotBeNil)
		})

		Convey("with trailing dashes", func() {
			_, err := NewVersion([]byte("5678.9-"))
			So(err, ShouldNotBeNil)
		})

		Convey("with too long parts", func() {
			_, err := NewVersion([]byte("100000000000007000000000000000070000000000000.0.0"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.0.0_alpha444444444444444444444444444444444444444"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.0.0_alpha-rc444444444444444444444444444444444444"))
			So(err, ShouldNotBeNil)
			_, err = NewVersion([]byte("1.0.0_alpha-rc1+build44444444444444444444444444444"))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestNextVersions(t *testing.T) {
	toStr := func(list []*Version) []string {
		ss := make([]string, len(list))
		for i := range list {
			ss[i] = fmt.Sprintf("%v", list[i])
		}
		return ss
	}

	Convey("NextVersions works with…", t, func() {

		ver := "1.0.0"
		Convey(ver, func() {
			ver, err := NewVersion([]byte(ver))
			So(err, ShouldBeNil)

			Convey("Without pre-releases", func() {
				next := toStr(ver.NextVersions(0, false))
				So(next, ShouldResemble, []string{
					"1.0.0-r1",
					"1.0.0-p1",
					"1.0.1",
					"1.1.0",
					"2.0.0",
				})
			})

			Convey("With some pre-releases", func() {
				next := toStr(ver.NextVersions(-2, false))
				So(next, ShouldResemble, []string{
					"1.0.0-r1",
					"1.0.0-p1",
					"1.0.1",
					"1.1.0-pre",
					"1.1.0-rc",
					"1.1.0",
					"2.0.0-pre",
					"2.0.0-rc",
					"2.0.0",
				})
			})

			Convey("With some pre-releases and numbers", func() {
				next := toStr(ver.NextVersions(-2, true))
				So(next, ShouldResemble, []string{
					"1.0.0-r1",
					"1.0.0-p1",
					"1.0.1",
					"1.1.0-pre1",
					"1.1.0-rc1",
					"1.1.0",
					"2.0.0-pre1",
					"2.0.0-rc1",
					"2.0.0",
				})
			})

			Convey("With all pre-releases", func() {
				next := toStr(ver.NextVersions(-4, false))
				So(next, ShouldResemble, []string{
					"1.0.0-r1",
					"1.0.0-p1",
					"1.0.1",
					"1.1.0-alpha",
					"1.1.0-beta",
					"1.1.0-pre",
					"1.1.0-rc",
					"1.1.0",
					"2.0.0-alpha",
					"2.0.0-beta",
					"2.0.0-pre",
					"2.0.0-rc",
					"2.0.0",
				})
			})
		})

		ver = "1.2.3"
		Convey(ver, func() {
			ver, err := NewVersion([]byte(ver))
			So(err, ShouldBeNil)

			Convey("Without pre-releases", func() {
				next := toStr(ver.NextVersions(0, false))
				So(next, ShouldResemble, []string{
					"1.2.3-r1",
					"1.2.3-p1",
					"1.2.4",
					"1.3.0",
					"2.0.0",
				})
			})

			Convey("With some pre-releases", func() {
				next := toStr(ver.NextVersions(-2, false))
				So(next, ShouldResemble, []string{
					"1.2.3-r1",
					"1.2.3-p1",
					"1.2.4",
					"1.3.0-pre",
					"1.3.0-rc",
					"1.3.0",
					"2.0.0-pre",
					"2.0.0-rc",
					"2.0.0",
				})
			})

			Convey("With all pre-releases", func() {
				next := toStr(ver.NextVersions(-4, false))
				So(next, ShouldResemble, []string{
					"1.2.3-r1",
					"1.2.3-p1",
					"1.2.4",
					"1.3.0-alpha",
					"1.3.0-beta",
					"1.3.0-pre",
					"1.3.0-rc",
					"1.3.0",
					"2.0.0-alpha",
					"2.0.0-beta",
					"2.0.0-pre",
					"2.0.0-rc",
					"2.0.0",
				})
			})
		})

		ver = "1.2.0-beta2"
		Convey(ver, func() {
			ver, err := NewVersion([]byte(ver))
			So(err, ShouldBeNil)

			Convey("With all pre-releases and numbers", func() {
				next := toStr(ver.NextVersions(-4, true))
				So(next, ShouldResemble, []string{
					"1.2.0-beta3",
					"1.2.0-pre1",
					"1.2.0-rc1",
					"1.2.0",
					"1.2.1",
					"1.3.0-alpha1",
					"1.3.0-beta1",
					"1.3.0-pre1",
					"1.3.0-rc1",
					"1.3.0",
					"2.0.0-alpha1",
					"2.0.0-beta1",
					"2.0.0-pre1",
					"2.0.0-rc1",
					"2.0.0",
				})
			})
		})

	})
}

func TestVersionOrder(t *testing.T) {

	Convey("Version 1.2.3-alpha4 should be…", t, func() {
		v1, _ := NewVersion([]byte("1.2.3-alpha4"))

		Convey("reasonably less than Version 1.2.3", func() {
			v2, _ := NewVersion([]byte("1.2.3"))
			So(v1.limitedLess(v2), ShouldBeTrue)
		})

		Convey("reasonably less than Version 1.2.3-alpha4.0.0.1", func() {
			v2, _ := NewVersion([]byte("1.2.3-alpha4.0.0.1"))
			So(v1.limitedLess(v2), ShouldBeTrue)
		})

		Convey("not reasonably less than 1.2.3-alpha4-p5", func() {
			v2, _ := NewVersion([]byte("1.2.3-alpha4-p5"))
			So(v1.limitedLess(v2), ShouldBeFalse)
		})
	})

}

func TestVersionAccessors(t *testing.T) {
	Convey("For version 1.2.3 we should have", t, func() {
		v, _ := NewVersion([]byte("1.2.3"))

		Convey("major equals 1", func() {
			So(v.Major(), ShouldEqual, 1)
		})

		Convey("minor equals 2", func() {
			So(v.Minor(), ShouldEqual, 2)
		})

		Convey("patch equals 3", func() {
			So(v.Patch(), ShouldEqual, 3)
		})
	})
}

func ExampleVersion_LimitedEqual_first() {
	var (
		pre, _ = NewVersion([]byte("1.0.0-pre"))
		rc, _  = NewVersion([]byte("1.0.0-rc"))
	)

	fmt.Println(pre.LimitedEqual(rc))
	// Output: false
}

func ExampleVersion_LimitedEqual_second() {
	var (
		a, _ = NewVersion([]byte("1.0.0-beta-pre3"))
		b, _ = NewVersion([]byte("1.0.0-beta-pre5"))
	)

	fmt.Println(a.LimitedEqual(b))
	// Output: true
}

func ExampleVersion_LimitedEqual_third() {
	var (
		regular, _ = NewVersion([]byte("1.0.0"))
		patched, _ = NewVersion([]byte("1.0.0-p1"))
	)

	fmt.Println(patched.LimitedEqual(regular))
	fmt.Println(regular.LimitedEqual(patched))

	// Output:
	// false
	// true
}

var verForBenchmarks = []byte("1.2.3-beta")
var benchV, benchErr = NewVersion(append(verForBenchmarks, '5'))

func BenchmarkSemverNewVersion(b *testing.B) {
	v, e := NewVersion(verForBenchmarks)

	for n := 0; n < b.N; n++ {
		v, e = NewVersion(verForBenchmarks)
	}
	benchV, benchErr = v, e
}

var compareResult = 5

func BenchmarkSemverCompare(b *testing.B) {
	v, _ := NewVersion(verForBenchmarks)
	r := Compare(benchV, v)

	for n := 0; n < b.N; n++ {
		r = Compare(benchV, v)
	}
	compareResult = r
}
