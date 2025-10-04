// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Width(t *testing.T) {
	for _, tc := range onelineTestCase {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEntry(tc.input)
			got := e.Width()
			assert.Equal(t, tc.width, got)
		})
	}
}
