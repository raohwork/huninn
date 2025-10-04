// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyle_Parse(t *testing.T) {
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
			expected: nil,
		},
		{
			name:     "reset code simplified",
			code:     "\x1b[m",
			prev:     &style{bold: true, fg: "31"},
			expected: nil,
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
			name:     "turn off bold and faint",
			code:     "\x1b[22m",
			prev:     &style{bold: true, faint: true},
			expected: nil,
		},
		{
			name:     "turn off bold and faint while other styles remain",
			code:     "\x1b[22m",
			prev:     &style{bold: true, faint: true, italic: true},
			expected: &style{italic: true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAnsiCode(tc.code, tc.prev)
			assert.Equal(t, tc.expected, got)
		})
	}
}
