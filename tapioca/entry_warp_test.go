// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_StyledWarps(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		width    int
		expected []string
	}{
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledWarps(tc.width)
			assert.Equal(t, tc.expected, got)
		})
	}

}
