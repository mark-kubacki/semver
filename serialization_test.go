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
				expect, _ := NewVersion("2.31.4")

				err := json.Unmarshal(in, &out)
				So(err, ShouldBeNil)
				So(out.Ver, ShouldResemble, expect)
			})
			Convey("even without quotes", func() {
				in := []byte(`{"ver": 2}`)
				var out struct{ Ver Version }
				expect, _ := NewVersion("v2")

				err := json.Unmarshal(in, &out)
				So(err, ShouldBeNil)
				So(out.Ver, ShouldResemble, expect)
			})
		})

		Convey("will be serialized correctly", func() {
			given, _ := NewVersion("2.31.4")
			t := struct{ Ver Version }{given}

			out, err := json.Marshal(&given)
			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `"2.31.4"`) // cast to 'string' for legibility

			out, err = json.Marshal(&t)
			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `{"Ver":"2.31.4"}`) // cast to 'string' for legibility
		})
	})
}
