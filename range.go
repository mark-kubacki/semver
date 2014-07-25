// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"errors"
	"strings"
)

type Range struct {
	lower       *Version
	equalsLower bool
	upper       *Version
	equalsUpper bool
}

func NewRange(str string) (*Range, error) {
	if str == "*" || str == "x" || str == "" {
		// an empty Range contains everything
		return new(Range), nil
	}
	if strings.HasSuffix(str, ".x") || strings.HasSuffix(str, ".*") {
		str = strings.TrimRight(str, ".x*")
	}
	if str[0] == '^' || str[0] == '~' {
		return newRangeByShortcut(str)
	}

	isNaturalRange := strings.ContainsAny(str, " –")
	if !isNaturalRange {
		switch strings.Count(str, ".") {
		case 1:
			return newRangeByShortcut("~" + str)
		case 0:
			return newRangeByShortcut("^" + str)
		}
	}

	vr := new(Range)
	if !isNaturalRange {
		err := vr.setBound(str)
		return vr, err
	}

	for _, delimiter := range []string{" - ", " – ", "–", " "} {
		if strings.Contains(str, delimiter) {
			parts := strings.Split(str, delimiter)
			if len(parts) == 2 {
				if parts[0][0] == '>' {
					vr.setBound(parts[0])
				} else {
					vr.setBound(">=" + parts[0])
				}
				if parts[1][0] == '<' {
					vr.setBound(parts[1])
				} else {
					vr.setBound("<=" + parts[1])
				}
				return vr, nil
			} else {
				return nil, errors.New("Range contains more than two elements.")
			}
		}
	}

	return nil, nil
}

func (r *Range) setBound(str string) error {
	t := strings.TrimLeft(str, "~=v<>^")
	num, err := NewVersion(t)
	if err != nil {
		return err
	}

	equalOk := strings.Contains(str, "=")
	if !strings.Contains(str, ">") {
		r.equalsUpper = equalOk
		r.upper = num
	}
	if !strings.Contains(str, "<") {
		r.equalsLower = equalOk
		r.lower = num
	}

	return nil
}

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
		return nil, errors.New("Unsupported shortcut notation for Range.")
	}

	return r, nil
}

func (r *Range) Contains(v *Version) bool {
	if v == nil {
		return false
	}

	if r.upper == r.lower {
		return r.lower.LimitedEqual(v)
	}

	return r.satisfiesLowerBound(v) && r.satisfiesUpperBound(v)
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
		equal = signDelta(r.upper.version, v.version, idxReleaseType) == 0
	}

	return v.limitedLess(r.upper) && !equal
}
