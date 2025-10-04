// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_StyledSubstring(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		start    int
		end      int
		expected string
	}{
		{
			name:     "simple substring",
			input:    "Hello, World!",
			start:    7,
			end:      12,
			expected: "World",
		},
		{
			name:     "wide characters",
			input:    "Hello, 你好!",
			start:    7,
			end:      9,
			expected: "你好",
		},
		{
			name:     "substring with ANSI codes",
			input:    "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m",
			start:    4,
			end:      7,
			expected: "and",
		},
		{
			name:     "cut through styled text at start",
			input:    "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m",
			start:    2,
			end:      5,
			expected: "\x1b[31md\x1b[0m a",
		},
		{
			name:     "cut through styled text at end",
			input:    "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m",
			start:    6,
			end:      10,
			expected: "d \x1b[32mGr\x1b[0m",
		},
		{
			name:     "cut at both ends",
			input:    "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m",
			start:    2,
			end:      10,
			expected: "\x1b[31md\x1b[0m and \x1b[32mGr\x1b[0m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.styledSubstring(tc.start, tc.end)
			assert.Equal(t, tc.expected, got)
		})
	}
}
