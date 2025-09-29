// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func arr[T any](values ...T) []T { return values }

func TestEntryWarpPositions(t *testing.T) {
	cases := []struct {
		name   string
		line   string
		width  int
		expect []int
		lines  int
	}{
		{
			name:   "short line",
			line:   "Hello, World!",
			width:  20,
			expect: arr(0),
			lines:  1,
		},
		{
			name:   "exact fit",
			line:   "Hello, World!",
			width:  13,
			expect: arr(0),
			lines:  1,
		},
		{
			name:   "multi line",
			line:   "Hello, World!",
			width:  5,
			expect: arr(0, 5, 10),
			lines:  3,
		},
		{
			name:   "wide characters",
			line:   "こんにちは、世界！",
			width:  10,
			expect: arr(0, 5),
			lines:  2,
		},
		{
			name:   "cut at wide character",
			line:   "こんにちは",
			width:  5,
			expect: arr(0, 2, 4),
			lines:  3,
		},
		{
			name:   "extreme narrow",
			line:   "Hello",
			width:  2,
			expect: arr(0, 2, 4),
			lines:  3,
		},
		{
			name:   "extreme narrow with wide chars",
			line:   "你好",
			width:  2,
			expect: arr(0, 1),
			lines:  2,
		},
		{
			name:   "incorrect width (<2)",
			line:   "hello, 你好",
			width:  1,
			expect: arr(0, 2, 4, 6, 7, 8),
			lines:  6,
		},
		{
			name:   "empty string",
			line:   "",
			width:  10,
			expect: nil,
			lines:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEntry(tc.line)
			got := e.warpPositions(tc.width)
			assert.Equal(t, tc.expect, got, "entry: %q, width: %d", tc.line, tc.width)
			lines := e.Lines(tc.width)
			assert.Equal(t, tc.lines, lines, "entry: %q, width: %d", tc.line, tc.width)
		})
	}
}

func TestEntryWarps(t *testing.T) {
	cases := []struct {
		name  string
		line  string
		width int
		want  []string
	}{
		{
			name:  "simple",
			line:  "Hello, World!",
			width: 10,
			want:  []string{"Hello, Wor", "ld!"},
		},
		{
			name:  "with spaces",
			line:  "Hello,    World! This is a test.",
			width: 13,
			want:  []string{"Hello,    Wor", "ld! This is a", " test."},
		},
		{
			name:  "wide characters",
			line:  "こんにちは、世界！",
			width: 10,
			want:  []string{"こんにちは", "、世界！"},
		},
		{
			name:  "cut at wide character",
			line:  "こんにちは",
			width: 7,
			want:  []string{"こんに", "ちは"},
		},
		{
			name:  "mixed characters",
			line:  "Hello, こんにちは World!",
			width: 15,
			want:  []string{"Hello, こんにち", "は World!"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewEntry(tc.line)
			got := m.Warps(tc.width)
			assert.Equal(t, tc.want, got)
		})
	}
}
