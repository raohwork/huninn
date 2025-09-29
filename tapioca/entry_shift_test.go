// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryShift(t *testing.T) {
	cases := []struct {
		name  string
		line  string
		begin int
		width int
		want  string
	}{
		{
			name:  "simple substring",
			line:  "Hello, world!",
			begin: 7,
			width: 5,
			want:  "world",
		},
		{
			name:  "start from beginning",
			line:  "Hello, world!",
			begin: 0,
			width: 5,
			want:  "Hello",
		},
		{
			name:  "width larger than remaining",
			line:  "Hello",
			begin: 2,
			width: 10,
			want:  "llo",
		},
		{
			name:  "begin at end",
			line:  "Hello",
			begin: 5,
			width: 3,
			want:  "",
		},
		{
			name:  "begin beyond end",
			line:  "Hello",
			begin: 10,
			width: 3,
			want:  "",
		},
		{
			name:  "empty string",
			line:  "",
			begin: 0,
			width: 5,
			want:  "",
		},
		{
			name:  "zero width",
			line:  "Hello",
			begin: 2,
			width: 0,
			want:  "ll",
		},
		{
			name:  "wide characters simple",
			line:  "世界",
			begin: 0,
			width: 2,
			want:  "世",
		},
		{
			name:  "wide characters full",
			line:  "世界",
			begin: 0,
			width: 4,
			want:  "世界",
		},
		{
			name:  "wide characters shifted by 1 - missing first half",
			line:  "世界",
			begin: 1,
			width: 10,
			want:  " 界",
		},
		{
			name:  "wide characters cut in middle",
			line:  "你好，世界",
			begin: 2,
			width: 5,
			want:  "好， ",
		},
		{
			name:  "wide characters shifted by 2 - second character",
			line:  "世界",
			begin: 2,
			width: 10,
			want:  "界",
		},
		{
			name:  "wide characters shifted by 3 - missing second half",
			line:  "世界",
			begin: 3,
			width: 10,
			want:  " ",
		},
		{
			name:  "mixed characters",
			line:  "Hello世界",
			begin: 5,
			width: 4,
			want:  "世界",
		},
		{
			name:  "mixed characters - cut at wide char",
			line:  "Hello世界",
			begin: 6,
			width: 10,
			want:  " 界",
		},
		{
			name:  "complex example from doc",
			line:  "Hello, world!",
			begin: 7,
			width: 5,
			want:  "world",
		},
		{
			name:  "wide char example from doc",
			line:  "世界",
			begin: 1,
			width: 10,
			want:  " 界",
		},
		{
			name:  "three wide characters",
			line:  "你好世",
			begin: 1,
			width: 10,
			want:  " 好世",
		},
		{
			name:  "three wide characters middle cut",
			line:  "你好世",
			begin: 3,
			width: 10,
			want:  " 世",
		},
		{
			name:  "ascii and wide mixed complex",
			line:  "a世b界c",
			begin: 2,
			width: 5,
			want:  " b界c",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEntry(tc.line)
			got := e.Shift(tc.begin, tc.width)
			// trim trailing spaces as it does not affect display
			// and makes testing easier
			got = strings.TrimRight(got, " ")
			want := strings.TrimRight(tc.want, " ")
			assert.Equal(t, want, got, "entry: %q, begin: %d, width: %d", tc.line, tc.begin, tc.width)
		})
	}
}
