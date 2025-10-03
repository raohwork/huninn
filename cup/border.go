// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

func DefaultBorderConfig() BorderConfig {
	return BorderConfig{
		Left:              true,
		Top:               true,
		Right:             true,
		Bottom:            true,
		VerticalLine:      '│',
		HorizontalLine:    '─',
		TopLeftCorner:     '┌',
		TopRightCorner:    '┐',
		BottomLeftCorner:  '└',
		BottomRightCorner: '┘',
	}
}

type BorderConfig struct {
	Left, Top, Right, Bottom bool
	VerticalLine             rune
	HorizontalLine           rune
	TopLeftCorner            rune
	TopRightCorner           rune
	BottomLeftCorner         rune
	BottomRightCorner        rune
}

func (bc *BorderConfig) size() (v, h int) {
	if bc.VerticalLine != 0 {
		v = tapioca.RuneWidth(bc.VerticalLine)
	}
	if bc.HorizontalLine != 0 {
		h = tapioca.RuneWidth(bc.HorizontalLine)
	}
	return
}

type BorderedBox struct {
	BorderConfig
	inner tea.Model

	// width caches, computed only once in init
	vLineWidth int // width of single vertical line rune
	hLineWidth int // width of single horizontal line rune
	// width of corners
	lt, lb, rt, rb int

	// cached values
	size int
	w, h int

	hasError bool
	wReserve int // width reserved for inner component
	hReserve int // height reserved for inner component

	// true if we have to left 1 char at right
	// ex: width 13 with wide character border
	reminder bool
}

func NewBorderedBox(inner tea.Model) *BorderedBox {
	return &BorderedBox{
		BorderConfig: DefaultBorderConfig(),
		inner:        inner,
	}
}

func (b *BorderedBox) Init() tea.Cmd {
	b.vLineWidth, b.hLineWidth = b.BorderConfig.size()
	b.lt = tapioca.RuneWidth(b.TopLeftCorner)
	b.lb = tapioca.RuneWidth(b.BottomLeftCorner)
	b.rt = tapioca.RuneWidth(b.TopRightCorner)
	b.rb = tapioca.RuneWidth(b.BottomRightCorner)
	return b.inner.Init()
}

func (b *BorderedBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tapioca.ResizeMsg:
		b.computeSize(msg.Width, msg.Height)
		b.inner, cmd = b.inner.Update(tapioca.ResizeMsg{
			Width:  b.wReserve,
			Height: b.hReserve,
		})
	default:
		b.inner, cmd = b.inner.Update(msg)
	}
	return b, cmd
}

func (b *BorderedBox) computeSize(width, height int) {
	b.size = width * height
	b.wReserve, b.hReserve = width, height
	if b.Left {
		b.wReserve -= b.vLineWidth
	}
	if b.Right {
		b.wReserve -= b.vLineWidth
	}
	if b.Top {
		b.hReserve -= b.hLineWidth
	}
	if b.Bottom {
		b.hReserve -= b.hLineWidth
	}

	b.hReserve = max(0, b.hReserve)
	b.wReserve = max(0, b.wReserve)

	// if width is odd number with wide character border, we have to leave 1 char at right
	b.reminder = b.wReserve%2 == 1 && b.hLineWidth == 2

	b.hasError = b.wReserve < 2 || b.hReserve < 1
}

func (b *BorderedBox) View() string {
	if b.hasError {
		return "too small"
	}

	buf := &strings.Builder{}
	buf.Grow(b.size)

	if b.Top {
		w := b.wReserve
		if b.Left {
			buf.WriteRune(b.TopLeftCorner)
		}
		for w >= b.hLineWidth {
			buf.WriteRune(b.HorizontalLine)
			w -= b.hLineWidth
		}
		if b.Right {
			buf.WriteRune(b.TopRightCorner)
		}
		if b.reminder {
			buf.WriteRune(' ')
		}
		buf.WriteRune('\n')
	}

	// render inner component
	innerView := strings.Split(strings.TrimRight(b.inner.View(), "\n"), "\n")
	for i := 0; i < b.hReserve; i++ {
		if b.Left {
			buf.WriteRune(b.VerticalLine)
		}
		buf.WriteString(innerView[i])
		if b.Right {
			buf.WriteRune(b.VerticalLine)
		}
		if b.reminder {
			buf.WriteRune(' ')
		}
		buf.WriteRune('\n')
	}

	if b.Bottom {
		// render bottom line
		w := b.wReserve
		if b.Left {
			buf.WriteRune(b.BottomLeftCorner)
		}
		for w >= b.hLineWidth {
			buf.WriteRune(b.HorizontalLine)
			w -= b.hLineWidth
		}
		if b.Right {
			buf.WriteRune(b.BottomRightCorner)
		}
		if b.reminder {
			buf.WriteRune(' ')
		}
	}

	return buf.String()
}
