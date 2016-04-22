// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package semver contains types and functions for
// parsing of Versions and (Version-)Ranges.
package semver

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Errors that are thrown when translating from a string.
var (
	ErrInvalidVersionString = errors.New("Given string does not resemble a Version")
	ErrTooMuchColumns       = errors.New("Version consists of too much columns")
)

type dotDelimitedNumber []int

func newDotDelimitedNumber(str string) (dotDelimitedNumber, error) {
	strSequence := strings.Split(str, ".")
	if len(strSequence) > 4 {
		return nil, ErrTooMuchColumns
	}
	numSequence := make(dotDelimitedNumber, 0, len(strSequence))
	for _, s := range strSequence {
		i, err := strconv.Atoi(s)
		if err != nil {
			return numSequence, err
		}
		numSequence = append(numSequence, i)
	}
	return numSequence, nil
}

// alpha = -4, beta = -3, pre = -2, rc = -1, common = 0, patch = 1
const (
	alpha = iota - 4
	beta
	pre
	rc
	common
	patch
)

const (
	idxReleaseType   = 4
	idxRelease       = 5
	idxSpecifierType = 9
	idxSpecifier     = 10
)

var releaseDesc = map[int]string{
	alpha: "alpha",
	beta:  "beta",
	pre:   "pre",
	rc:    "rc",
	patch: "p",
}

var releaseValue = map[string]int{
	"alpha": alpha,
	"beta":  beta,
	"pre":   pre,
	"":      pre,
	"rc":    rc,
	"p":     patch,
}

var verRegexp = regexp.MustCompile(`^(\d+(?:\.\d+){0,3})(?:([-_]alpha|[-_]beta|[-_]pre|[-_]rc|[-_]p|-)(\d+(?:\.\d+){0,3})?)?(?:([-_]alpha|[-_]beta|[-_]pre|[-_]rc|[-_]p|-)(\d+(?:\.\d+){0,3})?)?(?:(\+build)(\d*))?$`)

// Version represents a version:
// Columns consisting of up to four unsigned integers (1.2.4.99)
// optionally further divided into 'release' and 'specifier' (1.2-634.0-99.8).
type Version struct {
	// 0–3: version, 4: releaseType, 5–8: releaseVer, 9: releaseSpecifier, 10–14: specifier
	version [14]int
	build   int
}

// NewVersion translates the given string, which must be free of whitespace,
// into a single Version.
func NewVersion(str string) (*Version, error) {
	ver := &Version{}
	err := ver.Parse(str)
	return ver, err
}

// Parse reads a string into the given version, overwriting any existing values.
func (t *Version) Parse(str string) error {
	allMatches := verRegexp.FindAllStringSubmatch(str, -1)
	if allMatches == nil {
		return ErrInvalidVersionString
	}

	m := allMatches[0]

	// version
	n, err := newDotDelimitedNumber(m[1])
	if err != nil {
		return err
	}
	copy(t.version[:], n)

	// release
	if m[2] != "" {
		t.version[idxReleaseType] = releaseValue[strings.Trim(m[2], "-_")]
	}
	if m[3] != "" {
		n, err := newDotDelimitedNumber(m[3])
		if err != nil {
			return err
		}
		copy(t.version[idxRelease:], n)
	}

	// release specifier
	if m[4] != "" {
		t.version[idxSpecifierType] = releaseValue[strings.Trim(m[4], "-_")]
	}
	if m[5] != "" {
		n, err := newDotDelimitedNumber(m[5])
		if err != nil {
			return err
		}
		copy(t.version[idxSpecifier:], n)
	}

	// build
	if m[7] != "" {
		i, err := strconv.Atoi(m[7])
		if err != nil {
			return err
		}
		t.build = i
	}

	return nil
}

// signDelta returns the signum of the difference,
// which' precision can be limited by 'cuttofIdx'.
func signDelta(a, b [14]int, cutoffIdx int) int8 {
	for i := range a {
		if i >= cutoffIdx {
			return 0
		}
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// Less is a convenience function for sorting.
func (t *Version) Less(o *Version) bool {
	sd := signDelta(t.version, o.version, 15)
	return sd < 0 || (sd == 0 && t.build < o.build)
}

// limitedLess compares two Versions
// with a precision limited to version, (pre-)release type and (pre-)release version.
//
// Commutative.
func (t *Version) limitedLess(o *Version) bool {
	return signDelta(t.version, o.version, idxSpecifierType) < 0
}

// LimitedEqual returns true of two versions share the same prefix,
// which is the "actual version", (pre-)release type, and (pre-)release version.
// The exception are patch-levels, which are always equal.
//
// Use this, for example, to tell a beta from a regular version;
// or to accept a patched version as regular version.
func (t *Version) LimitedEqual(o *Version) bool {
	if t.version[idxReleaseType] == common && o.version[idxReleaseType] > common {
		return t.sharesPrefixWith(o)
	}
	return signDelta(t.version, o.version, idxSpecifierType) == 0
}

// IsAPreRelease is used to discriminate pre-releases.
func (t *Version) IsAPreRelease() bool {
	return t.version[idxReleaseType] < common
}

// sharesPrefixWith compares two Versions with a fixed limited precision.
//
// A 'prefix' is the major, minor, patch and revision number.
// For example: 1.2.3.4…
func (t *Version) sharesPrefixWith(o *Version) bool {
	return signDelta(t.version, o.version, idxReleaseType) == 0
}
