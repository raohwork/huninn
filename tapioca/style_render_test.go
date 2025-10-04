// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyle_Render(t *testing.T) {
	cases := []struct {
		name      string
		prev, cur *style
		expected  string
	}{
		{
			name:     "nil+nil",
			prev:     nil,
			cur:      nil,
			expected: "",
		},
		{
			name:     "empty+nil",
			prev:     &style{},
			cur:      nil,
			expected: "",
		},
		{
			name:     "nil+empty",
			prev:     nil,
			cur:      &style{},
			expected: "",
		},
		{
			name:     "empty+empty",
			prev:     &style{},
			cur:      &style{},
			expected: "",
		},
		{
			name:     "bold+nil",
			prev:     &style{bold: true},
			cur:      nil,
			expected: "\x1b[0m",
		},
		{
			name:     "bold+empty",
			prev:     &style{bold: true},
			cur:      nil,
			expected: "\x1b[0m",
		},
		{
			name:     "nil+bold",
			prev:     nil,
			cur:      &style{bold: true},
			expected: "\x1b[1m",
		},
		{
			name:     "empty+bold",
			prev:     nil,
			cur:      &style{bold: true},
			expected: "\x1b[1m",
		},
		{
			name:     "bold+bold",
			prev:     &style{bold: true},
			cur:      &style{bold: true},
			expected: "\x1b[0m\x1b[1m",
		},
		{
			name:     "bold+italic",
			prev:     &style{bold: true},
			cur:      &style{italic: true},
			expected: "\x1b[0m\x1b[3m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cur.Render(tc.prev)
			assert.Equal(t, tc.expected, got)
		})
	}
}
