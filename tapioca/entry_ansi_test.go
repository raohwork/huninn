// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEntryWithAnsiStyles(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string // Expected plain text output
	}{
		{
			name:     "no ansi codes",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "simple color code",
			input:    "Hello, \x1b[31mRed\x1b[0m World!",
			expected: "Hello, Red World!",
		},
		{
			name:     "multiple styles",
			input:    "\x1b[1;32mBold Green\x1b[0m Normal",
			expected: "Bold Green Normal",
		},
		{
			name:     "unsupported CSI sequences removed",
			input:    "Hello\x1b[2J\x1b[31mRed\x1b[H World",
			expected: "HelloRed World",
		},
		{
			name:     "256 color",
			input:    "\x1b[38;5;196mBright Red\x1b[0m",
			expected: "Bright Red",
		},
		{
			name:     "RGB color",
			input:    "\x1b[38;2;255;0;0mRGB Red\x1b[0m",
			expected: "RGB Red",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.String()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestStyledLines(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		width    int
		expected []string
	}{
		{
			name:  "style attribute reset issue",
			input: "\x1b[33m\x1b[3mabc\x1b[23mdef",
			width: 20,
			expected: []string{
				"\x1b[0m\x1b[33m\x1b[3mabc\x1b[0m\x1b[33mdef\x1b[m",
			},
		},
		{
			name:  "simple styled text",
			input: "\x1b[31mRed\x1b[0m Normal",
			width: 20,
			expected: []string{
				"\x1b[0m\x1b[31mRed\x1b[0m Normal\x1b[m",
			},
		},
		{
			name:  "styled text with wrapping",
			input: "\x1b[1mBold\x1b[0m Normal Text",
			width: 8,
			expected: []string{
				"\x1b[0m\x1b[1mBold\x1b[0m Nor\x1b[m",
				"\x1b[0mmal Text\x1b[m",
			},
		},
		{
			name:  "multiple styles across lines",
			input: "\x1b[31mRed\x1b[32mGreen\x1b[0mNormal",
			width: 6,
			expected: []string{
				"\x1b[0m\x1b[31mRed\x1b[0m\x1b[32mGre\x1b[m",
				"\x1b[0m\x1b[32men\x1b[0mNorm\x1b[m",
				"\x1b[0mal\x1b[m",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledLines(tc.width)
			assert.Equal(t, tc.expected, got)
			t.Log("Expected:")
			for i, line := range tc.expected {
				t.Logf("  Line %d: %s", i, line)
			}
			t.Log("Got:")
			for i, line := range got {
				t.Logf("  Line %d: %s", i, line)
			}
		})
	}
}

func TestStyledShift(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		startCol int
		width    int
		expected string
	}{
		{
			name:     "simple styled shift",
			input:    "\x1b[31mRed\x1b[0m Normal",
			startCol: 0,
			width:    6,
			expected: "\x1b[0m\x1b[31mRed\x1b[0m No\x1b[m",
		},
		{
			name:     "shift with style boundary",
			input:    "\x1b[31mRed\x1b[0m Normal",
			startCol: 2,
			width:    6,
			expected: "\x1b[0m\x1b[31md\x1b[0m Norm\x1b[m",
		},
		{
			name:     "shift beyond styled content",
			input:    "\x1b[31mRed\x1b[0m",
			startCol: 5,
			width:    5,
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := NewEntry(tc.input)
			got := entry.StyledShift(tc.startCol, tc.width)
			assert.Equal(t, tc.expected, got)
			t.Logf("Expected: %s", tc.expected)
			t.Logf("Got:      %s", got)
		})
	}
}

func TestParseAnsiCode(t *testing.T) {
	defaultStyle := &style{}

	cases := []struct {
		name     string
		code     string
		prev     *style
		expected *style
	}{
		{
			name:     "reset code",
			code:     "\x1b[0m",
			prev:     &style{bold: true, fg: "31"},
			expected: &style{},
		},
		{
			name:     "reset code simplified",
			code:     "\x1b[m",
			prev:     &style{bold: true, fg: "31"},
			expected: &style{},
		},
		{
			name:     "bold code",
			code:     "\x1b[1m",
			prev:     defaultStyle,
			expected: &style{bold: true},
		},
		{
			name:     "color code",
			code:     "\x1b[31m",
			prev:     defaultStyle,
			expected: &style{fg: "31"},
		},
		{
			name:     "multiple codes",
			code:     "\x1b[1;31m",
			prev:     defaultStyle,
			expected: &style{bold: true, fg: "31"},
		},
		{
			name:     "long multiple codes",
			code:     "\x1b[1;31;46m",
			prev:     defaultStyle,
			expected: &style{bold: true, fg: "31", bg: "46"},
		},
		{
			name:     "256 color",
			code:     "\x1b[38;5;196m",
			prev:     defaultStyle,
			expected: &style{fg: "38;5;196"},
		},
		{
			name:     "RGB color",
			code:     "\x1b[38;2;255;128;0m",
			prev:     defaultStyle,
			expected: &style{fg: "38;2;255;128;0"},
		},
		{
			name:     "turn off bold",
			code:     "\x1b[22m",
			prev:     &style{bold: true, faint: true},
			expected: &style{bold: false, faint: false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAnsiCode(tc.code, tc.prev)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestStyleString(t *testing.T) {
	cases := []struct {
		name     string
		style    *style
		expected string
	}{
		{
			name:     "nil style",
			style:    nil,
			expected: "\x1b[m",
		},
		{
			name:     "empty style",
			style:    &style{},
			expected: "\x1b[0m",
		},
		{
			name:     "bold only",
			style:    &style{bold: true},
			expected: "\x1b[0m\x1b[1m",
		},
		{
			name:     "color only",
			style:    &style{fg: "31"},
			expected: "\x1b[0m\x1b[31m",
		},
		{
			name:     "multiple attributes",
			style:    &style{bold: true, fg: "31", bg: "42"},
			expected: "\x1b[0m\x1b[31m\x1b[42m\x1b[1m",
		},
		{
			name: "all attributes",
			style: &style{
				fg: "31", bg: "42", bold: true, faint: true,
				italic: true, underline: true, blink: true,
				reverse: true, hidden: true, strike: true,
			},
			expected: "\x1b[0m\x1b[31m\x1b[42m\x1b[1m\x1b[2m\x1b[3m\x1b[4m\x1b[5m\x1b[7m\x1b[8m\x1b[9m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.String()
			assert.Equal(t, tc.expected, got)
		})
	}
}
