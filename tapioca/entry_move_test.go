// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_StyledMove(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		start, width int
		expected     string
	}{
		{
			name:     "Move styled text, no padding",
			input:    "\x1b[31mHello, World!\x1b[0m",
			start:    7,
			width:    3,
			expected: "\x1b[31mWor\x1b[0m",
		},
		{
			name:     "Move styled text, with prefix",
			input:    "\x1b[32mHello, World!\x1b[0m",
			start:    -3,
			width:    5,
			expected: "   \x1b[32mHe\x1b[0m",
		},
		{
			name:     "Move styled text, with suffix",
			input:    "\x1b[34mHello, World!\x1b[0m",
			start:    10,
			width:    5,
			expected: "\x1b[34mld!\x1b[0m  ",
		},
		{
			name:     "Move styled text, with both prefix and suffix",
			input:    "\x1b[35mHi\x1b[0m",
			start:    -2,
			width:    6,
			expected: "  \x1b[35mHi\x1b[0m  ",
		},
		// wide characters
		{
			name:     "Move styled wide text, no padding, no cutting",
			input:    "\x1b[31mこんにちは\x1b[0m", // "Hello" in Japanese
			start:    0,
			width:    6,
			expected: "\x1b[31mこんに\x1b[0m",
		},
		{
			name:     "Move styled wide text, with prefix, no cutting",
			input:    "\x1b[32mこんにちは\x1b[0m",
			start:    -2,
			width:    6,
			expected: "  \x1b[32mこん\x1b[0m",
		},
		{
			name:     "Move styled wide text, with suffix, no cutting",
			input:    "\x1b[34mこんにちは\x1b[0m",
			start:    6,
			width:    5,
			expected: "\x1b[34mちは\x1b[0m ",
		},
		{
			name:     "Move styled wide text, with prefix, cutting at end",
			input:    "\x1b[35mこんにちは\x1b[0m",
			start:    -1,
			width:    4,
			expected: " \x1b[35mこ\x1b[0m ",
		},
		{
			name:     "Move styled wide text, with suffix, cutting at start",
			input:    "\x1b[36mこんにちは\x1b[0m",
			start:    7,
			width:    4,
			expected: " \x1b[36mは\x1b[0m ",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			moved := entry.StyledMove(tc.start, tc.width)
			assert.Equal(t, tc.expected, moved)
		})
	}
}
