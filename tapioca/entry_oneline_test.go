// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// max 10 characters
var onelineTestCase = []struct {
	name     string
	input    string
	expected string
	width    int
}{
	{
		name:     "no style",
		input:    "simple",
		expected: "simple",
		width:    6,
	},
	{
		name:     "simple style (bold)",
		input:    "\x1b[1mbold",
		expected: "\x1b[1mbold\x1b[0m",
		width:    4,
	},
	{
		name:     "simple style with reset (bold)",
		input:    "\x1b[1mbold\x1b[0mnormal",
		expected: "\x1b[1mbold\x1b[0mnormal",
		width:    10,
	},
	{
		name:     "simple style with short reset (bold)",
		input:    "\x1b[1mbold\x1b[mnormal",
		expected: "\x1b[1mbold\x1b[0mnormal",
		width:    10,
	},
	{
		name:     "nested style (bold + red), need reset at end",
		input:    "\x1b[1mbold \x1b[31mred",
		expected: "\x1b[1mbold \x1b[0m\x1b[31m\x1b[1mred\x1b[0m",
		width:    8,
	},
	{
		name:     "nested style with reset (bold + red), no reset at end",
		input:    "\x1b[1mbold\x1b[31mred\x1b[0mnor",
		expected: "\x1b[1mbold\x1b[0m\x1b[31m\x1b[1mred\x1b[0mnor",
		width:    10,
	},
	{
		name: "complex nested style with resets",
		// n for normal, b for blue (fg 34), r for red (fg 31), uppercase for bold
		// nNbBrRn
		input: "n" +
			"\x1b[1mN" +
			"\x1b[0;34mb" +
			"\x1b[1;34mB" +
			"\x1b[0;31mr" +
			"\x1b[1m\x1b[31mR" +
			"\x1b[0mn",
		expected: "n" +
			"\x1b[1mN" +
			"\x1b[0m\x1b[34mb" +
			"\x1b[0m\x1b[34m\x1b[1mB" +
			"\x1b[0m\x1b[31mr" +
			"\x1b[0m\x1b[31m\x1b[1mR" +
			"\x1b[0mn",
		width: 7,
	},
}

func TestEntry_StyledString(t *testing.T) {
	for _, tc := range onelineTestCase {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledString()
			assert.Equal(t, tc.expected, got)
		})
	}
}
