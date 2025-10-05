// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestComponent_BasicRendering(t *testing.T) {
	// the component is resized to 10x3 by default

	cases := []struct {
		name     string
		msgs     []tea.Msg
		entries  []string
		hScroll  bool
		vScroll  bool
		expected []string
	}{
		{
			name:     "no scrolling, short string",
			msgs:     []tea.Msg{},
			entries:  []string{"Hello"},
			expected: []string{"Hello     ", "          ", "          "},
		},
		{
			name:     "no scrolling, long string",
			msgs:     []tea.Msg{},
			entries:  []string{"Hello, World!"},
			expected: []string{"Hello, Wor", "ld!       ", "          "},
		},
		{
			name:     "no scrolling, super long string",
			msgs:     []tea.Msg{},
			entries:  []string{"Hello, World! This is a test."},
			expected: []string{"Hello, Wor", "ld! This i", "s a test. "},
		},
		// horizontal scrolling
		{
			name:     "horizontal scrolling, long string",
			msgs:     []tea.Msg{ScrollRightMsg(1)},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"ello      ", "ello, Worl", "orld      "},
		},
		{
			name:     "horizontal scrolling, long string, scroll exactly to the end",
			msgs:     []tea.Msg{ScrollRightMsg(3)},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"lo        ", "lo, World!", "ld        "},
		},
		{
			name:     "horizontal scrolling, long string, scroll beyond the end",
			msgs:     []tea.Msg{ScrollRightMsg(5)},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"lo        ", "lo, World!", "ld        "},
		},
		{
			name:     "horizontal scrolling to end",
			msgs:     []tea.Msg{ScrollEndMsg{}},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"lo        ", "lo, World!", "ld        "},
		},
		{
			name:     "horizontal scrolling to end, then back 1 char",
			msgs:     []tea.Msg{ScrollEndMsg{}, ScrollLeftMsg(1)},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"llo       ", "llo, World", "rld       "},
		},
		{
			name:     "horizontal scrolling to end, then back beyond start",
			msgs:     []tea.Msg{ScrollEndMsg{}, ScrollLeftMsg(10)},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"Hello     ", "Hello, Wor", "World     "},
		},
		{
			name:     "horizontal scrolling to end, then scroll to beginning",
			msgs:     []tea.Msg{ScrollEndMsg{}, ScrollBeginMsg{}},
			entries:  []string{"Hello", "Hello, World!", "World"},
			hScroll:  true,
			expected: []string{"Hello     ", "Hello, Wor", "World     "},
		},
		// vertical scrolling with short lines
		{
			name:     "vertical scrolling, multiple short entries",
			msgs:     []tea.Msg{ScrollDownMsg(1)},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"Two       ", "Three     ", "Four      "},
		},
		{
			name:     "vertical scrolling, multiple short entries, scroll exactly to the end",
			msgs:     []tea.Msg{ScrollDownMsg(2)},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"Three     ", "Four      ", "Five      "},
		},
		{
			name:     "vertical scrolling, multiple short entries, scroll beyond the end",
			msgs:     []tea.Msg{ScrollDownMsg(5)},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"Three     ", "Four      ", "Five      "},
		},
		{
			name:     "vertical scrolling to end",
			msgs:     []tea.Msg{ScrollBottomMsg{}},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"Three     ", "Four      ", "Five      "},
		},
		{
			name:     "vertical scrolling to end, then back 1 entry",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(1)},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"Two       ", "Three     ", "Four      "},
		},
		{
			name:     "vertical scrolling to end, then back beyond start",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(10)},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"One       ", "Two       ", "Three     "},
		},
		{
			name:     "vertical scrolling to end, then scroll to top",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollTopMsg{}},
			entries:  []string{"One", "Two", "Three", "Four", "Five"},
			vScroll:  true,
			expected: []string{"One       ", "Two       ", "Three     "},
		},
		// vertical scrolling with long lines and warping
		{
			name:     "vertical scrolling, multiple long entries with wrapping",
			msgs:     []tea.Msg{ScrollDownMsg(1)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"TwoTwo    ", "ThreeThree", "Three     "},
		},
		{
			name:     "vertical scrolling, multiple long entries with wrapping, scroll exactly to the end",
			msgs:     []tea.Msg{ScrollDownMsg(3)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"Three     ", "FourFourFo", "urFour    "},
		},
		{
			name:     "vertical scrolling, multiple long entries with wrapping, scroll beyond the end",
			msgs:     []tea.Msg{ScrollDownMsg(10)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"Three     ", "FourFourFo", "urFour    "},
		},
		{
			name:     "vertical scrolling to end with wrapping",
			msgs:     []tea.Msg{ScrollBottomMsg{}},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"Three     ", "FourFourFo", "urFour    "},
		},
		{
			name:     "vertical scrolling to end with wrapping, then back 1 entry",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(1)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"ThreeThree", "Three     ", "FourFourFo"},
		},
		{
			name:     "vertical scrolling to end with wrapping, then back beyond start",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(10)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"One       ", "TwoTwo    ", "ThreeThree"},
		},
		{
			name:     "vertical scrolling to end with wrapping, then scroll to beginning",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollTopMsg{}},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll:  true,
			expected: []string{"One       ", "TwoTwo    ", "ThreeThree"},
		},
		// vertical scrolling with long lines and no warping
		{
			name:     "vertical scrolling, multiple long entries without wrapping",
			msgs:     []tea.Msg{ScrollDownMsg(1)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"TwoTwo    ", "ThreeThree", "FourFourFo"},
		},
		{
			name:     "vertical scrolling, multiple long entries without wrapping, scroll exactly to the end",
			msgs:     []tea.Msg{ScrollDownMsg(2)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"ThreeThree", "FourFourFo", "FiveFiveFi"},
		},
		{
			name:     "vertical scrolling, multiple long entries without wrapping, scroll beyond the end",
			msgs:     []tea.Msg{ScrollDownMsg(5)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"ThreeThree", "FourFourFo", "FiveFiveFi"},
		},
		{
			name:     "vertical scrolling to end with no wrapping",
			msgs:     []tea.Msg{ScrollBottomMsg{}},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"ThreeThree", "FourFourFo", "FiveFiveFi"},
		},
		{
			name:     "vertical scrolling to end with no wrapping, then back 1 entry",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(1)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"TwoTwo    ", "ThreeThree", "FourFourFo"},
		},
		{
			name:     "vertical scrolling to end with no wrapping, then back beyond start",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollUpMsg(10)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"One       ", "TwoTwo    ", "ThreeThree"},
		},
		{
			name:     "vertical scrolling to end with no wrapping, then scroll to beginning",
			msgs:     []tea.Msg{ScrollBottomMsg{}, ScrollTopMsg{}},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"One       ", "TwoTwo    ", "ThreeThree"},
		},
		// both horizontal and vertical scrolling
		{
			name:    "both horizontal and vertical scrolling with warp",
			msgs:    []tea.Msg{ScrollRightMsg(3), ScrollDownMsg(1)},
			entries: []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll: true,
			expected: []string{
				"TwoTwo    ",
				"ThreeThree",
				"Three     ",
			},
		},
		{
			name:     "both horizontal and vertical scrolling without warp",
			msgs:     []tea.Msg{ScrollRightMsg(3), ScrollDownMsg(1)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"Two       ", "eeThreeThr", "rFourFourF"},
		},
		{
			name:    "both vertical and horizontal scrolling with warp",
			msgs:    []tea.Msg{ScrollDownMsg(1), ScrollRightMsg(3)},
			entries: []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour"},
			vScroll: true,
			expected: []string{
				"TwoTwo    ",
				"ThreeThree",
				"Three     ",
			},
		},
		{
			name:     "both vertical and horizontal scrolling without warp",
			msgs:     []tea.Msg{ScrollDownMsg(1), ScrollRightMsg(3)},
			entries:  []string{"One", "TwoTwo", "ThreeThreeThree", "FourFourFourFour", "FiveFiveFiveFive"},
			vScroll:  true,
			hScroll:  true,
			expected: []string{"Two       ", "eeThreeThr", "rFourFourF"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := []ComponentOption{}
			if tc.hScroll {
				options = append(options, HorizontalScrollable())
			}
			if tc.vScroll {
				options = append(options, VerticalScrollable())
			}
			comp := NewComponent(10, options...)
			for _, entry := range tc.entries {
				comp.Append(entry)
			}
			comp.Update(ResizeMsg{Width: 10, Height: 3})
			for _, msg := range tc.msgs {
				comp.Update(msg)
			}

			// check expect
			if len(tc.expected) != 3 {
				t.Fatalf("expected must have exactly 3 lines, got %d", len(tc.expected))
			}
			for _, line := range tc.expected {
				if len(line) != 10 {
					t.Fatalf("each expected line must have exactly 10 characters, got %d in line '%s'", len(line), line)
				}
			}

			assert.Equal(t, strings.Join(tc.expected, "\n"), strings.TrimRight(comp.View(), "\n"))

			assert.Equal(t, "", IsThisTopping(ToppingTestSpec{
				Width:  10,
				Height: 3,
				Model:  comp,
			}))
		})
	}
}

func TestComponent_EdgeCase(t *testing.T) {
	t.Run("height larger than entries, no warp", func(t *testing.T) {
		comp := NewComponent(1, VerticalScrollable())
		comp.Append("One")
		comp.Update(ResizeMsg{Width: 5, Height: 2})
		assert.Equal(t, "One  \n     ", comp.View())
	})

	t.Run("height larger than entries, with warp", func(t *testing.T) {
		comp := NewComponent(1)
		comp.Append("OneTwo")
		comp.Update(ResizeMsg{Width: 3, Height: 3})
		assert.Equal(t, "One\nTwo\n   ", comp.View())
	})
}
