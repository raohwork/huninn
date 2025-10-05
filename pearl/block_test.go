// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	cases := []struct {
		width, height int
	}{
		{1, 1},
		{2, 1},
		{1, 2},
		{2, 2},
		{10, 3},
		{3, 10},
		{10, 30},
		{30, 10},
	}

	str := strings.Repeat("A", 100)

	for _, c := range cases {
		t.Run(fmt.Sprintf("%dx%d", c.width, c.height), func(t *testing.T) {
			b := NewBlock()
			sender := func(msg tea.Msg) {
				b.Update(msg)
			}
			setter := b.Setter(sender)
			setter(str)

			assert.Equal(t, "", tapioca.IsThisTopping(tapioca.ToppingTestSpec{
				Width:  c.width,
				Height: c.height,
				Model:  b,
			}))
		})
	}
}
