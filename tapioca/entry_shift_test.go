// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_StyledShift(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		start    int
		end      int
		expected string
	}{
		{
			name:     "simple shift, no styles",
			input:    "0123456789",
			start:    3,
			end:      6,
			expected: "345",
		},
		{
			name:     "shift with styles",
			input:    "012\x1b[31m345\x1b[0m6789", // "345" is red
			start:    3,
			end:      6,
			expected: "\x1b[31m345\x1b[0m",
		},
		{
			name:     "shift with styles, partial start",
			input:    "012\x1b[31m345\x1b[0m6789", // "345" is red
			start:    4,
			end:      7,
			expected: "\x1b[31m45\x1b[0m6",
		},
		{
			name:     "shift with styles, partial end",
			input:    "012\x1b[31m345\x1b[0m6789", // "345" is red
			start:    2,
			end:      5,
			expected: "2\x1b[31m34\x1b[0m",
		},
		{
			name:     "shift with styles, spanning multiple styles, reset at end",
			input:    "012\x1b[31m345\x1b[0m67\x1b[32m89\x1b[0m", // "345" is red, "89" is green
			start:    4,
			end:      9,
			expected: "\x1b[31m45\x1b[0m67\x1b[32m8\x1b[0m",
		},
		{
			name:     "shift with styles, spanning multiple styles, no reset at end",
			input:    "012\x1b[31m345\x1b[0m67\x1b[32m89", // "345" is red, "89" is green
			start:    4,
			end:      9,
			expected: "\x1b[31m45\x1b[0m67\x1b[32m8\x1b[0m",
		},
		// with wide characters
		{
			name:     "shift with wide characters, no styles",
			input:    "01二三四5六789",
			start:    3,
			end:      6,
			expected: "三四5",
		},
		{
			name:     "shift wide characters with styles",
			input:    "01二\x1b[31m三四5\x1b[0m六789", // "345" is red
			start:    3,
			end:      6,
			expected: "\x1b[31m三四5\x1b[0m",
		},
		{
			name:     "shift wide characters with styles, partial start",
			input:    "01二\x1b[31m三四5\x1b[0m六789", // "345" is red
			start:    4,
			end:      7,
			expected: "\x1b[31m四5\x1b[0m六",
		},
		{
			name:     "shift wide characters with styles, partial end",
			input:    "01二\x1b[31m三四5\x1b[0m六789", // "345" is red
			start:    2,
			end:      5,
			expected: "二\x1b[31m三四\x1b[0m",
		},
		{
			name:     "shift wide characters with styles, spanning multiple styles, reset at end",
			input:    "01二\x1b[31m三四5\x1b[0m六7\x1b[32m89\x1b[0m", // "345" is red, "89" is green
			start:    4,
			end:      9,
			expected: "\x1b[31m四5\x1b[0m六7\x1b[32m8\x1b[0m",
		},
		{
			name:     "shift wide characters with styles, spanning multiple styles, no reset at end",
			input:    "01二\x1b[31m三四5\x1b[0m六7\x1b[32m89", // "345" is red, "89" is green
			start:    4,
			end:      9,
			expected: "\x1b[31m四5\x1b[0m六7\x1b[32m8\x1b[0m",
		},
		{
			name:     "shift exceeding length",
			input:    "0123456789",
			start:    8,
			end:      15, // width = 7
			expected: "3456789",
		},
		{
			name:     "width > length",
			input:    "0123456789",
			start:    7,
			end:      50, // width = 43
			expected: "0123456789",
		},
		{
			name:     "start > length",
			input:    "0123456789",
			start:    15,
			end:      20, // width = 5
			expected: "56789",
		},
		{
			name:     "start is negative",
			input:    "0123456789",
			start:    -5,
			end:      3, // width = 8
			expected: "01234567",
		},
		{
			name:     "width is negative",
			input:    "0123456789",
			start:    5,
			end:      3, // width = -2
			expected: "5",
		},
		{
			name:     "empty input",
			input:    "",
			start:    0,
			end:      5,
			expected: "",
		},
		{
			name:     "empty input with styles",
			input:    "\x1b[31m\x1b[0m",
			start:    0,
			end:      5,
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledShift(tc.start, tc.end-tc.start)
			assert.Equal(t, tc.expected, got)
		})
	}
}
