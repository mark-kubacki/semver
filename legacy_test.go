// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
			ver := MustParse(ver)

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
			ver := MustParse(ver)

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
			ver := MustParse(ver)

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
