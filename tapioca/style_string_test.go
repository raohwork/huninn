// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyle_String(t *testing.T) {
	cases := []struct {
		name     string
		style    *style
		expected string
	}{
		{
			name:     "nil style",
			style:    nil,
			expected: "",
		},
		{
			name:     "empty style",
			style:    &style{},
			expected: "",
		},
		{
			name:     "bold only",
			style:    &style{bold: true},
			expected: "\x1b[1m",
		},
		{
			name:     "color only",
			style:    &style{fg: "31"},
			expected: "\x1b[31m",
		},
		{
			name:     "multiple attributes",
			style:    &style{bold: true, fg: "31", bg: "42"},
			expected: "\x1b[31m\x1b[42m\x1b[1m",
		},
		{
			name: "all attributes",
			style: &style{
				fg: "31", bg: "42", bold: true, faint: true,
				italic: true, underline: true, blink: true,
				reverse: true, hidden: true, strike: true,
			},
			expected: "\x1b[31m\x1b[42m\x1b[1m\x1b[2m\x1b[3m\x1b[4m\x1b[5m\x1b[7m\x1b[8m\x1b[9m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.String()
			assert.Equal(t, tc.expected, got)
		})
	}
}
