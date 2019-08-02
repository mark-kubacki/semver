// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver_test

import (
	"fmt"

	"blitznote.com/src/semver/v3"
)

func Example_version() {
	v1 := semver.MustParse("1.2.3-beta")
	v2 := semver.MustParse("2.0.0-alpha20140805.456-rc3+build1800")

	fmt.Println(v1.Less(v2))

	// Output: true
}

func Example_range() {
	v1 := semver.MustParse("1.2.3-beta")
	r1, _ := semver.NewRange([]byte("~1.2"))

	fmt.Println(r1.Contains(v1))
	fmt.Println(r1.IsSatisfiedBy(v1)) // pre-releases don't satisfy

	// Output:
	// true
	// false
}

func ExampleCompare() {
	v1 := semver.MustParse("v1")
	v2 := semver.MustParse("v2.0")
	v3 := semver.MustParse("v3.0.0")

	fmt.Println("Compare", v3, v2, "=", semver.Compare(v3, v2))
	fmt.Println("Compare", v2, v2, "=", semver.Compare(v2, v2))
	fmt.Println("Compare", v1, v2, "=", semver.Compare(v1, v2))

	// Output:
	// Compare 3.0.0 2.0.0 = 1
	// Compare 2.0.0 2.0.0 = 0
	// Compare 1.0.0 2.0.0 = -1
}

func ExampleRange_Contains_first() {
	v := semver.MustParse("1.4.3")
	r, _ := semver.NewRange([]byte("^1.2"))

	fmt.Println(r.Contains(v))

	// Output: true
}

func ExampleRange_Contains_second() {
	v := semver.MustParse("1.4.3")
	r, _ := semver.NewRange([]byte("1.2 <2.0.0"))

	fmt.Println(r.Contains(v))

	// Output: true
}

func ExampleRange_Contains_prerelases() {
	v := semver.MustParse("1.4.3-beta")
	r, _ := semver.NewRange([]byte("1.2 <2.0.0"))

	fmt.Println(r.Contains(v))

	// Output: true
}

func ExampleRange_GetLowerBoundary() {
	r, _ := semver.NewRange([]byte("^1.2"))
	fmt.Println(*r.GetLowerBoundary())

	// Output:
	// 1.2.0
}

func ExampleRange_GetUpperBoundary_first() {
	r, _ := semver.NewRange([]byte("1.2 <2.0.0"))
	fmt.Println(*r.GetUpperBoundary())

	// Output:
	// 2.0.0
}

func ExampleRange_GetUpperBoundary_second() {
	r, _ := semver.NewRange([]byte("~1.2.3"))
	fmt.Println(*r.GetUpperBoundary())

	// Output:
	// 1.3.0
}

func ExampleRange_IsSatisfiedBy_full() {
	v := semver.MustParse("1.2.3")
	r, _ := semver.NewRange([]byte("~1.2"))

	fmt.Println(r.IsSatisfiedBy(v))

	// Output: true
}

func ExampleRange_IsSatisfiedBy_prerelases() {
	// Unlike with Contains, this won't select prereleases.
	pre := semver.MustParse("1.2.3-beta")
	r, _ := semver.NewRange([]byte("~1.2"))

	fmt.Println(r.IsSatisfiedBy(pre))

	// Output: false
}

func ExampleMustParse() {
	v := semver.MustParse("v1.14")
	fmt.Println(v)

	// Output:
	// 1.14.0
}

func ExampleNewVersion() {
	for _, str := range []string{"v1.14", "6.0.2.1", "14b6"} {
		v, err := semver.NewVersion([]byte(str))
		fmt.Println(v, err)
	}

	// Output:
	// 1.14.0 <nil>
	// 6.0.2.1 <nil>
	// 14.0.0 Given string does not resemble a Version
}

func ExampleVersion_Bytes_first() {
	v := semver.MustParse("1.0")
	fmt.Println(v.Bytes())

	// Output:
	// [49]
}

func ExampleVersion_Bytes_second() {
	v := semver.MustParse("4.8")
	fmt.Println(v.Bytes())

	// Output:
	// [52 46 56]
}

func ExampleVersion_Bytes_minimal() {
	v := semver.MustParse("1.13beta")

	fmt.Println("Bytes()  =", string(v.Bytes()))
	fmt.Println("String() =", v.String())

	// Output:
	// Bytes()  = 1.13-beta
	// String() = 1.13.0-beta
}

func ExampleVersion_IsAPreRelease() {
	v, pre := semver.MustParse("1.12"), semver.MustParse("1.13beta")

	fmt.Println(v, "is a pre-release:", v.IsAPreRelease())
	fmt.Println(pre, "is a pre-release:", pre.IsAPreRelease())

	// Output:
	// 1.12.0 is a pre-release: false
	// 1.13.0-beta is a pre-release: true
}

func ExampleVersion_Less() {
	l, r := semver.MustParse("v2"), semver.MustParse("v3")
	fmt.Println(l.Less(r), ",", l, "is 'less' than", r)

	l, r = semver.MustParse("v1+build7"), semver.MustParse("v1+build9")
	fmt.Println(l.Less(r), ",", l, "is 'less' than", r)

	// Output:
	// true , 2.0.0 is 'less' than 3.0.0
	// true , 1.0.0+build7 is 'less' than 1.0.0+build9
}

func ExampleVersion_LimitedEqual_first() {
	// The version prefix does, but the first pre-release type does not match.
	pre := semver.MustParse("1.0.0-pre")
	rc := semver.MustParse("1.0.0-rc")

	fmt.Println(pre.LimitedEqual(rc))

	// Output: false
}

func ExampleVersion_LimitedEqual_second() {
	// The difference is beyond LimitedEqual's cutoff, so these "equal".
	a := semver.MustParse("1.0.0-beta-pre3")
	b := semver.MustParse("1.0.0-beta-pre5")

	fmt.Println(a.LimitedEqual(b))

	// Output: true
}

func ExampleVersion_LimitedEqual_third() {
	regular := semver.MustParse("1.0.0")
	patched := semver.MustParse("1.0.0-p1")

	// A patched version supposedly does more, so is more than; and its unequal.
	fmt.Println(patched.LimitedEqual(regular))
	// This will work because the regular version is a subset is in its subset.
	fmt.Println(regular.LimitedEqual(patched))

	// Output:
	// false
	// true
}

func ExampleVersion_Major() {
	v := semver.MustParse("v1.2.3")
	fmt.Println(v.Major())
	// Output: 1
}

func ExampleVersion_Minor() {
	v := semver.MustParse("v1.2.3")
	fmt.Println(v.Minor())
	// Output: 2
}

func ExampleVersion_Patch() {
	v := semver.MustParse("v1.2.3")
	fmt.Println(v.Patch())
	// Output: 3
}

func ExampleVersion_Scan() {
	a, b, c := new(semver.Version), new(semver.Version), new(semver.Version)
	errA := a.Scan("5.5.65")
	errB := b.Scan(int64(12))
	errC := c.Scan(-1)

	fmt.Println(a, errA)
	fmt.Println(b, errB)
	fmt.Println(c, errC)

	// Output:
	// 5.5.65 <nil>
	// 12.0.0 <nil>
	// 0.0.0 Cannot read this type into a Version
}

func ExampleVersion_String() {
	str := "v2.1"
	v := semver.MustParse(str)
	fmt.Println(str, "is", v.String(), "but as Bytes():", string(v.Bytes()))

	// Output:
	// v2.1 is 2.1.0 but as Bytes(): 2.1
}
