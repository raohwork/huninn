// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
	"github.com/stretchr/testify/assert"
)

func TestTaskList(t *testing.T) {
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

	for _, c := range cases {
		t.Run(fmt.Sprintf("%dx%d", c.width, c.height), func(t *testing.T) {
			l := NewTaskList()
			sender := func(msg tea.Msg) {
				l.Update(msg)
			}
			m := l.CreateManager(sender)
			m.AddTask("task1", "")
			m.AddTask("task2", "")
			m.AddTask("task3333333333333333333333333", "")
			m.AddTask("", "")

			assert.Equal(t, "", tapioca.IsThisTopping(tapioca.ToppingTestSpec{
				Width:  c.width,
				Height: c.height,
				Model:  l,
			}))
		})
	}
}
