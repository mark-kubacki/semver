// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package semver contains types and functions for
// parsing of Versions and (Version-)Ranges.
package semver

import "strconv"

// Errors that are thrown when translating from a string.
const (
	errInvalidVersionString InvalidStringValue = "Given string does not resemble a Version"
	errTooMuchColumns       InvalidStringValue = "Version consists of too much columns"
	errVersionStringLength  InvalidStringValue = "Version is too long"
	errInvalidBuildSuffix   InvalidStringValue = "Version has a '+' but no +buildNNN suffix"
)

// alpha = -4, beta = -3, pre = -2, rc = -1, common = 0, revision = 1, patch = 2
const (
	alpha = iota - 4
	beta
	pre
	rc
	common
	revision
	patch
)

const (
	idxReleaseType   = 4
	idxRelease       = 5
	idxSpecifierType = 9
	idxSpecifier     = 10
)

var releaseDesc = map[int]string{
	alpha:    "alpha",
	beta:     "beta",
	pre:      "pre",
	rc:       "rc",
	revision: "r",
	patch:    "p",
}

var releaseValue = map[string]int{
	"alpha": alpha,
	"beta":  beta,
	"pre":   pre,
	"":      pre,
	"rc":    rc,
	"r":     revision,
	"p":     patch,
}

// InvalidStringValue is returned as error when translating a string into type fail.
type InvalidStringValue string

// Error implements the error interface.
func (e InvalidStringValue) Error() string { return string(e) }

// IsInvalid satisfies a function IsInvalid().
func (e InvalidStringValue) IsInvalid() bool { return true }

// Version represents a version:
// Columns consisting of up to four unsigned integers (1.2.4.99)
// optionally further divided into 'release' and 'specifier' (1.2-634.0-99.8).
type Version struct {
	// 0–3: version, 4: releaseType, 5–8: releaseVer, 9: releaseSpecifier, 10–14: specifier
	version [14]int32
	build   int32
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
	var idx, toIdx, fieldNum, column int
	var strlen = len(str)

	for idx < strlen {
		r := str[idx]
		switch {
		case '0' <= r && r <= '9':
			if column == 4 {
				return errTooMuchColumns
			}
			column++
			for toIdx = idx + 1; toIdx < strlen; toIdx++ {
				p := str[toIdx]
				if !('0' <= p && p <= '9') {
					break
				}
			}

			if fieldNum == idxReleaseType || fieldNum == idxSpecifierType {
				fieldNum++
			}

			n, err := strconv.Atoi(str[idx:toIdx])
			if err != nil {
				return err
			}
			t.version[fieldNum] = int32(n)

			idx = toIdx
			fieldNum++
		case 'a' <= r && r <= 'z':
			column = 0
			for toIdx = idx + 1; toIdx < strlen; toIdx++ {
				p := str[toIdx]
				if !('a' <= p && p <= 'z') {
					break
				}
			}

			typ, known := releaseValue[str[idx:toIdx]]
			if !known {
				return errInvalidVersionString
			}
			switch {
			case fieldNum <= idxReleaseType:
				fieldNum = idxReleaseType
			case fieldNum <= idxSpecifierType:
				fieldNum = idxSpecifierType
			default:
				return errInvalidVersionString
			}
			t.version[fieldNum] = int32(typ)

			idx = toIdx
			fieldNum++
		case r == '.':
			idx++
		case r == '-' || r == '_':
			idx++
			switch {
			case fieldNum <= idxReleaseType:
				fieldNum = idxReleaseType
			case fieldNum <= idxSpecifierType:
				fieldNum = idxSpecifierType
			default:
				return errInvalidVersionString
			}
		case r == '+':
			if strlen < idx+7 || str[idx:idx+6] != "+build" {
				return errInvalidBuildSuffix
			}
			n, err := strconv.Atoi(str[idx+6:])
			if err != nil {
				return err
			}
			t.build = int32(n)
			return nil
		default:
			return errInvalidVersionString
		}

		if fieldNum > 14 {
			return errVersionStringLength
		}
	}

	return nil
}

// signDelta returns the signum of the difference,
// which' precision can be limited by 'cuttofIdx'.
func signDelta(a, b [14]int32, cutoffIdx int) int8 {
	//fmt.Println(a, b)
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

// Compare computes the difference between two Versions and returns its signum.
//
//   1  if a > b
//   0  if a == b
//   -1 if a < b
//
// The 'build' is not compared.
func Compare(a, b Version) int {
	return int(signDelta(a.version, b.version, 14))
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
