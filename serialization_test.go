// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerialization(t *testing.T) {
	Convey("Versions within JSON…", t, FailureContinues, func() {
		Convey("get parsed into structs", func() {
			Convey("if quoted", func() {
				in := []byte(`{"ver": "2.31.4"}`)
				var out struct{ Ver Version }
				expect, _ := NewVersion([]byte("2.31.4"))

				err := json.Unmarshal(in, &out)
				So(err, ShouldBeNil)
				So(out.Ver, ShouldResemble, expect)
			})
			Convey("even without quotes", func() {
				in := []byte(`{"ver": 2}`)
				var out struct{ Ver Version }
				expect, _ := NewVersion([]byte("v2"))

				err := json.Unmarshal(in, &out)
				So(err, ShouldBeNil)
				So(out.Ver, ShouldResemble, expect)
			})
		})

		Convey("will be serialized correctly", func() {
			for _, str := range []string{
				"2.31.4", "14.9", "1.5.3.1", "8",
				"8+build66",
				"1.5.1-3", "1.12-rc2",
				"0-0-0.0.0.4",
			} {
				given, err := NewVersion([]byte(str))
				So(err, ShouldBeNil)
				if err != nil {
					continue
				}
				t := struct{ Ver Version }{given}

				out, err := json.Marshal(&given)
				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, `"`+string(str)+`"`) // cast to 'string' for legibility

				out, err = json.Marshal(&t)
				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, `{"Ver":"`+string(str)+`"}`) // cast to 'string' for legibility
			}
		})
	})
}
