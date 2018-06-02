// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"strconv"
)

func numDecimalPlaces(n int32) int {
	var i int
	for i = 1; n > 9; i++ {
		n = n / 10
	}
	return i
}

// Serialize builds a minimal human-readable representation of this Version,
// and returns it as slice.
// Set |minPlaces| to how many columns the prefix must contain.
func (t *Version) serialize(minPlaces int, quoted bool) []byte {
	var idx, bytesNeeded int

	if quoted {
		bytesNeeded = 2
	}

	// Determine how much target space is needed (i.e. the string length).
	for idx = 0; idx <= len(t.version); idx += 5 {
		switch {
		case t.version[idx+3] != 0 || minPlaces >= 4:
			bytesNeeded += 1 + numDecimalPlaces(t.version[idx+3])
			fallthrough
		case t.version[idx+2] != 0 || minPlaces >= 3:
			bytesNeeded += 1 + numDecimalPlaces(t.version[idx+2])
			fallthrough
		case t.version[idx+1] != 0 || minPlaces >= 2:
			bytesNeeded += 1 + numDecimalPlaces(t.version[idx+1])
			fallthrough
		default:
			bytesNeeded += numDecimalPlaces(t.version[idx])
		}
		if t.version[idx+4] != 0 {
			bytesNeeded += 1 + len(releaseDesc[int(t.version[idx+4])])
		}

		remainderEmpty := true
		for j := idx + 5; j < len(t.version); j++ {
			if t.version[j] != 0 {
				remainderEmpty = false
				break
			}
		}
		if remainderEmpty {
			break
		}
	}
	if t.build != 0 {
		bytesNeeded += len("+build") + numDecimalPlaces(t.build)
	}

	// Build the string representation
	target := make([]byte, 0, bytesNeeded)

	if quoted {
		target = append(target, '"')
	}

	for idx = 0; idx < len(t.version); idx += 5 {
		switch {
		case t.version[idx+3] != 0 || minPlaces >= 4:
			target = strconv.AppendUint(target, uint64(t.version[idx]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+1]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+2]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+3]), 10)
		case t.version[idx+2] != 0 || minPlaces >= 3:
			target = strconv.AppendUint(target, uint64(t.version[idx]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+1]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+2]), 10)
		case t.version[idx+1] != 0 || minPlaces >= 2:
			target = strconv.AppendUint(target, uint64(t.version[idx]), 10)
			target = append(target, '.')
			target = strconv.AppendUint(target, uint64(t.version[idx+1]), 10)
		default:
			target = strconv.AppendUint(target, uint64(t.version[idx]), 10)
		}
		if t.version[idx+4] != 0 {
			target = append(target, '-')
			target = append(target, []byte(releaseDesc[int(t.version[idx+4])])...)
		}

		remainderEmpty := true
		for j := idx + 5; j < len(t.version); j++ {
			if t.version[j] != 0 {
				remainderEmpty = false
				break
			}
		}
		if remainderEmpty {
			break
		}
		minPlaces -= 5
	}
	if t.build != 0 {
		target = append(target, []byte("+build")...)
		target = strconv.AppendUint(target, uint64(t.build), 10)
	}
	if quoted {
		target = append(target, '"')
	}

	return target
}

// Bytes returns a slice with the minimal human-readable representation of this Version.
func (t *Version) Bytes() []byte {
	return t.serialize(0, false)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t *Version) MarshalBinary() ([]byte, error) {
	return t.serialize(0, false), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Version) UnmarshalBinary(b []byte) error {
	return t.Parse(string(b))
}

// String returns the string representation of t.
func (t *Version) String() string {
	return string(t.serialize(3, false))
}

// MarshalJSON implements the json.Marshaler interface.
func (t *Version) MarshalJSON() ([]byte, error) {
	return t.serialize(0, true), nil
}

// MarshalText implements the encoding.TestMarshaler interface.
func (t *Version) MarshalText() ([]byte, error) {
	return t.serialize(0, false), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Version) UnmarshalJSON(b []byte) error {
	if len(b) > 2 && b[0] == '"' || b[0] == '\'' || b[0] == '`' {
		// We can ignore the closing because the JSON engine will throw an error on any mismatch for us.
		return t.Parse(string(b[1 : len(b)-1]))
	}
	return t.Parse(string(b))
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Version) UnmarshalText(b []byte) error {
	return t.Parse(string(b))
}
