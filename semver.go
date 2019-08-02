// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"bytes"
)

// Errors that are thrown when translating from a string.
const (
	errInvalidVersionString InvalidStringValue = "Given string does not resemble a Version"
	errTooManyColumns       InvalidStringValue = "Version consists of too many columns"
	errVersionStringLength  InvalidStringValue = "Version is too long"
	errInvalidBuildSuffix   InvalidStringValue = "Version has a '+' but no +buildNNN suffix"
	errInvalidType          InvalidStringValue = "Cannot read this type into a Version"
	errOutOfBounds          InvalidStringValue = "The source representation does not fit into a Version"
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

var buildsuffix = []byte("+build")

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
	// 0–3: version, 4: releaseType, 5–8: releaseVer, 9: releaseSpecifier, 10–: specifier
	version [14]int32
	build   int32
}

// MustParse is NewVersion on strings, and panics on errors.
//
// This is a convenience function for a cloud plattform provider.
func MustParse(str string) Version {
	ver, err := NewVersion([]byte(str))
	if err != nil {
		panic(err.Error())
	}
	return ver
}

// NewVersion translates the given string, which must be free of whitespace,
// into a single Version.
func NewVersion(str []byte) (Version, error) {
	ver := Version{}
	err := (&ver).unmarshalText(str)
	return ver, err
}

// Parse reads a string into the given version, overwriting any existing values.
//
// Deprecated: Use the idiomatic UnmarshalText instead.
func (t *Version) Parse(str string) error {
	t.version = [14]int32{}
	t.build = 0

	return t.unmarshalText([]byte(str))
}

func isNumeric(ch byte) bool {
	return ((ch - '0') <= 9)
}

func isSmallLetter(ch byte) bool {
	// case insensitive: (ch | 0x20)
	return ((ch - 'a') <= ('z' - 'a'))
}

// atoui consumes up to n byte from b to convert them into |val|.
func atoui(b []byte) (n int, val uint32) {
	for ; n <= 10 && n < len(b); n++ {
		v := b[n] - '0' // see above 'isNumeric'
		if v > 9 {
			break
		}
		val = val*10 + uint32(v)
	}
	return
}

// unmarshalText implements the encoding.TextUnmarshaler interface,
// but assumes the data structure is pristine.
func (t *Version) unmarshalText(str []byte) error {
	var idx, fieldNum, column int
	var strlen = len(str)

	if strlen > 1 && str[idx] == 'v' {
		idx++
	}

	for idx < strlen {
		r := str[idx]
		switch {
		case r == '.':
			idx++
			column++
			if column >= 4 || idx >= strlen {
				return errTooManyColumns
			}
			fieldNum++
			fallthrough
		case isNumeric(r):
			idxDelta, n := atoui(str[idx:])
			if idxDelta == 0 || idxDelta >= 10 { // strlen(maxInt) is 10
				return errInvalidVersionString
			}
			t.version[fieldNum] = int32(n)

			idx += idxDelta
		case r == '-' || r == '_':
			idx++
			if idx < strlen && isNumeric(str[idx]) {
				column = 0
				switch {
				case fieldNum < idxReleaseType:
					fieldNum = idxReleaseType + 1
				case fieldNum < idxSpecifierType:
					fieldNum = idxSpecifierType + 1
				default:
					return errInvalidVersionString
				}
				continue
			}
			fallthrough
		case isSmallLetter(r):
			toIdx := idx + 1
			for ; toIdx < strlen && isSmallLetter(str[toIdx]); toIdx++ {
			}

			if toIdx > strlen {
				return errInvalidVersionString
			}
			typ, known := releaseValue[string(str[idx:toIdx])]
			if !known {
				return errInvalidVersionString
			}
			switch {
			case fieldNum < idxReleaseType:
				fieldNum = idxReleaseType
			case fieldNum < idxSpecifierType:
				fieldNum = idxSpecifierType
			default:
				return errInvalidVersionString
			}
			t.version[fieldNum] = int32(typ)
			if toIdx+1 < strlen && str[toIdx] == '.' {
				toIdx++
			}

			fieldNum++
			column = 0
			idx = toIdx
		case r == '+':
			if strlen < idx+len(buildsuffix)+1 || !bytes.Equal(str[idx:idx+len(buildsuffix)], buildsuffix) {
				return errInvalidBuildSuffix
			}
			idx += len(buildsuffix)
			idxDelta, n := atoui(str[idx:])
			if idxDelta > 9 || idx+idxDelta < strlen {
				return errInvalidBuildSuffix
			}
			t.build = int32(n)
			return nil
		default:
			return errInvalidVersionString
		}
	}

	return nil
}

