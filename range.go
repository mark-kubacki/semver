// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import "strings"

// Errors which can be encountered when parsing into a Range.
const (
	errUnsupportedShortcutNotation InvalidStringValue = "Unsupported shortcut notation for Range"
)

// Range is a subset of the universe of Versions: It can have a lower and upper boundary.
// For example, "1.2–2.0" is such a Range, with two boundaries.
type Range struct {
	lower       *Version
	equalsLower bool
	upper       *Version
	equalsUpper bool
}

// NewRange translates a string into a Range.
func NewRange(str string) (*Range, error) {
	if str == "*" || str == "x" || str == "" {
		// an empty Range contains everything
		return new(Range), nil
	}
	isNaturalRange := true
	if strings.HasSuffix(str, ".x") || strings.HasSuffix(str, ".*") {
		str = strings.TrimRight(str, ".x*")
		isNaturalRange = false
	}
	if str[0] == '^' || str[0] == '~' {
		return newRangeByShortcut(str)
	}

	var leftEnd, rightStart, leftDotCount int
	var upperBound, lowerBound bool = true, true
	for i, r := range str {
		if r == '.' {
			leftDotCount++
		}
		if r == ' ' || r == '–' || r == ',' {
			if leftEnd == 0 {
				leftEnd = i
			}
			rightStart = i
			continue
		} else {
			if rightStart != 0 {
				rightStart++
				if r != '-' {
					break
				}
			}
		}
		switch r {
		case '<', '≤':
			lowerBound = false
		case '>', '≥':
			upperBound = false
		}
	}

	isNaturalRange = isNaturalRange && leftEnd != rightStart && (len(str)-rightStart) > 0
	if !isNaturalRange {
		switch leftDotCount {
		case 1:
			return newRangeByShortcut("~" + str)
		case 0:
			return newRangeByShortcut("^" + str)
		}
	}
	vr := new(Range)
	if leftEnd == rightStart {
		err := vr.setBound(str, lowerBound, upperBound)
		return vr, err
	}

	if leftEnd == 0 {
		leftEnd = len(str)
	}
	if err := vr.setBound(str[:leftEnd], true, false); err != nil {
		return vr, err
	}
	if err := vr.setBound(str[rightStart:], false, true); err != nil {
		return vr, err
	}

	return vr, nil
}

func (r *Range) setBound(str string, isLower, isUpper bool) error {
	var versionStartIdx int
	for ; versionStartIdx < len(str); versionStartIdx++ {
		r := str[versionStartIdx]
		if '0' <= r && r <= '9' {
			goto startFound
		}
	}
	return errInvalidVersionString

startFound:
	num, err := NewVersion(str[versionStartIdx:])
	if err != nil {
		return err
	}

	prefix := str[:versionStartIdx]
	equalOk := versionStartIdx == 0 || strings.Contains(prefix, "=")
	if isUpper {
		r.equalsUpper = equalOk
		r.upper = num
	}
	if isLower {
		r.equalsLower = equalOk
		r.lower = num
	}

	return nil
}

// newRangeByShortcut covers the special case of Ranges whose boundaries
// are declared using prefixes.
func newRangeByShortcut(str string) (*Range, error) {
	t := strings.TrimLeft(str, "~^")
	num, err := NewVersion(t)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(t, "0.0.") {
		return NewRange(t)
	}

	r := new(Range)
	r.lower = num
	r.equalsLower = true
	r.upper = new(Version)

	switch {
	case strings.HasPrefix(t, "0."):
		r.upper.version[0] = r.lower.version[0]
		r.upper.version[1] = r.lower.version[1] + 1
	case str[0] == '^' || !strings.ContainsAny(t, "."):
		r.upper.version[0] = r.lower.version[0] + 1
	case str[0] == '~':
		r.upper.version[0] = r.lower.version[0]
		r.upper.version[1] = r.lower.version[1] + 1
	default:
		return nil, errUnsupportedShortcutNotation
	}

	return r, nil
}

// GetLowerBoundary translates a boundary into a Version.
func (r *Range) GetLowerBoundary() *Version {
	return r.lower
}

// GetUpperBoundary translates a boundary into a Version.
func (r *Range) GetUpperBoundary() *Version {
	return r.upper
}

// Contains returns true if a Version is inside this Range.
func (r *Range) Contains(v *Version) bool {
	if v == nil {
		return false
	}

	if r.upper == r.lower {
		return r.lower.LimitedEqual(v)
	}

	return r.satisfiesLowerBound(v) && r.satisfiesUpperBound(v)
}

// IsSatisfiedBy works like Contains,
// but rejects pre-releases if neither of the bounds is a pre-release.
//
// Use this in the context of pulling in packages because it follows the spirit of §9 SemVer.
// Also see https://github.com/npm/node-semver/issues/64
func (r *Range) IsSatisfiedBy(v *Version) bool {
	if !r.Contains(v) {
		return false
	}
	if v.IsAPreRelease() {
		if r.lower != nil && r.lower.IsAPreRelease() && r.lower.sharesPrefixWith(v) {
			return true
		}
		if r.upper != nil && r.upper.IsAPreRelease() && r.upper.sharesPrefixWith(v) {
			return true
		}
		return false
	}
	return true
}

func (r *Range) satisfiesLowerBound(v *Version) bool {
	if r.lower == nil {
		return true
	}

	equal := r.lower.LimitedEqual(v)
	if r.equalsLower && equal {
		return true
	}

	return r.lower.limitedLess(v) && !equal
}

func (r *Range) satisfiesUpperBound(v *Version) bool {
	if r.upper == nil {
		return true
	}

	equal := r.upper.LimitedEqual(v)
	if r.equalsUpper && equal {
		return true
	}

	if !r.equalsUpper && r.upper.version[idxReleaseType] == common {
		equal = r.upper.sharesPrefixWith(v)
	}

	return v.limitedLess(r.upper) && !equal
}

// Satisfies is a convenience function for former NodeJS developers
// which works on two strings.
func Satisfies(aVersion, aRange string) (bool, error) {
	v, err := NewVersion(aVersion)
	if err != nil {
		return false, err
	}
	r, err := NewRange(aRange)
	if err != nil {
		return false, err
	}

	return r.IsSatisfiedBy(v), nil
}
