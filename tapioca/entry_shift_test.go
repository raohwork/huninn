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
		// with wide characters
		{
			name:     "[wide] shift with wide characters, no styles",
			input:    "01三五67九", // position 2-3, 4-5 and 8-9 is wide characters
			start:    2,
			end:      6,
			expected: "三五",
		},
		{
			name:     "[wide] begin at partial wide character",
			input:    "01三五67九",
			start:    3,
			end:      6,
			expected: " 五",
		},
		{
			name:     "[wide] end at partial wide character",
			input:    "01三五67九",
			start:    2,
			end:      5,
			expected: "三 ",
		},
		{
			name:  "[wide] shift with wide characters and styles",
			input: "01\x1b[31m三五\x1b[0m67\x1b[32m九\x1b[0m", // "三五" is red, "九" is green
			//      01        2345       67        89
			start:    2,
			end:      7,
			expected: "\x1b[31m三五\x1b[0m6",
		},
		{
			name:     "[wide] shift with wide characters and styles, partial start",
			input:    "01\x1b[31m三五\x1b[0m67\x1b[32m九\x1b[0m", // "三五" is red, "九" is green
			start:    3,
			end:      9,
			expected: " \x1b[31m五\x1b[0m67 ",
		},
		{
			name:  "[wide] shift with wide characters and styles, partial start without reset at end",
			input: "01\x1b[31m三五\x1b[0m67\x1b[32m九", // "三五" is red, "九" is green
			//      01        2345       67        89
			start:    3,
			end:      9,
			expected: " \x1b[31m五\x1b[0m67 ",
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

func TestEntry_ComputeStartAndEndForShift(t *testing.T) {
	cases := []struct {
		name              string
		input             string
		startCol          int
		width             int
		expectedStartIdx  int
		expectedEndIdx    int
		expectedHasPrefix bool
		expectedHasSuffix bool
	}{
		{
			name:  "width exceeds length",
			input: "A一B二C三D",
			//      0123456789  position
			//      0 12 34 56  rune index
			startCol: 7,
			width:    4,
			// expected to show "三D"
			expectedStartIdx:  5,
			expectedEndIdx:    7,
			expectedHasPrefix: false,
			expectedHasSuffix: false,
		},
		{
			name:  "normal chars",
			input: "01三五七89", // (7 runes, 10 characters wide)
			//      0123456789  position
			//      01 2 3 456  rune index
			// expected to show "1三五七8"
			startCol:          1,
			width:             8,
			expectedStartIdx:  1, // 1
			expectedEndIdx:    6, // 9
			expectedHasPrefix: false,
			expectedHasSuffix: false,
		},
		{
			name:  "no cut",
			input: "01三五七89", // (7 runes, 10 characters wide)
			//      0123456789  position
			//      01 2 3 456  rune index
			// expected to show "三五七"
			startCol:          2,
			width:             6,
			expectedStartIdx:  2, // 三
			expectedEndIdx:    5, // 8
			expectedHasPrefix: false,
			expectedHasSuffix: false,
		},
		{
			name:  "cut both ends",
			input: "01三五七89", // (7 runes, 10 characters wide)
			// 	0123456789  position
			//      01 2 3 456  rune index
			// expected to show " 五 "
			startCol:          3,
			width:             4,
			expectedStartIdx:  3, // 五
			expectedEndIdx:    4, // 七
			expectedHasPrefix: true,
			expectedHasSuffix: true,
		},
		{
			name:  "cut at start",
			input: "01三五七89", // (7 runes, 10 characters wide)
			//      0123456789  position
			//      01 2 3 456  rune index
			// expected to show " 五七"
			startCol:          3,
			width:             5,
			expectedStartIdx:  3, // 五
			expectedEndIdx:    5, // 8
			expectedHasPrefix: true,
			expectedHasSuffix: false,
		},
		{
			name:  "cut at end",
			input: "01三五七89", // (7 runes, 10 characters wide)
			//      0123456789  position
			//      01 2 3 456  rune index
			// expected to show "三五 "
			startCol:          2,
			width:             5,
			expectedStartIdx:  2, // 三
			expectedEndIdx:    4, // 七
			expectedHasPrefix: false,
			expectedHasSuffix: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			offsets := entry.runeEndOffsets()
			startIdx, endIdx, hasPrefix, hasSuffix := computeStartAndEndForShift(offsets, tc.startCol, tc.width)
			assert.Equal(t, tc.expectedStartIdx, startIdx, "startIdx mismatch")
			assert.Equal(t, tc.expectedEndIdx, endIdx, "endIdx mismatch")
			assert.Equal(t, tc.expectedHasPrefix, hasPrefix, "hasPrefix mismatch")
			assert.Equal(t, tc.expectedHasSuffix, hasSuffix, "hasSuffix mismatch")
		})
	}
}