// signDelta returns the signum of the difference,
// which' precision can be limited by 'cuttofIdx'.
func signDelta(a, b [14]int32, cutoffIdx int) int8 {
	_ = a[0:cutoffIdx]
	for i := 0; i < len(a) && i < cutoffIdx; i++ {
		if a[i] == b[i] {
			continue
		}
		x := a[i] - b[i]
		return int8((x >> 31) - (-x >> 31))
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
	for i := 0; i < len(a.version); i++ {
		if a.version[i] == b.version[i] {
			continue
		}
		x := a.version[i] - b.version[i]
		return int((x >> 31) - (-x >> 31))
	}
	return 0
}

// Less is a convenience function for sorting.
func (t Version) Less(o Version) bool {
	sd := Compare(t, o)
	return sd < 0 || (sd == 0 && t.build < o.build)
}

// limitedLess compares two Versions
// with a precision limited to version, (pre-)release type and (pre-)release version.
//
// Commutative.
func (t Version) limitedLess(o Version) bool {
	return signDelta(t.version, o.version, idxSpecifierType) < 0
}

// LimitedEqual returns true of two versions share the same prefix,
// which is the "actual version", (pre-)release type, and (pre-)release version.
// The exception are patch-levels, which are always equal.
//
// Use this, for example, to tell a beta from a regular version;
// or to accept a patched version as regular version.
func (t Version) LimitedEqual(o Version) bool {
	if t.version[idxReleaseType] == common && o.version[idxReleaseType] > common {
		return t.sharesPrefixWith(o)
	}
	return signDelta(t.version, o.version, idxSpecifierType) == 0
}

// IsAPreRelease is used to discriminate pre-releases.
func (t Version) IsAPreRelease() bool {
	return t.version[idxReleaseType] < common
}

// sharesPrefixWith compares two Versions with a fixed limited precision.
//
// A 'prefix' is the major, minor, patch and revision number.
// For example: 1.2.3.4…
func (t Version) sharesPrefixWith(o Version) bool {
	return signDelta(t.version, o.version, idxReleaseType) == 0
}

// Major returns the major of a version.
func (t Version) Major() int {
	return int(t.version[0])
}

// Minor returns the minor of a version.
func (t Version) Minor() int {
	return int(t.version[1])
}

// Patch returns the patch of a version.
func (t Version) Patch() int {
	return int(t.version[2])
}

// NextVersions returns a list of possible next versions after t. For each of
// the three version points, pre-releases are given as options starting with
// the minimum release type (-4 <= 0), and those release types are numbered
// if numberedPre is true. Release types:
//
//   alpha: -4
//   beta:  -3
//   pre:   -2
//   rc:    -1
//   common: 0
//
// Thus, if you don't want any pre-release options, set minReleaseType to 0.
func (t *Version) NextVersions(minReleaseType int, numberedPre bool) []*Version {
	var next []*Version

	if minReleaseType < alpha || minReleaseType > common {
		return next
	}

	// if this is a pre-release, suggest next pre-releases or
	// common of same version
	for releaseType := t.version[idxReleaseType]; releaseType < common; releaseType++ {
		if releaseType == t.version[idxReleaseType] {
			if !numberedPre {
				continue
			}
			ver := *t
			ver.version[idxRelease]++
			next = append(next, &ver)
		} else {
			ver := *t
			ver.version[idxReleaseType] = releaseType
			if numberedPre {
				ver.version[idxRelease] = 1
			} else {
				ver.version[idxRelease] = 0
			}
			next = append(next, &ver)
		}
	}
	if t.version[idxReleaseType] < common {
		ver := *t
		ver.version[idxReleaseType] = common
		ver.version[idxRelease] = 0
		next = append(next, &ver)
	}

	// if the current version is at least common release type,
	// suggest patch or revision if not one of those already
	if t.version[idxReleaseType] == common ||
		t.version[idxReleaseType] == patch {
		ver := *t
		ver.version[idxReleaseType] = revision
		ver.version[idxRelease] = 1
		next = append(next, &ver)
	}
	if t.version[idxReleaseType] == common ||
		t.version[idxReleaseType] == revision {
		ver := *t
		ver.version[idxReleaseType] = patch
		ver.version[idxRelease] = 1
		next = append(next, &ver)
	}

	for i := idxReleaseType - 2; 0 <= i; i-- {
		// for each version point, iterate the release types within desired bounds
		for releaseType := int32(minReleaseType); releaseType <= common; releaseType++ {
			ver := *t
			ver.version[i]++
			for j := i + 1; j < len(ver.version); j++ {
				ver.version[j] = 0 // when incrementing, reset next points to 0
			}
			if i == 2 && releaseType < common {
				continue // patches seldom have pre-releases
			}
			ver.version[idxReleaseType] = releaseType
			if releaseType < common {
				if numberedPre {
					ver.version[idxRelease] = 1
				} else {
					ver.version[idxRelease] = 0
				}
			}
			next = append(next, &ver)
		}
	}

	return next
}
