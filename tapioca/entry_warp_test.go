// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var warpCases = []struct {
	name     string
	input    string
	width    int
	expected []string
}{
	{
		name:  "edge case: width greater than string length, with wide characters",
		input: "Hello, World! 你好，世界！",
		width: 90,
		expected: []string{
			"Hello, World! 你好，世界！",
		},
	},
	{
		name:  "edge case: zero width",
		input: "Hello, World!",
		width: 0,
		expected: []string{
			"H", "e", "l", "l", "o", ",", " ", "W", "o", "r", "l", "d", "!",
		},
	},
	{
		name:  "edge case: negative width",
		input: "Hello, World!",
		width: -1,
		expected: []string{
			"H", "e", "l", "l", "o", ",", " ", "W", "o", "r", "l", "d", "!",
		},
	},
	{
		name:  "edge case: width 1 with wide characters",
		input: "一二三",
		width: 1,
		expected: []string{
			"一", "二", "三",
		},
	},
	{
		name:  "edge case: width 1 with wide characters, with style",
		input: "一\x1b[31m二\x1b[m三",
		width: 1,
		expected: []string{
			"一", "\x1b[31m二\x1b[0m", "三",
		},
	},
	{
		name:  "no style, no warps",
		input: "Hello, World!",
		width: 20,
		expected: []string{
			"Hello, World!",
		},
	},
	{
		name:  "no style, with warps",
		input: "Hello, World! This is a test of the wrapping function.",
		width: 20,
		expected: []string{
			"Hello, World! This i",
			"s a test of the wrap",
			"ping function.",
		},
	},
	{
		name:  "with style, no warps",
		input: "\x1b[31mHello, World!\x1b[0m", // Red text
		width: 20,
		expected: []string{
			"\x1b[31mHello, World!\x1b[0m",
		},
	},
	{
		name:  "with style, with warps, with reset",
		input: "\x1b[32mHello, World! This is a test of the wrapping function.\x1b[0m", // Green text
		width: 20,
		expected: []string{
			"\x1b[32mHello, World! This i\x1b[0m",
			"\x1b[32ms a test of the wrap\x1b[0m",
			"\x1b[32mping function.\x1b[0m",
		},
	},
	{
		name:  "with style, with warps, without reset",
		input: "\x1b[32mHello, World! This is a test of the wrapping function.", // Green text
		width: 20,
		expected: []string{
			"\x1b[32mHello, World! This i\x1b[0m",
			"\x1b[32ms a test of the wrap\x1b[0m",
			"\x1b[32mping function.\x1b[0m",
		},
	},
	{
		name:  "with multiple styles, with warps",
		input: "\x1b[31mRed\x1b[34mBlue\x1b[0mNormal\x1b[32mGreen",
		width: 5,
		expected: []string{
			"\x1b[31mRed\x1b[0m\x1b[34mBl\x1b[0m",
			"\x1b[34mue\x1b[0mNor",
			"mal\x1b[32mGr\x1b[0m",
			"\x1b[32meen\x1b[0m",
		},
	},
	// wide characters
	{
		name:  "wide characters, with warps, no cut",
		input: "一二三四",
		width: 4,
		expected: []string{
			"一二",
			"三四",
		},
	},
	{
		name:  "wide characters, with warps, with cut",
		input: "一二三",
		width: 3,
		expected: []string{
			"一 ",
			"二 ",
			"三",
		},
	},
	{
		name:  "mixed characters, with warps, with cut",
		input: "A一B二C三D",
		width: 4,
		expected: []string{
			"A一B",
			"二C ",
			"三D",
		},
	},
	{
		name:  "mixed characters and styles, with warps, with cut",
		input: "\x1b[31mA一\x1b[34mB二C\x1b[m三D",
		width: 4,
		expected: []string{
			"\x1b[31mA一\x1b[0m\x1b[34mB\x1b[0m",
			"\x1b[34m二C\x1b[0m ",
			"三D",
		},
	},
}

func TestEntry_StyledWarps(t *testing.T) {
	for _, tc := range warpCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledWarps(tc.width)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestEntry_StyledBlock(t *testing.T) {
	for _, tc := range warpCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledBlock(tc.width)

			assert.Equal(t, len(tc.expected), len(got), "number of lines")
			for i := range tc.expected[:len(got)-1] {
				exp := NewEntry(tc.expected[i])
				act := NewEntry(got[i])
				assert.Equal(t, exp.StyledString(), act.StyledString(), "line#%d content", i)
				if tc.width > 1 {
					assert.Equal(t, tc.width, act.Width(), "line#%d width", i)
				} else {
					assert.Equal(t, exp.Width(), act.Width(), "line#%d width", i)
				}
			}

			// handle last line
			lastIdx := len(got) - 1
			assert.Equal(
				t,
				strings.TrimRight(tc.expected[lastIdx], " "),
				strings.TrimRight(got[lastIdx], " "),
				"last line content",
			)
			if tc.width > 1 {
				assert.Equal(t, tc.width, NewEntry(got[lastIdx]).Width(), "last line width")
			} else {
				assert.Equal(t, NewEntry(tc.expected[lastIdx]).Width(), NewEntry(got[lastIdx]).Width(), "last line width")
			}
		})
	}
}

func TestEntry_SubstringWidth(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		start    int
		end      int
		expected int
	}{
		{
			name:     "simple ASCII",
			input:    "Hello, World!",
			start:    0,
			end:      5,
			expected: 5,
		},
		{
			name:     "with style codes",
			input:    "\x1b[31mHello\x1b[0m, World!",
			start:    0,
			end:      5,
			expected: 5,
		},
		{
			name:     "wide characters",
			input:    "你好，世界！",
			start:    0,
			end:      2,
			expected: 4, // Each Chinese character is width 2
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.substringWidth(tc.start, tc.end)
			assert.Equal(t, tc.expected, got)
		})
	}
}
